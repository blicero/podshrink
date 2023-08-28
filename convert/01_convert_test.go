// /home/krylon/go/src/github.com/blicero/podshrink/convert/01_convert_test.go
// -*- mode: go; coding: utf-8; -*-
// Created on 28. 08. 2023 by Benjamin Walkenhorst
// (c) 2023 Benjamin Walkenhorst
// Time-stamp: <2023-08-28 21:00:50 krylon>

package convert

import "testing"

var conv *Converter

func TestCreateConverter(t *testing.T) {
	var err error

	if conv, err = New(); err != nil {
		conv = nil
		t.Fatalf("Cannot create Converter: %s",
			err.Error())
	}
} // func TestCreateConverter(t *testing.T)

func TestFindConverter(t *testing.T) {
	type testCase struct {
		input      string
		expect     string
		expectFail bool
	}

	if conv == nil {
		t.SkipNow()
	}

	var cases = []testCase{
		testCase{
			input:  "/data/Podcasts/bla.mp3",
			expect: "mpg123",
		},
		testCase{
			input:  "/data/tmp/BLA.OGG",
			expect: "ogg123",
		},
		testCase{
			input:  "/tmp/wer-das-liest-ist-doof.m4a",
			expect: "ffmpeg",
		},
		testCase{
			input:      "/home/myawesomeusername/Podcasts/some-file.txt",
			expectFail: true,
		},
	}

	for _, c := range cases {
		var cmd = conv.generateCommand(c.input, "bla")

		if len(cmd) == 0 {
			if !c.expectFail {
				t.Errorf("Did not find converter for %s",
					c.input)
			}
		} else if cmd[0] != c.expect {
			t.Errorf("Unexpected converter found for %s: %s (expected %s)",
				c.input,
				cmd[0],
				c.expect)
		}
	}
} // func TestFindConverter(t *testing.T)
