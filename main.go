package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
)

//TODO: Should mark header tags to avoid fetching those files
//TODO: Could also make Tag -> Array of attributes
var crawlTags = map[string]string{
	"link":   "href",
	"script": "src",
	"a":      "href",
}

const defaultDuration = 10
const defaultOutputFile = "stdout"
const defaultBufferSize = 10
const defaultNumWorkers = 16
const defaultStackSize = 1024

func main() {

	durationPointer := flag.Duration("d", defaultDuration, "Cralwing duration in seconds.")
	outputFilePointer := flag.String("f", defaultOutputFile, "output file | stdout.")
	workerBuffSizePointer := flag.Int("b", defaultBufferSize, "Worker input buffer size.")
	numWorkersPointer := flag.Int("w", defaultNumWorkers, "Worker input buffer size.")

	flag.Parse()

	duration := *durationPointer
	outputFile := *outputFilePointer
	bufferSize := *workerBuffSizePointer
	numWorkers := *numWorkersPointer
	domainNames := flag.Args()

	log.Printf("duration: %v, output: %v", duration, outputFile)
	log.Printf("Args %v", flag.Args())

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
	initNBatchesCrawler(c, domainNames, fe, p, fr, duration, s, numWorkers, bufferSize, sm)

	result, _ := c.Crawl()

	var f io.Writer
	var out *bufio.Writer
	var file os.File
	if outputFile == "stdout" {
		f = os.Stdout
	} else {
		file, err := os.Create(outputFile)
		if err != nil {
			log.Fatal("Error opening file.")
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
