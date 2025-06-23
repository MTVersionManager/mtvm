package shared

import (
	"errors"
	"testing"
)

func AssertIsNotFoundError(t *testing.T, err error, thing string, source Source) {
	if err == nil {
		t.Fatal("want error, got nil")
	}
	var notFoundError NotFoundError
	errors.As(err, &notFoundError)
	if !errors.As(err, &notFoundError) {
		t.Fatalf("want error to be NotFoundError, got error not containing NotFoundError")
	}
	if notFoundError.Thing != thing {
		t.Fatalf("want error to contain thing %v, got %v", thing, notFoundError.Thing)
	}
	if notFoundError.Source != source {
		t.Fatalf("want error to contain source %v, got %v", source, notFoundError.Source)
	}
}
