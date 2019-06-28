package site

import (
	"regexp"
	"strings"
	"time"

	"github.com/keltia/ssllabs"
)

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

func hasExpired(t int64) bool {
	return time.Now().After(time.Unix(fixTimestamp(t)))
}

func checkKey(cert ssllabs.Cert) bool {
	return (cert.KeySize == DefaultKeySize && cert.KeyAlg == DefaultAlg || cert.KeyAlg == DefaultECAlg)
}

func checkIssuer(cert ssllabs.Cert, ours *regexp.Regexp) string {
	if ours.MatchString(cert.IssuerSubject) {
		return "TRUE"
	}
	if (cert.Issues ^ 0x40) == 0 {
		return "SELF"
	}
	return "FALSE"
}

func checkHSTS(det ssllabs.EndpointDetails) int64 {
	if det.HstsPolicy.Status == "present" {
		return det.HstsPolicy.MaxAge
	}
	return -1
}
