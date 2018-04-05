// types.go

package main

import "time"

// TLSReport is one single run for all sites
type TLSReport struct {
	EngineVersion   string
	CriteriaVersion string
	Date            time.Time
	Sites           [][]string
}
