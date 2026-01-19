package shared

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotFoundError_Error(t *testing.T) {
	err := NotFoundError{
		Thing: "lorem",
		Source: Source{
			File:     "ipsum",
			Function: "dolor",
		},
	}
	assert.Equal(t, "ipsum:dolor could not find lorem", err.Error())
}
