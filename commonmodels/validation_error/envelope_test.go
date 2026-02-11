package validation_error_test

import (
	"encoding/json"
	"errors"
	"testing"

	ve "github.com/grasp-labs/ds-go-commonmodels/v3/commonmodels/validation_error"
	"github.com/stretchr/testify/assert"
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

func throwErr() error {
	env := ve.New()
	env.Append(ve.ValidationError{Field: "email", Message: "invalid", Loc: string(ve.Body), Code: "invalid"})
	return env
}

func TestErrorIs(t *testing.T) {
	err := throwErr()
	assert.True(t, errors.Is(err, ve.ErrValidation))
}

func TestErrorAs(t *testing.T) {
	err := throwErr()
	var errEnvelope *ve.ErrorEnvelope
	var isEnv = false
	if errors.As(err, &errEnvelope) {
		isEnv = true
	}
	assert.True(t, isEnv)
}

func TestEnvelopeCanBeRetrieved(t *testing.T) {
	err := throwErr()
	var env *ve.ErrorEnvelope
	if errors.As(err, &env) {
		if len(env.Details) != 1 {
			t.Fatalf("expected 1 validation error, got: %v", env)
		}
	} else {
		t.Fatalf("expected validation error envelope")
	}
}

func TestExtract(t *testing.T) {
	err := throwErr()
	env, ok := ve.Extract(err)
	if !ok {
		t.Fatalf("expected ok, got %v, %v", env, ok)
	}
	msg := env.Error()
	if msg != "validation error" {
		t.Fatalf("expected validation error, got %s", msg)
	}

}
