package main

import (
	"log"
	"testing"
	"time"
)

const defaultStackSize = 1024

func newBasicCrawlerWithDomainPolicy(userAgent string, domainNames []string, duration time.Duration) crawler {
	c := newBasicCrawler()
	p := newCheckDomainPolicy()
	initCheckDomainPolicy(p, domainNames)

	fe := defaultFetcher(p)
	fr := newStackFrontier(defaultStackSize)
	s := newInMemoryURLStore()
	initBasicCrawler(c, domainNames, fe, p, fr, duration, s)
	return c
}

func TestBasicCrawlerWithNoLinks(t *testing.T) {
	//TODO: Does not accept subdomains
	startMock()
	defer endMock()
	domainNames := []string{"domainGGC.com", "www.domainGGC.com"}
	setUpFakePage("http://www.domainGGC.com/", "testFiles/home.html")
	setUpFakePage("http://domainGGC.com/", "testFiles/home.html")
	setUpFakePage("http://www.domainGGC.com/page1/", "testFiles/page1.html")
	setUpFakePage("http://domainGGC.com/page1/", "testFiles/page1.html")
	oneSeconds := time.Duration(1) * time.Second
	c := newBasicCrawlerWithDomainPolicy("GGC", domainNames, oneSeconds)
	nilSitemap, error := c.crawl()
	if error != nil {
		t.Fail()
	}
	log.Printf("%s", nilSitemap)
}
