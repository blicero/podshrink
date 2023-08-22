// /home/krylon/go/src/github.com/blicero/podshrink/walker/01_walker_test.go
// -*- mode: go; coding: utf-8; -*-
// Created on 22. 08. 2023 by Benjamin Walkenhorst
// (c) 2023 Benjamin Walkenhorst
// Time-stamp: <2023-08-22 21:06:58 krylon>

package walker

import (
	"regexp"
	"testing"
)

func TestWalkerSimple(t *testing.T) {
	const expectedFileCount = 60
	var (
		err     error
		w       *Walker
		q       = make(chan string, 5)
		folders = []string{
			"testdata/folder1",
			"testdata/folder2",
			"testdata/folder3",
		}
	)

	if w, err = Create(testFileFn, q, folders); err != nil {
		t.Fatalf("Cannot create Walker: %s",
			err.Error())
	}

} // func TestWalkerSimple(t *testing.T)

//////////////////////////////////////////////////////////////////////////////
// Helpers ///////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////

var filenamePat = regexp.MustCompile(`(?i)[.](mp3|ogg)$`)

func testFileFn(path string) bool {
	var m = filenamePat.MatchString(path)
	return m
} // func testFileFn(path string) bool

func countResults(fileq <-chan string, q chan<- int) {
	var cnt = 0
	for range fileq {
		cnt++
	}

	q <- cnt
	close(q)
} // func countResults(fileq <-chan string, q chan<- int)
