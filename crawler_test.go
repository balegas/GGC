package main

import (
	"math"
	"testing"
	"time"
)

const defaultStackSize = 1024

func newBasicCrawlerWithNoPolicy(domainName string, duration time.Duration) crawler {
	c := newBasicCrawler()
	p := newAllAllowedPolicy()
	f := newStackFrontier(defaultStackSize)
	initBasicCrawler(c, domainName, p, f, duration)
	return c
}

func TestBasicCrawlerWithNoLinks(t *testing.T) {
	c := newBasicCrawlerWithNoPolicy("domainA.com", math.MaxInt64)
	_, error := c.crawl()
	if error != nil {
		t.Fail()
	}
}
