package shared

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSource_String(t *testing.T) {
	source := Source{
		File:     "lorem",
		Function: "ipsum",
	}
	expected := "lorem:ipsum"
	output := source.String()
	assert.Equal(t, expected, output)
}
