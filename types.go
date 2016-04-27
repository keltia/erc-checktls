// types.go

package main

import "time"

type EECSite struct {
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

type EECReport struct {
	Date time.Time
	Sites map[string]EECSite
}

// EECLine is used to hold a CSV-tobe line
type EECLine []interface{}

