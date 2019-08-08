// main.go

/*
This package implements reading the json from ssllabs-scan output
and generating a csv file.
*/
package main // import "github.com/keltia/erc-checktls"

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/keltia/cryptcheck"
	"github.com/keltia/observatory"
	"github.com/keltia/ssllabs"
	"github.com/pkg/errors"

	TLS "github.com/keltia/erc-checktls"
)

var (
	// MyName is obvious
	MyName = filepath.Base(os.Args[0])

	logLevel = 0
)

const (
	// MyVersion uses semantic versioning.
	MyVersion = "0.69.99"
)

// checkOutput checks whether we want to specify an output file
func checkOutput(fn string) *os.File {
	var err error

	OutputFH := os.Stdout

	// Open output file
	if fn != "" {
		verbose("Output file is %s\n", fn)

		if fn != "-" {
			fn = fn + "." + fType
			OutputFH, err = os.Create(fn)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating %s\n", fn)
				return nil
			}
		}
	}
	debug("output=%v\n", OutputFH)
	return OutputFH
}

// init is for pg connection and stuff
func init() {
	flag.Usage = Usage
}

func checkFlags(a []string) error {
	// Basic argument check
	if a == nil || len(a) != 1 {
		return fmt.Errorf("you must specify an input file!")
	}

	// Set logging level
	if fVerbose {
		logLevel = 1
	}

	if fDebug {
		fVerbose = true
		logLevel = 2
		debug("debug mode\n")
	}

	// Expect less perf is fJobs is too high (contention)
	if fJobs > runtime.NumCPU() {
		fmt.Fprintf(os.Stderr, "Warning, '-j %d' higher than %d, possible perf issue\n",
			fJobs, runtime.NumCPU())
	}

	return nil
}

func checkInput(file string) error {
	if file == "" {
		return errors.New("No file found\n")
	}

	_, err := os.Stat(file)
	return errors.Wrap(err, "checkInput")
}

// getResults read the JSON array generated and gone through jq
func getResults(file string) (res []byte, err error) {
	fh, err := os.Open(file)
	if err != nil {
		return res, errors.Wrapf(err, "can not open %s", file)
	}
	defer fh.Close()

	res, err = ioutil.ReadAll(fh)
	return res, errors.Wrapf(err, "can not read json %s", file)
}

// Most of the work is here
func realmain(args []string) int {
	// Announce ourselves
	fmt.Printf("%s version %s/j%d - Imirhil/%s SSLLabs/%s Mozilla/%s\n\n",
		filepath.Base(MyName), MyVersion, fJobs,
		cryptcheck.MyVersion, ssllabs.MyVersion, observatory.MyVersion)

	if err := checkFlags(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err.Error())
		return 1
	}

	// Open output
	var fhBase, fhSumm *os.File

	if fOutput == "-" || fOutput == "" {
		fhBase, fhSumm = os.Stdout, os.Stdout
	} else {
		var tag string

		if fAddDate {
			tag = "-" + makeDate()
		}
		fhBase = checkOutput(fOutput + tag)
		fhSumm = checkOutput(fOutput + "-summary" + tag)
	}

	// Check for output errors
	if fhBase == nil || fhSumm == nil {
		fmt.Fprintf(os.Stderr, "Can not open files.\n")
		return 1
	}

	// Initialise the library
	TLS.Init(TLS.Config{
		LogLevel:      logLevel,
		IgnoreMozilla: fIgnoreMozilla,
		IgnoreImirhil: fIgnoreImirhil,
	})

	file := args[0]
	if err := checkInput(file); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	raw, err := getResults(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't read %s: %v\n", file, err)
		return 1
	}

	// raw is the []byte array to be deserialized into Hosts
	allSites, err := ssllabs.ParseResults(raw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't parse %s: %v\n", file, err)
		return 1
	}

	if fCmdWild {
		str := displayWildcards(allSites)
		debug("str=%s\n", str)
		fmt.Fprintf(fhBase, "All wildcards certs:\n%s", str)
		return 0
	}

	// generate the final report & summary
	final, err := TLS.NewReport(allSites, fJobs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error analyzing report: %v\n", err)
		return 1
	}

	verbose("SSLabs engine: %s\n", final.SSLLabs)

	switch fType {
	case "csv":
		if err := final.WriteCSV(fhBase); err != nil {
			fmt.Fprintf(os.Stderr, "WriteCSV failed: %v\n", err)
			return 1
		}
		if err := final.WriteCSVSummary(fhSumm); err != nil {
			fmt.Fprintf(os.Stderr, "WriteCSVSummary failed: %v\n", err)
			return 1
		}
	case "html":
		if err := final.WriteHTML(fhBase); err != nil {
			fmt.Fprintf(os.Stderr, "WriteHTML failed: %v\n", err)
			return 1
		}
		if err = final.WriteHTMLSummary(fhSumm); err != nil {
			fmt.Fprintf(os.Stderr, "WriteHTML failed: %v\n", err)
			return 1
		}

	default:
		// XXX Early debugging
		fmt.Printf("%#v\n", final)
	}

	return 0
}

// main is the the starting point
func main() {
	flag.Parse()

	_ = realmain(flag.Args())
}
