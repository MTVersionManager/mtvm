package shared

import "testing"

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
	if output != expected {
		t.Fatalf("want %v, got %v", expected, output)
	}
}
