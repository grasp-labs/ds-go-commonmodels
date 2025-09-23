// Package commonmodels provides common persisted fields and helpers that can be
// embedded in domain entities. It standardizes IDs, tenancy, audit fields,
// status, and flexible metadata/tags handling across services.
package core

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	status "github.com/grasp-labs/ds-go-commonmodels/v2/commonmodels/enum/status"
	types "github.com/grasp-labs/ds-go-commonmodels/v2/commonmodels/types"
	verr "github.com/grasp-labs/ds-go-commonmodels/v2/commonmodels/validation_error"
	"github.com/grasp-labs/ds-go-commonmodels/v2/commonmodels/validators/email"
)

// Now returns the current time in UTC.
//
// Override in tests for deterministic timestamps, e.g.:
//
//	fixed := time.Date(2025, 8, 18, 12, 0, 0, 0, time.UTC)
//	old := Now
//	Now = func() time.Time { return fixed }
//	defer func() { Now = old }()
var Now = func() time.Time { return time.Now().UTC() }

// BaseModel defines a common set of fields shared by persisted entities.
// Embed this type in your structs to inherit IDs, tenancy, audit metadata,
// status, and JSONB-backed free-form metadata and tags.
//
// GORM notes:
//   - ID: uses Postgres pgcrypto's gen_random_uuid() by default.
//   - TenantID: indexed; Status+TenantID composite index for common filters.
//   - CreatedAt/ModifiedAt: auto-populated by GORM; also set in hooks.
type CoreModel struct {
	ID          uuid.UUID                      `gorm:"type:uuid;primaryKey" json:"id"`
	TenantID    uuid.UUID                      `gorm:"type:uuid" json:"tenant_id"`
	OwnerID     string                         `json:"owner_id"`
	Issuer      string                         `json:"issuer"`
	Name        string                         `json:"name"`
	Version     string                         `json:"version"`
	Description string                         `json:"description"`
	Status      status.Status                  `json:"status"`
	Metadata    types.JSONB[map[string]string] `gorm:"column:metadata;type:jsonb" json:"metadata"`
	Tags        types.JSONB[map[string]string] `gorm:"column:tags;type:jsonb"     json:"tags"`
	CreatedAt   time.Time                      `json:"created_at"`
	ModifiedAt  time.Time                      `json:"modified_at"`
	CreatedBy   string                         `json:"created_by"`
	ModifiedBy  string                         `json:"modified_by"`
}

// Validate checks required fields, status values, and JSONB shape.
//
// Call this after defaults have been applied (e.g., after Create/Touch or
// after GORM hooks). It returns nil if the model is valid.
func (b *CoreModel) Validate(loc string, code string) []verr.ValidationError {
	var errs []verr.ValidationError

	// local function helper creating and appending errors.
	req := func(field, msg string) {
		errs = append(errs, verr.ValidationError{Field: field, Message: msg, Loc: loc, Code: code})
	}

	if b.ID == uuid.Nil {
		req("id", "required")
	}
	if b.TenantID == uuid.Nil {
		req("tenant_id", "required")
	}
	if b.Name == "" {
		req("name", "required")
	}
	if b.CreatedBy == "" {
		req("created_by", "required")
	}
	if !email.IsEmailFormat(b.CreatedBy) {
		req("created_by", "not a valid email format")
	}
	if b.ModifiedBy == "" {
		req("modified_by", "required")
	}
	if !email.IsEmailFormat(b.ModifiedBy) {
		req("modified_by", "not a valid email format")
	}
	if b.CreatedAt.IsZero() {
		req("created_at", "required")
	}
	if b.ModifiedAt.IsZero() {
		req("modified_at", "required")
	}

	switch b.Status {
	case status.Active, status.Deleted, status.Suspended, status.Rejected, status.Draft:
		// ok
	default:
		req("status", fmt.Sprintf("invalid %q; expected one of active, deleted, suspended, rejected, draft", b.Status))
	}

	if err := b.Metadata.Validate(); err != nil {
		req("metadata", "invalid JSON structure")
	}
	if err := b.Tags.Validate(); err != nil {
		req("tags", "invalid JSON structure")
	}

	return errs
}

// Create applies create-time defaults and audit fields to Core model.
//
// Parameters:
//   - subject: authenticated principal creating the entity; copied to CreatedBy and ModifiedBy.
//   - issuer:  identity provider or system issuing the subject; copied to Issuer.
//   - tenantID: tenant scope; sets TenantID and ensures Tags["tenant_id"].
//
// Call this exactly once during entity creation, before persisting.
// GORM hooks will also set safe defaults if Create is not called, but Create
// allows you to propagate subject/issuer explicitly.
func (c *CoreModel) Create(subject, issuer string, tenantID uuid.UUID) {
	now := Now()

	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	c.TenantID = tenantID
	c.Issuer = issuer

	c.CreatedAt = now
	c.CreatedBy = subject

	c.Touch(subject) // sets ModifiedAt/ModifiedBy

	// Ensure tenant_id tag is present without mutating a nil map.
	if c.Tags.Data == nil {
		c.Tags.Data = map[string]string{}
	}
	if _, exists := c.Tags.Data["tenant_id"]; !exists {
		c.Tags.Data["tenant_id"] = tenantID.String()
	}
}

// Touch updates the modification audit fields.
//
// Parameters:
//   - subject: authenticated principal performing the update; copied to ModifiedBy.
//
// Call this prior to persisting updates (e.g., in your service layer). GORM
// also updates ModifiedAt automatically; this method ensures ModifiedBy too.
func (c *CoreModel) Touch(subject string) {
	now := Now()
	c.ModifiedAt = now
	c.ModifiedBy = subject
}

// BeforeCreate is a GORM hook that applies safe defaults for new rows.
// It does not set CreatedBy/ModifiedBy (no subject context in hooks).
func (b *CoreModel) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	if b.CreatedAt.IsZero() {
		b.CreatedAt = Now()
	}
	if b.ModifiedAt.IsZero() {
		b.ModifiedAt = b.CreatedAt
	}
	// Ensure tenant_id tag exists if TenantID is present.
	if b.Tags.Data == nil {
		b.Tags.Data = map[string]string{}
	}
	if b.TenantID != uuid.Nil {
		if _, ok := b.Tags.Data["tenant_id"]; !ok {
			b.Tags.Data["tenant_id"] = b.TenantID.String()
		}
	}
	return nil
}

// BeforeUpdate is a GORM hook that refreshes ModifiedAt on updates.
// (ModifiedBy should be set by callers via Touch to include subject.)
func (b *CoreModel) BeforeUpdate(tx *gorm.DB) error {
	b.ModifiedAt = Now()
	return nil
}
