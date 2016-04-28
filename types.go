// types.go

package main

import "time"

// TLSSite is for one site checked
type TLSSite struct {
	Site          string
	Contract      string
	Grade         string
	//CertScore     string		DEPRECATED
	//ProtocolScore string
	//KeyExchScore  string
	//StrengthScore string
	Key           string
	Sign          string
	Issuer        string
	Validity      string
	Path          string
	ChainIssues   string
	Protocols     string
	RC4           string
	Pfs           string
	OcspStapling  string
	Hsts          string
	Alpn          string
	Drown         string
	ImirhilScore  string
}

// TLSReport is one single run for all sites
type TLSReport struct {
	Date time.Time
	Sites []EECLine
}

// EECLine is used to hold a CSV-tobe line
type EECLine []string

