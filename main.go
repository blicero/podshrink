// /home/krylon/go/src/github.com/blicero/podshrink/main.go
// -*- mode: go; coding: utf-8; -*-
// Created on 28. 08. 2023 by Benjamin Walkenhorst
// (c) 2023 Benjamin Walkenhorst
// Time-stamp: <2023-08-31 21:00:08 krylon>

package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"runtime"

	"github.com/blicero/podshrink/common"
	"github.com/blicero/podshrink/convert"
	"github.com/blicero/podshrink/walker"
)

var audioPat = regexp.MustCompile("[.](?i:mp3|ogg|oga|mpga|m4a)$")

func isAudioFile(f string) bool {
	return audioPat.MatchString(f)
} // func isAudioFile(f string) bool

func main() {
	fmt.Printf("%s %s - built on %s\n",
		common.AppName,
		common.Version,
		common.BuildStamp,
	)

	var (
		err       error
		trav      *walker.Walker
		conv      *convert.Converter
		queue     chan string
		workerCnt int
	)

	flag.IntVar(&workerCnt, "workers", runtime.NumCPU(), "Number of worker goroutines to use in parallel")
	flag.StringVar(&convert.TmpDir, "tmp", convert.TmpDir, "Where to store temporary files")

	flag.Parse()

	queue = make(chan string, workerCnt)

	if trav, err = walker.Create(isAudioFile, queue, flag.Args()); err != nil {
		fmt.Printf("Cannot create tree walker: %s\n",
			err.Error())
		os.Exit(1)
	} else if conv, err = convert.New(workerCnt, queue); err != nil {
		fmt.Printf("Cannot create converter: %s\n",
			err.Error())
		os.Exit(1)
	}

	go trav.Run() // nolint: errcheck
	conv.Run()
}
