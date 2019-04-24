package site

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/keltia/ssllabs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var f = Flags{
	IgnoreImirhil: true,
	IgnoreMozilla: true,
	Contracts:     map[string]string{},
}

func Setup(t *testing.T) {
	Init(f)

	require.NoError(t, os.Unsetenv("http_proxy"))
	require.NoError(t, os.Unsetenv("https_proxy"))
	require.NoError(t, os.Unsetenv("all_proxy"))
}

func TestNewFromHost(t *testing.T) {
	ji, err := ioutil.ReadFile("../testdata/site.json")
	require.NoError(t, err)

	Setup(t)

	all, err := ssllabs.ParseResults(ji)
	require.NoError(t, err)

	tls := NewFromHost(all[0])
	require.NotEmpty(t, tls)
}

func TestNewFromHost2(t *testing.T) {
	Setup(t)

	tls := NewFromHost(ssllabs.Host{})
	require.Empty(t, tls)
}

func TestInit(t *testing.T) {
	Setup(t)

	assert.Empty(t, fnMozilla(ssllabs.Host{}))
	assert.Empty(t, fnImirhil(ssllabs.Host{}))
}

func TestInit1(t *testing.T) {
	Init(Flags{
		IgnoreImirhil: false,
		IgnoreMozilla: true,
		Contracts:     map[string]string{},
	})

	assert.Empty(t, fnMozilla(ssllabs.Host{}))
	g := fnImirhil(ssllabs.Host{})
	assert.NotEmpty(t, g)
	assert.Equal(t, "Z", g)
}

func TestInit2(t *testing.T) {
	Init(Flags{
		IgnoreImirhil: true,
		IgnoreMozilla: false,
		Contracts:     map[string]string{},
	})

	assert.Empty(t, fnMozilla(ssllabs.Host{}))
	assert.Empty(t, fnImirhil(ssllabs.Host{}))
}

func TestInit3(t *testing.T) {
	Init(Flags{
		IgnoreImirhil: false,
		IgnoreMozilla: false,
		Contracts:     map[string]string{},
	})

	assert.Empty(t, fnMozilla(ssllabs.Host{}))
	g := fnImirhil(ssllabs.Host{})
	assert.NotEmpty(t, g)
	assert.Equal(t, "Z", g)
}

func Test_HasExpiredTrue(t *testing.T) {
	tm := int64(1536423013000)
	assert.True(t, hasExpired(tm))
}

func Test_HasExpiredFalse(t *testing.T) {
	tm := int64(1855828800000)
	assert.False(t, hasExpired(tm))
}

func TestSweet32(t *testing.T) {
	ji, err := ioutil.ReadFile("../testdata/reallybad.json")
	require.NoError(t, err)

	Setup(t)

	bad, err := ssllabs.ParseResults(ji)
	require.NoError(t, err)

	assert.True(t, checkSweet32(bad[0].Endpoints[0].Details))
}

func TestCheckKey(t *testing.T) {
	ji, err := ioutil.ReadFile("../testdata/site.json")
	require.NoError(t, err)

	Setup(t)

	good, err := ssllabs.ParseResults(ji)
	require.NoError(t, err)

	assert.True(t, checkKey(good[0].Certs[0]))
}

func TestFindServerTypeEmpty(t *testing.T) {
	tt := findServerType(ssllabs.Host{})
	require.Equal(t, TypeHTTP, tt)
}

func TestFindServerType(t *testing.T) {
	ji, err := ioutil.ReadFile("../testdata/site.json")
	require.NoError(t, err)

	Setup(t)

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

	ji, err := ioutil.ReadFile("../testdata/ectl.json")
	require.NoError(t, err)

	Setup(t)

	all, err := ssllabs.ParseResults(ji)
	require.NoError(t, err)

	fIgnoreMozilla = false

	// Save & swap
	omoz, oirml := moz, irml
	moz, irml = fmoz, firml

	tt := findServerType(all[0])
	require.Equal(t, TypeHTTPSok, tt)
	fIgnoreImirhil = false

	moz, irml = omoz, oirml
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

func TestCheckHSTS_Empty(t *testing.T) {
	endp := &ssllabs.EndpointDetails{}

	require.Equal(t, int64(-1), checkHSTS(*endp))
}

func TestCheckHSTS_Present(t *testing.T) {
	endp := &ssllabs.EndpointDetails{
		HstsPolicy: ssllabs.HstsPolicy{MaxAge: 666, Status: "present"},
	}

	require.Equal(t, int64(666), checkHSTS(*endp))
}
