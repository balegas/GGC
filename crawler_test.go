package main

import (
	"log"
	"testing"
	"time"

	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

//TODO: Compare results of different Crawlers.
//TODO: Define test checks.
func TestCrawlersMock(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	domainNames := []string{"domainGGC.com", "www.domainGGC.com"}
	setUpFakePage("http://www.domainGGC.com/", "testFiles/home.html")
	setUpFakePage("http://domainGGC.com/", "testFiles/home.html")
	setUpFakePage("http://www.domainGGC.com/page1/", "testFiles/page1.html")
	setUpFakePage("http://domainGGC.com/page1/", "testFiles/page1.html")
	oneSeconds := time.Duration(1) * time.Second

	bC := newBasicCrawlerWithDomainPolicy("GGC", domainNames, oneSeconds)
	pcC := newProducerConsumerWithDomainPolicy("GGC", domainNames, oneSeconds)
	nbC := newNBatchesCrawlerWithDomainPolicy("GGC", domainNames, oneSeconds)

	C := []crawler{bC, pcC, nbC}

	for _, c := range C {
		nilSitemap, error := c.crawl()
		if error != nil {
			t.Fail()
		}
		log.Printf("%s", nilSitemap)
	}

}
