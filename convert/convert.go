// /home/krylon/go/src/github.com/blicero/podshrink/convert/convert.go
// -*- mode: go; coding: utf-8; -*-
// Created on 28. 08. 2023 by Benjamin Walkenhorst
// (c) 2023 Benjamin Walkenhorst
// Time-stamp: <2023-09-01 10:30:30 krylon>

// Package convert implements the conversion of various audio formats to opus.
package convert

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/blicero/podshrink/common"
	"github.com/blicero/podshrink/logdomain"
	"github.com/blicero/podshrink/meta"
)

var suffixPat = regexp.MustCompile("[.]([^.]+)$")

// TmpDir is the path where the waveform files are stored temporarily.
var TmpDir = "/data/ram"

// Converter wraps the state associated with converting audio files.
type Converter struct {
	log   *log.Logger
	alive atomic.Bool
	cnt   int
	fileQ <-chan string
	meta  *meta.Extractor
}

// New creates a new Converter
func New(cnt int, queue <-chan string) (*Converter, error) {
	var (
		err error
		c   = &Converter{
			cnt:   cnt,
			fileQ: queue,
		}
	)

	if c.log, err = common.GetLogger(logdomain.Converter); err != nil {
		return nil, err
	} else if c.meta, err = meta.NewExtractor(); err != nil {
		c.log.Printf("[ERROR] Cannot create Extractor: %s\n",
			err.Error())
		return nil, err
	}

	c.alive.Store(true)

	return c, nil
} // func New(cnt int, queue <-chan string) (*Converter, error)

func (c *Converter) Run() {
	var wg sync.WaitGroup

	for i := 1; i <= c.cnt; i++ {
		go c.worker(i, &wg)
		wg.Add(1)
	}

	wg.Wait()
}

func (c *Converter) worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	for file := range c.fileQ {
		var (
			err       error
			tags      *meta.FileMeta
			tmpfile   = filepath.Join(TmpDir, filepath.Base(file))
			decodeCmd = c.generateCommand(file, tmpfile)
		)

		if len(decodeCmd) == 0 {
			c.log.Printf("[INFO] Did not find decoder for %s\n",
				file)
			continue
		} else if err = c.execute(decodeCmd); err != nil {
			c.log.Printf("[ERROR] Failed to decode %s: %s\n",
				file,
				err.Error())
			os.Remove(tmpfile) // nolint: errcheck
			continue
		} else if tags, err = c.meta.ReadTags(file); err != nil {
			c.log.Printf("[ERROR] Cannot extract metadata from %s: %s\n",
				file,
				err.Error())
			continue
		} else if tags == nil {
			c.log.Printf("[ERROR] Failed to extract metadata from %s\n",
				file)
			continue
		}

		var opus = suffixPat.ReplaceAllString(file, ".opus")
		c.log.Printf("[DEBUG] Convert %s to %s\n",
			file,
			opus)

		var encodeCmd = []string{
			"opusenc",
			"--speech",
			"--title",
			tags.Title,
			"--album",
			tags.Album,
			"--tracknumber",
			strconv.Itoa(tags.Track),
			"--date",
			fmt.Sprintf("%04d", tags.Year),
		}

		if tags.Cover != "" {
			encodeCmd = append(encodeCmd,
				"--picture",
				tags.Cover)
		}

		encodeCmd = append(encodeCmd,
			tmpfile,
			opus)

		if err = c.execute(encodeCmd); err != nil {
			c.log.Printf("[ERROR] Failed to encode %s to %s: %s\n",
				file, opus,
				err.Error())
			os.Remove(opus) // nolint: errcheck
		} else {
			os.Remove(tmpfile) // nolint: errcheck
			os.Remove(file)    // nolint: errcheck
		}

	}
} // func (c *Converter) worker(id int)

func (c *Converter) execute(cmd []string) error {
	proc := exec.Command(cmd[0], cmd[1:]...)

	return proc.Run()
} // func (c *Converter) decode(cmd []string) error

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
