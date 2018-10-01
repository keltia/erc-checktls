// cli.go

package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/keltia/cryptcheck"
	"github.com/keltia/observatory"
	"github.com/keltia/ssllabs"
)

var (
	fDebug         bool
	fJobs          int
	fType          string
	fOutput        string
	fSummary       string
	fSiteName      string
	fIgnoreImirhil bool
	fIgnoreMozilla bool
	fVerbose       bool

	fCmdWild bool
)

const (
	cliUsage = `%s version %s - Imirhil/%s SSLLabs/%s Mozilla/%s

Usage: %s [-hvDIM] [-j n] [-t text|csv|html] [-s file] [-S site] [-o file] file[.json]
       %s [-vD] -wild file[.json]
`
)

// Usage string override.
var Usage = func() {
	fmt.Fprintf(os.Stderr, cliUsage, MyName,
		MyVersion, cryptcheck.MyVersion, ssllabs.MyVersion, observatory.MyVersion,
		MyName, MyName)
	flag.PrintDefaults()
}

func init() {
	flag.IntVar(&fJobs, "j", runtime.NumCPU(), "# of parallel jobs")
	flag.StringVar(&fOutput, "o", "-", "Save into file (default stdout)")
	flag.StringVar(&fSummary, "s", "summaries", "Save summary there")
	flag.StringVar(&fType, "t", "csv", "Type of report")
	flag.StringVar(&fSiteName, "S", "", "Display that site")
	flag.BoolVar(&fIgnoreImirhil, "I", false, "Do not fetch tls.imirhil.fr grade")
	flag.BoolVar(&fIgnoreMozilla, "M", false, "Do not fetch Mozilla Observatory data")
	flag.BoolVar(&fDebug, "D", false, "Debug mode")
	flag.BoolVar(&fVerbose, "v", false, "Verbose mode")

	flag.BoolVar(&fCmdWild, "wild", false, "Display wildcards")
}
