// /home/krylon/go/src/github.com/blicero/podshrink/walker/walker.go
// -*- mode: go; coding: utf-8; -*-
// Created on 18. 08. 2023 by Benjamin Walkenhorst
// (c) 2023 Benjamin Walkenhorst
// Time-stamp: <2023-08-22 18:54:41 krylon>

// Package walker implements the walking of directory trees.
package walker

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/blicero/podshrink/common"
	"github.com/blicero/podshrink/logdomain"
)

// Filter is a function that determines whether a file is of interest.
type Filter func(string) bool

// Walker wraps up the state for the tree traversal.
type Walker struct {
	folders []string
	log     *log.Logger
	fileQ   chan<- string
	filter  Filter
	root    string
}

// Create creates a new Walker.
func Create(filter Filter, queue chan<- string, folders []string) (*Walker, error) {
	var (
		err error
		w   = &Walker{
			folders: folders,
			fileQ:   queue,
			filter:  filter,
		}
	)

	if w.log, err = common.GetLogger(logdomain.Walker); err != nil {
		return nil, err
	}

	return w, nil
} // func Create(filter Filter, queue chan string, folders []string) (*Walker, error)

// Run executes the walker.
func (w *Walker) Run() error {
	defer func() {
		w.root = ""
		close(w.fileQ)
	}()

	for _, folder := range w.folders {
		var (
			err error
			f   fs.FS
		)
		w.root = folder
		f = os.DirFS(folder)

		if err = fs.WalkDir(f, "/", w.process); err != nil {
			return err
		}
	}

	return nil
} // func (w *Walker) Run() error

func (w *Walker) process(path string, dir fs.DirEntry, incoming error) error {
	if incoming != nil {
		w.log.Printf("[ERROR] Error processing %s: %s\n",
			path,
			incoming.Error())
		if dir.IsDir() {
			return fs.SkipDir
		}
	}

	if dir.IsDir() {
		return nil
	} else if !w.filter(path) {
		return nil
	}

	w.fileQ <- filepath.Join(w.root, path)

	return nil
} // func (w *Walker) process(path string, dir fs.DirEntry, incoming error) error
