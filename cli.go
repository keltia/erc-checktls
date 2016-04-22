// cli.go

package main

import "flag"

var (
	fSiteName string
	fVerbose bool
	fReallyVerbose bool
)

func init() {
	flag.StringVar(&fSiteName, "S", "", "Display that site")
	flag.BoolVar(&fVerbose, "v", false, "Verbose mode")
	flag.BoolVar(&fReallyVerbose, "V", false, "More verbose mode")
}
