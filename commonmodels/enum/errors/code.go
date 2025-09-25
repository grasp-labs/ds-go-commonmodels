// Enums and function for generating ideomatic
// responses accross service implementations.
//
// # Example
//
// HumanMessage(Required, "email")
// "email is required."
//
// HumanMessageLocale("nb", InvalidEmailFormat, "E-post")
// "E-post må være en gyldig e-postadresse."
//
// StatusFor(TooManyRequests) // 429
package errors

import (
	"fmt"
	"net/http"
)

// -----------------------------------------------------------------------------
// Machine-readable error codes
// -----------------------------------------------------------------------------

const (
	Internal           = "internal_error"
	Unauthorized       = "unauthorized"
	Forbidden          = "forbidden"
	NotFound           = "not_found"
	Conflict           = "conflict"
	BadRequest         = "bad_request"
	ValidationFailed   = "validation_failed"
	TooManyRequests    = "too_many_requests"
	Required           = "required"
	InvalidEmailFormat = "invalid_email_format"
	InvalidJSONFormat  = "invalid_json_format"
	InvalidStatus      = "invalid_status"
)

// -----------------------------------------------------------------------------
// Human readable error messages EN
// -----------------------------------------------------------------------------

var messagesEN = map[string]string{
	Internal:           "Something went wrong on our side. Please try again.",
	Unauthorized:       "You need to sign in to continue.",
	Forbidden:          "You don't have permission to perform this action.",
	NotFound:           "The requested resource was not found.",
	Conflict:           "The resource is in a conflicting state.",
	BadRequest:         "Your request could not be understood.",
	ValidationFailed:   "One or more fields failed validation.",
	TooManyRequests:    "Too many requests. Please slow down.",
	Required:           "%s is required.",
	InvalidEmailFormat: "%s must be a valid email address.",
	InvalidJSONFormat:  "Request body must be valid JSON.",
	InvalidStatus:      "The provided status %s is invalid. It should be one of active, deleted, suspended, rejected, draft",
}

// -----------------------------------------------------------------------------
// Human readable error messages NB
// -----------------------------------------------------------------------------

var messagesNB = map[string]string{
	Internal:           "Noe gikk galt hos oss. Prøv igjen.",
	Unauthorized:       "Du må være innlogget for å fortsette.",
	Forbidden:          "Du har ikke tilgang til denne handlingen.",
	NotFound:           "Forespurt ressurs ble ikke funnet.",
	Conflict:           "Ressursen er i konflikt.",
	BadRequest:         "Forespørselen kunne ikke forstås.",
	ValidationFailed:   "Ett eller flere felt feilet validering.",
	TooManyRequests:    "For mange forespørsler. Vent litt.",
	Required:           "%s er påkrevd.",
	InvalidEmailFormat: "%s må være en gyldig e-postadresse.",
	InvalidJSONFormat:  "Kroppen i forespørselen må være gyldig JSON.",
	InvalidStatus:      "Oppgitt statusverdi er ugyldig.",
}

// -----------------------------------------------------------------------------
// Message catalog
// -----------------------------------------------------------------------------

var catalogs = map[string]map[string]string{
	"en": messagesEN,
	"nb": messagesNB,
}

// HumanMessage returns a human-readable message for a machine code.
// args are passed to fmt.Sprintf to fill placeholders like %s.
//
// # Example usage
//
// // local function helper creating and appending errors.
//
//	req := func(field, code string) {
//		msg := errC.HumanMessageLocale(locale, code, field)
//		errs = append(errs, verr.ValidationError{Field: field, Message: msg, Loc: loc, Code: code})
//	}
//
// // Plain
// msg := errC.HumanMessageLocale("en", errC.Required, "id")
func HumanMessage(code string, args ...any) string {
	return HumanMessageLocale("en", code, args...)
}

// HumanMessageLocale lets you choose a locale (e.g., "en", "nb").
func HumanMessageLocale(locale, code string, args ...any) string {
	cat, ok := catalogs[locale]
	if !ok {
		cat = messagesEN
	}
	msg, ok := cat[code]
	if !ok {
		msg = messagesEN[Internal] // safe fallback
	}
	if len(args) > 0 {
		return fmt.Sprintf(msg, args...)
	}
	return msg
}

// (Optional) map codes to HTTP status for consistent responses.
var statusByCode = map[string]int{
	Internal:           http.StatusInternalServerError,
	Unauthorized:       http.StatusUnauthorized,
	Forbidden:          http.StatusForbidden,
	NotFound:           http.StatusNotFound,
	Conflict:           http.StatusConflict,
	BadRequest:         http.StatusBadRequest,
	ValidationFailed:   http.StatusBadRequest,
	TooManyRequests:    http.StatusTooManyRequests,
	Required:           http.StatusBadRequest,
	InvalidEmailFormat: http.StatusBadRequest,
	InvalidJSONFormat:  http.StatusBadRequest,
	InvalidStatus:      http.StatusBadRequest,
}

func StatusFor(code string) int {
	if s, ok := statusByCode[code]; ok {
		return s
	}
	return http.StatusInternalServerError
}
