// /home/krylon/go/src/github.com/blicero/podshrink/walker/01_walker_test.go
// -*- mode: go; coding: utf-8; -*-
// Created on 22. 08. 2023 by Benjamin Walkenhorst
// (c) 2023 Benjamin Walkenhorst
// Time-stamp: <2023-08-23 10:16:42 krylon>

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
		fileQ   = make(chan string, 5)
		cntQ    = make(chan int)
		fileCnt int
		folders = []string{
			"testdata/folder1",
			"testdata/folder2",
			"testdata/folder3",
		}
	)

	if w, err = Create(testFileFn, fileQ, folders); err != nil {
		t.Fatalf("Cannot create Walker: %s",
			err.Error())
	}

	go countResults(fileQ, cntQ)

	if err = w.Run(); err != nil {
		t.Errorf("Failed to walk directory trees: %s",
			err.Error())
	}

	fileCnt = <-cntQ

	if fileCnt != expectedFileCount {
		t.Errorf("Unexpected number of files emitted by Walker: %d (expected %d)",
			fileCnt,
			expectedFileCount)
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
