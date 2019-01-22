package main

import (
	"io/ioutil"
	"testing"

	"github.com/keltia/ssllabs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTLSSite(t *testing.T) {
	ji, err := ioutil.ReadFile("testdata/site.json")
	require.NoError(t, err)

	// Simulate
	fIgnoreMozilla = true
	fIgnoreImirhil = true

	all, err := ssllabs.ParseResults(ji)
	require.NoError(t, err)

	tls := NewTLSSite(all[0])
	require.NotEmpty(t, tls)
}

func TestNewTLSSite1(t *testing.T) {
	tls := NewTLSSite(ssllabs.Host{})
	require.Empty(t, tls)
}

func TestInitAPI(t *testing.T) {
	// Simulate
	fIgnoreMozilla = true
	fIgnoreImirhil = true
	initAPIs()
	assert.Empty(t, fnMozilla(ssllabs.Host{}))
	assert.Empty(t, fnImirhil(ssllabs.Host{}))
}

func TestInitAPI1(t *testing.T) {
	// Simulate
	fIgnoreMozilla = true
	fIgnoreImirhil = false
	initAPIs()
	assert.Empty(t, fnMozilla(ssllabs.Host{}))
	g := fnImirhil(ssllabs.Host{})
	assert.NotEmpty(t, g)
	assert.Equal(t, "Z", g)
}

func TestInitAPI2(t *testing.T) {
	// Simulate
	fIgnoreMozilla = false
	fIgnoreImirhil = true
	initAPIs()
	assert.Empty(t, fnMozilla(ssllabs.Host{}))
	assert.Empty(t, fnImirhil(ssllabs.Host{}))
}

func TestInitAPI3(t *testing.T) {
	// Simulate
	fIgnoreMozilla = false
	fIgnoreImirhil = false
	initAPIs()
	assert.Empty(t, fnMozilla(ssllabs.Host{}))
	g := fnImirhil(ssllabs.Host{})
	assert.NotEmpty(t, g)
	assert.Equal(t, "Z", g)
}

func TestTLSSite_HasExpiredTrue(t *testing.T) {
	tm := int64(1536423013000)
	assert.True(t, hasExpired(tm))
}

func TestTLSSite_HasExpiredFalse(t *testing.T) {
	tm := int64(1855828800000)
	assert.False(t, hasExpired(tm))
}

func TestDisplayWildcards(t *testing.T) {
	ji, err := ioutil.ReadFile("testdata/site.json")
	require.NoError(t, err)

	// Simulate
	fIgnoreMozilla = true
	fIgnoreImirhil = true

	all, err := ssllabs.ParseResults(ji)
	require.NoError(t, err)

	str := displayWildcards(all)
	assert.NotEmpty(t, str)
}

func TestSweet32(t *testing.T) {
	ji, err := ioutil.ReadFile("testdata/reallybad.json")
	require.NoError(t, err)

	// Simulate
	fIgnoreMozilla = true
	fIgnoreImirhil = true

	bad, err := ssllabs.ParseResults(ji)
	require.NoError(t, err)

	assert.True(t, checkSweet32(bad[0].Endpoints[0].Details))
}

const (
	mozURL = "https://http-observatory.security.mozilla.org/api/v1"
)

func TestFindServerTypeEmpty(t *testing.T) {
	tt := findServerType(ssllabs.Host{})
	require.Equal(t, TypeHTTP, tt)
}

func TestFindServerType(t *testing.T) {
	ji, err := ioutil.ReadFile("testdata/site.json")
	require.NoError(t, err)

	// Simulate
	fIgnoreMozilla = true
	fIgnoreImirhil = true

	all, err := ssllabs.ParseResults(ji)
	require.NoError(t, err)

	tt := findServerType(all[0])
	require.Equal(t, TypeHTTP, tt)
}
