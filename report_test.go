package TLS

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/keltia/ssllabs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/keltia/erc-checktls/site"
)

func TestNewReport(t *testing.T) {
	rep, err := NewReport([]ssllabs.Host{}, 1)
	require.Error(t, err)
	require.Nil(t, rep)
}

func TestNewReport2(t *testing.T) {
	ji, err := ioutil.ReadFile("testdata/site.json")
	require.NoError(t, err)

	// Simulate
	fIgnoreMozilla = true
	fIgnoreImirhil = true

	all, err := ssllabs.ParseResults(ji)
	require.NoError(t, err)

	sites, err := NewReport(all, 1)
	require.NoError(t, err)
	require.NotEmpty(t, sites)
}

func TestReport_ToCSV(t *testing.T) {
	var buf strings.Builder

	ji, err := ioutil.ReadFile("testdata/site.json")
	require.NoError(t, err)

	// Simulate
	fIgnoreMozilla = true
	fIgnoreImirhil = true

	all, err := ssllabs.ParseResults(ji)
	require.NoError(t, err)

	sites, err := NewReport(all, 1)
	require.NoError(t, err)

	err = sites.ToCSV(&buf)
	assert.NoError(t, err)
}

func TestReport_WriteCSV(t *testing.T) {
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

	err := r.WriteCSV(os.Stderr)
	assert.Error(t, err)

}

func TestReport_WriteCSV2(t *testing.T) {
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

	err := r.WriteCSV(os.Stderr)
	assert.Error(t, err)
}

func TestReport_WriteCSV3(t *testing.T) {
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

	fIgnoreImirhil = true
	fIgnoreMozilla = true

	final, err := NewReport(allSites, 1)
	require.NoError(t, err)

	final.cntrs = cntrs
	final.https = https

	null, err := os.OpenFile("/dev/null", os.O_WRONLY, 0666)

	err = final.WriteCSV(null)
	assert.NoError(t, err)
}

func TestReport_WriteCSVSummary(t *testing.T) {
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

	null, err := os.OpenFile("/dev/null", os.O_WRONLY, 0666)

	err = r.WriteCSVSummary(null)
	assert.Error(t, err)

}

func TestReport_WriteCSVSummary2(t *testing.T) {
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

	null, err := os.OpenFile("/dev/null", os.O_WRONLY, 0666)

	err = r.WriteCSVSummary(null)
	assert.Error(t, err)
}

func TestReport_WriteCSVSummary3(t *testing.T) {
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

	fIgnoreImirhil = true
	fIgnoreMozilla = true

	final, err := NewReport(allSites, 1)
	require.NoError(t, err)

	final.cntrs = cntrs
	final.https = https

	null, err := os.OpenFile("/dev/null", os.O_WRONLY, 0666)

	err = final.WriteCSVSummary(null)
	assert.NoError(t, err)
}

func TestReport_ColourMap(t *testing.T) {
	r := &Report{Sites: []site.TLSSite{}}
	tt := r.ColourMap("A")
	assert.NotEmpty(t, tt)
	assert.Empty(t, tt.Corrects)
}

func TestReport_ColourMap2(t *testing.T) {
	r := &Report{
		Sites: []site.TLSSite{
			{Type: TypeHTTPSok},
			{Type: TypeHTTP},
			{Type: TypeHTTPSnok},
		},
	}
	tt := r.ColourMap("A")
	assert.NotEmpty(t, tt)
	assert.NotEmpty(t, tt.Corrects)
	assert.EqualValues(t, map[string]int{"green": 1}, tt.Corrects)
	assert.Equal(t, 1, tt.Insecure)
	assert.Equal(t, 1, tt.ToFix)
}

func TestReport_GatherStats(t *testing.T) {
	ji, err := ioutil.ReadFile("testdata/site.json")
	require.NoError(t, err)

	// Simulate
	fIgnoreMozilla = true
	fIgnoreImirhil = true

	all, err := ssllabs.ParseResults(ji)
	require.NoError(t, err)

	r, err := NewReport(all, 1)
	require.NoError(t, err)
	require.NotEmpty(t, r)

	// Fake it
	r.Sites[0].Mozilla = "A+"

	t.Logf("r=%#v", r)
	t.Logf("r=%#v", r)
	assert.NotEmpty(t, r.cntrs)
	assert.EqualValues(t, map[string]int{"": 1, "A+": 1, "HSTS": 1, "Issues": 1, "OCSPStapling": 1, "PFS": 1, "Total": 1, "Z": 1}, r.cntrs)
}

