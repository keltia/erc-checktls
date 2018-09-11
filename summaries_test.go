package main

import (
	"strings"
	"testing"

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
