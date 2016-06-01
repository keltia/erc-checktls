// labs2pg.go

/*
This package implements reading the json from ssllabs-scan output
into our Pg database.
 */
package main

import (
	"flag"

	"github.com/keltia/erc-checktls/ssllabs"
//	"github.com/astaxie/beego/orm"
    _ "github.com/lib/pq" // import your used driver
	"fmt"
	"os"
	"encoding/csv"
	"log"
	"path/filepath"
)

var (
	contracts map[string]string
)

const (
	contractFile = "sites-list.csv"
	tlsVersion = "0.7.1"
)

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
	if fVerbose {
		fmt.Printf("%s version %s\n\n", filepath.Base(os.Args[0]), tlsVersion)
	}

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

	//fmt.Printf("all=%#v\n", allSites)

	// generate the final report
	final, err := NewTLSReport(allSites)

	if fType == "csv" {
		err := final.ToCSV(os.Stdout)
		if err != nil {
			log.Fatalf("Error can not generate CSV: %v", err)
		}
	} else {
		// XXX Early debugging
		fmt.Printf("%#v\n", final)
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
