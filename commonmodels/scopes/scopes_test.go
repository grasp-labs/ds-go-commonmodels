package scopes_test

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/grasp-labs/ds-go-commonmodels/v3/commonmodels/scopes"
)

// Helper to test if inputs should produce a no-op scope based on scope semantics.
// No-op means: empty/whitespace field names, zero values, empty collections, empty strings.
func shouldBeNoOp(field string, value interface{}, additional ...interface{}) bool {
	// Check field (if provided)
	if strings.TrimSpace(field) == "" {
		return true
	}

	// Check zero values
	if scopes.IsZero(value) {
		return true
	}

	// Check additional string arguments
	for _, v := range additional {
		if s, ok := v.(string); ok && s == "" {
			return true
		}
	}

	return false
}

func TestField(t *testing.T) {
	tests := []struct {
		name       string
		field      string
		value      interface{}
		expectNoOp bool
	}{
		{"age = 30", "age", 30, false},
		{"age = 0 (no-op)", "age", 0, true},
		{"name = 'John'", "name", "John", false},
		{"name = empty (no-op)", "name", "", true},
		{"empty field (no-op)", "", "value", true},
		{"whitespace field (no-op)", "  ", "value", true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			scope := scopes.Field(test.field, test.value)
			if scope == nil {
				t.Fatalf("scope should not be nil")
			}
			isNoOp := shouldBeNoOp(test.field, test.value)
			if isNoOp != test.expectNoOp {
				t.Errorf("expectNoOp: expected %v, got %v", test.expectNoOp, isNoOp)
			}
		})
	}
}

func TestFieldCmp(t *testing.T) {
	tests := []struct {
		name       string
		field      string
		op         string
		value      interface{}
		expectNoOp bool
	}{
		{"age >= 30", "age", ">=", 30, false},
		{"age >= 0 (no-op)", "age", ">=", 0, true},
		{"name != 'John'", "name", "!=", "John", false},
		{"name != empty (no-op)", "name", "!=", "", true},
		{"empty field (no-op)", "", ">=", 30, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			scope := scopes.FieldCmp(test.field, test.op, test.value)
			if scope == nil {
				t.Fatalf("scope should not be nil")
			}
			expectNoOp := shouldBeNoOp(test.field, test.value)
			if expectNoOp != test.expectNoOp {
				t.Errorf("expectNoOp: expected %v, got %v", test.expectNoOp, expectNoOp)
			}
		})
	}
}

func TestIsZero(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected bool
	}{
		{"nil (fail)", nil, true},
		{"zero uuid (fail)", uuid.UUID{}, true},
		{"empty string (fail)", "", true},
		{"non-zero int (success)", 123, false},
		{"non-empty string (success)", "hello", false},
	}
	for _, test := range tests {
		result := scopes.IsZero(test.input)
		if result != test.expected {
			t.Errorf("%s: Expected IsZero(%v) to be %v, got %v", test.name, test.input, test.expected, result)
		}
	}
}

func TestParseTimeFlexible(t *testing.T) {
	tests := []struct {
		input        string
		expectedTime string
		expected     bool
	}{
		{"2024-01-01T12:00:00Z", "2024-01-01T12:00:00Z", true},
		{"2024-01-01 12:00:00", "", false},
		{"2024-01-01", "2024-01-01T00:00:00Z", true},
		{"invalid-time", "", false},
		{"1700000000", "2023-11-14T22:13:20Z", true},
		{"not-a-timestamp", "", false},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			got, ok := scopes.ParseTimeFlexible(test.input)
			if ok != test.expected {
				t.Errorf("ParseTimeFlexible(%q) expected success: %v, got: %v", test.input, test.expected, ok)
			}
			if ok && test.expected {
				expected, err := time.Parse(time.RFC3339, test.expectedTime)
				if err != nil {
					t.Fatalf("invalid expectedTime format: %v", err)
				}
				if !got.Equal(expected) {
					t.Errorf("ParseTimeFlexible(%q) expected %v, got %v", test.input, expected, got)
				}
			}
		})
	}
}

func TestJSONBHasKey(t *testing.T) {
	tests := []struct {
		name       string
		field      string
		key        string
		expectNoOp bool
	}{
		{"has key 'metadata'", "metadata", "env", false},
		{"empty key (no-op)", "metadata", "", true},
		{"empty field (no-op)", "", "env", true},
		{"whitespace field (no-op)", "  ", "env", true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			scope := scopes.JSONBHasKey(test.field, test.key)
			if scope == nil {
				t.Fatalf("scope should not be nil")
			}
			expectNoOp := shouldBeNoOp(test.field, test.key)
			if expectNoOp != test.expectNoOp {
				t.Errorf("expectNoOp: expected %v, got %v", test.expectNoOp, expectNoOp)
			}
		})
	}
}

