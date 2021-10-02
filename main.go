package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/St0rmPetrel/chain-log-cat/processors/filetracker"
)

var (
	root   string
	patern string
	hours  int
)

const (
	rootFlagName = "dir"
	rootDefault  = "."

	paternFlagName = "regexp"
	paternDefault  = ""

	hoursFlagName = "h"
)

func init() {
	flag.StringVar(
		&root, rootFlagName, rootDefault,
		//TODO fill it
		"Usage of root flag",
	)
	flag.StringVar(
		&patern, paternFlagName, paternDefault,
		//TODO fill it
		"Usage of patern flag",
	)
	flag.IntVar(
		&hours, hoursFlagName, 1,
		//TODO fill it
		"Usage of hours flag",
	)
	flag.Parse()
}

func main() {
	ft := filetracker.New(root, patern, time.Duration(hours)*time.Hour)
	b, _ := ft.TrackChanges()
	fmt.Print("\n" + string(b))
}
