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

type stackFrontier struct {
	urlStack stack
}

var errEmptyFrontier = errors.New("no more urls in stack")

func newStackFrontier(baseSize uint) *stackFrontier {
	s := S.NewStack(1024)
	return &stackFrontier{urlStack: s}
}

func (f *stackFrontier) addURLString(url string) {
	f.urlStack.Push(url)
}

func (f *stackFrontier) nextURLString() (string, error) {
	ret, err := f.urlStack.Pop()
	if err != nil {
		return "", errEmptyFrontier
	}
	return ret.(string), err
}

func (f *stackFrontier) isEmpty() bool {
	return f.urlStack.Len() == 0
}
