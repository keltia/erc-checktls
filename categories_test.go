package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDisplayCategories(t *testing.T) {
	cntrs := map[string]int{
		"A": 666,
		"B": 0,
		"G": 1,
	}
	str := displayCategories(cntrs)
	assert.NotEmpty(t, str)
}

func TestSelectColours(t *testing.T) {
	var td = []struct{ in, out string }{
		{"A+", "green"},
		{"A", "green"},
		{"A-", "orange"},
		{"B", "orange"},
		{"C", "red"},
		{"D", "red"},
	}

	for _, e := range td {
		assert.EqualValues(t, e.out, selectColours(e.in))
	}
}
