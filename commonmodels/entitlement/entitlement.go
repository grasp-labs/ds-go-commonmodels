package entitlement

import (
	val_err "github.com/grasp-labs/ds-go-commonmodels/v2/commonmodels/validation_error"
)

type Entitlement struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	TenantId string `json:"tenant_id"`
}

func (e *Entitlement) Validate() []val_err.ValidationError {
	var errors []val_err.ValidationError
	if e.ID == "" {
		errors = append(errors, val_err.ValidationError{
			Field:   "id",
			Message: "Required.",
		})
	}
	if e.Name == "" {
		errors = append(errors, val_err.ValidationError{
			Field:   "name",
			Message: "Required.",
		})
	}
	if e.TenantId == "" {
		errors = append(errors, val_err.ValidationError{
			Field:   "tenant_id",
			Message: "Required.",
		})
	}
	return errors
}
