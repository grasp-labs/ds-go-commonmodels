package entitlement

import (
	err "github.com/grasp-labs/ds-go-commonmodels/v2/commonmodels/enum/errors"
	val_err "github.com/grasp-labs/ds-go-commonmodels/v2/commonmodels/validation_error"
)

type Entitlement struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	TenantId string `json:"tenant_id"`
}

func (e *Entitlement) Validate(locale string) []val_err.ValidationError {
	var errors []val_err.ValidationError
	if e.ID == "" {
		errors = append(errors, val_err.ValidationError{
			Field:   "id",
			Message: err.HumanMessageLocale(locale, err.Required, "id"),
			Loc:     "body",
			Code:    err.Required,
		})
	}
	if e.Name == "" {
		errors = append(errors, val_err.ValidationError{
			Field:   "name",
			Message: err.HumanMessageLocale(locale, err.Required, "name"),
			Loc:     "body",
			Code:    err.Required,
		})
	}
	if e.TenantId == "" {
		errors = append(errors, val_err.ValidationError{
			Field:   "tenant_id",
			Message: err.HumanMessageLocale(locale, err.Required, "tenant_id"),
			Loc:     "body",
			Code:    err.Required,
		})
	}
	return errors
}
