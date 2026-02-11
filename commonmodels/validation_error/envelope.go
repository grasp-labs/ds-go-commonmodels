package validation_error

import "errors"

// ErrValidation is a sentinel error used to identify validation failures.
//
// It allows callers to detect validation errors via:
//
//	errors.Is(err, validation_error.ErrValidation)
//
// The concrete error will typically be *ErrorEnvelope.
var ErrValidation = errors.New("validation error")

// ErrorEnvelope represents a structured validation error containing
// one or more field-level validation failures.
//
// It is designed to:
//
//   - Be returned as a standard `error` from repositories/services
//   - Be inspected using `errors.As`
//   - Serialize cleanly to JSON for HTTP responses
//
// The envelope allows multiple validation errors to be returned
// in a single response, enabling better UX by reporting all
// field issues at once.
//
// Example JSON response:
//
//	{
//	  "details": [
//	    {
//	      "field": "email",
//	      "message": "invalid email format",
//	      "loc": "body",
//	      "code": "invalid"
//	    }
//	  ]
//	}
//
// Typical usage in service layer:
//
//	env := validation_error.New()
//	env.Append(validation_error.ValidationError{
//	    Field:   "email",
//	    Message: "invalid email format",
//	    Loc:     "body",
//	    Code:    "invalid",
//	})
//	return env
//
// Typical usage in HTTP handler:
//
//	var ve *validation_error.ErrorEnvelope
//	if errors.As(err, &ve) {
//	    return c.JSON(http.StatusUnprocessableEntity, ve)
//	}
//
// ErrorEnvelope implements:
//
//   - error
//   - errors.Is support via ErrValidation
//   - errors.As compatibility
//
// This makes it fully idiomatic within Go's error handling model.
type ErrorEnvelope struct {
	// Details contains individual field-level validation errors.
	// It is always initialized as an empty slice (never nil).
	Details []ValidationError `json:"details"`
}

// New creates a new, empty ErrorEnvelope.
//
// The returned envelope has a non-nil Details slice to ensure
// consistent JSON output (empty array instead of null).
func New() *ErrorEnvelope {
	return &ErrorEnvelope{Details: make([]ValidationError, 0)}
}

// Error implements the error interface.
//
// It returns the sentinel validation error message.
// Use errors.As to access structured details.
func (e *ErrorEnvelope) Error() string { return ErrValidation.Error() }

// Unwrap enables errors.Is(err, ErrValidation) checks.
func (e *ErrorEnvelope) Unwrap() error { return ErrValidation }

// Is allows direct comparison using errors.Is.
//
//	errors.Is(err, validation_error.ErrValidation)
func (e *ErrorEnvelope) Is(target error) bool { return target == ErrValidation }

// Append adds a field-level ValidationError to the envelope.
//
// It is safe to call multiple times in order to accumulate
// multiple validation failures before returning.
func (e *ErrorEnvelope) Append(v ValidationError) {
	e.Details = append(e.Details, v)
}

// Extract convenience helper
//
//	if envelope, ok := ve.Extract(err); ok {
//	   fmt.Println(envelope.Details)
//	}
func Extract(err error) (*ErrorEnvelope, bool) {
	var envelope *ErrorEnvelope
	if errors.As(err, &envelope) {
		return envelope, true
	}
	return nil, false
}
