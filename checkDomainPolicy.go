package main

import (
	"net/url"
	"strings"
)

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
	// Not very efficient...
	parsedURL, _ := url.Parse(urlString)
	for _, domain := range p.domainNames {
		if strings.EqualFold(parsedURL.Hostname(), domain) {
			return true
		}
	}
	return false
}
