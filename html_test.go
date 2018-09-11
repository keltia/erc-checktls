package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBooleanT(t *testing.T) {
	assert.Equal(t, white("TRUE"), booleanT(true))
	assert.Equal(t, red("FALSE"), booleanT(false))
}

func TestBooleanF(t *testing.T) {
	assert.Equal(t, white("FALSE"), booleanF(false))
	assert.Equal(t, red("TRUE"), booleanF(true))
}

func TestRed(t *testing.T) {
	str := ""
	out := red(str)
	assert.Equal(t, `<td class=xl64 align=center></td>`, out)
}

func TestRed1(t *testing.T) {
	str := "foobar"
	out := red(str)
	assert.Equal(t, `<td class=xl64 align=center>foobar</td>`, out)
}
