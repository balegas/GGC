package main

import (
	"log"
	"time"

	S "github.com/junpengxiao/Stack"
)

// SharedNothigCrawler uses multiple go routines to crawl the internet.
// Each go routine receives a single urls to process and returns a slice of
// found urls.
// The coordinator filters received urls.
type sharedNothingCrawler struct {
	crawlerInternals
	maxWorkers  int
	currWorkers int
	workerTT    time.Duration
}

type response struct {
	workerID int
	curls    []string
}

func newSharedNothingCrawler() *sharedNothingCrawler {
	return &sharedNothingCrawler{}
}

func newResponse(workerID int, curls []string) *response {
	return &response{workerID, curls}
}

func initSharedNothingCrawler(c *sharedNothingCrawler, seed []string, fet fetcher,
	rules accessPolicy, uf urlFrontier, duration time.Duration, s urlStore,
	maxRoutines int, bufferSize int, workerTT time.Duration, sm sitemap) {
	initCommonAttributes(&c.crawlerInternals, seed, fet, rules, uf, duration, s, sm)
	c.maxWorkers = maxRoutines
	c.workerTT = workerTT
	c.currWorkers = 0
}

// Spawns workers to crawl webpages while there are urls left in the frontier.
// Workers return a batch of urls ready to put on the frontier.
// While no more workers can be spawned, add results to the frontier.

// NOTE: Might underperform if results queue gets full and frontier gets empty.
// The objective is to ensure is that the frontier is sufficiently occupied, to
// ensure that the workers always have enough work. Can tune number of threads
// and the buffer size of channels.
func (c *sharedNothingCrawler) Crawl() (sitemap, error) {
	foundURLs := 0

	freeWorkers := S.NewStack(uint(c.maxWorkers))
	for i := 0; i < c.maxWorkers; i++ {
		freeWorkers.Push(i)
	}

	newURLsC := make(chan *response, c.maxWorkers)
	controlC := make(chan interface{})

	workers := c.spawnRoutines(newURLsC, controlC)

	//init
	for !c.frontier.isEmpty() && freeWorkers.Len() != 0 {
		nextWorker, _ := freeWorkers.Pop()
		nextURL, _ := c.frontier.nextURLString()
		c.sitemap.addURL(nextURL)
		log.Printf("Pushing URL %v to worker %v", nextURL, nextWorker.(int))
		workers[nextWorker.(int)] <- nextURL

	}

	select {
	case newResponse := <-newURLsC:
		{
			log.Printf("Received Response %v", newResponse.workerID)
			//Process incoming
			for _, curl := range newResponse.curls {
				if !c.seen(curl) && c.canProcess(curl) {
					foundURLs++
					c.store.put(curl, []byte{})
					c.frontier.addURLString(curl)
				}
			}
			//Start new task

			//Check finished
			if c.frontier.isEmpty() || c.isTimeout() {
				//TODO: Add remaining to sitemap
				log.Printf("Finished. Found %v urls.", foundURLs)
				close(controlC)
				return c.sitemap, nil
			}
			nextURL, _ := c.frontier.nextURLString()
			c.sitemap.addURL(nextURL)
			workers[newResponse.workerID] <- nextURL

		}
	}

	close(newURLsC)
	//close(signalC)
	return c.sitemap, nil
}

func (c *sharedNothingCrawler) spawnRoutines(newURLsC chan *response, controlC chan interface{}) []chan string {
	pendingURLsSlice := make([]chan string, c.maxWorkers)

	for i := range pendingURLsSlice {
		pendingURLsSlice[i] = make(chan string)
		go c.processURLs(i, pendingURLsSlice[i], newURLsC, controlC)
	}
	return pendingURLsSlice
}

func (c *sharedNothingCrawler) processURLs(myID int, inputC chan string,
	outputC chan *response, controlC chan interface{}) {
	select {
	case newCURL := <-inputC:
		{
			log.Printf("Worker %v received URL", myID)
			curls := make([]string, 0)

			nextURL, _ := toURL(newCURL)
			newURLs, _, _ := c.findURLLinksGetBody(nextURL)
			log.Printf("NEW URLs %v", newURLs)
			for _, url := range newURLs {
				curli, _ := getCanonicalURLString(url, nextURL)
				curls = append(curls, curli)
			}
			log.Printf("Going to respond %v", myID)
			outputC <- newResponse(myID, curls)
			time.Sleep(c.workerTT)
		}
	case <-controlC:
		log.Printf("Closing worker %v", myID)
		close(inputC)
	}

}
