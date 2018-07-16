// main.go

/*
This package implements reading the json from ssllabs-scan output
and generating a csv file.
*/
package main // import "github.com/keltia/erc-checktls"

import (
	"flag"

	"encoding/csv"
	"github.com/keltia/cryptcheck"
	"github.com/keltia/erc-checktls/ssllabs"
	"log"
	"os"
	"path/filepath"
)

var (
	MyName = filepath.Base(os.Args[0])

	contracts map[string]string

	logLevel = 0
)

const (
	contractFile = "sites-list.csv"
	// MyVersion uses semantic versioning.
	MyVersion = "0.20.0"
)

type Context struct {
	proxyauth string
}

// getContract retrieve the site's contract from the DB
func readContractFile(file string) (contracts map[string]string, err error) {
	var (
		fh *os.File
	)

	_, err = os.Stat(file)
	if err != nil {
		return
	}

	if fh, err = os.Open(file); err != nil {
		return
	}
	defer fh.Close()

	all := csv.NewReader(fh)
	allSites, err := all.ReadAll()

	contracts = make(map[string]string)
	for _, site := range allSites {
		contracts[site[0]] = site[1]
	}
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

	// We need that for the reports
	contracts, err = readContractFile(contractFile)
	if err != nil {
		log.Fatalf("Error: can not read contract file %s: %v", contractFile, err)
	}

	// Set logging level
	if fVerbose {
		logLevel = 1
	}

	if fDebug {
		logLevel = 2
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
		categoryCounts(allSites)
		if fType == "csv" {
			categoriesCSV(cntrs, os.Stdout)
		} else {
			log.Printf("%s\n", displayCategories(cntrs))
		}
	}
}
