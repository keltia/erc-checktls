package main

import (
	"io/ioutil"
	"testing"

	"github.com/keltia/ssllabs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetResults(t *testing.T) {
	ji, err := ioutil.ReadFile("testdata/site.json")
	require.NoError(t, err)

	buf, err := getResults("testdata/site.json")
	require.NoError(t, err)

	assert.Equal(t, ji, buf)
}

func TestGetResultsNothing(t *testing.T) {
	buf, err := getResults("testdata/site.nowhere")
	require.Error(t, err)
	require.Empty(t, buf)
}

func TestTLSReport_HTTPCountsEmpty(t *testing.T) {
	r := &TLSReport{
		https: map[string]int{},
		Sites: []TLSSite{},
	}

	r.httpCounts()
	assert.Empty(t, r.cntrs)
	assert.Empty(t, r.https)
}

func TestTLSReport_HTTPCountsReport(t *testing.T) {
	ji, err := ioutil.ReadFile("testdata/site.json")
	require.NoError(t, err)

	// Simulate
	fIgnoreMozilla = true
	fIgnoreImirhil = true

	all, err := ssllabs.ParseResults(ji)
	require.NoError(t, err)

	r, err := NewTLSReport(all)
	require.NoError(t, err)
	require.NotEmpty(t, r)

	// Fake it
	r.Sites[0].Mozilla = "A+"

	t.Logf("r=%#v", r)
	r.httpCounts()
	t.Logf("r=%#v", r)
	assert.NotEmpty(t, r.cntrs)
	assert.EqualValues(t, map[string]int{"": 1, "A+": 1, "HSTS": 1, "Issues": 1, "OCSP": 1, "PFS": 1, "Total": 1, "Z": 1}, r.cntrs)
}

func TestTLSReport_HTTPCountsReport_1(t *testing.T) {
	ji, err := ioutil.ReadFile("testdata/site.json")
	require.NoError(t, err)

	// Simulate
	fIgnoreMozilla = true
	fIgnoreImirhil = true

	all, err := ssllabs.ParseResults(ji)
	require.NoError(t, err)

	r, err := NewTLSReport(all)
	require.NoError(t, err)
	require.NotEmpty(t, r)

	// Fake it
	r.Sites[0].Mozilla = "H"

	r.httpCounts()
	assert.NotEmpty(t, r.cntrs)
	assert.EqualValues(t, map[string]int{"": 1, "A+": 1, "HSTS": 1, "Issues": 1, "OCSP": 1, "PFS": 1, "Total": 1, "Z": 1}, r.cntrs)
}

func TestTLSReport_HTTPCountsReport_2(t *testing.T) {
	ji, err := ioutil.ReadFile("testdata/reallybad.json")
	require.NoError(t, err)

	// Simulate
	fIgnoreMozilla = true
	fIgnoreImirhil = true

	all, err := ssllabs.ParseResults(ji)
	require.NoError(t, err)

	r, err := NewTLSReport(all)
	require.NoError(t, err)
	require.NotEmpty(t, r)

	// Fake it
	r.Sites[0].Mozilla = "H"

	r.httpCounts()
	assert.NotEmpty(t, r.cntrs)
	assert.EqualValues(t, map[string]int{"A+": 1, "HSTS": 1, "Issues": 1, "OCSP": 1, "PFS": 1, "Sweet32": 1, "Total": 1}, r.cntrs)
}
