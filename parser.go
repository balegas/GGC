package main

import (
	"io"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// Function to retrieve a list of attribute values in HTML5 tags.
// Example: getAllTagAttr(map[string]string{"a": "href",}, File) retrieves
// all values of attribute href for all "a" tags in the document.

// TODO: This parser implementation uses a html tokenizer. Alternatively, we could
// use regex expressions. Need benchmarks to compare results.

// TODO: It is potentially more efficient to return only new links. The design
// of the code is less clear that way. Need benchmarks to evaluate benefits.
func getAllTagAttr(tagAttr map[string]string, content io.Reader) []string {
	var found []string
	parser := html.NewTokenizer(content)
	token := parser.Next()

	for token != html.ErrorToken {

		t := parser.Token()

		if searchAttr, ok := tagAttr[t.Data]; ok {
			for _, a := range t.Attr {
				if a.Key == searchAttr {
					found = append(found, a.Val)
				}
			}
		}
		token = parser.Next()
	}

	return found
}

func getCanonicalURLString(urlString string, parentURL *url.URL) (string, error) {
	//TODO: transform parameters to path segments; order arguments by index lex. order
	//Full path

	//Remove Anchors #
	lastIndexOf := strings.LastIndex(urlString, "#")
	for lastIndexOf > 0 {
		urlString = urlString[:len(urlString)-(len(urlString)-lastIndexOf)]
		lastIndexOf = strings.LastIndex(urlString, "#")
	}

	if strings.Index(urlString, "http://") == 0 || strings.Index(urlString,
		"https://") == 0 {
		return urlString, nil
	}

	if strings.Index(urlString, "/") == 0 {
		return parentURL.Scheme + "://" + parentURL.Hostname() + urlString, nil
	}

	//Relative path
	return parentURL.Scheme + "://" + parentURL.Hostname() + "/" + urlString, nil

}

func toURL(urlString string) (*url.URL, error) {
	//TODO: after transforming urlParams in getCanonicalURLString,
	// transform them back to URL.
	return url.Parse(urlString)
}
