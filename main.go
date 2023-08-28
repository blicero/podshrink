// /home/krylon/go/src/github.com/blicero/podshrink/main.go
// -*- mode: go; coding: utf-8; -*-
// Created on 28. 08. 2023 by Benjamin Walkenhorst
// (c) 2023 Benjamin Walkenhorst
// Time-stamp: <2023-08-28 21:03:18 krylon>

package main

import (
	"fmt"

	"github.com/blicero/podshrink/common"
)

func main() {
	fmt.Printf("%s %s - built on %s\n",
		common.AppName,
		common.Version,
		common.BuildStamp,
	)
}
