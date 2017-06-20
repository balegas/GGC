package main

import "errors"

//FIXME: How to use unexported types properly?
type stackFrontier struct {
	urlStack []string
}

var errEmptyFrontier = errors.New("no more urls in stack")

func newStackFrontier(baseSize uint) *stackFrontier {
	return &stackFrontier{urlStack: make([]string, 0, baseSize)}
}

func (f *stackFrontier) addURL(url string) {
	f.urlStack = append(f.urlStack, url)
}

//pop the top item out, if stack is empty, will return ErrEmptyStack decleared above
func (f *stackFrontier) nextURL() (string, error) {
	if !f.isEmpty() {
		URL := f.urlStack[len(f.urlStack)-1]
		f.urlStack = f.urlStack[:len(f.urlStack)-1]
		return URL, nil
	}
	return "", errEmptyFrontier
}

func (f *stackFrontier) isEmpty() bool {
	return len(f.urlStack) == 0
}
