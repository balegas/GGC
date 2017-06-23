package main

import (
	"io"
	"net/http"
	"net/url"
)

type htmlParser interface {
	findURLs() []string
}

type urlFrontier interface {
	// Returns the next url to process and ok,
	// or an error if there are no urls left.
	// Order depends on implementation.
	nextURLString() (string, error)
	// Add new url to the frontier.
	addURLString(string)
	// Check if there is any url in the frontier.
	isEmpty() bool
	size() int
}

type urlStore interface {
	// Returns true if values does not exist or replaced the previous value.
	// Return false on error.
	put(k string, content []byte) bool
	// Returns the stored page (can be full byte content or just the hash)
	// return true if value exists, false otherwise.
	get(k string) ([]byte, bool)
}

type accessPolicyChecker interface {
	// check wether an URL can be accessed or not.
	checkURL(url string) bool
}

type fetcher interface {
	getURLContent(url *url.URL) (*http.Response, error)
}

type crawler interface {
	crawl() (sitemap, error)
}

//Stores urls from a domains and can print them.
type sitemap interface {
	// Add an url for printing
	addURL() string // any advantage with URL data type?
	// Print the sitemap to a Writer. (TODO: Check if it allows print to console
	// and file)
	printSiteMap(io.Writer)
	numberOfLinks() int
}
