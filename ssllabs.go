// ssllabs.go

/*
SSLLabs-related functions.
 */
package main

import "log"

// Display for one report
func (rep *LabsReport) Display() {
	host := rep.Host
	grade := rep.Endpoints[0].Grade
	details := rep.Endpoints[0].Details
	log.Printf("Looking at %s… — grade %s", host, grade)
	if fVerbose {
		log.Printf("  Ciphers: %d", details.Suites.len())
	} else if fReallyVerbose {
		for _, cipher := range details.Suites.List {
			log.Printf("  %s: %d bits", cipher.Name, cipher.CipherStrength)
		}
	}
}

