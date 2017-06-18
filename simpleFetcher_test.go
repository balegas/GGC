package main

import (
	"io/ioutil"
	"net/url"
	"reflect"
	"testing"

	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

func TestGet(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	location := "https://domainA.com/pageA"
	content := []byte("Hello!")

	httpmock.RegisterResponder("GET", location,
		httpmock.NewBytesResponder(200, content))

	f := defaultFetcher("domainA.com")
	url, _ := url.Parse(location)

	resp, _ := f.getURLContent(*url)
	body, _ := ioutil.ReadAll(resp.Body)

	if !(reflect.DeepEqual(body, content)) {
		t.Fail()
	}
}

func TestRedirect(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	content := []byte("Hello!")

	domainA := "domainA.com"
	locationA := "https://" + domainA + "/pageA"
	locationB := "https://" + domainA + "/pageB"

	httpmock.RegisterResponder("GET", locationA,
		NewRedirectResponder(locationB))
	httpmock.RegisterResponder("GET", locationB,
		httpmock.NewBytesResponder(200, content))

	f := defaultFetcher(domainA)
	url, _ := url.Parse(locationA)

	resp, _ := f.getURLContent(*url)
	body, _ := ioutil.ReadAll(resp.Body)
	if !(reflect.DeepEqual(body, content)) {
		t.Fail()
	}
}

func TestInfiniteRedirect(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	location := "https://domainA.com/pageA"

	httpmock.RegisterResponder("GET", location,
		NewRedirectResponder(location))

	f := defaultFetcher("domainA.com")
	url, _ := url.Parse(location)

	_, err := f.getURLContent(*url)

	if err != errorRedirection {
		t.Fail()
	}
}

func TestDifferentDomainRedirect(t *testing.T) {

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	domainA := "domainA.com"
	locationA := "https://" + domainA + "/pageA"
	locationB := "https://domainB.com/pageA"

	httpmock.RegisterResponder("GET", locationA,
		NewRedirectResponder(locationB))

	f := defaultFetcher(domainA)
	url, _ := url.Parse(locationA)

	_, err := f.getURLContent(*url)

	if err != errorDomain {
		t.Fail()
	}
}

func TestError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	location := "https://domainA.com/pageA"

	f := defaultFetcher("domainA.com")
	url, _ := url.Parse(location)

	_, err := f.getURLContent(*url)

	if err != errorFetching {
		t.Fail()
	}
}
