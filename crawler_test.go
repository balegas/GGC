package main

import (
	"log"
	"testing"
	"time"
)

/*
func TestCrawlerOnWeb(t *testing.T) {
	domainNames := []string{"gatofedorento.blogspot.pt", "www.gatofedorento.blogspot.pt"}
	TenSeconds := time.Duration(10) * time.Second
	c := newBasicCrawlerWithDomainPolicy("GGC", domainNames, TenSeconds)
	nilSitemap, error := c.crawl()
	if error != nil {
		t.Fail()
	}
	log.Printf("%s", nilSitemap)
}
*/

func TestBasicCrawlerMock(t *testing.T) {
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
