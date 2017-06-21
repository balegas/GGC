package main

import (
	"errors"
	"log"
	"net"
	"net/http"
	"net/url"
)

var errorFetching = errors.New("Error fetching URL")
var errorDomain = errors.New("Reading outside domain")
var errorRedirection = errors.New("Maximum redirections reached")

//var ErrorLocation = errors.New("Location header is empty")

type simpleFetcher struct {
	rules           accessPolicyChecker
	httpClient      *http.Client
	ipAddress       net.IP
	maxRedirections int
}

func defaultFetcher(rules accessPolicyChecker) simpleFetcher {
	httpClient := &http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}}

	return simpleFetcher{rules, httpClient, nil, 10}
}

func (f simpleFetcher) getURLContent(url *url.URL) (*http.Response, error) {
	//convert url back to normal form in case it was transformed.

	var err error
	redirections := 0
	nextLocation := url
	nextLocationCanonical, _ := getCanonicalURLString(nextLocation.String(), url)

	if !f.rules.checkURL(nextLocationCanonical) {
		log.Fatal("Should not try to get a forbidden url.")
	}

	for redirections < f.maxRedirections {
		resp, eR := f.httpClient.Get(nextLocation.String())

		if eR != nil {
			err = errorFetching
			break
		}

		if resp.StatusCode >= 300 && resp.StatusCode < 400 {
			location, _ := resp.Location()

			if !f.rules.checkURL(location.String()) {
				err = errorDomain
				break
			} else {
				redirections++
				nextLocation = location
				//log.Printf("Redirect to %v", nextLocation)
				continue
			}
		}

		// TODO: Need support for other codes?
		if resp.StatusCode == 200 {
			return resp, nil
		}
	}
	if redirections >= f.maxRedirections {
		err = errorRedirection
	}
	return nil, err
}

func main() {}
