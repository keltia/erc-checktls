package main

import (
	"io/ioutil"
	"testing"

	"github.com/keltia/ssllabs"
	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestDisplayCategories(t *testing.T) {
	cntrs := map[string]int{
		"A": 666,
		"B": 0,
		"G": 1,
	}
	str := displayCategories(cntrs)
	assert.NotEmpty(t, str)
}

func TestHTTPCountsNil(t *testing.T) {
	cntrs := httpCounts(nil)
	assert.Empty(t, cntrs)
}

func TestHTTPCountsEmpty(t *testing.T) {
	cntrs := httpCounts(&TLSReport{})
	assert.NotEmpty(t, cntrs)
	assert.EqualValues(t, map[string]int{"Total": 0, "Broken": 0}, cntrs)
}

func TestHTTPCountsReport(t *testing.T) {
	ji, err := ioutil.ReadFile("testdata/site.json")
	require.NoError(t, err)

	// Simulate
	fIgnoreMozilla = true
	fIgnoreImirhil = true

	all, err := ssllabs.ParseResults(ji)
	require.NoError(t, err)

	sites, err := NewTLSReport(all)
	require.NoError(t, err)
	require.NotEmpty(t, sites)

	// Fake it
	sites.Sites[0].Mozilla = "A+"

	cntrs := httpCounts(sites)
	assert.NotEmpty(t, cntrs)
	assert.EqualValues(t, map[string]int{"A+": 1, "Total": 1, "Broken": 1}, cntrs)
}

func TestHTTPCountsReport_1(t *testing.T) {
	ji, err := ioutil.ReadFile("testdata/site.json")
	require.NoError(t, err)

	// Simulate
	fIgnoreMozilla = true
	fIgnoreImirhil = true

	all, err := ssllabs.ParseResults(ji)
	require.NoError(t, err)

	sites, err := NewTLSReport(all)
	require.NoError(t, err)
	require.NotEmpty(t, sites)

	// Fake it
	sites.Sites[0].Mozilla = "H"

	cntrs := httpCounts(sites)
	assert.NotEmpty(t, cntrs)
	assert.EqualValues(t, map[string]int{"H": 1, "Total": 0, "Broken": 1}, cntrs)
}
