// types.go

package main

import (
	"time"
)

// TLSReport is one single run for all sites
type TLSReport struct {
	SSLLabs string
	Date    time.Time
	Sites   []TLSSite
}

// TLSSite is a summary for each site
type TLSSite struct {
	Name     string
	Contract string

	Grade      string
	CryptCheck string
	Mozilla    string

	DefKey bool
	DefSig bool
	DefCA  bool

	IsExpired  bool
	PathIssues bool

	Protocols string
	PFS       bool

	OCSP    bool
	HSTS    bool
	Sweet32 bool
}
