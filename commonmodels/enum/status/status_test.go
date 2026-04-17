package status_test

import (
	"testing"

	st "github.com/grasp-labs/ds-go-commonmodels/v3/commonmodels/enum/status"
)

// Test for validating status values.
func TestValidationError_ValidStatus(t *testing.T) {
	locs := []string{"inactive", "active", "deleted", "suspended", "rejected", "draft", "closed"}

	for _, l := range locs {
		_, ok := st.ValidStatus[st.Status(l)]
		if !ok {
			t.Fatalf("expected %s to be ok, got error", l)
		}
	}
}

// Test for validating jobstatus values
func TestValidationError_ValidJobStatus(t *testing.T) {
	locs := []string{"new", "queued", "running", "completed", "failed", "cancelled"}

	for _, l := range locs {
		_, ok := st.ValidJobStatus[st.JobStatus(l)]
		if !ok {
			t.Fatalf("expected %s to be ok, got error", l)
		}
	}
}
