package main

import "time"

type allAllowed struct {
	userAgeng string
	thinkTime time.Duration
}

func newAllAllowedPolicy() *allAllowed {
	return &allAllowed{"GGC", 0}
}

func (r *allAllowed) checkURL(url string) bool {
	return true
}
