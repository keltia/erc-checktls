package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/gobuffalo/packr"
	"github.com/keltia/ssllabs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTLSReport_ToHTML(t *testing.T) {
	var buf strings.Builder

	ji, err := ioutil.ReadFile("testdata/site.json")
	require.NoError(t, err)

	// Simulate
	fIgnoreMozilla = true
	fIgnoreImirhil = true

	all, err := ssllabs.ParseResults(ji)
	require.NoError(t, err)

	sites, err := NewTLSReport(all)
	require.NoError(t, err)

	raw, err := ioutil.ReadFile("files/templ.html")
	tmpl := string(raw)
	require.NoError(t, err)
	require.NotEmpty(t, tmpl)

	err = sites.ToHTML(&buf, tmpl)
	assert.NoError(t, err)
}

func TestWriteHTML2(t *testing.T) {
	cntrs := map[string]int{
		"A": 666,
		"B": 42,
		"F": 1,
	}

	https := map[string]int{
		"A":  666,
		"B+": 37,
		"F":  42,
	}

	r := &TLSReport{}
	err := r.WriteHTML(os.Stderr, cntrs, https)
	assert.Error(t, err)
}

func TestWriteHTML3(t *testing.T) {
	cntrs := map[string]int{
		"A": 666,
		"B": 42,
		"F": 1,
	}

	https := map[string]int{
		"A":  666,
		"B+": 37,
		"F":  42,
	}

	file := "testdata/site.json"
	raw, err := getResults(file)
	require.NoError(t, err)

	allSites, err := ssllabs.ParseResults(raw)
	require.NoError(t, err)

	box := packr.NewBox("./files")
	tmpls, err = loadTemplates(box)
	require.NoError(t, err)

	fIgnoreImirhil = true
	fIgnoreMozilla = true

	final, err := NewTLSReport(allSites)
	null, err := os.OpenFile("/dev/null", os.O_WRONLY, 0666)
	require.NoError(t, err)

	err = final.WriteHTML(null, cntrs, https)
	assert.NoError(t, err)
}
