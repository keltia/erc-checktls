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

	IsExpired bool

	Protocols string

	RC4     bool
	PFS     bool
	OCSP    bool
	HSTS    bool
	ALPN    bool
	Drown   bool
	Sweet32 bool
}
