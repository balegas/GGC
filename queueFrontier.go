package main

//

type queueFrontier struct {
	urlQueue []interface{}
}

func newQueueFrontier(baseSize uint) *queueFrontier {
	q := make([]interface{}, 0, baseSize)
	return &queueFrontier{urlQueue: q}
}

func (f *queueFrontier) addURLString(url string) {
	f.urlQueue = append(f.urlQueue, url)
}

func (f *queueFrontier) nextURLString() (string, error) {
	if len(f.urlQueue) > 0 {
		ret := f.urlQueue[0]
		f.urlQueue = f.urlQueue[1:]
		return ret.(string), nil
	}
	return "", errEmptyFrontier
}

func (f *queueFrontier) isEmpty() bool {
	return len(f.urlQueue) <= 0
}
