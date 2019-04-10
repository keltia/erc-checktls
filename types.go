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

	// Statistics
	cntrs map[string]int
	https map[string]int
}

const (
	// 1, 2, 3 are the main categories 1=green, 2=yellow, 3=red
	CatGreen = 1 + iota
	CatYellow
	CatRed

	// 1 is for correct https w/ redirection, 2 is https&http, 3 is http only
	TypeError = 1 + iota
	TypeHTTPSok
	TypeHTTPSnok
	TypeHTTP
)

// TLSSite is a summary for each site
type TLSSite struct {
	Name     string
	Contract string

	Grade      string
	CryptCheck string
	Mozilla    string

	DefKey bool
	DefSig bool
	DefCA  string

	IsExpired  bool
	PathIssues bool

	Protocols string
	PFS       bool

	OCSP    bool
	HSTS    int64
	Sweet32 bool

	Type    int
	CatHTTP int
	CatTLS  int
}
