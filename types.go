// types.go

package main

import "time"

// TLSReport is one single run for all sites
type TLSReport struct {
	Date  time.Time
	Sites [][]string
}
