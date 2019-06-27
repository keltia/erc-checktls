package main

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func TestCheckOutput_2(t *testing.T) {
	fh := checkOutput("/nonexistent")
	assert.Nil(t, fh)
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

func TestCheckInput(t *testing.T) {
	err := checkInput("")
	assert.Error(t, err)
}

func TestCheckInput2(t *testing.T) {
	err := checkInput("/nonexistent")
	assert.Error(t, err)
}

func TestCheckInput3(t *testing.T) {
	err := checkInput("testdata/site.json")
	assert.NoError(t, err)
}

func TestCheckInput4(t *testing.T) {
	file := "testdata/site.json"
	require.NoError(t, os.Chmod(file, 0600))
	err := checkInput(file)
	assert.NoError(t, err)
	require.NoError(t, os.Chmod(file, 0644))
}

func TestRealmain(t *testing.T) {
	ret := realmain([]string{})
	assert.Equal(t, 1, ret)
}

func TestRealmain2(t *testing.T) {
	ret := realmain([]string{"/dev/null"})
	assert.Equal(t, 1, ret)
}

func TestRealmain3(t *testing.T) {
	ret := realmain([]string{"/nonexistent"})
	assert.Equal(t, 1, ret)
}

func TestRealmain4(t *testing.T) {
	ret := realmain([]string{"testdata/site.json"})
	assert.Equal(t, 0, ret)
}

func TestRealmain5(t *testing.T) {
	fType = "html"
	ret := realmain([]string{"testdata/site.json"})
	assert.Equal(t, 0, ret)
	fType = ""
}

func TestRealmain6(t *testing.T) {
	fType = "html"
	fOutput = "/nonexistent"
	ret := realmain([]string{"testdata/site.json"})
	assert.Equal(t, 1, ret)
	fType = ""
	fOutput = ""
}

func TestRealmain7(t *testing.T) {
	file := "testdata/site.json"
	require.NoError(t, os.Chmod(file, 0600))
	ret := realmain([]string{file})
	assert.Equal(t, 0, ret)
	require.NoError(t, os.Chmod(file, 0644))
}

func TestRealmain8(t *testing.T) {
	fType = "csv"
	ret := realmain([]string{"testdata/site.json"})
	assert.Equal(t, 0, ret)
	fType = ""
}

func TestRealmain9(t *testing.T) {
	fType = "csv"
	fOutput = "/nonexistent"
	ret := realmain([]string{"testdata/site.json"})
	assert.Equal(t, 1, ret)
	fType = ""
	fOutput = ""
}

func TestRealmain10(t *testing.T) {
	fType = "csv"
	fOutput = "/dev/null"
	ret := realmain([]string{"testdata/emptysite.json"})
	assert.Equal(t, 0, ret)
	fType = ""
	fOutput = ""
}

func TestWild(t *testing.T) {
	fCmdWild = true
	ret := realmain([]string{"testdata/site.json"})
	assert.Equal(t, 0, ret)
	fCmdWild = false
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
