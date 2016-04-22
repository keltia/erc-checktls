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
	KeyScore      int
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

// NewEECReport is
func NewEECReport(r LabsReport) (e *EECReport, err error) {
	contract, err := getContract(r.Host)
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
func getContract(name string) (contract string, err error) {
	return
}
