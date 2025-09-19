package validation_error_test

import (
	"testing"

	ve "github.com/grasp-labs/ds-go-commonmodels/v2/commonmodels/validation_error"
)

// Test for validating Locations.
func TestValidationError_ValidLocations(t *testing.T) {
	locs := []string{"query", "body"}

	for _, l := range locs {
		_, ok := ve.ValidLocations[ve.Location(l)]
		if !ok {
			t.Fatalf("expected %s to be ok, got error", l)
		}
	}
}
