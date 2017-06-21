package main

//

import (
	"errors"

	S "github.com/junpengxiao/Stack"
)

type stack interface {
	Push(value interface{})
	Pop() (interface{}, error)
	Len() int
}

//FIXME: How to use unexported types properly?
type stackFrontier struct {
	urlStack stack
}

var errEmptyFrontier = errors.New("no more urls in stack")

func newStackFrontier(baseSize uint) *stackFrontier {
	s := S.NewStack(1024)
	return &stackFrontier{urlStack: s}
}

func (f *stackFrontier) addURL(url string) {
	f.urlStack.Push(url)
}

func (f *stackFrontier) nextURL() (string, error) {
	ret, err := f.urlStack.Pop()
	return ret.(string), err
}

func (f *stackFrontier) isEmpty() bool {
	return f.urlStack.Len() <= 0
}
