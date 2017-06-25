package main

import (
	"bytes"
	"io"

	rbt "github.com/emirpasic/gods/trees/redblacktree"
)

type orderedTreeSitemap struct {
	prefixTree *rbt.Tree
}

func newOrderedTreeSitemap() *orderedTreeSitemap {
	return &orderedTreeSitemap{}
}

func initOrderedTreeSitemap(s *orderedTreeSitemap) {
	s.prefixTree = rbt.NewWithStringComparator()
}

func (s *orderedTreeSitemap) addURL(curl string) {
	s.prefixTree.Put(curl, struct{}{})
}

func (s *orderedTreeSitemap) printSitemap(out io.Writer) {
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

func (s *orderedTreeSitemap) numberOfLinks() int {
	return s.prefixTree.Size()
}
