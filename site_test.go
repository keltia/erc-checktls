package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTLSSite_HasExpiredTrue(t *testing.T) {
	tm := int64(1536423013000)
	assert.True(t, hasExpired(tm))
}

func TestTLSSite_HasExpiredFalse(t *testing.T) {
	tm := int64(1855828800000)
	assert.False(t, hasExpired(tm))
}
