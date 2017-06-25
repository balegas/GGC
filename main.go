package main

import (
	"log"
	"os"
	"strconv"
	"time"
)

//TODO: Should mark header tags to avoid fetching those files
var crawlTags = map[string]string{
	// Header tags
	"link":   "href",
	"script": "source",
	// Body tags
	"a": "href",
	// others?
}

const urlChanBufferSize = 10

func main() {

	domainNames := os.Args[1 : len(os.Args)-1]
	durationInt, _ := strconv.Atoi(os.Args[len(os.Args)-1])
	duration := time.Duration(durationInt) * time.Second

	log.Printf("DomainNames: %v, duration: %v", domainNames, durationInt)

	//c := newBasicCrawler()
	c := newProducerConsumerCrawler()
	//c := newNBatchesCrawler()

	p := newCheckDomainPolicy()
	initCheckDomainPolicy(p, domainNames)

	fe := defaultFetcher(p)
	fr := newQueueFrontier(defaultStackSize)
	s := newInMemoryURLStore()

	//initBasicCrawler(c, domainNames, fe, p, fr, duration, s)
	initProducerConsumerCrawler(c, domainNames, fe, p, fr, duration, s)
	//initNBatchesCrawler(c, domainNames, fe, p, fr, duration, s, 4)

	nilSitemap, _ := c.Crawl()
	log.Printf("%s", nilSitemap)
}
