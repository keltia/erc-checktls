// categories.go

package main

import (
	"fmt"
	"github.com/keltia/ssllabs"
	tw "github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"io"
)

var (
	keys = []string{
		"A+",
		"A",
		"A-",
		"B",
		"C",
		"D",
		"E",
		"F",
		"T",
		"X",
		"Z",
		"Total",
		"RC4",
		"OCSP",
		"HSTS",
		"PFS",
		"Sweet32",
	}
)

func categoryCounts(reports []ssllabs.Host) (cntrs map[string]int) {
	cntrs = make(map[string]int)

	baddies := 0
	broken := 0
	reals := 0

	for _, r := range reports {
		if r.Endpoints != nil {
			endp := r.Endpoints[0]
			det := endp.Details

			if r.Endpoints[0].Grade != "" && r.Endpoints[0].Grade != "Z" {
				reals++
			} else {
				baddies++
			}
			cntrs[r.Endpoints[0].Grade]++
			if det.ForwardSecrecy >= 2 {
				cntrs["PFS"]++
			}
			if checkSweet32(det) {
				cntrs["Sweet32"]++
			}
			if det.SupportsRC4 {
				cntrs["RC4"]++
			}
			if det.OcspStapling {
				cntrs["OCSP"]++
			}
			if det.HstsPolicy.Status == "present" {
				cntrs["HSTS"]++
			}
		} else {
			broken++
		}
	}
	cntrs["Total"] = reals
	cntrs["X"] = broken
	cntrs["Z"] = baddies
	return cntrs
}

func displayCategories(cntrs map[string]int) string {
	str := ""
	for _, k := range keys {
		str = str + fmt.Sprintf("%s:%d ", k, cntrs[k])
	}
	return str
}

func writeSummary(cntrs map[string]int, w io.Writer) (err error) {
	table := tw.NewWriter(w)
	table.SetHeader(keys)
	table.SetAlignment(tw.ALIGN_CENTER)

	line := []string{}
	for _, c := range keys {
		if v, ok := cntrs[c]; ok {
			line = append(line, fmt.Sprintf("%d", v))
		} else {
			line = append(line, "0")
		}
	}

	table.Append(line)
	table.Render()

	return errors.Wrap(err, "table")
}
