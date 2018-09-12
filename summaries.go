package main

import (
	"bytes"
	"fmt"
	"io"

	tw "github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
)

type cmap struct {
	elem string
	fn   func(string) string
}

var ctlsmap = []cmap{
	{"A+", green},
	{"A", green},
	{"A-", yellow},
	{"B", orange},
	{"C", orange},
	{"D", orange},
	{"E", orange},
	{"F", orange},
	{"T", orange},
	{"X", orange},
	{"Z", orange},
	{"Total", white},
	{"Issues", white},
	{"PFS", white},
	{"OCSP", white},
	{"HSTS", white},
	{"Sweet32", white},
}

var httpmap = []cmap{
	{"A+", green},
	{"A", green},
	{"A-", yellow},
	{"B-", orange},
	{"B", orange},
	{"B-", orange},
	{"C+", orange},
	{"C", orange},
	{"C-", orange},
	{"D+", red},
	{"D", red},
	{"D-", red},
	{"E+", red},
	{"E", red},
	{"E-", red},
	{"F+", red},
	{"F", red},
	{"F-", red},
	{"T", red},
	{"X", red},
	{"Z", red},
	{"Total", white},
	{"Broken", white},
}

func writeSummary(w io.Writer, keys []string, cntrs map[string]int) (err error) {
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

func writeHTMLSummary(w io.Writer, keys []cmap, cntrs map[string]int) (err error) {
	var str bytes.Buffer

	if len(keys) == 0 || len(cntrs) == 0 {
		return nil
	}
	for _, cm := range keys {
		str.WriteString(cm.fn(fmt.Sprintf("%d", cntrs[cm.elem])))
		str.WriteString("\n")
	}
	w.Write(str.Bytes())
	return nil
}
