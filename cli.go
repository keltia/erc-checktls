// cli.go

package main

import "flag"

var (
	fVerbose bool
)

func init() {
	flag.BoolVar(&fVerbose, "v", false, "Verbose mode")
}
