package main

import (
	"log"
	"testing"
	"time"

	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

func TestCrawlersMock(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	domainURLs := []string{
		"http://www.domainGGC.com/",
		"http://www.domainGGC.com/nonExistingPage.html",
		"http://www.domainGGC.com/page1.html",
		"http://www.domainGGC.com/page2.html",
		"http://www.domainGGC.com/page3.html",
		"http://www.domainGGC.com/source.js",
		"http://www.domainGGC.com/theme.css",
	}

	domainNames := []string{"www.domainGGC.com"}
	setUpFakePage("http://www.domainGGC.com/", "testFiles/home.html")
	setUpFakePage("http://www.domainGGC.com/page1.html", "testFiles/page1.html")
	setUpFakePage("http://www.domainGGC.com/page2.html", "testFiles/page2.html")
	setUpFakePage("http://www.domainGGC.com/page3.html", "testFiles/page3.html")
	setUpFakePage("http://www.domainGGC.com/source.js", "testFiles/source.js")
	setUpFakePage("http://www.domainGGC.com/theme.css", "testFiles/theme.css")

	oneSeconds := time.Duration(1) * time.Second

	bC := NewBasicCrawlerWithDomainPolicy("GGC", domainNames, oneSeconds)
	pcC := NewProducerConsumerWithDomainPolicy("GGC", domainNames, oneSeconds)
	nbC := NewNBatchesCrawlerWithDomainPolicy("GGC", domainNames, oneSeconds)

	C := []Crawler{bC, pcC, nbC}

	for _, c := range C {
		sm, error := c.Crawl()
		if error != nil {
			t.Fail()
		}
		links := sm.getOrderedLinks()

		if len(links) != len(domainURLs) {
			t.Errorf("Number of URLs found does not match.")
		}

		log.Printf("%v", links)

		for i, l := range links {
			if l != domainURLs[i] {
				t.Errorf("URL different from expected: %v C: %v, E: %v", i, l, domainURLs[i])
			}

		}

	}

}
