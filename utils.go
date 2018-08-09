package main

import (
	"fmt"
	"os"
)

// debug displays only if fDebug is set
func debug(str string, a ...interface{}) {
	if fDebug {
		fmt.Fprintf(os.Stderr, str, a...)
	}
}

// verbose displays only if fVerbose is set
func verbose(str string, a ...interface{}) {
	if fVerbose {
		fmt.Printf(str, a...)
	}
}

// fatalf is like log.Fatalf()
func fatalf(str string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, str, a...)
	os.Exit(1)
}
