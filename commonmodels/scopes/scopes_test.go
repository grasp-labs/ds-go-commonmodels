package scopes_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/grasp-labs/ds-go-commonmodels/v3/commonmodels/scopes"
)

func TestField(t *testing.T) {
	tests := []struct {
		name    string
		field   string
		value   interface{}
		success bool
	}{
		{"age = 30 (success)", "age", 30, true},
		{"age = 0 (fail)", "age", 0, false},
		{"name = 'John' (success)", "name", "John", true},
		{"name = empty (fail)", "name", "", false},
	}
	for _, test := range tests {
		scope := scopes.Field(test.field, test.value)
		if (scope != nil) != test.success {
			t.Errorf("%s: Expected success=%v, got success=%v", test.name, test.success, scope != nil)
		}
	}
}

func TestFieldComp(t *testing.T) {
	tests := []struct {
		name    string
		field   string
		op      string
		value   interface{}
		success bool
	}{
		{"age >= 30 (success)", "age", ">=", 30, true},
		{"age >= 0 (fail)", "age", ">=", 0, false},
		{"name != 'John' (success)", "name", "!=", "John", true},
		{"name != empty (fail)", "name", "!=", "", false},
	}
	for _, test := range tests {
		scope := scopes.FieldCmp(test.field, test.op, test.value)
		if (scope != nil) != test.success {
			t.Errorf("%s: Expected success=%v, got success=%v", test.name, test.success, scope != nil)
		}
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
		_, ok := scopes.ParseTimeFlexible(test.input)
		if ok != test.expected {
			t.Errorf("ParseTimeFlexible(%q) expected success: %v, got: %v", test.input, test.expected, ok)
		}
	}
}

func TestJSONBHasKey(t *testing.T) {
	tests := []struct {
		name    string
		field   string
		key     string
		success bool
	}{
		{"has key 'metadata' (success)", "metadata", "metadata", true},
		{"has key 'attributes' (success)", "attributes", "attributes", true},
		{"empty key (fail)", "metadata", "", false},
	}
	for _, test := range tests {
		scope := scopes.JSONBHasKey(test.field, test.key)
		if (scope != nil) != test.success {
			t.Errorf("%s: Expected success=%v, got success=%v", test.name, test.success, scope != nil)
		}
	}
}

func TestJSONBContainsValue(t *testing.T) {
	tests := []struct {
		name    string
		field   string
		value   string
		success bool
	}{
		{"metadata contains value 'active' (success)", "metadata", "active", true},
		{"attributes contains value 'enabled' (success)", "attributes", "enabled", true},
		{"empty column (fail)", "", "active", false},
		{"empty value (fail)", "metadata", "", false},
	}
	for _, test := range tests {
		scope := scopes.JSONBContainsValue(test.field, test.value)
		if (scope != nil) != test.success {
			t.Errorf("%s: Expected success=%v, got success=%v", test.name, test.success, scope != nil)
		}
	}
}

func TestJSONBHasKeyValue(t *testing.T) {
	tests := []struct {
		name    string
		field   string
		key     string
		value   string
		success bool
	}{
		{"has key 'metadata' with value 'active' (success)", "metadata", "status", "active", true},
		{"has key 'attributes' with value 'enabled' (success)", "attributes", "feature", "enabled", true},
		{"empty key (fail)", "metadata", "", "active", false},
		{"empty value (fail)", "attributes", "feature", "", false},
	}
	for _, test := range tests {
		scope := scopes.JSONBContains(test.field, test.key, test.value)
		if (scope != nil) != test.success {
			t.Errorf("%s: Expected success=%v, got success=%v", test.name, test.success, scope != nil)
		}
	}
}

func TestJSONBContainsAll(t *testing.T) {
	tests := []struct {
		name    string
		field   string
		keys    []string
		success bool
	}{
		{"contains all keys 'metadata' (success)", "metadata", []string{"status", "type"}, true},
		{"contains all keys 'attributes' (success)", "attributes", []string{"feature", "version"}, true},
		{"empty keys (fail)", "metadata", []string{}, false},
	}
	for _, test := range tests {
		keysMap := make(map[string]string)
		for _, key := range test.keys {
			keysMap[key] = "value" // add dummy value to avoid empty value
		}
		scope := scopes.JSONBContainsAll(test.field, keysMap)
		if (scope != nil) != test.success {
			t.Errorf("%s: Expected success=%v, got success=%v", test.name, test.success, scope != nil)
		}
	}
}
