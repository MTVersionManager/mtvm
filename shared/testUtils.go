package shared

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func AssertIsNotFoundError(t mock.TestingT, err error, thing string, source Source) {
	require.NotNil(t, err, "want not nil error, got nil")
	var notFoundError NotFoundError
	require.ErrorAs(t, err, &notFoundError)
	assert.Equalf(t, thing, notFoundError.Thing, "want error to contain thing %v, got %v", thing, notFoundError.Thing)
	assert.Equalf(t, source, notFoundError.Source, "want error to contain source %v, got %v", source, notFoundError.Source)
}
