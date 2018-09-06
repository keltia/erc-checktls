// report.go

/*
This file contains func for generating the report
*/
package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/atotto/encoding/csv"
	"github.com/keltia/cryptcheck"
	"github.com/keltia/observatory"
	"github.com/keltia/ssllabs"
	"github.com/pkg/errors"
)

var (
	yesno = map[bool]string{
		true:  "YES",
		false: "NO",
	}

	fnImirhil func(site ssllabs.Host) string
	fnMozilla func(site ssllabs.Host) string
)

const (
	DefaultKeySize = 2048
	DefaultAlg     = "RSA"
	DefaultIssuer  = "GlobalSign Organization Validation CA - SHA256 - G2"
	DefaultSig     = "SHA256withRSA"
)

// Private functions

func init() {
	if !fIgnoreImirhil {
		cnf := cryptcheck.Config{
			Log:     logLevel,
			Refresh: fRefresh,
		}
		client := cryptcheck.NewClient(cnf)

		fnImirhil = func(site ssllabs.Host) string {
			score, err := client.GetScore(site.Host)
			if err != nil {
				verbose("can not get cryptcheck score: %v\n", err)
			}
			return score
		}
	} else {
		fnImirhil = func(site ssllabs.Host) string {
			return ""
		}
	}

	if !fIgnoreMozilla {
		cnf := observatory.Config{
			Log: logLevel,
		}
		moz, err := observatory.NewClient(cnf)
		if err != nil {
			fmt.Fprintf(os.Stderr, "can not create observatory client: %v", err)
		}

		fnMozilla = func(site ssllabs.Host) string {
			score, err := moz.GetGrade(site.Host)
			if err != nil {
				verbose("can not get Mozilla score: %v\n", err)
			}
			return score
		}
	} else {
		fnMozilla = func(site ssllabs.Host) string {
			return ""
		}
	}
}

func fixTimestamp(ts int64) (int64, int64) {
	return ts / 1000, ts % 1000
}

func checkSweet32(det ssllabs.EndpointDetails) (yes bool) {
	if len(det.Suites) != 0 {
		ciphers := det.Suites[0].List
		for _, cipher := range ciphers {
			if strings.Contains(cipher.Name, "DES") {
				return true
			}
		}
	}
	return false
}

func getGrade(site ssllabs.Host, fn func(site ssllabs.Host) string) string {
	return fn(site)
}

func getSSLablsVersion(site ssllabs.Host) string {
	debug("%#v", site)
	return fmt.Sprintf("%s/%s", site.EngineVersion, site.CriteriaVersion)
}

// NewTLSReport generates everything we need for display/export
func NewTLSReport(reports []ssllabs.Host) (e *TLSReport, err error) {
	if len(reports) == 0 {
		return nil, fmt.Errorf("empty list")
	}

	e = &TLSReport{
		Date:    time.Now(),
		SSLLabs: getSSLablsVersion(reports[0]),
	}

	verbose("%d sites found.\n", len(reports))

	// Now analyze each site
	for _, site := range reports {
		var current TLSSite

		if site.Endpoints == nil {
			verbose("Site %s has no endpoint\n", site.Host)
			current = TLSSite{
				Name:     site.Host,
				Contract: contracts[site.Host],
			}
		} else {
			endp := site.Endpoints[0]
			det := endp.Details

			fmt.Printf("  Host: %s\n", site.Host)

			protos := []string{}
			for _, p := range det.Protocols {
				protos = append(protos, fmt.Sprintf("%sv%s", p.Name, p.Version))
			}

			// FIll in all details
			current = TLSSite{
				Name:       site.Host,
				Contract:   contracts[site.Host],
				Grade:      fmt.Sprintf("%s/%s", endp.Grade, endp.GradeTrustIgnored),
				CryptCheck: getGrade(site, fnImirhil),
				Mozilla:    getGrade(site, fnMozilla),
				Protocols:  strings.Join(protos, ","),
				RC4:        det.SupportsRC4,
				PFS:        det.ForwardSecrecy >= 2,
				OCSP:       det.OcspStapling,
				HSTS:       det.HstsPolicy.Status == "present",
				ALPN:       det.SupportsAlpn,
				Drown:      det.DrownVulnerable,
				Sweet32:    checkSweet32(det),
			}

			// Handle case where we have a DNS entry but no connection
			if len(site.Certs) != 0 {
				cert := site.Certs[0]
				current.DefKey = cert.KeySize == DefaultKeySize && cert.KeyAlg == DefaultAlg

				current.DefCA = cert.IssuerLabel == DefaultIssuer
				current.DefSig = cert.SigAlg == DefaultSig
				current.IsExpired = time.Now().After(time.Unix(fixTimestamp(cert.NotAfter)))
			}

			/*
				// make space
				var siteData []string

				// [7] = path
				siteData = append(siteData, fmt.Sprintf("%d", len(det.Chain.Certs)))

				// [8] = issues
				siteData = append(siteData, fmt.Sprintf("%d", det.Chain.Issues))
			*/
		}
		e.Sites = append(e.Sites, current)
	}
	return
}

func displayWildcards(all []ssllabs.Host) string {
	var buf strings.Builder

	fmt.Fprint(&buf, "")
	// Now analyze each site
	for _, site := range all {
		debug("site=%s\n", site.Host)
		if site.Endpoints != nil {
			// If Certs is empty, we could not connect
			if len(site.Certs) != 0 {
				cert := site.Certs[0]
				debug("  cert=%#v\n", cert)
				debug("  CN=%s\n", cert.Subject)

				if strings.HasPrefix(cert.Subject, "CN=*") {
					debug("adding %s\n", site.Host)
					buf.WriteString(fmt.Sprintf("  %-35s %-16s CN=%v \n", site.Host, site.Endpoints[0].IPAddress, cert.CommonNames))
				}
			}
		}
	}
	return buf.String()
}

// ToCSV output a CSV file from a report
func (r *TLSReport) ToCSV(w io.Writer) (err error) {
	wh := csv.NewWriter(w)
	debug("%v\n", r.Sites)
	if err = wh.WriteStructHeader(r.Sites[0]); err != nil {
		return errors.Wrap(err, "can not write csv header")
	}

	err = wh.WriteStructAll(r.Sites)
	return errors.Wrap(err, "can not write csv file")
}
