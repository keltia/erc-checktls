package TLS

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDebug(t *testing.T) {
	debug("false\n")
	logLevel = 2
	debug("true\n")
	logLevel = 0
}

func TestVerbose(t *testing.T) {
	verbose("false\n")
	logLevel = 1
	verbose("true\n")
	logLevel = 0
}

func TestMakeDate(t *testing.T) {
	str := makeDate()
	assert.NotEmpty(t, str)
}
