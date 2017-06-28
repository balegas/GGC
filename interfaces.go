package main

import (
	"io"
	"net/http"
	"net/url"
)

// htmlParser parses html pages and returns links to other files found in
// that page.
type htmlParser interface {
	findURLs() []string
}

// urlFrontier manages the links that are found by the crawler.
type urlFrontier interface {
	// Returns the next url to process
	// or an error if there are no urls left.
	// Order of results depends on implementation.
	nextURLString() (string, error)
	// Add new url to the frontier.
	addURLString(string)
	// Check if there are urls left in the frontier.
	isEmpty() bool
	// Returns the number of URLs left in the frontier.
	size() int
}

// urlStore manages the links that are found by the crawler.
type urlStore interface {
	// Returns true if values does not exist or replaced the previous value.
	// Return false on error.
	put(k string, content []byte) bool
	// Returns the stored page (can be full byte content or just the hash).
	// return true if value exists, false otherwise.
	get(k string) ([]byte, bool)
}

// Checks if an url meets the accessPolicy that is in place.
// E.g. of policies: url within a domain; robots.txt file rules.
type accessPolicy interface {
	// check wether an URL can be accessed or not.
	checkURL(url string) bool
}

// fetcher fetches urls using http protocol.
// Error is an HTML response code, or a negative number (fetcher errors).
type fetcher interface {
	getURLContent(url *url.URL) (*http.Response, int)
}

// Crawler has a minimal interface
type Crawler interface {
	Crawl() (sitemap, error)
}

//Stores urls from a domains and can print them.
//TODO: interface might change.
type sitemap interface {
	// Add an url to the sitemap.
	addURL(curl string)
	printSitemap(io.Writer)
	numberOfLinks() int
	getOrderedLinks() []string
}
