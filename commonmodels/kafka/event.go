package events

import (
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/grasp-labs/ds-go-commonmodels/commonmodels/types"
	verr "github.com/grasp-labs/ds-go-commonmodels/commonmodels/validation_error"
	"github.com/grasp-labs/ds-go-commonmodels/commonmodels/validators/email"
	"github.com/grasp-labs/ds-go-commonmodels/commonmodels/validators/uri"
)

// ValidationErrors is a collection of field-level validation errors.
// The element type ValidationError is defined in errors.go.
type ValidationErrors []verr.ValidationError

// MD5 hash validation: Ensure it’s 32 hex chars.
var md5Re = regexp.MustCompile(`^[a-fA-F0-9]{32}$`)

// Event define rigidly the data requirement of sending messages to
// DS Event Stream platform
type Event struct {
	ID                uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	SessionID         uuid.UUID `gorm:"type:uuid" json:"session_id"`
	RequestID         uuid.UUID `gorm:"type:uuid" json:"request_id"`
	TenantID          uuid.UUID `gorm:"type:uuid" json:"tenant_id"`
	OwnerID           *string   `json:"owner_id,omitempty"`
	EventType         string    `json:"event_type"`
	EventSource       string    `json:"event_source"`
	EventSourceURI    *string   `json:"event_source_uri,omitempty"`
	AffectedEntityURI *string   `json:"affected_entity_uri,omitempty"`
	Message           *string   `json:"message,omitempty"`

	//  Body has to be JSON
	Body      *types.JSONB[map[string]any]     `gorm:"type:jsonb" json:"body,omitempty"`
	BodyURI   *string                          `json:"body_uri,omitempty"`
	Metadata  types.JSONB[[]map[string]string] `gorm:"column:metadata;type:jsonb" json:"metadata"`
	Tags      types.JSONB[map[string]any]      `gorm:"type:jsonb" json:"tags,omitempty"`
	Timestamp time.Time                        `json:"timestamp"`
	CreatedBy string                           `json:"created_by"`
	MD5Hash   string                           `json:"md5_hash"`
}

// Validate checks required fields, status values, and JSONB shape.
//
// Call this after defaults have been applied.
// It returns nil if the model is valid.
func (e *Event) Validate() ValidationErrors {
	var errs ValidationErrors

	// local function helper creating and appending errors.
	req := func(field, msg string) { errs = append(errs, verr.ValidationError{Field: field, Message: msg}) }

	if e.ID == uuid.Nil {
		req("id", "required")
	}
	if e.SessionID == uuid.Nil {
		req("session_id", "required")
	}
	if e.RequestID == uuid.Nil { // remove if RequestID is optional
		req("request_id", "required")
	}
	if e.TenantID == uuid.Nil {
		req("tenant_id", "required")
	}
	if e.EventType == "" {
		req("event_type", "required")
	}
	if e.EventSource == "" {
		req("event_source", "required")
	}

	// body/body_uri rule: require at least one non-empty; optionally forbid both
	bodyEmpty := e.Body == nil || len(e.Body.Data) == 0
	bodyURINilOrEmpty := e.BodyURI == nil || strings.TrimSpace(*e.BodyURI) == ""
	if bodyEmpty && bodyURINilOrEmpty {
		req("body", "body and body_uri cannot both be empty")
		req("body_uri", "body and body_uri cannot both be empty")
	}

	if ve := uri.ValidateURI("event_source_uri", e.EventSourceURI, false); ve != nil {
		errs = append(errs, *ve)
	}
	if ve := uri.ValidateURI("affected_entity_uri", e.AffectedEntityURI, false); ve != nil {
		errs = append(errs, *ve)
	}
	if ve := uri.ValidateURI("body_uri", e.BodyURI, false); ve != nil {
		errs = append(errs, *ve)
	}
	if e.Timestamp.IsZero() {
		req("timestamp", "required")
	}
	if e.CreatedBy == "" {
		req("created_by", "required")
	} else if !email.IsEmailFormat(e.CreatedBy) {
		req("created_by", "not a valid email format")
	}
	if e.OwnerID != nil && strings.TrimSpace(*e.OwnerID) == "" {
		req("owner_id", "cannot be empty when provided")
	}

	if e.MD5Hash == "" || !md5Re.MatchString(e.MD5Hash) {
		req("md5_hash", "must be a 32-char hex MD5")
	}
	if e.Body != nil {
		if err := e.Body.Validate(); err != nil {
			req("body", "invalid JSON structure")
		}
	}
	if err := e.Tags.Validate(); err != nil {
		req("tags", "invalid JSON structure")
	}
	if err := e.Metadata.Validate(); err != nil {
		req("metadata", "invalid JSON structure")
	}

	return errs
}
