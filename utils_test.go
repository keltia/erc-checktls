package main

import "testing"

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
