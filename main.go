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
func checkOutput(fOutput string) (fOutputFH *os.File) {
	var err error

	fOutputFH = os.Stdout

	// Open output file
	if fOutput != "" {
		verbose("Output file is %s\n", fOutput)

		if fOutput != "-" {
			fOutputFH, err = os.Create(fOutput)
			if err != nil {
				fatalf("Error creating %s\n", fOutput)
			}
		}
	}
	debug("output=%v\n", fOutputFH)
	return
}

// init is for pg connection and stuff
func init() {
	flag.Usage = Usage
	flag.Parse()
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

// main is the the starting point
func main() {
	// Announce ourselves
	fmt.Printf("%s version %s/j%d - Imirhil/%s SSLLabs/%s Mozilla/%s\n\n",
		filepath.Base(os.Args[0]), MyVersion, fJobs,
		cryptcheck.MyVersion, ssllabs.MyVersion, observatory.MyVersion)

	err := checkFlags(flag.Args())
	if err != nil {
		fatalf("Error: %v", err.Error())
	}

	file := flag.Arg(0)

	raw, err := getResults(file)
	if err != nil {
		fatalf("Can't read %s: %v", file, err.Error())
	}

	// raw is the []byte array to be deserialized into Hosts
	allSites, err := ssllabs.ParseResults(raw)
	if err != nil {
		fatalf("Can't parse %s: %v", file, err.Error())
	}

	err = loadResources(resourcesPath)
	if err != nil {
		fatalf("Can't load resources %s: %v", resourcesPath, err)
	}

	// Open output file
	fOutputFH := checkOutput(fOutput)

	if fCmdWild {
		str := displayWildcards(allSites)
		debug("str=%s\n", str)
		fmt.Fprintf(fOutputFH, "All wildcards certs:\n%s", str)
		os.Exit(0)
	}

	// generate the final report & summary
	final, err := NewTLSReport(allSites)
	if err != nil {
		fatalf("error analyzing report: %v", err)
	}

	// Gather statistics for summaries
	cntrs := categoryCounts(allSites)
	https := httpCounts(final)

	verbose("SSLabs engine: %s\n", final.SSLLabs)

	switch fType {
	case "csv":
		err = WriteCSV(fOutputFH, final, cntrs, https)
		if err != nil {
			fatalf("WriteCSV failed: %v", err)
		}
	case "html":
		err = WriteHTML(fOutputFH, final, cntrs, https)
		if err != nil {
			fatalf("WriteHTML failed: %v", err)
		}
	default:
		// XXX Early debugging
		fmt.Printf("%#v\n", final)
		fmt.Printf("%s\n", displayCategories(cntrs))

	}
}
