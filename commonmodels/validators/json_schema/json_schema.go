// Package json_schema provides helpers to validate JSON documents against
// JSON Schema (draft-2020-12) using github.com/xeipuuv/gojsonschema.
//
// This package returns []validation_error.ValidationError for consistency
// with the rest of your APIs:
//
//   - On success (document is valid): returns nil
//   - On schema violations: returns one item per violation
//   - On validator/internal errors: returns one item with Field = NoneFieldError
package json_schema

import (
	"fmt"
	"strings"

	verr "github.com/grasp-labs/ds-go-commonmodels/v2/commonmodels/validation_error"
	"github.com/xeipuuv/gojsonschema"
)

// NoneFieldError is used for errors not tied to a specific JSON path
// (e.g., validator failures or root-level issues).
const NoneFieldError = "none_field_error"

// ValidateAgainstSchema validates the JSON document (docBytes) against the given
// JSON Schema (jsonSchema). It returns nil if the document is valid.
// On failure, it returns a slice of ValidationError with messages formatted as
// "(<keyword>): <description>" and Field set to the JSON path (or NoneFieldError).
//
// Typical usage:
//
//	b, _ := json.Marshal(myEvent)
//	if errs := json_schema.ValidateAgainstSchema(b, EventJSONSchema); len(errs) > 0 {
//	    // surface errs to the user as your API's standard validation errors
//	}
func ValidateAgainstSchema(docBytes []byte, jsonSchema []byte, loc string, code string) []verr.ValidationError {
	schemaLoader := gojsonschema.NewBytesLoader(jsonSchema)
	docLoader := gojsonschema.NewBytesLoader(docBytes)

	res, err := gojsonschema.Validate(schemaLoader, docLoader)
	if err != nil {
		return []verr.ValidationError{{
			Field:   NoneFieldError,
			Message: "schema validator error: " + err.Error(),
			Loc:     loc,
			Code:    code,
		}}
	}

	if res.Valid() {
		return nil
	}

	errs := make([]verr.ValidationError, 0, len(res.Errors()))
	for _, issue := range res.Errors() {
		field := issue.Field()
		if field == "(root)" || field == "" {
			field = NoneFieldError
		}
		kw := normalizeKeyword(issue.Type())
		errs = append(errs, verr.ValidationError{
			Field:   field,
			Message: fmt.Sprintf("(%s): %s", kw, issue.Description()),
			Loc:     loc,
			Code:    code,
		})
	}
	return errs
}

// normalizeKeyword maps library-specific issue types to canonical JSON Schema
// keywords to produce stable, user-friendly messages.
func normalizeKeyword(t string) string {
	switch strings.ToLower(t) {
	case "array_min_properties", "object_min_properties", "number_min_properties":
		return "minProperties"
	case "array_min_items":
		return "minItems"
	case "array_max_items":
		return "maxItems"
	case "string_min_length":
		return "minLength"
	case "string_max_length":
		return "maxLength"
	case "number_minimum":
		return "minimum"
	case "number_maximum":
		return "maximum"
	case "string_pattern":
		return "pattern"
	case "additional_property_not_allowed":
		return "additionalProperties"
	case "invalid_type":
		return "type"
	// common passthroughs
	case "required", "format", "enum", "const", "type", "pattern":
		return t
	default:
		return t
	}
}
