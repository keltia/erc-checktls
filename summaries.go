package main

import (
	"bytes"
	"fmt"
	"io"
	"text/template"
	"time"

	tw "github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
)

const summariesT = "summaries.html"

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

func writeOneRow(keys []cmap, cntrs map[string]int) string {
	var str bytes.Buffer

	// Contruct the row with all possible keys
	for _, cm := range keys {
		str.WriteString(cm.fn(fmt.Sprintf("%d", cntrs[cm.elem])))
		str.WriteString("\n")
	}
	return str.String()
}

func writeHTMLSummary(w io.Writer, cntrs, https map[string]int) (err error) {
	tm, ok := tmpls[summariesT]
	if !ok {
		debug("%s: %s", summariesT, tm)
		return fmt.Errorf("tmpl is empty")
	}

	t := template.Must(template.New(summariesT).Parse(tm))
	if t == nil {
		return fmt.Errorf("bad template")
	}

	if len(cntrs) == 0 || len(https) == 0 {
		return
	}
	date := time.Now().Format("2006-01") + "-01"
	data := struct {
		Date string
		TLS  string
		HTTP string
	}{
		Date: date,
		TLS:  writeOneRow(ctlsmap, cntrs),
		HTTP: writeOneRow(httpmap, https),
	}
	err = t.ExecuteTemplate(w, summariesT, data)
	if err != nil {
		return errors.Wrap(err, "writeHTMLSummary")
	}

	return nil
}
