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

func newProducerConsumerWithDomainPolicy(userAgent string, domainNames []string, duration time.Duration) crawler {
	c := newProducerConsumerCrawler()
	p := newCheckDomainPolicy()
	initCheckDomainPolicy(p, domainNames)

	fe := defaultFetcher(p)
	fr := newStackFrontier(defaultStackSize)
	s := newInMemoryURLStore()
	initProducerConsumerCrawler(c, domainNames, fe, p, fr, duration, s)
	return c
}

func newNBatchesCrawlerWithDomainPolicy(userAgent string, domainNames []string, duration time.Duration) crawler {
	c := newNBatchesCrawler()
	p := newCheckDomainPolicy()
	initCheckDomainPolicy(p, domainNames)

	fe := defaultFetcher(p)
	fr := newStackFrontier(defaultStackSize)
	s := newInMemoryURLStore()
	initNBatchesCrawler(c, domainNames, fe, p, fr, duration, s, 4)
	return c
}
