// categories.go

package main

import (
	"encoding/csv"
	"fmt"
	"github.com/keltia/erc-checktls/ssllabs"
	"io"
)

var (
	cntrs map[string]int

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
	}
)

func categoryCounts(reports *ssllabs.LabsReports) {
	cntrs = make(map[string]int)
	for _, r := range *reports {
		if r.Endpoints != nil {
			cntrs[r.Endpoints[0].Grade]++
		}
	}
}

func displayCategories(cntrs map[string]int) string {
	str := ""
	for _, k := range keys {
		str = str + fmt.Sprintf("%s:%d ", k, cntrs[k])
	}
	return str
}

func categoriesCSV(cntrs map[string]int, w io.Writer) (err error) {
	res := make([][]string, 2)

	hdrs := []string{}
	line := []string{}
	for _, c := range keys {
		hdrs = append(hdrs, c)
		if v, ok := cntrs[c]; ok {
			line = append(line, fmt.Sprintf("%d", v))
		} else {
			line = append(line, "0")
		}
	}
	res[0] = hdrs
	res[1] = line

	wh := csv.NewWriter(w)
	err = wh.WriteAll(res)
	return
}
