package main

import "log"

// debug displays only if fDebug is set
func debug(str string, a ...interface{}) {
	if fDebug {
		log.Printf(str, a...)
	}
}

// debug displays only if fVerbose is set
func verbose(str string, a ...interface{}) {
	if fVerbose {
		log.Printf(str, a...)
	}
}
