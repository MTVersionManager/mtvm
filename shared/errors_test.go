package shared

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNotFoundError_Error(t *testing.T) {
	err := NotFoundError{
		Thing: "lorem",
		Source: Source{
			File:     "ipsum",
			Function: "dolor",
		},
	}
	expected := "ipsum:dolor could not find lorem"
	output := err.Error()
	assert.Equal(t, expected, output)
}
