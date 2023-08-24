// /home/krylon/go/src/github.com/blicero/podshrink/logdomain/logdomain.go
// -*- mode: go; coding: utf-8; -*-
// Created on 16. 08. 2023 by Benjamin Walkenhorst
// (c) 2023 Benjamin Walkenhorst
// Time-stamp: <2023-08-24 16:05:20 krylon>

package logdomain

//go:generate stringer -type=ID

// ID is an id...
type ID uint8

// These constants represent the pieces of the application that need to log stuff.
const (
	Common ID = iota
	Walker
	Extractor
)

// AllDomains returns a slice of all the valid values for ID.
func AllDomains() []ID {
	return []ID{
		Common,
		Walker,
		Extractor,
	}
} // func AllDomains() []ID
