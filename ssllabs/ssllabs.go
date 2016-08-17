// ssllabs.go

/*
Package ssllabs contains SSLLabs-related functions.
*/
package ssllabs

import (
	"encoding/json"
	"log"
)

// Display for one report
func (rep *LabsReport) String() {
	host := rep.Host
	grade := rep.Endpoints[0].Grade
	//details := rep.Endpoints[0].Details
	log.Printf("Looking at %s — grade %s", host, grade)
	/*	if fVerbose {
			log.Printf("  Ciphers: %d", details.Suites.len())
		} else if fReallyVerbose {
			for _, cipher := range details.Suites.List {
				log.Printf("  %s: %d bits", cipher.Name, cipher.CipherStrength)
			}
		} */
}

// Insert for one report
func (rep *LabsReport) Insert() {

}

// ParseResults unmarshals the json payload
func ParseResults(content []byte) (*LabsReports, error) {
	var data LabsReports

	err := json.Unmarshal(content, &data)
	return &data, err
}

// InsertResults saves all reports
func (reports *LabsReports) InsertResults() (err error) {
	rep := *reports
	for _, report := range rep {
		report.Insert()
	}
	return
}

// InsertOneResult imports only one site report
func (reports *LabsReports) InsertOneResult(site string) (err error) {
	rep := *reports
	for _, report := range rep {
		if report.Host != site {
			continue
		}
		report.Insert()
	}
	return
}
