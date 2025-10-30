package usage

import (
	"time"

	"github.com/google/uuid"
	err "github.com/grasp-labs/ds-go-commonmodels/v3/commonmodels/enum/errors"
	"github.com/grasp-labs/ds-go-commonmodels/v3/commonmodels/enum/status"
	val_err "github.com/grasp-labs/ds-go-commonmodels/v3/commonmodels/validation_error"
)

type UsageEntry struct {
	ID             uuid.UUID           `json:"id"`
	TenantID       uuid.UUID           `json:"tenant_id"`
	OwnerID        *string             `json:"owner_id"` // pointer to allow null
	ProductID      uuid.UUID           `json:"product_id"`
	MemoryMB       int16               `json:"memory_mb"`
	StartTimestamp time.Time           `json:"start_timestamp"`
	EndTimestamp   time.Time           `json:"end_timestamp"`
	Duration       float64             `json:"duration"`
	Status         status.Status       `json:"status"`
	Metadata       []map[string]string `json:"metadata"`
	Tags           map[string]string   `json:"tags"`
	CreatedAt      time.Time           `json:"created_at"`
	CreatedBy      string              `json:"created_by"`
}

func (u *UsageEntry) Validate(locale string) []val_err.ValidationError {
	var errors []val_err.ValidationError
	if u.ID == uuid.Nil {
		errors = append(errors, val_err.ValidationError{
			Field:   "id",
			Message: err.HumanMessageLocale(locale, err.Required, "id"),
			Loc:     "body",
			Code:    err.Required,
		})
	}
	if u.TenantID == uuid.Nil {
		errors = append(errors, val_err.ValidationError{
			Field:   "tenant_id",
			Message: err.HumanMessageLocale(locale, err.Required, "tenant_id"),
			Loc:     "body",
			Code:    err.Required,
		})
	}
	if u.ProductID == uuid.Nil {
		errors = append(errors, val_err.ValidationError{
			Field:   "product_id",
			Message: err.HumanMessageLocale(locale, err.Required, "product_id"),
			Loc:     "body",
			Code:    err.Required,
		})
	}
	if u.MemoryMB <= 0 {
		errors = append(errors, val_err.ValidationError{
			Field:   "memory_mb",
			Message: err.HumanMessageLocale(locale, err.ValidationFailed, "memory_mb"),
			Loc:     "body",
			Code:    err.ValidationFailed,
		})
	}
	if u.StartTimestamp.IsZero() {
		errors = append(errors, val_err.ValidationError{
			Field:   "start_timestamp",
			Message: err.HumanMessageLocale(locale, err.ValidationFailed, "start_timestamp"),
			Loc:     "body",
			Code:    err.ValidationFailed,
		})
	}
	if u.EndTimestamp.IsZero() {
		errors = append(errors, val_err.ValidationError{
			Field:   "end_timestamp",
			Message: err.HumanMessageLocale(locale, err.ValidationFailed, "end_timestamp"),
			Loc:     "body",
			Code:    err.ValidationFailed,
		})
	}

	if !u.StartTimestamp.IsZero() && !u.EndTimestamp.IsZero() {
		if u.EndTimestamp.Before(u.StartTimestamp) || u.EndTimestamp.Equal(u.StartTimestamp) {
			errors = append(errors, val_err.ValidationError{
				Field:   "end_timestamp",
				Message: err.HumanMessageLocale(locale, err.ValidationFailed, "end_timestamp"),
				Loc:     "body",
				Code:    err.ValidationFailed,
			})
		}
	}

	if u.Duration <= 0 {
		errors = append(errors, val_err.ValidationError{
			Field:   "duration",
			Message: err.HumanMessageLocale(locale, err.ValidationFailed, "duration"),
			Loc:     "body",
			Code:    err.ValidationFailed,
		})
	}

	if u.Status == "" {
		errors = append(errors, val_err.ValidationError{
			Field:   "status",
			Message: err.HumanMessageLocale(locale, err.Required, "status"),
			Loc:     "body",
			Code:    err.Required,
		})
	}
	if u.Tags == nil {
		errors = append(errors, val_err.ValidationError{
			Field:   "tags",
			Message: err.HumanMessageLocale(locale, err.Required, "tags"),
			Loc:     "body",
			Code:    err.Required,
		})
	}
	if u.CreatedAt.IsZero() {
		errors = append(errors, val_err.ValidationError{
			Field:   "created_at",
			Message: err.HumanMessageLocale(locale, err.ValidationFailed, "created_at"),
			Loc:     "body",
			Code:    err.ValidationFailed,
		})
	}
	if u.CreatedBy == "" {
		errors = append(errors, val_err.ValidationError{
			Field:   "created_by",
			Message: err.HumanMessageLocale(locale, err.Required, "created_by"),
			Loc:     "body",
			Code:    err.Required,
		})
	}
	return errors
}
