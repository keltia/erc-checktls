// cli.go

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/keltia/cryptcheck"
	"github.com/keltia/observatory"
	"github.com/keltia/ssllabs"
)

var (
	fDebug         bool
	fType          string
	fOutput        string
	fSummary       string
	fSiteName      string
	fIgnoreImirhil bool
	fIgnoreMozilla bool
	fVerbose       bool
	fReallyVerbose bool

	fCmdWild bool
)

const (
	cliUsage = `%s version %s - Imirhil/%s SSLLabs/%s Mozilla/%s

Usage: %s [-hvIMV] [-t text|csv|html] [-s file] [-o file] [-wild] file[.json]

`
)

// Usage string override.
var Usage = func() {
	fmt.Fprintf(os.Stderr, cliUsage, MyName,
		MyVersion, cryptcheck.MyVersion, ssllabs.MyVersion, observatory.MyVersion,
		MyName)
	flag.PrintDefaults()
}

func init() {
	flag.StringVar(&fOutput, "o", "-", "Save into file (default stdout)")
	flag.StringVar(&fSummary, "s", "summary", "Save summary there")
	flag.StringVar(&fType, "t", "csv", "Type of report")
	flag.StringVar(&fSiteName, "S", "", "Display that site")
	flag.BoolVar(&fIgnoreImirhil, "I", false, "Do not fetch tls.imirhil.fr grade")
	flag.BoolVar(&fIgnoreMozilla, "M", false, "Do not fetch Mozilla Observatory data")
	flag.BoolVar(&fDebug, "D", false, "Debug mode")
	flag.BoolVar(&fVerbose, "v", false, "Verbose mode")
	flag.BoolVar(&fReallyVerbose, "V", false, "More verbose mode")

	flag.BoolVar(&fCmdWild, "wild", false, "Display wildcards")
}
