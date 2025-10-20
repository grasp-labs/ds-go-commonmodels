package audit

import (
	"time"

	"github.com/google/uuid"

	val_err "github.com/grasp-labs/ds-go-commonmodels/v2/commonmodels/validation_error"
)

type AuditEntry struct {
	// Identity
	TenantID uuid.UUID `json:"tenant_id"` // From JWT or request context
	Subject  string    `json:"subject"`   // Identity from JWT 'sub'
	Jti      uuid.UUID `json:"jti"`       // Unique JWT identifier

	// Operation
	HTTPMethod string    `json:"http_method"` // "POST", "PATCH", etc.
	Resource   string    `json:"resource"`    // E.g. "target", "result"
	ResourceID uuid.UUID `json:"resource_id"` // Optional

	// Data (diff/state)
	Payload any `json:"payload"` // Partial or full state (optional)

	// Request Context
	SourceIP  string    `json:"source_ip"`  // optional
	UserAgent string    `json:"user_agent"` // optional
	Timestamp time.Time `json:"timestamp"`  // RFC3339 preferred

	// Originating API (Self-awareness)
	Service     string    `json:"service"`        // e.g., "target-api", "report-api"
	Endpoint    string    `json:"endpoint"`       // e.g., "/v1/targets/{id}"
	FullURL     string    `json:"full_url"`       // e.g., "https://api.example.com/v1/targets/abc"
	ID          uuid.UUID `json:"id"`             // Powerful for cross-trace using `request_id`
	Correlation string    `json:"correlation_id"` // Optional: cross-service trace
}

func (a *AuditEntry) Validate() []val_err.ValidationError {
	var errors []val_err.ValidationError
	if a.Subject == "" {
		errors = append(errors, val_err.ValidationError{
			Field:   "subject",
			Message: "Required.",
		})
	}

	if a.TenantID == uuid.Nil {
		errors = append(errors, val_err.ValidationError{
			Field:   "tenant_id",
			Message: "Required.",
		})
	}

	if a.HTTPMethod == "" {
		errors = append(errors, val_err.ValidationError{
			Field:   "http_method",
			Message: "Required.",
		})
	}

	if a.Resource == "" {
		errors = append(errors, val_err.ValidationError{
			Field:   "resource",
			Message: "Required.",
		})
	}

	if a.Timestamp.IsZero() {
		errors = append(errors, val_err.ValidationError{
			Field:   "timestamp",
			Message: "Required.",
		})
	}

	return errors
}
