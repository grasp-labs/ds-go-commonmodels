package validation_error_test

import (
	"encoding/json"
	"testing"

	ve "github.com/grasp-labs/ds-go-commonmodels/v3/commonmodels/validation_error"
)

func TestEnvelope_WrapValidationErrorAndMarshal(t *testing.T) {
	var errors []ve.ValidationError

	errors = append(errors, ve.ValidationError{
		Field:   "foo",
		Message: "bar",
		Loc:     "header",
		Code:    "",
	})

	details := ve.ErrorEnvelope{
		Details: errors,
	}
	jsonBytes, err := json.Marshal(details)
	if err != nil {
		t.Fatalf("did not expect json marshal error, err: %v", err)
	}
	jsonStr := string(jsonBytes)
	expected := `{"details":[{"field":"foo","message":"bar","loc":"header","code":""}]}`
	if jsonStr != expected {
		t.Fatalf("failed to marshal struct, got %s, expected %s", jsonStr, expected)
	}
}

func TestEnvelope_MarshalNilDetails(t *testing.T) {

	details := ve.ErrorEnvelope{
		Details: nil,
	}
	jsonBytes, err := json.Marshal(details)
	if err != nil {
		t.Fatalf("did not expect json marshal error, err: %v", err)
	}
	jsonStr := string(jsonBytes)
	expected := `{"details":null}`
	if jsonStr != expected {
		t.Fatalf("failed to marshal struct, got %s, expected %s", jsonStr, expected)
	}
}

func TestEnvelope_MarshalEmptySlice(t *testing.T) {
	details := ve.ErrorEnvelope{
		Details: []ve.ValidationError{},
	}
	jsonBytes, err := json.Marshal(details)
	if err != nil {
		t.Fatalf("did not expect json marshal error, err: %v", err)
	}
	jsonStr := string(jsonBytes)
	expected := `{"details":[]}`
	if jsonStr != expected {
		t.Fatalf("failed to marshal struct, got %s, expected %s", jsonStr, expected)
	}
}
