package main

import (
	"fmt"
	"strings"
	"text/template"
)

func grade(val string) string {
	switch val {
	case "A+":
		fallthrough
	case "A":
		return green(val)
	case "A-":
		return yellow(val)
	case "B+":
		fallthrough
	case "B":
		fallthrough
	case "B-":
		return orange(val)
	case "C+":
		fallthrough
	case "C":
		fallthrough
	case "C-":
		fallthrough
	case "D+":
		fallthrough
	case "D":
		fallthrough
	case "D-":
		fallthrough
	case "E+":
		fallthrough
	case "E":
		fallthrough
	case "E-":
		fallthrough
	case "F":
		return red(val)
	}
	return white("&nbsp;")
}

func booleanT(val bool) string {
	if val {
		return white("TRUE")
	}
	return red("FALSE")
}

func booleanF(val bool) string {
	if val {
		return red("TRUE")
	}
	return white("FALSE")
}

func proto(val string) string {
	switch val {
	case "TLSv1.2":
		return green(val)
	case "TLSv1.1,TLSv1.2":
		fallthrough
	case "TLSv1.0,TLSv1.1,TLSv1.2":
		return yellow(val)
	case "SSLv3.0,TLSv1.0":
		return red(val)
	}
	return white(val)
}

const (
	hstsNone = -1 // red

	hstsLow  = 86400 * 90  // orange
	hstsGood = 86400 * 180 // yellow
	// green
)

func hsts(age int64) string {
	if age == hstsNone {
		return red("NO")
	}

	if age >= 0 && age < hstsLow {
		return orange(fmt.Sprintf("%d", age))
	}

	if age >= hstsLow && age < hstsGood {
		return yellow(fmt.Sprintf("%d", age))
	}
	return green(fmt.Sprintf("%d", age))
}

func servertype(t int) string {
	if t == TypeHTTPSok {
		return green("HTTPS")
	} else if t == TypeHTTPSnok {
		return orange("MIXED")
	} else {
		return red("HTTP")
	}
}

func red(str string) string {
	var buf strings.Builder

	t, _ := template.New("red").Parse(`<td class=xl64 align=center>{{.}}</td>`)
	t.Execute(&buf, str)
	return buf.String()
}

func yellow(str string) string {
	var buf strings.Builder

	t, _ := template.New("yellow").Parse(`<td class=xl631 align=center>{{.}}</td>`)
	t.Execute(&buf, str)
	return buf.String()

}

func orange(str string) string {
	var buf strings.Builder

	t, _ := template.New("orange").Parse(`<td class=xl63 align=center>{{.}}</td>`)
	t.Execute(&buf, str)
	return buf.String()

}

func green(str string) string {
	var buf strings.Builder

	t, _ := template.New("green").Parse(`<td class=xl65 align=center>{{.}}</td>`)
	t.Execute(&buf, str)
	return buf.String()
}

func white(str string) string {
	var buf strings.Builder

	t, _ := template.New("white").Parse(`<td class=xl661 align=center>{{.}}</td>`)
	t.Execute(&buf, str)
	return buf.String()
}

func text(str string) string {
	var buf strings.Builder

	t, _ := template.New("white").Parse(`<td height=21 style='height:16.0pt'>{{.}}</td>`)
	t.Execute(&buf, str)
	return buf.String()
}
