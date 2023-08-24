// /home/krylon/go/src/github.com/blicero/podshrink/meta/metadata.go
// -*- mode: go; coding: utf-8; -*-
// Created on 24. 08. 2023 by Benjamin Walkenhorst
// (c) 2023 Benjamin Walkenhorst
// Time-stamp: <2023-08-24 17:24:26 krylon>

// Package meta implements extracting metadata from audio files.
package meta

import (
	"log"

	"github.com/blicero/podshrink/common"
	"github.com/blicero/podshrink/logdomain"
)

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

func (e *Extractor) ReadTags(path string) {
}
