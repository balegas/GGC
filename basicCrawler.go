package main

import (
	"bytes"
	"io"
	"log"
	"net/url"
	"time"
)

var crawlTags = map[string]string{
	// Header tags
	"link":   "href",
	"script": "source",
	// Body tags
	"a": "href",
	// others?
}

// basicCrawler is a simple web crawler that is single-threaded, uses a
// stack-based url frontier, and applies no strategies to choose the urls to
// visit.
type basicCrawler struct {
	finishTime time.Time // make default time = math.MaxInt64
	fetcher    fetcher
	rules      accessPolicyChecker
	frontier   urlFrontier
	store      urlStore
}

func newBasicCrawler() *basicCrawler {
	return &basicCrawler{}
}

func initBasicCrawler(c *basicCrawler, seed []string, fet fetcher, rules accessPolicyChecker, uf urlFrontier, duration time.Duration, s urlStore) {
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

func (c *basicCrawler) isTimeout() bool {
	log.Printf("Checked timeout %v", c.finishTime.Before(time.Now()))
	return c.finishTime.Before(time.Now())
}

// Crawl a webdomain
func (c *basicCrawler) crawl() (sitemap, error) {
	// Check if you're doint pointers right in initbasicCrawler
	var s sitemap
	for !c.frontier.isEmpty() && !c.isTimeout() {
		curl, err := c.frontier.nextURLString()
		if err != nil {
			log.Fatal("Error dequeuing.")
		}
		if c.canProcess(curl) {
			nextURL, _ := toURL(curl)
			newURLs, body, err := c.findURLLinksGetBody(nextURL)
			log.Printf("newURLS %s", newURLs)
			if err != nil {
				log.Printf("Error processing page: %s", err)
				//Mark as visited
				c.store.put(curl, bytes.NewReader([]byte{}))
				continue
			} else {
				c.store.put(curl, body)
				for _, u := range newURLs {
					curli, _ := getCanonicalURLString(u, nextURL)
					if !c.visited(curli) && c.canProcess(curli) {
						log.Printf("Added to frontier %s", curli)
						c.frontier.addURLString(curli)
					}
				}
			}
			c.markProcessed(curl)
		}
	}
	log.Printf("Finished")
	return s, nil
}
func (c *basicCrawler) canProcess(curl string) bool {
	//TODO: Add cache verification
	return c.rules.checkURL(curl)
}

func (c *basicCrawler) visited(curl string) bool {
	//TODO: Integrate URL Store
	return false
}

func (c *basicCrawler) storeURL(curl string, content string) bool {
	//TODO: Integrate URL Store.
	//TODO: Hash content
	return false
}

func (c *basicCrawler) findURLLinksGetBody(url *url.URL) ([]string, io.Reader, error) {
	content, err := c.fetcher.getURLContent(url)
	//TODO: push content/or content hash to store
	if err != nil {
		log.Printf("error fetching content from url: %s : %s", url, err)
		return nil, nil, err
	}
	// Reading the value twice :/
	return getAllTagAttr(crawlTags, content.Body), content.Body, nil
}

func (c *basicCrawler) markProcessed(curl string) []string {
	//TODO
	return nil
}