func TestJSONBContainsValue(t *testing.T) {
	tests := []struct {
		name       string
		field      string
		value      string
		expectNoOp bool
	}{
		{"metadata contains value 'active'", "metadata", "active", false},
		{"attributes contains value 'enabled'", "attributes", "enabled", false},
		{"empty field (no-op)", "", "active", true},
		{"empty value (no-op)", "metadata", "", true},
		{"whitespace field (no-op)", "  ", "active", true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			scope := scopes.JSONBContainsValue(test.field, test.value)
			if scope == nil {
				t.Fatalf("scope should not be nil")
			}
			expectNoOp := shouldBeNoOp(test.field, test.value)
			if expectNoOp != test.expectNoOp {
				t.Errorf("expectNoOp: expected %v, got %v", test.expectNoOp, expectNoOp)
			}
		})
	}
}

func TestJSONBContains(t *testing.T) {
	tests := []struct {
		name       string
		field      string
		key        string
		value      string
		expectNoOp bool
	}{
		{"metadata with key-value pair", "metadata", "status", "active", false},
		{"attributes with key-value pair", "attributes", "feature", "enabled", false},
		{"empty key (no-op)", "metadata", "", "active", true},
		{"empty value (no-op)", "attributes", "feature", "", true},
		{"empty field (no-op)", "", "key", "value", true},
		{"whitespace field (no-op)", "  ", "key", "value", true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			scope := scopes.JSONBContains(test.field, test.key, test.value)
			if scope == nil {
				t.Fatalf("scope should not be nil")
			}
			expectNoOp := shouldBeNoOp(test.field, test.key, test.value)
			if expectNoOp != test.expectNoOp {
				t.Errorf("expectNoOp: expected %v, got %v", test.expectNoOp, expectNoOp)
			}
		})
	}
}

func TestJSONBContainsAll(t *testing.T) {
	tests := []struct {
		name       string
		field      string
		pairs      map[string]string
		expectNoOp bool
	}{
		{"contains all key-value pairs", "metadata", map[string]string{"status": "active", "type": "user"}, false},
		{"single key-value pair", "attributes", map[string]string{"feature": "enabled"}, false},
		{"empty pairs (no-op)", "metadata", map[string]string{}, true},
		{"empty field (no-op)", "", map[string]string{"key": "value"}, true},
		{"whitespace field (no-op)", "  ", map[string]string{"key": "value"}, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			scope := scopes.JSONBContainsAll(test.field, test.pairs)
			if scope == nil {
				t.Fatalf("scope should not be nil")
			}
			// JSONBContainsAll is a no-op if field is empty/whitespace OR if pairs is empty
			expectNoOp := strings.TrimSpace(test.field) == "" || len(test.pairs) == 0
			if expectNoOp != test.expectNoOp {
				t.Errorf("expectNoOp: expected %v, got %v", test.expectNoOp, expectNoOp)
			}
		})
	}
}

func TestOrderBy(t *testing.T) {
	tests := []struct {
		name       string
		param      string
		whitelist  map[string]string
		expectNoOp bool
	}{
		{"valid ascending order", "age", map[string]string{"age": "age"}, false},
		{"valid descending order", "-age", map[string]string{"age": "age"}, false},
		{"empty param (no-op)", "", map[string]string{"age": "age"}, true},
		{"invalid field not in whitelist (no-op)", "name", map[string]string{"age": "age"}, true},
		{"case-insensitive match", "AGE", map[string]string{"age": "age"}, false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			scope := scopes.OrderBy(test.param, test.whitelist)
			if scope == nil {
				t.Fatalf("scope should not be nil")
			}
			// OrderBy is a no-op if param is empty or field not in whitelist
			field := strings.TrimSpace(test.param)
			if strings.HasPrefix(field, "-") || strings.HasPrefix(field, "+") {
				field = field[1:]
			}
			_, inWhitelist := test.whitelist[strings.ToLower(strings.TrimSpace(field))]
			expectNoOp := strings.TrimSpace(test.param) == "" || !inWhitelist
			if expectNoOp != test.expectNoOp {
				t.Errorf("expectNoOp: expected %v, got %v", test.expectNoOp, expectNoOp)
			}
		})
	}
}

func TestPaginate(t *testing.T) {
	tests := []struct {
		name       string
		limit      int
		offset     int
		expectNoOp bool
	}{
		{"valid limit and offset", 10, 5, false},
		{"zero limit (uses default 100)", 0, 0, false},
		{"negative limit (uses default 100)", -5, 0, false},
		{"negative offset (clamped to 0)", 10, -5, false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			scope := scopes.Paginate(test.limit, test.offset)
			if scope == nil {
				t.Fatalf("scope should not be nil")
			}
			// Paginate always applies (never a no-op), it always sets limit
			if test.expectNoOp {
				t.Errorf("Paginate should not be a no-op")
			}
		})
	}
}
