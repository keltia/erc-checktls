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
		host := report.Host
		grade := report.Endpoints[0].Grade
		log.Printf("Looking at %s… — grade %s", host, grade)
	}
	return
}


