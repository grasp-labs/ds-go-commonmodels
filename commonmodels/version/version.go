package version

import verr "github.com/grasp-labs/ds-go-commonmodels/v3/commonmodels/validation_error"

type Version struct {
	Version string `json:"version"`
}

func (v *Version) Validate() []verr.ValidationError {
	var errors []verr.ValidationError
	if v.Version == "" {
		errors = append(errors, verr.ValidationError{
			Field:   "version",
			Message: "Required.",
		})
	}
	return errors
}
