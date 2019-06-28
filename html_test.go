package TLS

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/gobuffalo/packr/v2"
	"github.com/keltia/ssllabs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTLSReport_ToHTML(t *testing.T) {
	var buf strings.Builder

	Init(Config{
		IgnoreMozilla: true,
		IgnoreImirhil: true,
	})

	ji, err := ioutil.ReadFile("testdata/site.json")
	require.NoError(t, err)

	all, err := ssllabs.ParseResults(ji)
	require.NoError(t, err)

	fJobs = 1

	sites, err := NewReport(all)
	require.NoError(t, err)
	require.NotEmpty(t, sites)

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

	r := &Report{}
	r.cntrs = cntrs
	r.https = https

	err := r.WriteHTML(os.Stderr)
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
	raw, err := ioutil.ReadFile(file)
	require.NoError(t, err)

	allSites, err := ssllabs.ParseResults(raw)
	require.NoError(t, err)

	box := packr.New("test", "./files")
	tmpls, err = loadTemplates(box)
	require.NoError(t, err)

	fIgnoreImirhil = true
	fIgnoreMozilla = true

	final, err := NewReport(allSites)
	require.NoError(t, err)

	final.cntrs = cntrs
	final.https = https

	null, err := os.OpenFile("/dev/null", os.O_WRONLY, 0666)
	require.NoError(t, err)

	err = final.WriteHTML(null)
	assert.NoError(t, err)
}
