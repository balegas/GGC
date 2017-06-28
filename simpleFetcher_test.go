package main

import (
	"io/ioutil"
	"net/url"
	"reflect"
	"testing"

	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

const domainA = "domainA.com"
const domainB = "domainB.com"
const locationA = "https://" + domainA + "/pageA"
const locationA1 = "https://" + domainA + "/pageB"
const locationB = "https://" + domainB + "/pageB"

func TestGet(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	content := []byte("Hello!")

	httpmock.RegisterResponder("GET", locationA,
		httpmock.NewBytesResponder(200, content))

	p := newCheckDomainPolicy()
	initCheckDomainPolicy(p, []string{domainA})

	f := defaultFetcher(p)

	urlA, _ := url.Parse(locationA)
	resp, _ := f.getURLContent(urlA)
	body, _ := ioutil.ReadAll(resp.Body)

	if !(reflect.DeepEqual(body, content)) {
		t.Fail()
	}
}

func TestRedirect(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	content := []byte("Hello!")

	httpmock.RegisterResponder("GET", locationA,
		NewRedirectResponder(locationA1))
	httpmock.RegisterResponder("GET", locationA1,
		httpmock.NewBytesResponder(200, content))

	p := newCheckDomainPolicy()
	initCheckDomainPolicy(p, []string{domainA})

	f := defaultFetcher(p)

	urlA, _ := url.Parse(locationA)
	resp, _ := f.getURLContent(urlA)
	body, _ := ioutil.ReadAll(resp.Body)
	if !(reflect.DeepEqual(body, content)) {
		t.Fail()
	}
}

func TestInfiniteRedirect(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", locationA,
		NewRedirectResponder(locationA))

	p := newCheckDomainPolicy()
	initCheckDomainPolicy(p, []string{domainA})

	f := defaultFetcher(p)

	urlA, _ := url.Parse(locationA)
	_, err := f.getURLContent(urlA)

	if err != errorMaxRedirections {
		t.Fail()
	}
}

func TestDifferentDomainRedirect(t *testing.T) {

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", locationA,
		NewRedirectResponder(locationB))

	p := newCheckDomainPolicy()
	initCheckDomainPolicy(p, []string{domainA})

	f := defaultFetcher(p)

	urlA, _ := url.Parse(locationA)
	_, err := f.getURLContent(urlA)

	if err != errorForbidden {
		t.Fail()
	}
}

func TestError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	p := newCheckDomainPolicy()
	initCheckDomainPolicy(p, []string{domainA})

	f := defaultFetcher(p)

	urlA, _ := url.Parse(locationA)
	_, err := f.getURLContent(urlA)

	if err != errorFetching {
		t.Fail()
	}
}
