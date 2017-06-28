package main

import (
	"net/url"
	"strings"
)

//Access policy that checks if an URL is within a domain.
//TODO: Does not accept subdomains
type checkSubDomainPolicy struct {
	domainNames []string
}

func newCheckSubDomainPolicy() *checkSubDomainPolicy {
	return &checkSubDomainPolicy{}
}

func initCheckSubDomainPolicy(p *checkSubDomainPolicy, domainNames []string) {
	p.domainNames = domainNames
}

func (p *checkSubDomainPolicy) checkURL(urlString string) bool {
	// TODO: avoid url.Parse.
	parsedURL, _ := url.Parse(urlString)
	for _, domain := range p.domainNames {
		if strings.HasSuffix(parsedURL.Hostname(), domain) {
			return true
		}
	}
	return false
}
