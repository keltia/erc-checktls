package main

import (
	"bytes"
	"path/filepath"

	"github.com/atotto/encoding/csv"
	"github.com/gobuffalo/packr"
	"github.com/pkg/errors"
)

const (
	contractFile  = "sites-list.csv"
	resourcesPath = "./files"
)

type Templ map[string]string

// Load all .html files into an array
func loadTemplates(box packr.Box) (Templ, error) {
	list := map[string]string{}

	err := box.Walk(func(s string, file packr.File) error {
		ext := filepath.Ext(s)
		if ext == ".html" {
			t := box.String(s)
			list[filepath.Base(s)] = t
		}
		return nil
	})
	if err != nil {
		debug("got an error")
		return Templ{}, errors.Wrap(err, "loadTemplates")
	}

	return list, nil
}

// getContract retrieve the site's contract from the DB
func readContractFile(box packr.Box) (contracts map[string]string, err error) {
	debug("reading contracts\n")
	cf := box.Bytes(contractFile)
	fh := bytes.NewBuffer(cf)

	all := csv.NewReader(fh)
	allSites, err := all.ReadAll()
	if err != nil {
		return nil, errors.Wrap(err, "ReadAll")
	}

	contracts = make(map[string]string)
	for _, site := range allSites {
		contracts[site[0]] = site[1]
	}
	err = nil
	return
}

// Load all resources
func loadResources(path string) error {
	var err error

	// We embed the file now
	box := packr.NewBox(path)

	// We need that for the reports
	contracts, err = readContractFile(box)
	if err != nil {
		return errors.Wrapf(err, "readContractFile/%s", contractFile)
	}

	tmpls, err = loadTemplates(box)
	return errors.Wrapf(err, "loadTemplates/%s", path)
}
