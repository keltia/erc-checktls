package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestWriteSummary(t *testing.T) {
	var buf strings.Builder

	err := writeSummary(&buf, []string{}, nil)
	require.NoError(t, err)
}

func TestWriteSummary1(t *testing.T) {
	var buf strings.Builder

	keys := []string{"A"}
	cntrs := map[string]int{
		"A": 666,
		"B": 0,
	}
	err := writeSummary(&buf, keys, cntrs)
	require.NoError(t, err)
}

func TestWriteSummary2(t *testing.T) {
	var buf strings.Builder

	keys := []string{"C"}
	cntrs := map[string]int{
		"A": 666,
		"B": 0,
	}
	err := writeSummary(&buf, keys, cntrs)
	require.NoError(t, err)
}

func TestWriteHTMLSummaryEmpty(t *testing.T) {
	var buf strings.Builder

	cntrs := map[string]int{}

	err := writeHTMLSummary(&buf, ctlsmap, cntrs)
	assert.NoError(t, err)
	assert.Empty(t, buf)
}

func TestWriteHTMLSummary(t *testing.T) {
	var buf strings.Builder

	cntrs := map[string]int{
		"A": 666,
		"B": 42,
	}

	err := writeHTMLSummary(&buf, ctlsmap, cntrs)
	assert.NoError(t, err)
	assert.NotEmpty(t, buf)
	t.Logf("buf=%s", buf.String())
}

func TestWriteHTMLSummary_1(t *testing.T) {
	var buf strings.Builder

	cntrs := map[string]int{
		"A": 666,
		"B": 42,
		"F": 1,
	}

	err := writeHTMLSummary(&buf, httpmap, cntrs)
	assert.NoError(t, err)
	assert.NotEmpty(t, buf)
	t.Logf("buf=%s", buf.String())
}
