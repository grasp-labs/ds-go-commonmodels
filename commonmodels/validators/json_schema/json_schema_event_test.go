package json_schema_test

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	ecode "github.com/grasp-labs/ds-go-commonmodels/v3/commonmodels/enum/errors"
	events "github.com/grasp-labs/ds-go-commonmodels/v3/commonmodels/kafka"
	"github.com/grasp-labs/ds-go-commonmodels/v3/commonmodels/types"
	verr "github.com/grasp-labs/ds-go-commonmodels/v3/commonmodels/validation_error"
	js "github.com/grasp-labs/ds-go-commonmodels/v3/commonmodels/validators/json_schema"
)

var EventJSONSchema = []byte(`
{
    "$schema": "https://json-schema.org/draft/2020-12/schema",
    "$id": "https://grasp-labs.com/schemas/events/event.json",
    "title": "Event",
    "type": "object",
    "additionalProperties": false,
    "properties": {
      "id":                { "type": "string", "format": "uuid" },
      "session_id":        { "type": "string", "format": "uuid" },
      "request_id":        { "type": "string", "format": "uuid" },
      "tenant_id":         { "type": "string", "format": "uuid" },
      "owner_id":          { "type": "string", "minLength": 1 },
  
      "event_type":        { "type": "string", "minLength": 1 },
      "event_source":      { "type": "string", "minLength": 1 },
  
      "event_source_uri":    { "type": "string", "format": "uri" },
      "affected_entity_uri": { "type": "string", "format": "uri" },
  
      "message":           { "type": "string" },
  
      "payload": {
        "type": "object",
        "minProperties": 1,
        "additionalProperties": true
      },
      "payload_uri": {
        "type": "string",
        "format": "uri",
        "pattern": "\\\\S"
      },
  
      "metadata": {
        "type": "object",
        "additionalProperties": true
      },
  
      "tags": {
        "type": "object",
        "additionalProperties": true
      },
  
      "timestamp":        { "type": "string", "format": "date-time" },
      "created_by":       { "type": "string", "format": "email" },
      "md5_hash":         { "type": "string", "pattern": "^[A-Fa-f0-9]{32}$" },
	  "context": {
        "type": "object",
        "minProperties": 1,
        "additionalProperties": true
      },
	  "context_uri": {
        "type": "string",
        "format": "uri",
        "pattern": "\\\\S"
      }
    },
    "required": [
      "id",
      "session_id",
      "request_id",
      "tenant_id",
      "event_type",
      "event_source",
      "timestamp",
      "created_by",
      "md5_hash"
    ],
    "anyOf": [
      { "required": ["payload"] },
      { "required": ["payload_uri"] }
    ]
  }
`)

func strp(s string) *string { return &s }
func validMD5() string      { return "d41d8cd98f00b204e9800998ecf8427e" } // 32 hex

func newValidEvent() events.Event {
	return events.Event{
		ID:                uuid.New(),
		SessionID:         uuid.New(),
		RequestID:         uuid.New(),
		TenantID:          uuid.New(),
		OwnerID:           nil,
		EventType:         "created",
		EventSource:       "unit-test",
		EventSourceURI:    strp("https://example.com/source"),
		AffectedEntityURI: strp("https://example.com/entity"),
		Message:           strp("hello"),
		Payload:           &types.JSONB[map[string]any]{Data: map[string]any{"k": "v"}},
		PayloadURI:        nil,
		Metadata:          types.JSONB[map[string]string]{Data: map[string]string{"key": "value"}},
		Tags:              types.JSONB[map[string]string]{Data: map[string]string{"env": "test"}},
		Timestamp:         time.Now().UTC(),
		CreatedBy:         "dev@example.com",
		MD5Hash:           validMD5(),
		Context:           nil,
		ContextURI:        nil,
	}
}

func marshal(t *testing.T, v any) []byte {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return b
}

func TestEvent_JSONSchema_Valid(t *testing.T) {
	ev := newValidEvent()
	if err := js.ValidateAgainstSchema(marshal(t, ev), EventJSONSchema); err != nil {
		t.Fatalf("expected schema to accept valid event, got error:\n%s", err)
	}
}

