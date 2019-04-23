// main.go

/*
This package implements reading the json from ssllabs-scan output
and generating a csv file.
*/
package main // import "github.com/keltia/erc-checktls"

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/keltia/cryptcheck"
	"github.com/keltia/observatory"
	"github.com/keltia/ssllabs"
	"github.com/pkg/errors"
)

var (
	// MyName is obvious
	MyName = filepath.Base(os.Args[0])

	contracts map[string]string
	tmpls     map[string]string

	logLevel = 0
)

const (
	// MyVersion uses semantic versioning.
	MyVersion = "0.63.0"
)

// checkOutput checks whether we want to specify an output file
func checkOutput(fOutput string) *os.File {
	var err error

	OutputFH := os.Stdout

	// Open output file
	if fOutput != "" {
		verbose("Output file is %s\n", fOutput)

		if fOutput != "-" && fOutput != "" {
			OutputFH, err = os.Create(fOutput)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating %s\n", fOutput)
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
	return nil
}

func checkInput(file string) error {
	if file == "" {
		return errors.New("No file found\n")
	}

	_, err := os.Stat(file)
	return errors.Wrap(err, "checkInput")
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

	contracts, tmpls, err = loadResources()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't load resources %s: %v\n", resourcesPath, err)
		return 1
	}

	// Open output file
	OutputFH := checkOutput(fOutput)
	if OutputFH == nil {
		fmt.Fprintf(os.Stderr, "error output: %v\n", err)
		return 1
	}

	if fCmdWild {
		str := displayWildcards(allSites)
		debug("str=%s\n", str)
		fmt.Fprintf(OutputFH, "All wildcards certs:\n%s", str)
		return 0
	}

	// generate the final report & summary
	final, err := NewTLSReport(allSites)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error analyzing report: %v\n", err)
		return 1
	}

	verbose("SSLabs engine: %s\n", final.SSLLabs)

	switch fType {
	case "csv":
		if err := final.WriteCSV(OutputFH); err != nil {
			fmt.Fprintf(os.Stderr, "WriteCSV failed: %v\n", err)
			return 1
		}
	case "html":
		if err := final.WriteHTML(OutputFH); err != nil {
			fmt.Fprintf(os.Stderr, "WriteHTML failed: %v\n", err)
			return 1
		}
	default:
		// XXX Early debugging
		fmt.Printf("%#v\n", final)
		fmt.Printf("%s\n", displayCategories(final.cntrs))
	}
	return 0
}

// main is the the starting point
func main() {
	flag.Parse()

	_ = realmain(flag.Args())
}
