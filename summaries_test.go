package main

import (
	"strings"
	"testing"

	"github.com/gobuffalo/packr/v2"
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
	var (
		buf strings.Builder
		err error
	)

	box := packr.New("files", "./files")
	tmpls, err = loadTemplates(box)
	require.NoError(t, err)

	cntrs := map[string]int{}
	https := map[string]int{}

	err = writeHTMLSummary(&buf, cntrs, https)
	assert.NoError(t, err)
	assert.Empty(t, buf)
}

func TestWriteHTMLSummaryEmptyT(t *testing.T) {
	var (
		buf strings.Builder
		err error
	)

	tmpls = map[string]string{}

	cntrs := map[string]int{}
	https := map[string]int{}

	err = writeHTMLSummary(&buf, cntrs, https)
	assert.Error(t, err)
	assert.Empty(t, buf)
}

func TestWriteHTMLSummary(t *testing.T) {
	var (
		buf strings.Builder
		err error
	)

	box := packr.New("files", "./files")
	tmpls, err = loadTemplates(box)
	require.NoError(t, err)

	cntrs := map[string]int{
		"A": 666,
		"B": 1,
	}

	https := map[string]int{
		"A": 666,
		"F": 42,
	}

	err = writeHTMLSummary(&buf, cntrs, https)
	assert.NoError(t, err)
	assert.NotEmpty(t, buf)
	t.Logf("buf=%s", buf.String())
}

func TestWriteHTMLSummary_1(t *testing.T) {
	var (
		buf strings.Builder
		err error
	)

	box := packr.New("files", "./files")
	tmpls, err = loadTemplates(box)
	require.NoError(t, err)

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

	err = writeHTMLSummary(&buf, cntrs, https)
	assert.NoError(t, err)
	assert.NotEmpty(t, buf)
	t.Logf("buf=%s", buf.String())
}
