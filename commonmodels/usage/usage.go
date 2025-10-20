package usage

import (
	"time"

	"github.com/google/uuid"
	"github.com/grasp-labs/ds-go-commonmodels/v2/commonmodels/enum/status"
	val_err "github.com/grasp-labs/ds-go-commonmodels/v2/commonmodels/validation_error"
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

func (u *UsageEntry) Validate() []val_err.ValidationError {
	var errors []val_err.ValidationError
	if u.ID == uuid.Nil {
		errors = append(errors, val_err.ValidationError{
			Field:   "id",
			Message: "Required.",
		})
	}
	if u.TenantID == uuid.Nil {
		errors = append(errors, val_err.ValidationError{
			Field:   "tenant_id",
			Message: "Required.",
		})
	}
	if u.ProductID == uuid.Nil {
		errors = append(errors, val_err.ValidationError{
			Field:   "product_id",
			Message: "Required.",
		})
	}
	if u.MemoryMB <= 0 {
		errors = append(errors, val_err.ValidationError{
			Field:   "memory_mb",
			Message: "Must be greater than 0.",
		})
	}
	if u.StartTimestamp.IsZero() {
		errors = append(errors, val_err.ValidationError{
			Field:   "start_timestamp",
			Message: "Required.",
		})
	}
	if u.EndTimestamp.IsZero() {
		errors = append(errors, val_err.ValidationError{
			Field:   "end_timestamp",
			Message: "Required.",
		})
	}

	if !u.StartTimestamp.IsZero() && !u.EndTimestamp.IsZero() {
		if u.EndTimestamp.Before(u.StartTimestamp) || u.EndTimestamp.Equal(u.StartTimestamp) {
			errors = append(errors, val_err.ValidationError{
				Field:   "end_timestamp",
				Message: "Must be after start_timestamp.",
			})
		}
	}

	if u.Duration <= 0 {
		errors = append(errors, val_err.ValidationError{
			Field:   "duration",
			Message: "Required.",
		})
	}
	if u.Status == "" {
		errors = append(errors, val_err.ValidationError{
			Field:   "status",
			Message: "Required.",
		})
	}
	if u.Tags == nil {
		errors = append(errors, val_err.ValidationError{
			Field:   "tags",
			Message: "Required.",
		})
	}
	if u.CreatedAt.IsZero() {
		errors = append(errors, val_err.ValidationError{
			Field:   "created_at",
			Message: "Required.",
		})
	}
	if u.CreatedBy == "" {
		errors = append(errors, val_err.ValidationError{
			Field:   "created_by",
			Message: "Required.",
		})
	}
	return errors
}
