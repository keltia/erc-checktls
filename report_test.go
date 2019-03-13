package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/keltia/ssllabs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTLSReport(t *testing.T) {
	rep, err := NewTLSReport([]ssllabs.Host{})
	require.Error(t, err)
	require.Nil(t, rep)
}

func TestNewTLSReport2(t *testing.T) {
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
}

func TestTLSReport_ToCSV(t *testing.T) {
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

	err = sites.ToCSV(&buf)
	assert.NoError(t, err)
}

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

func TestTLSReport_WriteCSV(t *testing.T) {
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
	err := r.WriteCSV(os.Stderr, cntrs, https)
	assert.Error(t, err)

}

func TestWriteCSV2(t *testing.T) {
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
	err := r.WriteCSV(os.Stderr, cntrs, https)
	assert.Error(t, err)
}

func TestWriteCSV3(t *testing.T) {
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

	fIgnoreImirhil = true
	fIgnoreMozilla = true

	final, err := NewTLSReport(allSites)
	null, err := os.OpenFile("/dev/null", os.O_WRONLY, 0666)
	require.NoError(t, err)

	err = final.WriteCSV(null, cntrs, https)
	assert.NoError(t, err)
}
