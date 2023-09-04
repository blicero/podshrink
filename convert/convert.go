// /home/krylon/go/src/github.com/blicero/podshrink/convert/convert.go
// -*- mode: go; coding: utf-8; -*-
// Created on 28. 08. 2023 by Benjamin Walkenhorst
// (c) 2023 Benjamin Walkenhorst
// Time-stamp: <2023-09-04 10:50:47 krylon>

// Package convert implements the conversion of various audio formats to opus.
package convert

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
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
			err      error
			tags     *meta.FileMeta
			opusfile string
		)

		opusfile = suffixPat.ReplaceAllString(file, ".opus")

		var convertCmd = []string{
			"nice",
			"ffmpeg",
			"-i",
			file,
			opusfile,
		}

		if tags, err = c.meta.ReadTags(file); err != nil {
			c.log.Printf("[ERROR] Cannot extract metadata from %s: %s\n",
				file,
				err.Error())
			continue
		} else if tags == nil {
			c.log.Printf("[ERROR] Failed to extract data from %s\n",
				file)
			continue
		} else if err = c.execute(convertCmd); err != nil {
			c.log.Printf("[ERROR] Cannot convert %s to %s: %s\n",
				filepath.Base(file),
				filepath.Base(opusfile),
				err.Error())
			os.Remove(opusfile) // nolint: errcheck
			continue
		}

		// Tag the file.
		var tagCmd = []string{
			"opustags",
			"--in-place",
			"--add",
			fmt.Sprintf("ARTIST=%s", tags.Artist),
			"--add",
			fmt.Sprintf("ALBUM=%s", tags.Album),
			"--add",
			fmt.Sprintf("TITLE=%s", tags.Title),
			"--add",
			fmt.Sprintf("DATE=%d", tags.Year),
			"--add",
			fmt.Sprintf("TRACK=%d", tags.Track),
		}

		if tags.Cover != "" {
			tagCmd = append(tagCmd,
				"--set-cover",
				tags.Cover,
			)
		}

		tagCmd = append(tagCmd, opusfile)

		c.log.Printf("[DEBUG] Execute command: %v\n",
			tagCmd)

		if err = c.execute(tagCmd); err != nil {
			c.log.Printf("[ERROR] Cannot apply tags to %s: %s\n",
				opusfile,
				err.Error())
			continue
		}

		os.Remove(file) // nolint: errcheck
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