func TestEvent_JSONSchema_Failures(t *testing.T) {
	t.Run("missing_payload_and_payload_uri", func(t *testing.T) {
		ev := newValidEvent()
		ev.Payload = nil
		ev.PayloadURI = nil // both absent -> violates anyOf
		validationErrors := js.ValidateAgainstSchema(marshal(t, ev), EventJSONSchema)

		foundF := 0
		foundM := 0
		for _, e := range validationErrors {
			if e.Field == js.NoneFieldError {
				foundF++
				if strings.Contains(e.Message, "anyOf") || strings.Contains(e.Message, "required") {
					foundM++
				}
				// Assert Loc and Code are set correctly
				if e.Loc != string(verr.Body) {
					t.Errorf("Expected Loc to be %q, got %q", string(verr.Body), e.Loc)
				}
				if e.Code != ecode.ValidationFailed {
					t.Errorf("Expected Code to be %q, got %q", ecode.ValidationFailed, e.Code)
				}
			}

		}
		if foundF != 2 || foundM != 2 {
			t.Fatalf("Expected two Fields and Messages, got %d %d, err: %v", foundF, foundM, validationErrors)
		}
	})

	t.Run("empty_body_object", func(t *testing.T) {
		ev := newValidEvent()
		ev.Payload = &types.JSONB[map[string]any]{Data: map[string]any{}} // present but empty -> minProperties:1
		validationErrors := js.ValidateAgainstSchema(marshal(t, ev), EventJSONSchema)

		for i, e := range validationErrors {
			switch i {
			case 0:
				if e.Field != "payload" {
					t.Fatalf("Expected field to be payload, got %s", e.Field)
				}
				if e.Message != "(minProperties): Must have at least 1 properties" {
					t.Fatalf("Expected message '(minProperties): Must have at least 1 properties', got %s", e.Message)
				}
				// Assert Loc and Code are set correctly
				if e.Loc != string(verr.Body) {
					t.Errorf("Expected Loc to be %q, got %q", string(verr.Body), e.Loc)
				}
				if e.Code != ecode.ValidationFailed {
					t.Errorf("Expected Code to be %q, got %q", ecode.ValidationFailed, e.Code)
				}
			}
		}
	})

	t.Run("bad_uris", func(t *testing.T) {
		ev := newValidEvent()
		bad := "::::not-a-uri"
		ev.EventSourceURI = &bad
		ev.AffectedEntityURI = &bad
		ev.Payload = nil
		ev.PayloadURI = &bad
		validationErrors := js.ValidateAgainstSchema(marshal(t, ev), EventJSONSchema)
		foundF := 0
		foundM := 0
		for _, e := range validationErrors {
			if e.Field == "payload_uri" {
				foundF++
				if strings.Contains(e.Message, "pattern") || strings.Contains(e.Message, "format") {
					foundM++
				}
				if e.Loc != string(verr.Body) {
					t.Errorf("Expected Loc to be %q, got %q", string(verr.Body), e.Loc)
				}
				if e.Code != ecode.ValidationFailed {
					t.Errorf("Expected Code to be %q, got %q", ecode.ValidationFailed, e.Code)
				}
			}

			if e.Field == "affected_entity_uri" {
				foundF++
				if strings.Contains(e.Message, "format") {
					foundM++
				}
				// Assert Loc and Code for all errors
				if e.Loc != string(verr.Body) {
					t.Errorf("Expected Loc to be %q, got %q", string(verr.Body), e.Loc)
				}
				if e.Code != ecode.ValidationFailed {
					t.Errorf("Expected Code to be %q, got %q", ecode.ValidationFailed, e.Code)
				}
			}

			if e.Field == "event_source_uri" {
				foundF++
				if strings.Contains(e.Message, "format") {
					foundM++
				}
				// Assert Loc and Code for all errors
				if e.Loc != string(verr.Body) {
					t.Errorf("Expected Loc to be %q, got %q", string(verr.Body), e.Loc)
				}
				if e.Code != ecode.ValidationFailed {
					t.Errorf("Expected Code to be %q, got %q", ecode.ValidationFailed, e.Code)
				}
			}

		}
		if foundF != 4 || foundM != 4 {
			t.Fatalf("Expected two Fields and Messages, got %d %d, err: %v", foundF, foundM, validationErrors)
		}

	})

	t.Run("bad_email_and_md5", func(t *testing.T) {
		ev := newValidEvent()
		ev.CreatedBy = "not-an-email"
		ev.MD5Hash = "abc"
		validationErrors := js.ValidateAgainstSchema(marshal(t, ev), EventJSONSchema)
		foundF := 0
		foundM := 0
		for _, e := range validationErrors {
			if e.Field == "created_by" {
				foundF++
				if strings.Contains(e.Message, "pattern") || strings.Contains(e.Message, "format") {
					foundM++
				}
				// Assert Loc and Code for all errors
				if e.Loc != string(verr.Body) {
					t.Errorf("Expected Loc to be %q, got %q", string(verr.Body), e.Loc)
				}
				if e.Code != ecode.ValidationFailed {
					t.Errorf("Expected Code to be %q, got %q", ecode.ValidationFailed, e.Code)
				}
			}

			if e.Field == "md5_hash" {
				foundF++
				if strings.Contains(e.Message, "format") {
					foundM++
				}
				// Assert Loc and Code for all errors
				if e.Loc != string(verr.Body) {
					t.Errorf("Expected Loc to be %q, got %q", string(verr.Body), e.Loc)
				}
				if e.Code != ecode.ValidationFailed {
					t.Errorf("Expected Code to be %q, got %q", ecode.ValidationFailed, e.Code)
				}
			}
		}
	})

	t.Run("validator_error", func(t *testing.T) {
		// Test with an invalid schema to trigger the validator error path
		invalidSchema := []byte(`{invalid json`)
		ev := newValidEvent()
		validationErrors := js.ValidateAgainstSchema(marshal(t, ev), invalidSchema)

		if len(validationErrors) != 1 {
			t.Fatalf("Expected exactly one validation error, got %d", len(validationErrors))
		}

		e := validationErrors[0]
		// Validator errors should have Field = NoneFieldError and proper Loc and Code
		if e.Field != js.NoneFieldError {
			t.Errorf("Expected Field to be %q, got %q", js.NoneFieldError, e.Field)
		}
		if e.Loc != string(verr.Body) {
			t.Errorf("Expected Loc to be %q, got %q", string(verr.Body), e.Loc)
		}
		if e.Code != ecode.ValidationFailed {
			t.Errorf("Expected Code to be %q, got %q", ecode.ValidationFailed, e.Code)
		}
		if !strings.Contains(e.Message, "schema validator error") {
			t.Errorf("Expected message to contain 'schema validator error', got %q", e.Message)
		}
	})
}
