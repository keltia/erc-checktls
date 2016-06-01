// cli.go

package main

import "flag"

var (
	fType          string
	fSiteName      string
	fIgnoreImirhil bool
	fVerbose       bool
	fReallyVerbose bool
)

func init() {
	flag.StringVar(&fType, "t", "text", "Type of report")
	flag.StringVar(&fSiteName, "S", "", "Display that site")
	flag.BoolVar(&fIgnoreImirhil, "I", false, "Do not fetch tls.imirhil.fr grade")
	flag.BoolVar(&fVerbose, "v", false, "Verbose mode")
	flag.BoolVar(&fReallyVerbose, "V", false, "More verbose mode")
}
