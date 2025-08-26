package uri

import (
	"net/url"
	"strings"

	verr "github.com/grasp-labs/ds-go-commonmodels/commonmodels/validation_error"
)

// required=false means: empty/nil is allowed.
func ValidateURI(field string, v *string, required bool) *verr.ValidationError {
	if v == nil || strings.TrimSpace(*v) == "" {
		if required {
			return &verr.ValidationError{Field: field, Message: "required"}
		}
		return nil
	}

	s := strings.TrimSpace(*v)
	u, err := url.ParseRequestURI(s)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return &verr.ValidationError{Field: field, Message: "invalid URI"}
	}
	return nil
}
