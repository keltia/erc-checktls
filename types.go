// types.go

package main

import (
	"time"

	"github.com/keltia/erc-checktls/site"
)

// TLSReport is one single run for all sites
type TLSReport struct {
	SSLLabs string
	Date    time.Time
	Sites   []site.TLSSite

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
