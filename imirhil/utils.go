package imirhil

import "log"

// debug displays only if fVerbose is set
func verbose(str string, a ...interface{}) {
    if ctx.verbose {
        log.Printf(str, a...)
    }
}
