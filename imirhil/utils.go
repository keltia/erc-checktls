package imirhil

import "log"

// debug displays only if fDebug is set
func debug(str string, a ...interface{}) {
	if ctx.debug {
		log.Printf(str, a...)
	}
}

// debug displays only if fVerbose is set
func verbose(str string, a ...interface{}) {
	if ctx.verbose {
		log.Printf(str, a...)
	}
}
