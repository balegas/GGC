package main

import (
	"log"
	"time"
)

// basicCrawler is a simple web crawler that is single-threaded, uses a
// stack-based url frontier, and applies no strategies to choose the urls to
// visit.
type basicCrawler struct {
	domainName string    // default seed
	finishTime time.Time // make default time = math.MaxInt64
	fetcher    fetcher
	rules      accessPolicyChecker
	frontier   urlFrontier
}

func newBasicCrawler() *basicCrawler {
	return &basicCrawler{}
}

func initBasicCrawler(c *basicCrawler, domainName string, rules accessPolicyChecker, f urlFrontier, duration time.Duration) {
	c.domainName = domainName
	c.rules = rules
	c.finishTime = time.Now().Add(duration)
	c.frontier = f
	c.frontier.addURL("http://" + domainName + "/") // Causes redirect if https.

}

func (c *basicCrawler) isTimeout() bool {
	return c.finishTime.Before(time.Now())
}

// Crawl a webdomain
func (c *basicCrawler) crawl() (sitemap, error) {
	// Check if you're doint pointers right in initbasicCrawler
	var s sitemap
	for !c.frontier.isEmpty() || c.isTimeout() {
		nextURL, err := c.frontier.nextURL()
		if err != nil {
			log.Fatal("Error dequeuing.")
		}
		if c.canProcess(nextURL) {
			findURLLinks(nextURL)
			markProcessed(nextURL)
		}
	}
	return s, nil
}
func (c *basicCrawler) canProcess(url string) bool {
	// TODO: implement robots.txt policy
	// Check not visited.
	return true
}

func findURLLinks(url string) []string {
	//TODO
	return nil
}

func markProcessed(url string) []string {
	//TODO
	return nil
}
