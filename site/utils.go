package site

import (
	"fmt"
	"os"
	"time"
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

// makeDate for month tagging.
func makeDate() string {
	return time.Now().Format("2006-01") + "-01"
}
