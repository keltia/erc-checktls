// main.go

/*
This package implements reading the json from ssllabs-scan output
and generating a csv file.
*/
package main // import "github.com/keltia/erc-checktls"

import (
	"bytes"
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gobuffalo/packr"
	"github.com/keltia/cryptcheck"
	"github.com/keltia/ssllabs"
	"github.com/pkg/errors"
)

var (
	// MyName is obvious
	MyName = filepath.Base(os.Args[0])

	contracts map[string]string

	logLevel = 0
)

const (
	contractFile = "sites-list.csv"
	// MyVersion uses semantic versioning.
	MyVersion = "0.50.0-ssllabs"
)

// getContract retrieve the site's contract from the DB
func readContractFile(box packr.Box) (contracts map[string]string, err error) {
	cf := box.Bytes(contractFile)
	fh := bytes.NewBuffer(cf)

	all := csv.NewReader(fh)
	allSites, err := all.ReadAll()
	if err != nil {
		return nil, errors.Wrap(err, "ReadAll")
	}

	contracts = make(map[string]string)
	for _, site := range allSites {
		contracts[site[0]] = site[1]
	}
	err = nil
	return
}

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
	return
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

// init is for pg connection and stuff
func init() {
	flag.Usage = Usage
	flag.Parse()

	// Basic argument check
	if len(flag.Args()) != 1 {
		fatalf("Error: you must specify an input file!")
	}

	// Announce ourselves
	fmt.Printf("%s version %s - Imirhil %s SSLLabs %s\n\n", filepath.Base(os.Args[0]),
		MyVersion, cryptcheck.Version(), "v3")

}

// main is the the starting point
func main() {

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

	// We embed the file now
	box := packr.NewBox("./files")

	// We need that for the reports
	contracts, err = readContractFile(box)
	if err != nil {
		fatalf("Error: can not read contract file %s: %v", contractFile, err)
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

	if fType == "csv" {
		if err = final.ToCSV(fOutputFH); err != nil {
			fatalf("Error can not generate CSV: %v", err)
		}
		fmt.Printf("\nTLS Summary\n")
		if err = writeSummary(tlsKeys, cntrs, os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "can not generate TLS summary: %v", err)
		}
		fmt.Printf("\nHTTP Summary\n")
		if err = writeSummary(httpKeys, https, os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "can not generate HTTP summary: %v", err)
		}
	} else {
		// XXX Early debugging
		fmt.Printf("%#v\n", final)
		fmt.Printf("%s\n", displayCategories(cntrs))
	}
}
