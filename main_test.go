package main

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/keltia/ssllabs"

	"github.com/gobuffalo/packr"

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

func TestCheckOutput(t *testing.T) {
	fh := checkOutput("")
	assert.NotEmpty(t, fh)
	assert.EqualValues(t, os.Stdout, fh)
}

func TestCheckOutput_1(t *testing.T) {
	temp, err := ioutil.TempDir("", "test")
	require.NoError(t, err)

	defer os.RemoveAll(temp)

	fn := path.Join(temp, "foo.out")
	fh := checkOutput(fn)
	assert.NotEmpty(t, fh)

	fi, err := os.Stat(fn)
	assert.NoError(t, err)
	assert.NotNil(t, fi)
	assert.NotEmpty(t, fi)
}

func TestWriteCSV(t *testing.T) {
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

	err := WriteCSV(os.Stderr, nil, cntrs, https)
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

	err := WriteCSV(os.Stderr, &TLSReport{}, cntrs, https)
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

	err = WriteCSV(null, final, cntrs, https)
	assert.NoError(t, err)
}

func TestWriteHTML(t *testing.T) {
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

	err := WriteHTML(os.Stderr, nil, cntrs, https)
	assert.Error(t, err)
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

	err := WriteHTML(os.Stderr, &TLSReport{}, cntrs, https)
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

	err = WriteHTML(null, final, cntrs, https)
	assert.NoError(t, err)
}

func TestCheckFlags_Empty(t *testing.T) {
	err := checkFlags([]string{})
	require.Error(t, err)
}

func TestCheckFlags_Nil(t *testing.T) {
	err := checkFlags(nil)
	require.Error(t, err)
}

func TestCheckFlags_Good(t *testing.T) {
	err := checkFlags([]string{"foo"})
	require.NoError(t, err)
}

func TestCheckFlags_GoodVerbose(t *testing.T) {
	fVerbose = true
	err := checkFlags([]string{"foo"})
	require.NoError(t, err)
	assert.Equal(t, 1, logLevel)
	fVerbose = false
}

func TestCheckFlags_GoodDebug(t *testing.T) {
	fDebug = true
	err := checkFlags([]string{"foo"})
	require.NoError(t, err)
	assert.Equal(t, 2, logLevel)
	fDebug = false
}
