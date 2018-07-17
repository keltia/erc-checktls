// report.go

/*
This file contains func for generating the report
*/
package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/keltia/cryptcheck"
	"github.com/keltia/erc-checktls/obs"
	"github.com/keltia/erc-checktls/ssllabs"
	"github.com/pkg/errors"
)

var (
	headersLine = []string{
		"Site",
		"Contract",
		"Grade",
		"Key",
		"Signature",
		"Issuer",
		"Validity",
		"Path",
		"Issues",
		"Protocols",
		"RC4?",
		"PFS?",
		"OCSP?",
		"HSTS?",
		"ALPN?",
		"Drown?",
		"Ciphers",
		"cryptcheck",
		"Sweet32",
		"Robot",
		"Mozilla",
	}

	yesno = map[bool]string{
		true:  "YES",
		false: "NO",
	}
)

// Private functions

// getResults read the JSON array generated and gone through jq
func getResults(file string) (res []byte, err error) {
	fh, err := os.Open(file)
	if err != nil {
		return res, errors.Wrapf(err, "can not open %s", file)
	}
	defer fh.Close()

	res, err = ioutil.ReadAll(fh)
	return res, errors.Wrapf(err, "can not read json %s", file)
}

func fixTimestamp(ts int64) (int64, int64) {
	return ts / 1000, ts % 1000
}

func checkSweet32(det ssllabs.LabsEndpointDetails) (yes bool) {
	ciphers := det.Suites.List
	for _, cipher := range ciphers {
		if strings.Contains(cipher.Name, "DES") {
			return true
		}
	}
	return false
}

// Public functions

// NewTLSReport generates everything we need for display/export
func NewTLSReport(ctx *Context, reports *ssllabs.LabsReports) (e *TLSReport, err error) {
	var (
		client *cryptcheck.Client
		moz    *obs.Client
	)

	e = &TLSReport{
		Date:  time.Now(),
		Sites: make([][]string, len(*reports)+1),
	}

	if !fIgnoreImirhil {
		cnf := cryptcheck.Config{
			Log:     logLevel,
			Refresh: fRefresh,
		}
		client = cryptcheck.NewClient(cnf)
	}

	if !fIgnoreMozilla {
		cnf := obs.Config{
			Log:     logLevel,
			Refresh: fRefresh,
		}
		moz = obs.NewClient(cnf)
	}

	verbose("%d sites found.", len(*reports))
	// First add the headers line
	e.Sites[0] = headersLine

	// Now analyze each site
	for i, site := range *reports {
		// Hack to get SSLLabs report engiveVersion & CriteriaVersion
		if e.CriteriaVersion == "" || e.EngineVersion == "" {
			e.CriteriaVersion = site.CriteriaVersion
			e.EngineVersion = site.EngineVersion
		}

		if site.Endpoints == nil {
			log.Printf("Site %s has no endpoint", site.Host)
			continue
		}
		endp := site.Endpoints[0]
		det := endp.Details
		cert := endp.Details.Cert

		verbose("  Host: %s", site.Host)
		// make space
		var siteData []string

		// [0] = site
		siteData = append(siteData, site.Host)

		// [1] = contract
		siteData = append(siteData, contracts[site.Host])

		// [2] = grade
		siteData = append(siteData, fmt.Sprintf("%s/%s", endp.Grade, endp.GradeTrustIgnored))

		// [3] = key
		siteData = append(siteData, fmt.Sprintf("%s %d bits",
			det.Key.Alg,
			det.Key.Size))

		// [4] = signature
		siteData = append(siteData, det.Cert.SigAlg)

		// [5] = issuer
		siteData = append(siteData, det.Cert.IssuerLabel)

		// [6] = validity
		siteData = append(siteData, time.Unix(fixTimestamp(cert.NotAfter)).String())

		// [7] = path
		siteData = append(siteData, fmt.Sprintf("%d", len(det.Chain.Certs)))

		// [8] = issues
		siteData = append(siteData, fmt.Sprintf("%d", det.Chain.Issues))

		// [9] = protocols
		protos := []string{}
		for _, p := range det.Protocols {
			protos = append(protos, fmt.Sprintf("%sv%s", p.Name, p.Version))
		}
		siteData = append(siteData, strings.Join(protos, ","))

		// [10] = RC4
		siteData = append(siteData, yesno[det.SupportsRC4])

		// [11] = PFS
		// 0 = NO
		// 1 = with some browsers but not the reference ones
		// 2 = with modern browsers
		// 4 = with most browsers (ROBUST)
		siteData = append(siteData, yesno[det.ForwardSecrecy >= 2])

		// [12] = OCSP Stapling
		siteData = append(siteData, yesno[det.OcspStapling])

		// [13] = HSTS
		siteData = append(siteData, yesno[det.HstsPolicy.Status == "present"])

		// [14] = ALPN
		siteData = append(siteData, yesno[det.SupportsAlpn])

		// [15] = Drown vuln
		siteData = append(siteData, yesno[det.DrownVulnerable])

		// [16] = # of ciphers
		siteData = append(siteData, fmt.Sprintf("%d", len(det.Suites.List)))

		// [17] = cryptcheck score unless ignored
		if !fIgnoreImirhil {
			score, err := client.GetScore(site.Host)
			if err != nil {
				verbose("can not get cryptcheck score: %v", err)
			}
			siteData = append(siteData, score)
		} else {
			siteData = append(siteData, "")
		}

		// [18] = include 64-bytes ciphers?
		siteData = append(siteData, yesno[checkSweet32(det)])

		// [19] = Robot Attack, return of the Oracle?
		siteData = append(siteData, "NO")

		// [20] = Mozilla Observatory score unless ignored
		if !fIgnoreMozilla {
			grade, err := moz.GetGrade(site.Host)
			if err != nil {
				verbose("can not get Mozilla grade: %v", err)
			}
			siteData = append(siteData, grade)
		} else {
			siteData = append(siteData, "")
		}

		e.Sites[i+1] = siteData
	}
	return
}

// ToCSV output a CSV file from a report
func (r *TLSReport) ToCSV(w io.Writer) (err error) {
	wh := csv.NewWriter(w)
	verbose("%v\n", r.Sites)
	err = wh.WriteAll(r.Sites)
	return errors.Wrap(err, "can not write csv file")
}

/* Display for one report
func (rep *ssllabs.LabsReport) String() {
	host := rep.Host
	grade := rep.Endpoints[0].Grade
	details := rep.Endpoints[0].Details
	log.Printf("Looking at %s/%s — grade %s", host, contracts[host], grade)
	if fVerbose {
		log.Printf("  Ciphers: %d", details.Suites.len())
	} else if fReallyVerbose {
		for _, cipher := range details.Suites.List {
			log.Printf("  %s: %d bits", cipher.Name, cipher.CipherStrength)
		}
	}
}*/
