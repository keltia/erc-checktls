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
	fAddDate       bool
	fDebug         bool
	fJobs          int
	fType          string
	fOutput        string
	fSiteName      string
	fIgnoreImirhil bool
	fIgnoreMozilla bool
	fVerbose       bool

	fCmdWild bool
)

const (
	cliUsage = `%s version %s - Cryptcheck/%s SSLLabs/%s Mozilla/%s

Usage: %s [-dhvDIM] [-j n] [-t text|csv|html] [-S site] [-o file] file.json
       %s [-vD] -wild file[.json]

Default output: stdout for both complete & summary.
Use -o file to get file.<fmt> and file-summary.<fmt>.
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
	flag.BoolVar(&fAddDate, "date", false, "Add date to output filename")
	flag.IntVar(&fJobs, "j", runtime.NumCPU(), "# of parallel jobs")
	flag.StringVar(&fOutput, "o", "-", "Save into files (default stdout)")
	flag.StringVar(&fType, "t", "csv", "Type of report")
	flag.StringVar(&fSiteName, "S", "", "Display that site")
	flag.BoolVar(&fIgnoreImirhil, "I", false, "Do not fetch tls.imirhil.fr grade")
	flag.BoolVar(&fIgnoreMozilla, "M", false, "Do not fetch Mozilla Observatory data")
	flag.BoolVar(&fDebug, "D", false, "Debug mode")
	flag.BoolVar(&fVerbose, "v", false, "Verbose mode")

	flag.BoolVar(&fCmdWild, "wild", false, "Display wildcards")
}
