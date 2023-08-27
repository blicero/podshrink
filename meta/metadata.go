// /home/krylon/go/src/github.com/blicero/podshrink/meta/metadata.go
// -*- mode: go; coding: utf-8; -*-
// Created on 24. 08. 2023 by Benjamin Walkenhorst
// (c) 2023 Benjamin Walkenhorst
// Time-stamp: <2023-08-27 19:07:28 krylon>

// Package meta implements extracting metadata from audio files.
package meta

import (
	"log"
	"os"

	"github.com/blicero/podshrink/common"
	"github.com/blicero/podshrink/logdomain"
	"github.com/dhowden/tag"
)

// FileMeta contains the metadata for a single audio file.
type FileMeta struct {
	Path    string
	Artist  string
	Title   string
	Album   string
	Year    int
	Track   int
	Cover   string
	Comment string
}

// Extractor wraps up the state associated with
type Extractor struct {
	log *log.Logger
}

func NewExtractor() (*Extractor, error) {
	var (
		err error
		ex  = &Extractor{}
	)

	if ex.log, err = common.GetLogger(logdomain.Extractor); err != nil {
		return nil, err
	}

	return ex, nil
} // func NewExtractor() (*Extractor, error)

// ReadTags tries to extract the metadata from the given file.
func (e *Extractor) ReadTags(path string) (*FileMeta, error) {
	var (
		fh  *os.File
		err error
		rdr tag.Metadata
	)

	if fh, err = os.Open(path); err != nil {
		e.log.Printf("[ERROR] Cannot open %s: %s\n",
			path,
			err.Error())
		return nil, err
	}

	defer fh.Close() // nolint: errcheck

	if rdr, err = tag.ReadFrom(fh); err != nil {
		e.log.Printf("[ERROR] Cannot read metadata from %s: %s\n",
			path,
			err.Error())
		return nil, err
	}

	tnum, _ := rdr.Track()

	var m = &FileMeta{
		Path:    path,
		Artist:  rdr.Artist(),
		Title:   rdr.Title(),
		Album:   rdr.Album(),
		Year:    rdr.Year(),
		Track:   tnum,
		Comment: rdr.Comment(),
	}

	// Cover!!!
	var cover *tag.Picture

	if cover = rdr.Picture(); cover != nil {

	}

	return m, nil
} // func (e *Extractor) ReadTags(path string) (FileMeta, error)
