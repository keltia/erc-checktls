package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDebug(t *testing.T) {
	debug("false\n")
	fDebug = true
	debug("true\n")
	fDebug = false
}

func TestVerbose(t *testing.T) {
	verbose("false\n")
	fVerbose = true
	verbose("true\n")
	fVerbose = false
}

func TestMakeDate(t *testing.T) {
	str := makeDate()
	assert.NotEmpty(t, str)
}
