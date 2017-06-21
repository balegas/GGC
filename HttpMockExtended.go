package main

import (
	"net/http"
	"strconv"

	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

// NewRedirectResponder creates a Responder with an empty body, a 301 code,
// and a location header.
func NewRedirectResponder(location string) httpmock.Responder {
	return httpmock.ResponderFromResponse(NewRedirectResponse(location))
}

// NewRedirectResponse creates an *http.Response with an empty body, a 301 code,
// and a location header.
func NewRedirectResponse(location string) *http.Response {
	return &http.Response{
		Status:     strconv.Itoa(301),
		StatusCode: 301,
		Header:     http.Header{"Location": {location}},
	}
}

func startMock() {
	httpmock.Activate()
}

func endMock() {
	httpmock.DeactivateAndReset()
}
