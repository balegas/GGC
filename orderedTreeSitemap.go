package main

import (
	"bytes"
	"io"

	rbt "github.com/emirpasic/gods/trees/redblacktree"
)

type OrderedTreeSitemap struct {
	prefixTree *rbt.Tree
}

func newOrderedTreeSitemap() *OrderedTreeSitemap {
	return &OrderedTreeSitemap{}
}

func initOrderedTreeSitemap(s *OrderedTreeSitemap) {
	s.prefixTree = rbt.NewWithStringComparator()
}

func (s *OrderedTreeSitemap) addURL(curl string) {
	s.prefixTree.Put(curl, struct{}{})
}

func (s *OrderedTreeSitemap) printSitemap(out io.Writer) {
	var buffer bytes.Buffer
	buffer.Write([]byte("<body>\n\t<ul>\n"))
	for _, k := range s.prefixTree.Keys() {
		buffer.Write([]byte("\t\t<li><a href=\""))
		buffer.Write([]byte(k.(string)))
		buffer.Write([]byte("\" >"))
		buffer.Write([]byte(k.(string)))
		buffer.Write([]byte("</li>\n"))

	}
	buffer.Write([]byte("\t</ul>\n</body>\n"))
	buffer.WriteTo(out)

}

func (s *OrderedTreeSitemap) getOrderedLinks() []string {
	orderedKeys := make([]string, 0, s.prefixTree.Size())
	for _, k := range s.prefixTree.Keys() {
		orderedKeys = append(orderedKeys, k.(string))
	}
	return orderedKeys
}

func (s *OrderedTreeSitemap) numberOfLinks() int {
	return s.prefixTree.Size()
}
