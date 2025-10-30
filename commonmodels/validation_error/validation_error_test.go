package validation_error_test

import (
	"testing"

	ve "github.com/grasp-labs/ds-go-commonmodels/v3/commonmodels/validation_error"
)

// Test for validating Locations.
func TestValidationError_ValidLocations(t *testing.T) {
	locs := []string{"query", "body", "path", "header"}

	for _, l := range locs {
		_, ok := ve.ValidLocations[ve.Location(l)]
		if !ok {
			t.Fatalf("expected %s to be ok, got error", l)
		}
	}
}

// Test for validating Locations.
func TestValidationError_InvalidLocations(t *testing.T) {
	locs := []string{"not", "NOT", "PATH"}

	for _, l := range locs {
		_, ok := ve.ValidLocations[ve.Location(l)]
		if ok {
			t.Fatalf("expected %s not to be ok, got error", l)
		}
	}
}
