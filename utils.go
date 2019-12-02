package TLS

import (
	"fmt"
	"os"
	"time"

	"github.com/keltia/erc-checktls/site"
)

// debug displays only if fDebug is set
func debug(str string, a ...interface{}) {
	if logLevel >= 2 {
		fmt.Fprintf(os.Stderr, str, a...)
	}
}

// verbose displays only if fVerbose is set
func verbose(str string, a ...interface{}) {
	if logLevel >= 1 {
		fmt.Printf(str, a...)
	}
}

// makeDate for month tagging.
func makeDate() string {
	return time.Now().Format("2006-01") + "-01"
}

// ByAlphabet is for sorting
type ByAlphabet []site.TLSSite

func (a ByAlphabet) Len() int           { return len(a) }
func (a ByAlphabet) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByAlphabet) Less(i, j int) bool { return a[i].Name < a[j].Name }
