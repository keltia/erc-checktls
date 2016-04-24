// results.go

package main

import (
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
	defer fh.Close()

	res, err = ioutil.ReadAll(fh)
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
		report.Display()
	}
	return
}

