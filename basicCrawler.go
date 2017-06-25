package main

import (
	"io/ioutil"
	"log"
	"time"
)

// basicCrawler is a single-threaded web crawler with support for generic
// urlFrontier, url fetcher and cache, and access policy rules. It processes
// urls sequentially.
type basicCrawler struct {
	crawlerInternals
}

func newBasicCrawler() *basicCrawler {
	return &basicCrawler{}
}

func initBasicCrawler(c *basicCrawler, seed []string, fet fetcher,
	rules accessPolicy, uf urlFrontier, duration time.Duration, s urlStore) {
	initCommonAttributes(&c.crawlerInternals, seed, fet, rules, uf, duration, s)
}

func (c *crawlerInternals) Crawl() (sitemap, error) {
	var s sitemap
	foundURLs := 0

	for !c.frontier.isEmpty() && !c.isTimeout() {
		curl, err := c.frontier.nextURLString()
		if err != nil {
			log.Fatal("Error dequeuing.")
		}
		if c.canProcess(curl) {
			nextURL, _ := toURL(curl)
			newURLs, body, err := c.findURLLinksGetBody(nextURL)
			receivedURLs := len(newURLs)
			if err != nil {
				log.Printf("Error processing page: %s", err)
				c.storeURL(curl, []byte{})
				continue
			} else {
				bodyInBytes, _ := ioutil.ReadAll(body)
				c.storeURL(curl, bodyInBytes)
				for _, u := range newURLs {
					curli, _ := getCanonicalURLString(u, nextURL)
					if c.canProcess(curli) && !c.seen(curli) {
						foundURLs++
						c.storeURL(curli, []byte{})
						c.frontier.addURLString(curli)

					}
				}
			}
			log.Printf("Received %v URLS.", receivedURLs)
		}
	}
	log.Printf("Finished. Found %v urls.", foundURLs)
	return s, nil
}
