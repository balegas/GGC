package main

import (
	"log"
	"testing"
)

func TestGetTags(t *testing.T) {
	file, err := readFileFromDisk("testFiles/simple_html.html")

	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}

	attrs := map[string]string{
		"a": "href",
	}
	result := getAllTagAttr(attrs, file)

	log.Printf("RESULT: %v", result)
}
