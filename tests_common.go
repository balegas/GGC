package main

import (
  "bytes"
  "io"
  "io/ioutil"
)

func readFileFromDisk(filename, rootPath string) (io.Reader, error) {
	bytesRead, err := ioutil.ReadFile(rootPath + filename)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(bytesRead), nil
}
