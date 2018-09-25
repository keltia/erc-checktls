package main

import (
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/pkg/errors"
)

type htmlvars struct {
	Name       string
	Contract   string
	Grade      string
	Cryptcheck string
	Mozilla    string
	DefKey     string
	DefSig     string
	DefCA      string
	IsExpired  string
	Issues     string
	Protocols  string
	PFS        string
	OCSP       string
	HSTS       string
	Sweet32    string
}

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
	hstsNone = -1				// red

	hstsLow  = 86400 * 90		// orange
	hstsGood = 86400 * 180		// yellow
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

func (r *TLSReport) ToHTML(w io.Writer, tmpl string) error {
	var (
		err error
	)

	debug("ToHTML\n")
	debug("tmpl=%s\n", tmpl)
	t := template.Must(template.New("html-report").Parse(tmpl))

	debug("t=%v\n", t)
	debug("Sites=%#v\n", r.Sites)
	Sites := []htmlvars{}
	for _, s := range r.Sites {
		h := htmlvars{
			Name:     text(s.Name),
			Contract: text(s.Contract),

			Grade:      grade(s.Grade),
			Cryptcheck: grade(s.CryptCheck),
			Mozilla:    grade(s.Mozilla),

			DefKey:    booleanT(s.DefKey),
			DefSig:    booleanT(s.DefSig),
			DefCA:     booleanT(s.DefCA),
			IsExpired: booleanF(s.IsExpired),
			Issues:    booleanF(s.PathIssues),

			Protocols: proto(s.Protocols),

			PFS:     booleanT(s.PFS),
			OCSP:    booleanT(s.OCSP),
			HSTS:    hsts(s.HSTS),
			Sweet32: booleanF(s.Sweet32),
		}
		Sites = append(Sites, h)
		debug("h=%#v\n", h)
	}
	htmlVars := struct {
		Date    string
		Version string
		Sites   []htmlvars
	}{makeDate(), r.SSLLabs, Sites}
	err = t.ExecuteTemplate(w, "html-report", htmlVars)
	return errors.Wrap(err, "can not write HTML file")
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
