// /home/krylon/go/src/github.com/blicero/podshrink/convert/convert.go
// -*- mode: go; coding: utf-8; -*-
// Created on 28. 08. 2023 by Benjamin Walkenhorst
// (c) 2023 Benjamin Walkenhorst
// Time-stamp: <2023-08-28 19:26:16 krylon>

package convert

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var suffixPat = regexp.MustCompile("[.]([^.]+$")

func generateCommand(in, out string) []string {
	var match []string

	if match = suffixPat.FindStringSubmatch(in); match == nil {
		return nil
	}

	var suffix = match[1]

	switch strings.ToLower(suffix) {
	case "mp3":
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
	}

	fmt.Fprintf(os.Stderr, "Cannot find decoder for %s\n", suffix)
	return nil
} // func generateCommand(in, out string) []string
