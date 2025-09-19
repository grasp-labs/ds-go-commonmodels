package status_test

import (
	"testing"

	st "github.com/grasp-labs/ds-go-commonmodels/v2/commonmodels/enum/status"
)

// Test for validating Locations.
func TestValidationError_ValidStatus(t *testing.T) {
	locs := []string{"active", "deleted", "suspended", "rejected", "draft", "closed"}

	for _, l := range locs {
		_, ok := st.ValidStatus[st.Status(l)]
		if !ok {
			t.Fatalf("expected %s to be ok, got error", l)
		}
	}
}
