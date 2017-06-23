package main

import (
	"io"
	"io/ioutil"
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

// basicCrawler is a single-threaded web crawler with support for generic
// urlFrontier, url fetcher and cache, and access policy rules.
type basicCrawler struct {
	finishTime time.Time
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

// Crawl a webdomain
func (c *basicCrawler) crawl() (sitemap, error) {
	var s sitemap
	for !c.frontier.isEmpty() && !c.isTimeout() {
		curl, err := c.frontier.nextURLString()
		log.Printf("NEXT  %s", curl)
		if err != nil {
			log.Fatal("Error dequeuing.")
		}
		if c.canProcess(curl) {
			nextURL, _ := toURL(curl)
			newURLs, body, err := c.findURLLinksGetBody(nextURL)
			receivedURLs := len(newURLs)
			foundURLs := 0
			//log.Printf("newURLS %s", newURLs)
			if err != nil {
				log.Printf("Error processing page: %s", err)
				c.storeURL(curl, []byte{})
				continue
			} else {
				bodyInBytes, _ := ioutil.ReadAll(body)
				c.storeURL(curl, bodyInBytes)
				for _, u := range newURLs {
					//Must make url cannonical before checking that can be processed.
					curli, _ := getCanonicalURLString(u, nextURL)
					if c.canProcess(curli) && !c.seen(curli) {
						foundURLs++
						c.storeURL(curli, []byte{})
						c.frontier.addURLString(curli)

					}
				}
			}
			log.Printf("Received %v URLS. New: %v", receivedURLs, foundURLs)
		}
	}
	log.Printf("Finished")
	return s, nil
}

// Checks if access policy allows this URL.
func (c *basicCrawler) canProcess(curl string) bool {
	// This check is being done in two diff. places, but seems more efficient
	// this way
	return c.rules.checkURL(curl)
}

// Checks if url has been seen (might have not been processed yet)
func (c *basicCrawler) seen(curl string) bool {
	if _, exists := c.store.get(curl); exists {
		return true

	}
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

func (c *basicCrawler) storeURL(curl string, body []byte) {
	c.store.put(curl, body)
}

func (c *basicCrawler) isTimeout() bool {
	return c.finishTime.Before(time.Now())
}

/*
//TODO: Make a single main function
func main() {

	domainNames := os.Args[1:]
	TenSeconds := time.Duration(20) * time.Second

	c := newBasicCrawler()
	p := newCheckDomainPolicy()
	initCheckDomainPolicy(p, domainNames)

	fe := defaultFetcher(p)
	fr := newStackFrontier(defaultStackSize)
	s := newInMemoryURLStore()
	initBasicCrawler(c, domainNames, fe, p, fr, TenSeconds, s)

	nilSitemap, _ := c.crawl()
	log.Printf("%s", nilSitemap)
}
*/
