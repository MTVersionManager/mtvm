package shared

import "testing"

func TestSource_String(t *testing.T) {
	source := Source{
		File:     "lorem",
		Function: "ipsum",
	}
	expected := "lorem:ipsum"
	output := source.String()
	if output != expected {
		t.Fatalf("want %v, got %v", expected, output)
	}
}
