package main

import (
	"fmt"
	"io"
	"text/template"

	"github.com/pkg/errors"
)

type htmlvars struct {
	Name       string
	Contract   string
	Grade      string
	Cryptcheck string
	Mozilla    string
	Redir      string
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

			Redir:     servertype(s.Type),
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

func WriteHTML(w io.Writer, final *TLSReport, cntrs, https map[string]int) error {
	var err error

	debug("WriteHTML")
	if final == nil {
		return fmt.Errorf("nil final")
	}
	if len(final.Sites) == 0 {
		return fmt.Errorf("empty final")
	}

	debug("tmpls=%v\n", tmpls)
	if err = final.ToHTML(w, tmpls["templ.html"]); err != nil {
		return errors.Wrap(err, "Can not write HTML")
	}
	// Generate colour map
	//cm := final.ColourMap()
	if fSummary != "" {
		fn := fSummary + "-" + makeDate() + ".html"
		verbose("HTML summary: %s\n", fn)
		w = checkOutput(fn)
		err = writeHTMLSummary(w, cntrs, https)
	}
	return err
}

