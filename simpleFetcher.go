package main

import (
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

//var errorFetching = errors.New("error fetching URL")
//var errorDomain = errors.New("reading outside domain")
//var errorRedirection = errors.New("maximum redirections reached")
//var errorTooManyRequests = errors.New("error 429, Please throttle requests")
//var errorOther = errors.New("other response code")

var defaultMaxRedirection = 10
var defaultTimeoutMillis = 5000

var errorForbidden = -1
var errorFetching = -2
var errorMaxRedirections = -3

type simpleFetcher struct {
	rules           accessPolicy
	httpClient      *http.Client
	ipAddress       net.IP
	maxRedirections int
}

func defaultFetcher(rules accessPolicy) simpleFetcher {
	timeout := time.Duration(defaultTimeoutMillis) * time.Millisecond
	httpClient := &http.Client{Timeout: timeout, CheckRedirect: func(req *http.Request,
		via []*http.Request) error {
		return http.ErrUseLastResponse
	}}

	return simpleFetcher{rules, httpClient, nil, defaultMaxRedirection}
}

func (f simpleFetcher) getURLContent(url *url.URL) (*http.Response, int) {
	redirections := 0
	nextLocation := url
	nextLocationCanonical, _ := getCanonicalURLString(nextLocation.String(), url)

	if !f.rules.checkURL(nextLocationCanonical) {
		log.Fatal("should never try to get a forbidden url")
	}

	for redirections < f.maxRedirections {
		resp, eR := f.httpClient.Get(nextLocation.String())

		if eR != nil {
			return nil, errorFetching
		}

		if resp.StatusCode >= 300 && resp.StatusCode < 400 {
			location, _ := resp.Location()

			// TODO: ONLY ACCEPTS status code 200. Must check other codes.
			if !f.rules.checkURL(location.String()) {
				return nil, errorForbidden
			}
			redirections++
			nextLocation = location
			continue
		} else if resp.StatusCode == 200 {
			return resp, 200
		} else {
			return nil, resp.StatusCode
		}
	}
	if redirections >= f.maxRedirections {
		return nil, errorMaxRedirections
	}

	log.Fatal("should not enter here")
	return nil, errorMaxRedirections
}
