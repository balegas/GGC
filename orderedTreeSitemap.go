package main

import (
	"bytes"
	"io"
	"sync"

	rbt "github.com/emirpasic/gods/trees/redblacktree"
)

type orderedTreeSitemap struct {
	mutex      sync.Mutex
	prefixTree *rbt.Tree
}

func newOrderedTreeSitemap() *orderedTreeSitemap {
	return &orderedTreeSitemap{}
}

func initOrderedTreeSitemap(s *orderedTreeSitemap) {
	s.prefixTree = rbt.NewWithStringComparator()
}

func (s *orderedTreeSitemap) addURL(curl string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.prefixTree.Put(curl, struct{}{})
}

func (s *orderedTreeSitemap) printSitemap(out io.Writer) {
	var buffer bytes.Buffer
	buffer.Write([]byte("<body>\n\t<ul>\n"))

	s.mutex.Lock()
	defer s.mutex.Unlock()

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

func (s *orderedTreeSitemap) getOrderedLinks() []string {
	orderedKeys := make([]string, 0, s.prefixTree.Size())
	for _, k := range s.prefixTree.Keys() {
		orderedKeys = append(orderedKeys, k.(string))
	}
	return orderedKeys
}

func (s *orderedTreeSitemap) numberOfLinks() int {
	return s.prefixTree.Size()
}
