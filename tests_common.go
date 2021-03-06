package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"time"

	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

func readFileFromDisk(filename string) (io.Reader, error) {
	bytesRead, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(bytesRead), nil
}

func setUpFakePage(pageLocation, pageFile string) {
	if bytesRead, err := ioutil.ReadFile(pageFile); err == nil {
		httpmock.RegisterResponder("GET", pageLocation,
			httpmock.NewBytesResponder(200, bytesRead))
		return
	}
	log.Fatalf("Error loading pageFile %v", pageFile)

}

// NewBasicCrawlerWithDomainPolicy cretes a basicCrawler with acces policy based
// on domain names.
func NewBasicCrawlerWithDomainPolicy(userAgent string, domainNames []string,
	duration time.Duration) Crawler {
	c := newBasicCrawler()
	p := newCheckDomainPolicy()
	initCheckDomainPolicy(p, domainNames)

	fe := defaultFetcher(p)
	fr := newStackFrontier(defaultStackSize)
	s := newInMemoryURLStore()
	sm := newOrderedTreeSitemap()

	initOrderedTreeSitemap(sm)
	initBasicCrawler(c, domainNames, fe, p, fr, duration, s, sm)
	return c
}

// NewProducerConsumerWithDomainPolicy cretes a producerConsumerCrawler with
// acces policy based on domain names.
func NewProducerConsumerWithDomainPolicy(userAgent string, domainNames []string,
	duration time.Duration) Crawler {
	c := newProducerConsumerCrawler()
	p := newCheckDomainPolicy()
	initCheckDomainPolicy(p, domainNames)

	fe := defaultFetcher(p)
	fr := newStackFrontier(defaultStackSize)
	s := newInMemoryURLStore()
	sm := newOrderedTreeSitemap()

	initOrderedTreeSitemap(sm)
	initProducerConsumerCrawler(c, domainNames, fe, p, fr, duration, s,
		defaultBufferSize, sm)
	return c
}

// NewNBatchesCrawlerWithDomainPolicy cretes a nBatchesCrawler with acces policy
// based on domain names.
func NewNBatchesCrawlerWithDomainPolicy(userAgent string, domainNames []string,
	duration time.Duration) Crawler {
	c := newNBatchesCrawler()
	p := newCheckDomainPolicy()
	initCheckDomainPolicy(p, domainNames)

	fe := defaultFetcher(p)
	fr := newStackFrontier(defaultStackSize)
	s := newInMemoryURLStore()
	sm := newOrderedTreeSitemap()

	initOrderedTreeSitemap(sm)
	initNBatchesCrawler(c, domainNames, fe, p, fr, duration, s, defaultNumWorkers,
		defaultBufferSize, defaultWorkerThinkTime, sm)
	return c
}

func NewSharedNothingCrawlerWithDomainPolicy(userAgent string, domainNames []string,
	duration time.Duration) Crawler {
	c := newSharedNothingCrawler()
	p := newCheckDomainPolicy()
	initCheckDomainPolicy(p, domainNames)

	fe := defaultFetcher(p)
	fr := newStackFrontier(defaultStackSize)
	s := newInMemoryURLStore()
	sm := newOrderedTreeSitemap()

	initOrderedTreeSitemap(sm)
	initSharedNothingCrawler(c, domainNames, fe, p, fr, duration, s, defaultNumWorkers,
		defaultBufferSize, defaultWorkerThinkTime, sm)
	return c
}
