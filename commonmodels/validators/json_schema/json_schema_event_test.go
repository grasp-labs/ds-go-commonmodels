package json_schema_test

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	events "github.com/grasp-labs/ds-go-commonmodels/commonmodels/kafka"
	"github.com/grasp-labs/ds-go-commonmodels/commonmodels/types"
	js "github.com/grasp-labs/ds-go-commonmodels/commonmodels/validators/json_schema"
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
  
      "body": {
        "type": "object",
        "minProperties": 1,
        "additionalProperties": true
      },
      "body_uri": {
        "type": "string",
        "format": "uri",
        "pattern": "\\\\S"
      },
  
      "metadata": {
        "type": "array",
        "items": {
          "type": "object",
          "additionalProperties": { "type": "string" }
        }
      },
  
      "tags": {
        "type": "object",
        "additionalProperties": true
      },
  
      "timestamp":        { "type": "string", "format": "date-time" },
      "created_by":       { "type": "string", "format": "email" },
      "md5_hash":         { "type": "string", "pattern": "^[A-Fa-f0-9]{32}$" }
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
      { "required": ["body"] },
      { "required": ["body_uri"] }
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
		Body:              &types.JSONB[map[string]any]{Data: map[string]any{"k": "v"}},
		BodyURI:           nil,
		Metadata:          types.JSONB[[]map[string]string]{Data: []map[string]string{{"key": "value"}}},
		Tags:              types.JSONB[map[string]string]{Data: map[string]string{"env": "test"}},
		Timestamp:         time.Now().UTC(),
		CreatedBy:         "dev@example.com",
		MD5Hash:           validMD5(),
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
	t.Run("missing_body_and_body_uri", func(t *testing.T) {
		ev := newValidEvent()
		ev.Body = nil
		ev.BodyURI = nil // both absent -> violates anyOf
		validationErrors := js.ValidateAgainstSchema(marshal(t, ev), EventJSONSchema)

		foundF := 0
		foundM := 0
		for _, e := range validationErrors {
			if e.Field == "none_field_error" {
				foundF++
				if strings.Contains(e.Message, "anyOf") || strings.Contains(e.Message, "required") {
					foundM++
				}
			}

		}
		if foundF != 2 || foundM != 2 {
			t.Fatalf("Expected two Fields and Messages, got %d %d, err: %v", foundF, foundM, validationErrors)
		}
	})

	t.Run("empty_body_object", func(t *testing.T) {
		ev := newValidEvent()
		ev.Body = &types.JSONB[map[string]any]{Data: map[string]any{}} // present but empty -> minProperties:1
		validationErrors := js.ValidateAgainstSchema(marshal(t, ev), EventJSONSchema)

		for i, e := range validationErrors {
			switch i {
			case 0:
				if e.Field != "body" {
					t.Fatalf("Expected field to be body, got %s", e.Field)
				}
				if e.Message != "(minProperties): Must have at least 1 properties" {
					t.Fatalf("Expected message '(minProperties): Must have at least 1 properties', got %s", e.Message)
				}
			}
		}
	})

	t.Run("bad_uris", func(t *testing.T) {
		ev := newValidEvent()
		bad := "::::not-a-uri"
		ev.EventSourceURI = &bad
		ev.AffectedEntityURI = &bad
		ev.Body = nil
		ev.BodyURI = &bad
		validationErrors := js.ValidateAgainstSchema(marshal(t, ev), EventJSONSchema)
		foundF := 0
		foundM := 0
		for _, e := range validationErrors {
			if e.Field == "body_uri" {
				foundF++
				if strings.Contains(e.Message, "pattern") || strings.Contains(e.Message, "format") {
					foundM++
				}
			}

			if e.Field == "affected_entity_uri" {
				foundF++
				if strings.Contains(e.Message, "format") {
					foundM++
				}
			}

			if e.Field == "event_source_uri" {
				foundF++
				if strings.Contains(e.Message, "format") {
					foundM++
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
			}

			if e.Field == "md5_hash" {
				foundF++
				if strings.Contains(e.Message, "format") {
					foundM++
				}
			}
		}
	})
}
