package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"testing"
)

func TestGetTags(t *testing.T) {
	file, err := readFileFromDisk("testFiles/simple_html.html", "")

	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}

	attrs := map[string]string{
		"a": "href",
	}
	result := getAllTagAttr(attrs, file)

	log.Printf("RESULT: %v", result)
}

func readFileFromDisk(filename, rootPath string) (io.Reader, error) {
	bytesRead, err := ioutil.ReadFile(rootPath + filename)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(bytesRead), nil
}
