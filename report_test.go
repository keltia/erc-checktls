package main

import (
	"io/ioutil"
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
