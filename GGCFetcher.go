package main

import (
	"errors"
	"net"
	"net/http"
	"net/url"
	"strings"
)

var errorFetching = errors.New("Error fetching URL")
var errorDomain = errors.New("Reading outside domain")
var errorRedirection = errors.New("Maximum redirections reached")

//var ErrorLocation = errors.New("Location header is empty")

type simpleFetcher struct {
	domainName      string // must keep domain name to analyze redirects
	httpClient      *http.Client
	ipAddress       net.IP
	maxRedirections int
}

func defaultFetcher(domainName string) simpleFetcher {
	httpClient := &http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}}

	return simpleFetcher{domainName, httpClient, nil, 10}
}

//TODO: NEED TO CHECK DOMAIN
func (f simpleFetcher) getURLContent(url url.URL) (*http.Response, error) {
	var err error
	redirections := 0
	nextLocation := &url

	for redirections < f.maxRedirections {
		resp, eR := f.httpClient.Get(nextLocation.String())

		if eR != nil {
			err = errorFetching
			break
		}

		if resp.StatusCode >= 300 && resp.StatusCode < 400 {
			location, _ := resp.Location()
			if !strings.EqualFold(location.Hostname(), f.domainName) {
				err = errorDomain
				break
			} else {
				redirections++
				nextLocation = location
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
