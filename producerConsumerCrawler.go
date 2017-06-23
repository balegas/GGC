package main

import (
	"io"
	"log"
	"net/url"
	"time"
)

// ProducerConsumer Crawler has two cooperatingg goroutines.
// One manages the urlFrontier and the other processes the requests.
// Its a simpler version of the multiple-workers approach and main purpose
// is to exercise the concurrency primitives of the language.
type producerConsumerCrawler struct {
	finishTime time.Time
	fetcher    fetcher
	rules      accessPolicyChecker
	frontier   urlFrontier
	store      urlStore
}

func newProducerConsumerCrawler() *producerConsumerCrawler {
	return &producerConsumerCrawler{}
}

func initProducerConsumerCrawler(c *producerConsumerCrawler, seed []string, fet fetcher, rules accessPolicyChecker, uf urlFrontier, duration time.Duration, s urlStore) {
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
}

// Crawl function reads and writes to frontier in isolation.
func (c *producerConsumerCrawler) crawl() (sitemap, error) {
	var s sitemap
	foundURLs := 0
	finishedBatch := true

	newURLsC := make(chan string, urlChanBufferSize)
	signalC := make(chan bool)

	for (!c.frontier.isEmpty() || !finishedBatch) && !c.isTimeout() {

		//Fill the channel for processing
		pendingURLsC := make(chan string, urlChanBufferSize)
		c.enqueueMultiple(pendingURLsC)
		finishedBatch := false

		go c.processURLs(pendingURLsC, newURLsC, signalC)

		// Wait for a batch to be complete before enqueueing new.
		// Process new urls meanwhile.
		for !finishedBatch {
			select {
			// New URL arrived
			case curli := <-newURLsC:
				foundURLs++
				c.frontier.addURLString(curli)

			case <-signalC:
				finishedBatch = true
			}
		}

	}

	log.Printf("Finished. Found %v urls.", foundURLs)
	close(newURLsC)
	close(signalC)
	return s, nil
}

func (c *producerConsumerCrawler) enqueueMultiple(pendingURLsC chan string) {
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
func (c *producerConsumerCrawler) processURLs(pendingURLsC, newURLsC chan string,
	signalC chan bool) {
	for curl := range pendingURLsC {
		//log.Printf("NEXT  %s", curl)
		if c.canProcess(curl) {
			nextURL, _ := toURL(curl)
			newURLs, _, err := c.findURLLinksGetBody(nextURL)

			if err != nil {
				c.storeURL(curl, []byte{})
				continue
			} else {
				// bodyInBytes, _ := ioutil.ReadAll(body)
				// c.storeURL(curl, bodyInBytes)
				for _, url := range newURLs {
					curli, _ := getCanonicalURLString(url, nextURL)
					if c.canProcess(curli) && !c.seen(curli) {
						c.storeURL(curli, []byte{})
						newURLsC <- curli
					}
				}
				log.Printf("Received %v URLS.", len(newURLs))
			}
		}
	}
	signalC <- true
}

// Checks if access policy allows this URL.
func (c *producerConsumerCrawler) canProcess(curl string) bool {
	// This check is being done in two diff. places, but seems more efficient
	// this way
	return c.rules.checkURL(curl)
}

// Checks if url has been seen (might have not been processed yet)
func (c *producerConsumerCrawler) seen(curl string) bool {
	if _, exists := c.store.get(curl); exists {
		return true

	}
	return false
}

func (c *producerConsumerCrawler) findURLLinksGetBody(url *url.URL) ([]string, io.Reader, error) {
	content, err := c.fetcher.getURLContent(url)
	//TODO: push content/or content hash to store
	if err != nil {
		return nil, nil, err
	}
	// Reading the value twice :/
	return getAllTagAttr(crawlTags, content.Body), content.Body, nil
}

func (c *producerConsumerCrawler) storeURL(curl string, body []byte) {
	c.store.put(curl, body)
}

func (c *producerConsumerCrawler) isTimeout() bool {
	return c.finishTime.Before(time.Now())
}
