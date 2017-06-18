package main

import (
	"io"
	"net/url"

	"golang.org/x/net/html"
)

// Function to retrieve a list of attribute values in HTML5 tags.
// Example: getAllTagAttr(map[string]string{"a": "href",}, File) retrieves
// all values of attribute href for all "a" tags in the document.
// TODO: Add support for multiple attributes -- not necessary for this exercise.
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

func getCanonicalURL(urlString string) (*url.URL, error) {
	//TODO: transform parameters to path; order arguments by index lex. order
	return url.Parse(urlString)

}
