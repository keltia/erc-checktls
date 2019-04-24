// categories.go

package main

import (
	"fmt"
)

var (
	tlsKeys = []string{
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
		"Issues",
		"PFS",
		"OCSPStapling",
		"HSTS",
		"Sweet32",
	}
	httpKeys = []string{
		"A+",
		"A",
		"A-",
		"B-",
		"B",
		"B-",
		"C+",
		"C",
		"C-",
		"D+",
		"D",
		"D-",
		"E+",
		"E",
		"E-",
		"F+",
		"F",
		"F-",
		"T",
		"X",
		"Z",
		"Total",
		"Broken",
	}
)

func displayCategories(cntrs map[string]int) string {
	str := ""
	for _, k := range tlsKeys {
		str = str + fmt.Sprintf("%s:%d ", k, cntrs[k])
	}
	return str
}

func selectColours(grade string) string {
	switch grade {
	case "A+":
		fallthrough
	case "A":
		return "green"
	case "A-":
		fallthrough
	case "B+":
		fallthrough
	case "B":
		fallthrough
	case "B-":
		return "orange"
	}
	return "red"
}
