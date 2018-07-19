// main.go

/*
This package implements reading the json from ssllabs-scan output
and generating a csv file.
*/
package main // import "github.com/keltia/erc-checktls"

import (
	"flag"

	"bytes"
	"encoding/csv"
	"github.com/gobuffalo/packr"
	"github.com/keltia/cryptcheck"
	"github.com/keltia/erc-checktls/ssllabs"
	"github.com/pkg/errors"
	"log"
	"os"
	"path/filepath"
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
	MyVersion = "0.23.0"
)

type Context struct {
	proxyauth string
}

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
		if fVerbose {
			log.Printf("Output file is %s\n", fOutput)
		}

		if fOutput != "-" {
			fOutputFH, err = os.Create(fOutput)
			if err != nil {
				log.Fatalf("Error creating %s\n", fOutput)
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
	// set default database
	//orm.RegisterDataBase("default", "postgres", "roberto", 30)
}

// main is the the starting point
func main() {
	flag.Usage = Usage
	flag.Parse()

	// Announce ourselves
	verbose("%s version %s - Imirhil %s\n\n", filepath.Base(os.Args[0]),
		MyVersion, cryptcheck.Version())

	// Initiase context
	ctx := &Context{}

	// Basic argument check
	if len(flag.Args()) != 1 {
		log.Fatalf("Error: you must specify an input file!")
	}

	file := flag.Arg(0)

	raw, err := getResults(file)
	if err != nil {
		log.Fatalf("Can't read %s: %v", file, err.Error())
	}

	// raw is the []byte array to be deserialized into LabsReports
	allSites, err := ssllabs.ParseResults(raw)
	if err != nil {
		log.Fatalf("Can't parse %s: %v", file, err.Error())
	}

	// We embed the file now
	box := packr.NewBox("./files")

	// We need that for the reports
	contracts, err = readContractFile(box)
	if err != nil {
		log.Fatalf("Error: can not read contract file %s: %v", contractFile, err)
	}

	// Set logging level
	if fVerbose {
		logLevel = 1
	}

	if fDebug {
		fVerbose = true
		logLevel = 2
		debug("debug mode")
	}

	//fmt.Printf("all=%#v\n", allSites)

	// generate the final report
	final, err := NewTLSReport(ctx, allSites)

	verbose("SSLabs engine: %s/%s", final.EngineVersion, final.CriteriaVersion)

	// Open output file
	fOutputFH := checkOutput(fOutput)

	if fType == "csv" {
		err := final.ToCSV(fOutputFH)
		if err != nil {
			log.Fatalf("Error can not generate CSV: %v", err)
		}
	} else {
		// XXX Early debugging
		debug("%#v\n", final)
	}
	if fVerbose {
		cntrs := categoryCounts(allSites)
		if fType == "csv" {
			categoriesCSV(cntrs, os.Stdout)
		} else {
			log.Printf("%s\n", displayCategories(cntrs))
		}
	}
}
