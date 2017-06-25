package main

import (
	"io"
	"log"
	"net/url"
	"time"
)

//Common atributes of a crawler.
type crawlerInternals struct {
	finishTime time.Time
	fetcher    fetcher
	rules      accessPolicy
	frontier   urlFrontier
	store      urlStore
	sitemap    sitemap
}

func initCommonAttributes(c *crawlerInternals, seed []string,
	fet fetcher, rules accessPolicy, uf urlFrontier, duration time.Duration,
	s urlStore, sm sitemap) {
	c.rules = rules
	c.finishTime = time.Now().Add(duration)
	c.fetcher = fet
	c.frontier = uf
	for _, domain := range seed {
		domainURL, _ := url.Parse("http://" + domain + "/") // Causes redirect if https.
		curl, _ := getCanonicalURLString("/", domainURL)
		c.frontier.addURLString(curl)
	}
	c.store = s
	c.sitemap = sm
}

// Checks if access policy allows this URL.
func (c *crawlerInternals) canProcess(curl string) bool {
	return c.rules.checkURL(curl)
}

// Checks if url has been added to cache. I.e. it has been visited, or is in
// urlFrontier
func (c *crawlerInternals) seen(curl string) bool {
	if _, exists := c.store.get(curl); exists {
		return true

	}
	return false
}

// Store url in cache.
func (c *crawlerInternals) storeURL(curl string, body []byte) {
	c.store.put(curl, body)
}

// Check execution timeout.
func (c *crawlerInternals) isTimeout() bool {
	return c.finishTime.Before(time.Now())
}

// Find urls in a page and returns the body of the document.
func (c *crawlerInternals) findURLLinksGetBody(url *url.URL) ([]string,
	io.Reader, error) {
	content, err := c.fetcher.getURLContent(url)
	if err != nil {
		log.Printf("error fetching content from url: %s : %s", url, err)
		return nil, nil, err
	}
	return getAllTagAttr(crawlTags, content.Body), content.Body, nil
}

func (c *crawlerInternals) printSitemap(s sitemap, out io.Writer) {
	s.printSitemap(out)
}
