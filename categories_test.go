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
