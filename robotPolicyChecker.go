package main

import (
	"io"
	"time"
)

type robotsRules struct {
	pathsPerm map[string]bool //maps paths -> allowed/disallowed
	userAgeng string
	thinkTime time.Duration
}

func (r *robotsRules) new() *robotsRules {
	return &robotsRules{make(map[string]bool), "GGC", 0}
}

// TODO: ensure ioReader is the correct type
func (r *robotsRules) newFromRobotsFile(robotsFile io.Reader, userAgent string) *robotsRules {
	return &robotsRules{make(map[string]bool), userAgent, 0}
	// set thinktime based on file, or value from environment, or 0.
	// process paths and mark allowed/disallowed.
}

// TODO 1: Must convert to canonical paths

// TODO 2: Add support for a * operator in subpaths: e.g. /home/*/something

// NOTE 1: Must learn more about corner cases like
// is it possible to access some file in a folder, but
// prevent access to any other file in the same folder?

func (r *robotsRules) checkURL(url string) bool {
	//TODO implement policy
	return false
}

//Alternative approach, assuming that the size of robots.txt is larger than
//Most common URLS.
//func (r *robotRules) checkURL2(path) bool{
//TODO: Get all	url 'logical' prefixes, check if prefix is in pathsPrefix.
// Requires canonical links to be more effective (also needs to be done).
//}
