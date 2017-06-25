package main

import (
	"log"
	"time"
)

// NBatches Crawler uses multiple go routines to crawl the internet.
// Each go routine receives a batch of urls to process and returns a slice of
// found urls.
// Go routines share a cache of visited urls (must be thread-safe).
type nBatchesCrawler struct {
	crawlerInternals
	maxWorkers  int
	currWorkers int
}

func newNBatchesCrawler() *nBatchesCrawler {
	return &nBatchesCrawler{}
}

func initNBatchesCrawler(c *nBatchesCrawler, seed []string, fet fetcher,
	rules accessPolicy, uf urlFrontier, duration time.Duration, s urlStore,
	maxRoutines int) {
	initCommonAttributes(&c.crawlerInternals, seed, fet, rules, uf, duration, s)
	c.maxWorkers = maxRoutines
	c.currWorkers = 0
}

// Spawns workers to crawl webpages while there are urls left in the frontier.
// Workers return a batch of urls ready to put on the frontier.
// While no more workers can be spawned, add results to the frontier.

// NOTE: Might underperform if results queue gets full and frontier gets empty.
// The objective is to ensure is that the frontier is sufficiently occupied, to
// ensure that the workers always have enough work. Can tune number of threads
// and the buffer size of channels.
func (c *nBatchesCrawler) Crawl() (sitemap, error) {
	var s sitemap
	foundURLs := 0

	newURLsC := make(chan []string, urlChanBufferSize)
	signalC := make(chan bool)

	for (!c.frontier.isEmpty() || c.currWorkers != 0) && !c.isTimeout() {

		c.spawnRoutines(newURLsC, signalC)

		finishedBatch := false
		for !finishedBatch {
			select {
			case curls := <-newURLsC:
				for _, ci := range curls {
					c.frontier.addURLString(ci)
					foundURLs++
				}
			case <-signalC:
				finishedBatch = true
				c.currWorkers--
				log.Printf("FINISHED ROUTINE; CURRENT WORKERS: %v; SIZE OF FRONTIER: %v",
					c.currWorkers, c.frontier.size())
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

func (c *nBatchesCrawler) processURLs(pendingURLsC chan string,
	newURLsC chan []string, signalC chan bool) {
	visitedURLs := 0
	found := 0
	for curl := range pendingURLsC {
		//log.Printf("NEXT  %s", curl)
		if c.canProcess(curl) {
			visitedURLs++
			nextURL, _ := toURL(curl)
			newURLs, _, err := c.findURLLinksGetBody(nextURL)

			if err != nil {
				c.storeURL(curl, []byte{})
				continue
			} else {
				// bodyInBytes, _ := ioutil.ReadAll(body)
				// c.storeURL(curl, bodyInBytes)
				filteredCurl := make([]string, 0)
				for _, url := range newURLs {
					found++
					curli, _ := getCanonicalURLString(url, nextURL)
					if c.canProcess(curli) && !c.seen(curli) {
						c.storeURL(curli, []byte{})
						filteredCurl = append(filteredCurl, curli)
					}
				}
				newURLsC <- filteredCurl
				log.Printf("Received %v URLS.", len(newURLs))
			}
		}
	}
	log.Printf("BATCH RESULTS: VISITED: %v; FOUND: %v", visitedURLs, found)
	signalC <- true
}
