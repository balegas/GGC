package main

import (
	"net/url"
	"strings"
)

//Access policy that checks if an URL is within a domain.
//TODO: Does not accept subdomains
type checkDomainPolicy struct {
	domainNames []string
}

func newCheckDomainPolicy() *checkDomainPolicy {
	return &checkDomainPolicy{}
}

func initCheckDomainPolicy(p *checkDomainPolicy, domainNames []string) {
	p.domainNames = domainNames
}

func (p *checkDomainPolicy) checkURL(urlString string) bool {
	// TODO: avoid url.Parse.
	parsedURL, _ := url.Parse(urlString)
	for _, domain := range p.domainNames {
		if strings.EqualFold(parsedURL.Hostname(), domain) {
			return true
		}
	}
	return false
}
