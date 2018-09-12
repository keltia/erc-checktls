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
	contractFile = "sites-list.csv"
	htmlTemplate = "templ.html"
	// MyVersion uses semantic versioning.
	MyVersion = "0.60.0"
)

// getContract retrieve the site's contract from the DB
func readContractFile(box packr.Box) (contracts map[string]string, err error) {
	debug("reading contracts\n")
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
	debug("output=%v\n", fOutputFH)
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
}

func checkFlags() {
	// Basic argument check
	if len(flag.Args()) != 1 {
		fatalf("Error: you must specify an input file!")
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
}

// main is the the starting point
func main() {
	// Announce ourselves
	fmt.Printf("%s version %s - Imirhil/%s SSLLabs/%s Mozilla/%s\n\n", filepath.Base(os.Args[0]),
		MyVersion, cryptcheck.MyVersion, ssllabs.MyVersion, observatory.MyVersion)

	checkFlags()

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

	tmpls, err = loadTemplates(box)
	if err != nil {
		fatalf("Error: can not read HTML templates from 'files/': %v", err)
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
		if err = final.ToCSV(fOutputFH); err != nil {
			fatalf("Error can not generate CSV: %v", err)
		}
		fmt.Printf("\nTLS Summary\n")
		if err = writeSummary(os.Stdout, tlsKeys, cntrs); err != nil {
			fmt.Fprintf(os.Stderr, "can not generate TLS summary: %v", err)
		}
		fmt.Printf("\nHTTP Summary\n")
		if err = writeSummary(os.Stdout, httpKeys, https); err != nil {
			fmt.Fprintf(os.Stderr, "can not generate HTTP summary: %v", err)
		}
	case "html":
		tmpls, err := loadTemplates(box)
		if err != nil {
			fatalf("Error: can not read HTML templates from 'files/': %v", err)
		}
		debug("tmpls=%v\n", tmpls)
		if err := final.ToHTML(fOutputFH, tmpls["templ"]); err != nil {
			fatalf("Can not write HTML: %v", err)
		}
		if fSummary != "" {
			fOutputFH = checkOutput(filepath.Join(fSummary, ".html"))
			if err := writeHTMLSummary(fOutputFH, cntrs, https); err != nil {
				fatalf("summary failed: %v\n", err)
			}
		}
	default:
		// XXX Early debugging
		fmt.Printf("%#v\n", final)
		fmt.Printf("%s\n", displayCategories(cntrs))

	}
}
