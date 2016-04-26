// ssllabs.go

/*
SSLLabs-related functions.
*/
package ssllabs

import (
	"log"
	"encoding/json"

	"github.com/keltia/erc-checktls/imirhil"
)

// Display for one report
func (rep *LabsReport) Display() {
	host := rep.Host
	grade := rep.Endpoints[0].Grade
	details := rep.Endpoints[0].Details
	if fIgnoreImirhil {
		imirhil := imirhil.GetScore(host)
		log.Printf("Looking at %s/%s — grade %s/%s", host, contracts[host], grade, imirhil)
	} else {
		log.Printf("Looking at %s/%s — grade %s", host, contracts[host], grade)
	}
	if fVerbose {
		log.Printf("  Ciphers: %d", details.Suites.len())
	} else if fReallyVerbose {
		for _, cipher := range details.Suites.List {
			log.Printf("  %s: %d bits", cipher.Name, cipher.CipherStrength)
		}
	}
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
