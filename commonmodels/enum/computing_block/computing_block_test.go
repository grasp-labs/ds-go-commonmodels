package computingblock_test

import (
	"encoding/json"
	"testing"

	cb "github.com/grasp-labs/ds-go-commonmodels/commonmodels/enum/computing_block"
)

// test-only helper (kept here to avoid changing the package API)
func isValid(c cb.ComputingBlock) bool {
	switch c {
	case cb.Workflow, cb.Pipeline, cb.Clone:
		return true
	default:
		return false
	}
}

func TestComputingBlock_StringValues(t *testing.T) {
	tests := []struct {
		name string
		got  cb.ComputingBlock
		want string
	}{
		{"workflow", cb.Workflow, "workflow"},
		{"pipeline", cb.Pipeline, "pipeline"},
		{"clone", cb.Clone, "clone"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.got) != tt.want {
				t.Fatalf("string(%q) != %q", tt.got, tt.want)
			}
		})
	}
}

func TestComputingBlock_JSONRoundTrip(t *testing.T) {
	vals := []cb.ComputingBlock{cb.Workflow, cb.Pipeline, cb.Clone}
	for _, v := range vals {
		t.Run(string(v), func(t *testing.T) {
			b, err := json.Marshal(v)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			var out cb.ComputingBlock
			if err := json.Unmarshal(b, &out); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			if out != v {
				t.Fatalf("round-trip mismatch: got %q want %q", out, v)
			}
			if !isValid(out) {
				t.Fatalf("round-trip produced invalid value: %q", out)
			}
		})
	}
}

func TestComputingBlock_InvalidValues(t *testing.T) {
	// Zero value should be invalid
	var zero cb.ComputingBlock
	if isValid(zero) {
		t.Fatalf("zero value should be invalid, got %q", zero)
	}

	// Arbitrary string should unmarshal but be considered invalid by our helper.
	var got cb.ComputingBlock
	if err := json.Unmarshal([]byte(`"not-a-valid-kind"`), &got); err != nil {
		t.Fatalf("unmarshal invalid: %v", err)
	}
	if isValid(got) {
		t.Fatalf("expected invalid value after unmarshal, got %q", got)
	}
}
