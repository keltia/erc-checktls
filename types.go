// types.go

package TLS

import (
	"time"

	"github.com/keltia/erc-checktls/site"
)

// Config is for setting internal flags
type Config struct {
	Jobs          int
	LogLevel      int
	IgnoreMozilla bool
	IgnoreImirhil bool
}

// Report is one single run for all sites
type Report struct {
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