func TestReport_GatherStats_1(t *testing.T) {
	ji, err := ioutil.ReadFile("testdata/site.json")
	require.NoError(t, err)

	// Simulate
	fIgnoreMozilla = true
	fIgnoreImirhil = true

	all, err := ssllabs.ParseResults(ji)
	require.NoError(t, err)

	r, err := NewReport(all, 1)
	require.NoError(t, err)
	require.NotEmpty(t, r)

	// Fake it
	r.Sites[0].Mozilla = "H"

	assert.NotEmpty(t, r.cntrs)
	assert.EqualValues(t, map[string]int{"": 1, "A+": 1, "HSTS": 1, "Issues": 1, "OCSPStapling": 1, "PFS": 1, "Total": 1, "Z": 1}, r.cntrs)
}

func TestReport_GatherStats_2(t *testing.T) {
	ji, err := ioutil.ReadFile("testdata/reallybad.json")
	require.NoError(t, err)

	// Simulate
	fIgnoreMozilla = true
	fIgnoreImirhil = true

	all, err := ssllabs.ParseResults(ji)
	require.NoError(t, err)

	r, err := NewReport(all, 1)
	require.NoError(t, err)
	require.NotEmpty(t, r)

	// Fake it
	r.Sites[0].Mozilla = "H"

	assert.NotEmpty(t, r.cntrs)
	assert.EqualValues(t, map[string]int{"A+": 1, "HSTS": 1, "Issues": 1, "OCSPStapling": 1, "PFS": 1, "Sweet32": 1, "Total": 1}, r.cntrs)
}

func TestReport_GatherStats_Null(t *testing.T) {
	ji, err := ioutil.ReadFile("testdata/null.json")
	require.NoError(t, err)

	// Simulate
	fIgnoreMozilla = true
	fIgnoreImirhil = true

	all, err := ssllabs.ParseResults(ji)
	require.NoError(t, err)

	r, err := NewReport(all, 1)
	require.NoError(t, err)
	assert.NotEmpty(t, r)

	good := map[string]int{"X": 1}

	assert.NotEmpty(t, r.cntrs)
	assert.EqualValues(t, good, r.cntrs)
}

func TestTLSReport_GatherStats_Full(t *testing.T) {
	ji, err := ioutil.ReadFile("testdata/site.json")
	require.NoError(t, err)

	// Simulate
	fIgnoreMozilla = true
	fIgnoreImirhil = true

	all, err := ssllabs.ParseResults(ji)
	require.NoError(t, err)

	r, err := NewReport(all, 1)
	require.NoError(t, err)
	require.NotEmpty(t, r)

	r.Sites[0].Mozilla = "H"

	r.GatherStats(r.Sites[0])

	assert.NotEmpty(t, r.cntrs)
	assert.EqualValues(t, map[string]int{"": 1, "A+": 2, "HSTS": 2, "Issues": 2, "OCSPStapling": 2, "PFS": 2, "Total": 2, "Z": 1}, r.cntrs)
	assert.NotEmpty(t, r.https)
	assert.EqualValues(t, 1, r.https["Bad"])
}

func TestTLSReport_GatherStats_Full1(t *testing.T) {
	ji, err := ioutil.ReadFile("testdata/site.json")
	require.NoError(t, err)

	// Simulate
	fIgnoreMozilla = true
	fIgnoreImirhil = true

	all, err := ssllabs.ParseResults(ji)
	require.NoError(t, err)

	r, err := NewReport(all, 1)
	require.NoError(t, err)
	require.NotEmpty(t, r)

	r.Sites[0].Mozilla = "C+"

	r.GatherStats(r.Sites[0])

	assert.NotEmpty(t, r.cntrs)
	assert.EqualValues(t, map[string]int{"": 1, "A+": 2, "HSTS": 2, "Issues": 2, "OCSPStapling": 2, "PFS": 2, "Total": 2, "Z": 1}, r.cntrs)
	assert.NotEmpty(t, r.https)
	assert.EqualValues(t, 1, r.https["Total"])
}
