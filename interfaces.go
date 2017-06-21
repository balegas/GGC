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
}

type urlStore interface {
	// Returns true if values does not exist in cache or is different.
	put(k string, v io.Reader) bool
	// Returns the value of the key in bytes and ok, or empty []byte and an error.
	get(string) (io.Reader, error)
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
