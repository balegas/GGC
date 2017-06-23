package main

import (
	"io"
	"log"
	"net/url"
	"time"
)

// NBatches Crawler uses go routines to crawl the internet.
// Each go routine receives a batch of urls to process and returns a slice of
// found urls.
// Go routines share a cache of visited urls, that must be thread-safe.
type nBatchesCrawler struct {
	finishTime  time.Time
	fetcher     fetcher
	rules       accessPolicyChecker
	frontier    urlFrontier
	store       urlStore
	maxWorkers  int
	currWorkers int
}

func newNBatchesCrawler() *nBatchesCrawler {
	return &nBatchesCrawler{}
}

func initNBatchesCrawler(c *nBatchesCrawler, seed []string, fet fetcher, rules accessPolicyChecker, uf urlFrontier, duration time.Duration, s urlStore, maxRoutines int) {
	c.rules = rules
	c.finishTime = time.Now().Add(duration)
	c.fetcher = fet
	c.frontier = uf
	for _, domain := range seed {
		domainURL, _ := url.Parse("http://" + domain + "/")
		curl, _ := getCanonicalURLString("/", domainURL)
		c.frontier.addURLString(curl) // Causes redirect if https.
	}
	c.store = s
	c.maxWorkers = maxRoutines
	c.currWorkers = 0
}

// Spawns workers to crawl webpages while there are urls left in the frontier.
// Workers return a batch of urls ready to put on the frontier.
// While no more workers can be spawned, add results to the frontier.

// NOTE: Might underperform if results queue gets full and frontier gets empty.
// The objective is to ensure is that the frontier is sufficiently occupied, to
// ensure that the workers always have enough work.
func (c *nBatchesCrawler) crawl() (sitemap, error) {
	var s sitemap
	var foundURLs int32
	foundURLs = 0

	newURLsC := make(chan []string, urlChanBufferSize)
	signalC := make(chan bool)

	for (!c.frontier.isEmpty() || c.currWorkers != 0) && !c.isTimeout() {

		//Spawn multiple routines
		c.spawnRoutines(newURLsC, signalC)

		finishedBatch := false
		for !finishedBatch {
			select {
			// New URL arrived
			case curls := <-newURLsC:
				for _, ci := range curls {
					c.frontier.addURLString(ci)
					foundURLs++
				}
			case <-signalC:
				finishedBatch = true
				c.currWorkers--
				log.Printf("FINISHED ROUTINE; CURRENT WORKERS: %v; SIZE OF FRONTIER: %v", c.currWorkers, c.frontier.size())
			}
		}

	}

	log.Printf("Finished. Found %v urls.", foundURLs)
	close(newURLsC)
	close(signalC)
	return s, nil
}

func (c *nBatchesCrawler) spawnRoutines(newURLsC chan []string, signalC chan bool) {
	for !c.frontier.isEmpty() && c.currWorkers < c.maxWorkers {
		pendingURLsC := make(chan string, urlChanBufferSize)
		c.enqueueMultiple(pendingURLsC)
		go c.processURLs(pendingURLsC, newURLsC, signalC)
		c.currWorkers++
	}
}

func (c *nBatchesCrawler) enqueueMultiple(pendingURLsC chan string) {
	filled := false
	for !c.frontier.isEmpty() && !filled {
		url, err := c.frontier.nextURLString()
		if err != nil {
			log.Fatal("Error dequeuing.")
		}

		select {
		case pendingURLsC <- url:
		default:
			// Filled queue. Put value back and continue
			c.frontier.addURLString(url)
			filled = true

		}
	}
	close(pendingURLsC)
}

//TODO: Consider make routine independent of Crawler
func (c *nBatchesCrawler) processURLs(pendingURLsC chan string, newURLsC chan []string,
	signalC chan bool) {
	visitedURLs := 0
	found := 0
	for curl := range pendingURLsC {
		visitedURLs++
		//log.Printf("NEXT  %s", curl)
		if c.canProcess(curl) {
			nextURL, _ := toURL(curl)
			newURLs, _, err := c.findURLLinksGetBody(nextURL)

			if err != nil {
				//log.Printf("Error processing page: %s", err)
				c.storeURL(curl, []byte{})
				continue
			} else {
				// bodyInBytes, _ := ioutil.ReadAll(body)
				// c.storeURL(curl, bodyInBytes)
				filteredCurl := make([]string, 0)
				for _, url := range newURLs {
					found++
					//log.Printf("Push new %s", url)
					curli, _ := getCanonicalURLString(url, nextURL)
					if c.canProcess(curli) && !c.seen(curli) {
						c.storeURL(curli, []byte{})
						filteredCurl = append(filteredCurl, curli)
					}
				}
				newURLsC <- filteredCurl
				//log.Printf("Received %v URLS.", len(newURLs))
			}
		}
	}
	log.Printf("VISITED: %v; FOUND: %v", visitedURLs, found)
	signalC <- true
}

// Checks if access policy allows this URL.
func (c *nBatchesCrawler) canProcess(curl string) bool {
	// This check is being done in two diff. places, but seems more efficient
	// this way
	return c.rules.checkURL(curl)
}

// Checks if url has been seen (might have not been processed yet)
func (c *nBatchesCrawler) seen(curl string) bool {
	if _, exists := c.store.get(curl); exists {
		return true

	}
	return false
}

func (c *nBatchesCrawler) findURLLinksGetBody(url *url.URL) ([]string, io.Reader, error) {
	content, err := c.fetcher.getURLContent(url)
	//TODO: push content/or content hash to store
	if err != nil {
		//log.Printf("error fetching content from url: %s : %s", url, err)
		return nil, nil, err
	}
	// Reading the value twice :/
	return getAllTagAttr(crawlTags, content.Body), content.Body, nil
}

func (c *nBatchesCrawler) storeURL(curl string, body []byte) {
	c.store.put(curl, body)
}

func (c *nBatchesCrawler) isTimeout() bool {
	return c.finishTime.Before(time.Now())
}
