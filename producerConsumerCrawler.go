package main

import (
	"log"
	"time"
)

// ProducerConsumer is a Web crawler that runs independent routines for finding
// new URLs and managing the urlFrontier.
// Its a simpler version of the nBatchesCrawler. The purpose of developing this
// was to to exercise the concurrency primitives of the language.
type producerConsumerCrawler struct {
	crawlerInternals
}

func newProducerConsumerCrawler() *producerConsumerCrawler {
	return &producerConsumerCrawler{}
}

func initProducerConsumerCrawler(c *producerConsumerCrawler, seed []string,
	fet fetcher, rules accessPolicy, uf urlFrontier,
	duration time.Duration, s urlStore, sm sitemap) {
	initCommonAttributes(&c.crawlerInternals, seed, fet, rules, uf, duration, s, sm)
}

// Crawl function reads and writes to frontier in isolation.
func (c *producerConsumerCrawler) Crawl() (sitemap, error) {
	foundURLs := 0
	finishedBatch := true

	newURLsC := make(chan string, urlChanBufferSize)
	signalC := make(chan bool)

	for (!c.frontier.isEmpty() || !finishedBatch) && !c.isTimeout() {

		pendingURLsC := make(chan string, urlChanBufferSize)
		c.enqueueMultiple(pendingURLsC)
		finishedBatch := false

		go c.processURLs(pendingURLsC, newURLsC, signalC)

		// Wait for a batch to be complete before enqueueing new.
		// Process new urls meanwhile.
		for !finishedBatch {
			select {
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
	return c.sitemap, nil
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
			c.frontier.addURLString(url)
			filled = true

		}
	}
	close(pendingURLsC)
}

func (c *producerConsumerCrawler) processURLs(pendingURLsC, newURLsC chan string,
	signalC chan bool) {
	visitedURLs := 0
	found := 0

	for curl := range pendingURLsC {
		if c.canProcess(curl) {
			c.sitemap.addURL(curl)
			visitedURLs++
			nextURL, _ := toURL(curl)
			newURLs, _, err := c.findURLLinksGetBody(nextURL)

			if err != nil {
				c.storeURL(curl, []byte{})
				continue
			} else {
				// bodyInBytes, _ := ioutil.ReadAll(body)
				// c.storeURL(curl, bodyInBytes)
				for _, url := range newURLs {
					found++
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
	log.Printf("BATCH RESULTS: VISITED: %v; FOUND: %v", visitedURLs, found)
	signalC <- true
}
