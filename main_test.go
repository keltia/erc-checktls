package main

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

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

func TestReadContractFile(t *testing.T) {
	// We embed the file now
	box := packr.NewBox("./files")

	cntrs, err := readContractFile(box)
	assert.NoError(t, err)
	assert.NotEmpty(t, cntrs)
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
