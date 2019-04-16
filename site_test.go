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

func TestCheckKey(t *testing.T) {
	ji, err := ioutil.ReadFile("testdata/site.json")
	require.NoError(t, err)

	// Simulate
	fIgnoreMozilla = true
	fIgnoreImirhil = true

	good, err := ssllabs.ParseResults(ji)
	require.NoError(t, err)

	assert.True(t, checkKey(good[0].Certs[0]))
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
	require.NotEmpty(t, all)
	require.NotEmpty(t, all[0].CertHostnames)

	tt := findServerType(all[0])
	require.Equal(t, TypeHTTP, tt)
}

type Fmoz struct{}

func (f *Fmoz) GetGrade(site string) (string, error) {
	return "A+", nil
}

func (f *Fmoz) IsHTTPSonly(site string) (bool, error) {
	return true, nil
}

type Firml struct{}

func (c *Firml) GetScore(site string) (string, error) {
	return "A+", nil
}

func TestFindServerType2(t *testing.T) {
	var (
		fmoz  *Fmoz
		firml *Firml
	)

	ji, err := ioutil.ReadFile("testdata/ectl.json")
	require.NoError(t, err)

	// Simulate
	fIgnoreMozilla = true
	fIgnoreImirhil = true

	all, err := ssllabs.ParseResults(ji)
	require.NoError(t, err)

	fIgnoreMozilla = false

	omoz := moz
	oirml := irml
	moz = fmoz
	irml = firml

	tt := findServerType(all[0])
	require.Equal(t, TypeHTTPSok, tt)
	fIgnoreImirhil = false
	moz = omoz
	irml = oirml
}

func TestCheckIssuer_Ok(t *testing.T) {
	cert := new(ssllabs.Cert)
	cert.IssuerSubject = "foo GlobalSign bar"

	require.Equal(t, "TRUE", checkIssuer(*cert, DefaultIssuer))
}

func TestCheckIssuer_StillOk(t *testing.T) {
	cert := new(ssllabs.Cert)
	cert.IssuerSubject = "foo GlobalSign bar"
	cert.Issues = 65

	require.Equal(t, "TRUE", checkIssuer(*cert, DefaultIssuer))
}

func TestCheckIssuer_NotOk(t *testing.T) {
	cert := new(ssllabs.Cert)
	cert.IssuerSubject = "foo DigiCert bar"

	require.Equal(t, "FALSE", checkIssuer(*cert, DefaultIssuer))
}

func TestCheckIssuer_Self(t *testing.T) {
	cert := new(ssllabs.Cert)
	cert.IssuerSubject = "foo ECTL bar"
	cert.Issues = 64

	t.Logf("issues=%d", cert.Issues&0x40)
	require.Equal(t, "SELF", checkIssuer(*cert, DefaultIssuer))
}
