package main

import (
	"bufio"
	"io"
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

	domainNames := os.Args[1 : len(os.Args)-2]
	durationInt, _ := strconv.Atoi(os.Args[len(os.Args)-2])
	duration := time.Duration(durationInt) * time.Second
	outputFile := os.Args[len(os.Args)-1]

	log.Printf("DomainNames: %v, duration: %v", domainNames, durationInt)

	//c := newBasicCrawler()
	//c := newProducerConsumerCrawler()
	c := newNBatchesCrawler()

	p := newCheckDomainPolicy()
	initCheckDomainPolicy(p, domainNames)

	fe := defaultFetcher(p)
	fr := newQueueFrontier(defaultStackSize)
	s := newInMemoryURLStore()
	sm := newOrderedTreeSitemap()
	initOrderedTreeSitemap(sm)

	//initBasicCrawler(c, domainNames, fe, p, fr, duration, s, sm)
	//initProducerConsumerCrawler(c, domainNames, fe, p, fr, duration, s, sm)
	initNBatchesCrawler(c, domainNames, fe, p, fr, duration, s, 4, sm)

	result, _ := c.Crawl()

	var f io.Writer
	var out *bufio.Writer
	var file os.File
	if outputFile == "stdout" {
		f = os.Stdout
	} else {
		file, err := os.Create(outputFile)
		if err != nil {
			log.Print("File cannot be created, output to stdout.")
			f = os.Stdout
		} else {
			f = file

		}

	}
	out = bufio.NewWriter(f)
	result.printSitemap(out)
	if f != nil {
		file.Close()
	}

}
