// /home/krylon/go/src/github.com/blicero/podshrink/meta/01_meta_test.go
// -*- mode: go; coding: utf-8; -*-
// Created on 28. 08. 2023 by Benjamin Walkenhorst
// (c) 2023 Benjamin Walkenhorst
// Time-stamp: <2023-08-28 17:24:04 krylon>

package meta

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMetaExtractor(t *testing.T) {
	const tdir = "testdata"
	var (
		err   error
		e     *Extractor
		dirh  *os.File
		files []string
	)

	if dirh, err = os.Open(tdir); err != nil {
		t.Fatalf("Cannot open directory %s: %s",
			tdir,
			err.Error())
	}

	defer dirh.Close()

	if files, err = dirh.Readdirnames(-1); err != nil {
		t.Fatalf("Cannot read from directory %s: %s",
			tdir,
			err.Error())
	} else if e, err = NewExtractor(); err != nil {
		t.Fatalf("Cannot create Extractor: %s",
			err.Error())
	}

	for _, name := range files {
		var (
			m     *FileMeta
			fname = filepath.Join(tdir, name)
		)

		if m, err = e.ReadTags(fname); err != nil {
			t.Errorf("Cannot read tags from %s: %s",
				name,
				err.Error())
			continue
		}

		t.Logf("%s -> %s - %s - %04d %s (%d)",
			m.Path,
			m.Artist,
			m.Album,
			m.Track,
			m.Title,
			m.Year)

		if m.Cover != "" {
			t.Logf("Got cover image for %s: %s",
				name,
				m.Cover)
		}
	}
} // func TestMetaExtractor(t *testing.T)
