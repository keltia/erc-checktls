package main

import (
	"io/ioutil"
	"testing"

	"github.com/keltia/ssllabs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
