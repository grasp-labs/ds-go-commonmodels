package validation_error

type Location string

const (
	Header Location = "header"
	Query  Location = "query"
	Body   Location = "body"
)

// The ValidationError is a common model used to
// structure error responses across our internal apis.
//
// Field: The path to the exception.
//
//	Example: data.owner.id
//
// Message: Human readable message to decipher what happened.
//
//	Example: Missing owner id
//
// Loc: Location of where the field resides
//
//	Example: body
//
// Code: The error code.
//
//	Example: required.missing
//	Example: length.min / length.max
//	Example: type.integer
//
// See ErrorEnvelope for returning validation errors to client.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Loc     string `json:"loc"`
	Code    string `json:"code"`
}

var ValidLocations = map[Location]struct{}{
	Header: {},
	Query:  {},
	Body:   {},
}
