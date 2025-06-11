package plugin

import (
	"errors"
	"testing"
)

func TestNoInstalledVersion(t *testing.T) {
	_, err := InstalledVersion("loremIpsum")
	if err == nil {
		t.Fatal("Want error, got nil")
	} else if !errors.Is(err, ErrNotFound) {
		t.Fatal("Want error to contain ErrNotFound, got error not containing ErrNotFound")
	}
}
