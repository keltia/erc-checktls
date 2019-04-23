package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/keltia/ssllabs"
	"github.com/pkg/errors"
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

func getSSLablsVersion(site ssllabs.Host) string {
	debug("%#v", site)
	return fmt.Sprintf("%s/%s", site.EngineVersion, site.CriteriaVersion)
}

func (r *TLSReport) categoryCounts(s ssllabs.Host) {
	if s.Endpoints != nil && len(s.Endpoints) != 0 {
		endp := s.Endpoints[0]
		det := endp.Details

		if s.Endpoints[0].Grade != "" && s.Endpoints[0].Grade != "Z" {
			r.cntrs["Total"]++
		} else {
			r.cntrs["Z"]++
		}
		r.cntrs[s.Endpoints[0].Grade]++
		if det.ForwardSecrecy >= 2 {
			r.cntrs["PFS"]++
		}
		if checkSweet32(det) {
			r.cntrs["Sweet32"]++
		}
		if len(det.CertChains) == 0 ||
			det.CertChains[0].Issues != 0 {
			r.cntrs["Issues"]++
		}
		if det.OcspStapling {
			r.cntrs["OCSP"]++
		}
		if det.HstsPolicy.Status == "present" {
			r.cntrs["HSTS"]++
		}
	} else {
		r.cntrs["X"]++
	}
}

func (r *TLSReport) httpCounts() {
	for _, s := range r.Sites {
		if s.Mozilla != "" {
			if s.Mozilla >= "G" {
				r.https["Bad"]++
			} else {
				r.https["Total"]++
			}
			r.https[s.Mozilla]++
		} else {
			r.https["Broken"]++
		}
	}
	return
}
