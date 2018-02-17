package imirhil

import "log"

// debug displays only if fDebug is set
func debug(str string, a ...interface{}) {
	if ctx.level >= 2 {
		log.Printf(str, a...)
	}
}

// debug displays only if fVerbose is set
func verbose(str string, a ...interface{}) {
	if ctx.level >= 1 {
		log.Printf(str, a...)
	}
}
