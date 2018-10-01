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
