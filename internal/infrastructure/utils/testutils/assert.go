package testutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func AssertError(t *testing.T, expected string, err error) {
	if expected == "" {
		assert.NoError(t, err)
	} else {
		assert.EqualError(t, err, expected)
	}
}
