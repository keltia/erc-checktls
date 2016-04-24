// eecreport.go

/*
This file contains our EEC-specific data/func
 */
package main

import (
	"os"
	"encoding/csv"
)

// EECReport is the data we want to extract
type EECReport struct {
	Site          string
	Contract      string
	Grade         string
	CertScore     int
	ProtocolScore int
	KeyExchScore  int
	StrengthScore int
	Key           string
	Sign          string
	Issuer        string
	Validity      int
	Path          int
	ChainIssues   string
	Protocols     []string
	RC4           bool
	Pfs           bool
	OcspStapling  bool
	Hsts          bool
	Alpn          bool
	ImirhilScore  string
}

// EECLine is used to hold a CSV-tobe line
type EECLine []interface{}

// NewEECReport is
func NewEECReport(r LabsReport) (e *EECReport, err error) {
	contract := contracts[r.Host]
	e = &EECReport{
		Site: r.Host,
		Contract: contract,
		Grade: r.Endpoints[0].Grade,
	}

	return
}

// toLine groups part of the data into a single array
func (r *EECReport) toLine() {

}

// ToCSV generate a CSV file from a given report
func (r *EECReport) ToCSV() {

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

	for _, site := range allSites {
		contracts[site[0]] = site[1]
	}
	return
}

