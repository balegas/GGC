package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"time"

	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

const defaultStackSize = 1024

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
	initBasicCrawler(c, domainNames, fe, p, fr, duration, s)
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
	initProducerConsumerCrawler(c, domainNames, fe, p, fr, duration, s)
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
	initNBatchesCrawler(c, domainNames, fe, p, fr, duration, s, 4)
	return c
}
