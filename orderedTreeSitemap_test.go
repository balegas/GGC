package main

import (
	"math/rand"
	"testing"
)

func TestSitemap(t *testing.T) {

	expectedResults := [5]string{
		"a",
		"b",
		"c",
		"d",
		"e",
	}

	randomOrder := make([]string, len(expectedResults))
	perm := rand.Perm(len(expectedResults))
	for i, v := range perm {
		randomOrder[v] = expectedResults[i]
	}

	sitemap := newOrderedTreeSitemap()
	initOrderedTreeSitemap(sitemap)

	for _, e := range randomOrder {
		sitemap.addURL(e)
	}

	for i, elem := range sitemap.getOrderedLinks() {
		if expectedResults[i] != elem {
			t.Errorf("Element out of order.")
		}
	}

}
