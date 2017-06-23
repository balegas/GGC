package main

import (
	"log"
	"os"
	"strconv"
	"time"
)

//Should mark header tags to avoid fetching those files
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
	//c := newProducerConsumerCrawler()
	c := newNBatchesCrawler()

	p := newCheckDomainPolicy()
	initCheckDomainPolicy(p, domainNames)

	fe := defaultFetcher(p)
	fr := newQueueFrontier(defaultStackSize)
	s := newInMemoryURLStore()

	//initBasicCrawler(c, domainNames, fe, p, fr, duration, s)
	//initProducerConsumerCrawler(c, domainNames, fe, p, fr, duration, s)
	initNBatchesCrawler(c, domainNames, fe, p, fr, duration, s, 4)

	nilSitemap, _ := c.crawl()
	log.Printf("%s", nilSitemap)
}

/*
func main() {
	startMock()
	defer endMock()
	domainNames := []string{"domainGGC.com", "www.domainGGC.com"}
	setUpFakePage("http://www.domainGGC.com/", "testFiles/home.html")
	setUpFakePage("http://domainGGC.com/", "testFiles/home.html")
	setUpFakePage("http://www.domainGGC.com/page1/", "testFiles/page1.html")
	setUpFakePage("http://domainGGC.com/page1/", "testFiles/page1.html")
	oneSeconds := time.Duration(1) * time.Second

	//bC := newBasicCrawlerWithDomainPolicy("GGC", domainNames, oneSeconds)
	pcC := newProducerConsumerWithDomainPolicy("GGC", domainNames, oneSeconds)

	C := []crawler{pcC}

	//C := []crawler{bC, pcC}

	for _, c := range C {
		nilSitemap, error := c.crawl()
		if error != nil {

		}
		log.Printf("%s", nilSitemap)
	}

}
*/
