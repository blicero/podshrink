// /home/krylon/go/src/github.com/blicero/podshrink/convert/convert.go
// -*- mode: go; coding: utf-8; -*-
// Created on 28. 08. 2023 by Benjamin Walkenhorst
// (c) 2023 Benjamin Walkenhorst
// Time-stamp: <2023-08-28 21:00:59 krylon>

// Package convert implements the conversion of various audio formats to opus.
package convert

import (
	"log"
	"regexp"
	"strings"

	"github.com/blicero/podshrink/common"
	"github.com/blicero/podshrink/logdomain"
)

var suffixPat = regexp.MustCompile("[.]([^.]+)$")

// Converter wraps the state associated with converting audio files.
type Converter struct {
	log *log.Logger
}

// New creates a new Converter
func New() (*Converter, error) {
	var (
		err error
		c   = &Converter{}
	)

	if c.log, err = common.GetLogger(logdomain.Converter); err != nil {
		return nil, err
	}

	return c, nil
} // func New() (*Converter, error)

func (c *Converter) generateCommand(in, out string) []string {
	var match []string

	if match = suffixPat.FindStringSubmatch(in); match == nil {
		return nil
	}

	var suffix = match[1]

	switch strings.ToLower(suffix) {
	case "mp3", "mpga":
		return []string{
			"mpg123",
			"-q",
			"--no-control",
			"-o",
			"wav",
			"-w",
			out,
			in,
		}
	case "ogg", "oga":
		return []string{
			"ogg123",
			"-q",
			"-d",
			"wav",
			"-f",
			out,
			in,
		}
	case "m4a":
		return []string{
			"ffmpeg",
			"-i",
			in,
			out,
		}
	default:
		c.log.Printf("[ERROR] Cannot find decoder for %s\n", suffix)
		return nil
	}
} // func (c *Converter) generateCommand(in, out string) []string
