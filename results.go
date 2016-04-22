// results.go

package main

import (
	"log"
	"os"
	"io/ioutil"
	"encoding/json"
)

// getResults read the JSON array generated and gone through jq
func getResults(file string) (res []byte, err error) {
	fh, err := os.Open(file)
	if err != nil {
		return
	}

	res, err = ioutil.ReadAll(fh)
	if err != nil {
		return
	}

	return
}

// parseResults unmarshals the json payload
func parseResults(content []byte) (rep *[]LabsReport, err error) {
	var data []LabsReport

	err = json.Unmarshal(content, &data)
	rep = &data
	return
}

// insertResults saves all reports
func insertResults(reports *[]LabsReport) (err error) {
	rep := *reports
	for _, report := range rep {
		if fSiteName != "" {
			if report.Host != fSiteName {
				continue
			}
		}
		displayReport(report)
	}
	return
}

// displayReport displays one report
func displayReport(report LabsReport) {
	host := report.Host
	grade := report.Endpoints[0].Grade
	details := report.Endpoints[0].Details
	log.Printf("Looking at %s… — grade %s", host, grade)
	if fVerbose {
		log.Printf("  Ciphers: %d", details.Suites.len())
	} else if fReallyVerbose {
		for _, cipher := range details.Suites.List {
			log.Printf("  %s: %d bits", cipher.Name, cipher.CipherStrength)
		}
	}
}

