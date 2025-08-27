// event_test.go
package event_test

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	events "github.com/grasp-labs/ds-go-commonmodels/commonmodels/kafka"
	"github.com/grasp-labs/ds-go-commonmodels/commonmodels/types"
)

func strp(s string) *string { return &s }
func validMD5() string      { return "d41d8cd98f00b204e9800998ecf8427e" } // 32 hex

func newValidEvent() events.Event {
	return events.Event{
		ID:                uuid.New(),
		SessionID:         uuid.New(),
		RequestID:         uuid.New(),
		TenantID:          uuid.New(),
		EventType:         "created",
		EventSource:       "unit-test",
		EventSourceURI:    strp("https://example.com/source"),
		AffectedEntityURI: strp("https://example.com/entity"),
		Body:              &types.JSONB[map[string]any]{Data: map[string]any{"k": "v"}},
		BodyURI:           nil,
		Metadata:          types.JSONB[[]map[string]string]{Data: []map[string]string{{"m": "1"}}},
		Tags:              types.JSONB[map[string]any]{Data: map[string]any{"env": "test"}},
		Timestamp:         time.Now().UTC(),
		CreatedBy:         "dev@example.com",
		MD5Hash:           validMD5(),
	}
}

func hasErr(errs events.ValidationErrors, field, contains string) bool {
	for _, e := range errs {
		if e.Field == field && strings.Contains(e.Message, contains) {
			return true
		}
	}
	return false
}

func TestEventValidate_OK(t *testing.T) {
	ev := newValidEvent()
	if errs := ev.Validate(); len(errs) != 0 {
		t.Fatalf("expected no errors, got: %+v", errs)
	}
}

func TestEventValidate_MissingRequireds(t *testing.T) {
	ev := events.Event{} // zero value
	errs := ev.Validate()
	must := []struct {
		field string
		part  string
	}{
		{"id", "required"},
		{"session_id", "required"},
		{"request_id", "required"},
		{"tenant_id", "required"},
		{"event_type", "required"},
		{"event_source", "required"},
		{"timestamp", "required"},
		{"created_by", "required"},
		{"md5_hash", "32-char hex"},
		{"body", "cannot both be empty"},
		{"body_uri", "cannot both be empty"},
	}
	for _, m := range must {
		if !hasErr(errs, m.field, m.part) {
			t.Errorf("expected error %q containing %q, got: %+v", m.field, m.part, errs)
		}
	}
}

func TestEventValidate_BodyVsBodyURI(t *testing.T) {
	// Both empty -> error on both fields
	ev := newValidEvent()
	ev.Body = nil
	ev.BodyURI = nil
	errs := ev.Validate()
	if !hasErr(errs, "body", "cannot both be empty") || !hasErr(errs, "body_uri", "cannot both be empty") {
		t.Fatalf("expected both body/body_uri emptiness errors, got: %+v", errs)
	}

	// Only Body present (non-empty) -> OK
	ev = newValidEvent()
	ev.BodyURI = nil
	if errs := ev.Validate(); len(errs) != 0 {
		t.Fatalf("body-only should be valid, got: %+v", errs)
	}

	// Only BodyURI present (non-empty) -> OK
	ev = newValidEvent()
	ev.Body = nil
	ev.BodyURI = strp("https://example.com/payload")
	if errs := ev.Validate(); len(errs) != 0 {
		t.Fatalf("body_uri-only should be valid, got: %+v", errs)
	}
}

func TestEventValidate_InvalidURIs(t *testing.T) {
	ev := newValidEvent()
	bad := "::::not-a-valid-uri"
	ev.EventSourceURI = &bad
	ev.AffectedEntityURI = &bad
	ev.BodyURI = &bad

	errs := ev.Validate()
	for _, f := range []string{"event_source_uri", "affected_entity_uri", "body_uri"} {
		if !hasErr(errs, f, "invalid URI") {
			t.Errorf("expected %s invalid URI error, got: %+v", f, errs)
		}
	}
}

func TestEventValidate_EmailAndMD5(t *testing.T) {
	ev := newValidEvent()
	ev.CreatedBy = "not-an-email"
	ev.MD5Hash = "abc"

	errs := ev.Validate()
	if !hasErr(errs, "created_by", "not a valid email format") {
		t.Errorf("expected created_by format error, got: %+v", errs)
	}
	if !hasErr(errs, "md5_hash", "32-char hex") {
		t.Errorf("expected md5_hash length/hex error, got: %+v", errs)
	}
}

func TestEventValidate_JSONBStructureErrors(t *testing.T) {
	ev := newValidEvent()

	// Make Body JSON invalid for encoding/json by inserting a channel value.
	ch := make(chan int)
	ev.Body = &types.JSONB[map[string]any]{Data: map[string]any{"bad": ch}}
	// Provide BodyURI so presence rule isn't tripped
	ev.BodyURI = strp("https://example.com/payload")
	// Make Tags invalid too
	ev.Tags = types.JSONB[map[string]any]{Data: map[string]any{"also_bad": ch}}

	errs := ev.Validate()
	if !hasErr(errs, "body", "invalid JSON structure") {
		t.Errorf("expected body invalid JSON structure, got: %+v", errs)
	}
	if !hasErr(errs, "tags", "invalid JSON structure") {
		t.Errorf("expected tags invalid JSON structure, got: %+v", errs)
	}
}

func TestEventValidate_OwnerIDProvidedButEmpty(t *testing.T) {
	ev := newValidEvent()
	ev.OwnerID = strp("   ")

	errs := ev.Validate()
	if !hasErr(errs, "owner_id", "cannot be empty") {
		t.Errorf("expected owner_id empty error, got: %+v", errs)
	}
}
