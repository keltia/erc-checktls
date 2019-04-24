package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/keltia/ssllabs"
	"github.com/pkg/errors"

	"github.com/keltia/erc-checktls/site"
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

func (r *TLSReport) GatherStats(s site.TLSSite) {
	if !s.Empty {
		if s.Grade != "" && s.Grade != "Z" {
			r.cntrs["Total"]++
		} else {
			r.cntrs["Z"]++
		}
		r.cntrs[s.Grade]++
		if s.PFS {
			r.cntrs["PFS"]++
		}
		if s.Sweet32 {
			r.cntrs["Sweet32"]++
		}
		if s.Issues {
			r.cntrs["Issues"]++
		}
		if s.OCSPStapling {
			r.cntrs["OCSPStapling"]++
		}
		if s.HSTS > 0 {
			r.cntrs["HSTS"]++
		}

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
	} else {
		r.cntrs["X"]++
	}
}
