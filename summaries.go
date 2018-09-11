package main

import (
	"fmt"
	"io"

	tw "github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
)

func writeSummary(keys []string, cntrs map[string]int, w io.Writer) (err error) {
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
