// /home/krylon/go/src/github.com/blicero/podshrink/common/common.go
// -*- mode: go; coding: utf-8; -*-
// Created on 16. 08. 2023 by Benjamin Walkenhorst
// (c) 2023 Benjamin Walkenhorst
// Time-stamp: <2023-09-04 19:19:16 krylon>

package common

import (
	"crypto/sha512"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/blicero/podshrink/logdomain"
	"github.com/hashicorp/logutils"
	"github.com/odeke-em/go-uuid"
)

//go:generate ./build_time_stamp.pl

// Debug indicates whether to emit additional log messages and perform
// additional sanity checks.
// Version is the version number to display.
// AppName is the name of the application.
// TimestampFormat is the format string used to render datetime values.
// HeartBeat is the interval for worker goroutines to wake up and check
// their status.
const (
	Debug                    = true
	Version                  = "0.1.0"
	AppName                  = "podshrink"
	TimestampFormat          = "2006-01-02 15:04:05"
	TimestampFormatMinute    = "2006-01-02 15:04"
	TimestampFormatSubSecond = "2006-01-02 15:04:05.0000 MST"
	TimestampFormatDate      = "2006-01-02"
	HeartBeat                = time.Millisecond * 500
	RCTimeout                = time.Millisecond * 10
	Interval                 = time.Second * 120
	NetName                  = "unixpacket"
	BufferSize               = 65536 // 64 KiB
)

// LogLevels are the names of the log levels supported by the logger.
var LogLevels = []logutils.LogLevel{
	"TRACE",
	"DEBUG",
	"INFO",
	"WARN",
	"ERROR",
	"CRITICAL",
	"CANTHAPPEN",
	"SILENT",
}

// PackageLevels defines minimum log levels per package.
var PackageLevels = make(map[logdomain.ID]logutils.LogLevel, len(LogLevels))

// MinLogLevel is the mininum log level all loggers forward.
const MinLogLevel = "TRACE"

// EncJSON is the MIME type used for JSON payloads.
const EncJSON = "application/json"

func init() {
	for _, id := range logdomain.AllDomains() {
		PackageLevels[id] = MinLogLevel
	}
} // func init()

// SuffixPattern is a regular expression that matches the suffix of a file name.
// For "text.txt", it should match ".txt" and capture "txt".
var SuffixPattern = regexp.MustCompile("([.][^.]+)$")

// BaseDir is the folder where all application-specific files (database,
// log files, etc) are stored.
// LogPath is the file to the log path.
// CoverDir is the directory to temporarily store cover images.
var (
	BaseDir  = filepath.Join(os.Getenv("HOME"), ".local", "share", AppName)
	LogPath  = filepath.Join(BaseDir, AppName+".log")
	CoverDir = filepath.Join(BaseDir, "covers")
)

// SetBaseDir sets the BaseDir and related variables.
func SetBaseDir(path string) error {
	fmt.Printf("Setting BaseDir to %s\n", path)

	BaseDir = path
	LogPath = filepath.Join(BaseDir, AppName+".log")
	CoverDir = filepath.Join(BaseDir, "covers")

	if err := InitApp(); err != nil {
		fmt.Printf("Error initializing application environment: %s\n", err.Error())
		return err
	}

	return nil
} // func SetBaseDir(path string)

// InitApp performs some basic preparations for the application to run.
// Currently, this means creating the BASE_DIR folder.
func InitApp() error {
	var err error

	if err = os.Mkdir(BaseDir, 0755); err != nil {
		if !os.IsExist(err) {
			msg := fmt.Sprintf("Error creating BaseDir %s: %s", BaseDir, err.Error())
			return errors.New(msg)
		}
	}

	if err = os.Mkdir(CoverDir, 0755); err != nil {
		if !os.IsExist(err) {
			msg := fmt.Sprintf("Error creating BufferPath %s: %s", CoverDir, err.Error())
			return errors.New(msg)
		}
	}

	return nil
} // func InitApp() error

// GetLogger Tries to create a named logger instance and return it.
// If the directory to hold the log file does not exist, try to create it.
func GetLogger(dom logdomain.ID) (*log.Logger, error) {
	var (
		err     error
		logName string
	)

	if err = InitApp(); err != nil {
		return nil, fmt.Errorf("Error initializing application environment: %s", err.Error())
	}

	logName = fmt.Sprintf("%s.%s",
		AppName,
		dom)

	fmt.Printf("Creating Logger for %s\n", dom)

	var logfile *os.File
	logfile, err = os.OpenFile(LogPath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		msg := fmt.Sprintf("Error opening log file: %s\n", err.Error())
		fmt.Println(msg)
		return nil, errors.New(msg)
	}

	filter := &logutils.LevelFilter{
		Levels:   LogLevels,
		MinLevel: MinLogLevel,
		Writer:   io.MultiWriter(os.Stdout, logfile),
	}

	logger := log.New(filter, logName, log.Ldate|log.Ltime|log.Lshortfile)

	if logger == nil {
		fmt.Fprintf(
			os.Stderr,
			"log.New returned nil! Why?\n",
		)
	}

	return logger, nil
} // func GetLogger(name string) (*log.logger, error)

// GetUUID returns a randomized UUID
func GetUUID() string {
	return uuid.NewRandom().String()
} // func GetUUID() string

// TimeEqual returns true if the two timestamps are less than one second apart.
func TimeEqual(t1, t2 time.Time) bool {
	var delta = t1.Sub(t2)

	if delta < 0 {
		delta = -delta
	}

	return delta < time.Second
} // func TimeEqual(t1, t2 time.Time) bool

// GetChecksum computes the SHA512 checksum of the given data.
func GetChecksum(data []byte) (string, error) {
	var err error
	var hash = sha512.New()

	if _, err = hash.Write(data); err != nil {
		fmt.Fprintf( // nolint: errcheck
			os.Stderr,
			"Error computing checksum: %s\n",
			err.Error(),
		)
		return "", err
	}

	var checkSumBinary = hash.Sum(nil)
	var checkSumText = fmt.Sprintf("%x", checkSumBinary)

	return checkSumText, nil
} // func getChecksum(data []byte) (string, error)
