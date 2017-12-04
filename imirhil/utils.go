package imirhil

import "log"

// debug displays only if fVerbose is set
func verbose(str string, a ...interface{}) {
    if ctx.fVerbose {
        log.Printf(str, a...)
    }
}
