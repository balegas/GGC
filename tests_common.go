package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"

	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

func readFileFromDisk(filename string) (io.Reader, error) {
	bytesRead, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(bytesRead), nil
}

func setUpFakePage(pageLocation, pageFile string) {
	if bytesRead, err := ioutil.ReadFile(pageFile); err == nil {
		httpmock.RegisterResponder("GET", pageLocation,
			httpmock.NewBytesResponder(200, bytesRead))
		return
	}
	log.Fatalf("Error loading pageFile %v", pageFile)

}
