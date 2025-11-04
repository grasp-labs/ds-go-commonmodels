// HTTP Error common
//
// JSON method of struct is omnitted
// for the purpose of leaving commonmodels
// as slim as possible.
package httperror

import (
	"errors"
	"net/http"
	"strconv"

	errC "github.com/grasp-labs/ds-go-commonmodels/v3/commonmodels/enum/errors"
)

// -----------------------------------------------------------------------------
// HTTPError type
// -----------------------------------------------------------------------------

// HTTPError is the canonical error type returned to clients.
//
// Every error response has the shape:
//
//	{
//	  "error": {
//	    "code": "bad_request",
//	    "message": "invalid payload",
//	    "request_id": "3c640e85-75b3-4e0b-84c3-1b8427a64e23",
//	    "recoverable": true,
//	    "retry_after": 120
//	  }
//	}
//
// Fields:
//   - Code:       a stable, machine-readable identifier.
//   - Message:    a human-readable description safe to show clients.
//   - RequestID:  always present, injected by middleware (never omitted).
//   - Recoverable: (optional) indicates if the error is likely to be transient and if retrying might succeed.
//   - RetryAfter:  (optional) seconds to wait before retrying indicating the client should wait before retrying.
//
// Internal-only:
//   - cause:   the underlying Go error (not serialized).
//   - status:  HTTP status code to return.
type HTTPError struct {
	Code        string `json:"code"`
	Message     string `json:"message"`
	RequestID   string `json:"request_id"`
	Recoverable bool   `json:"recoverable"`
	RetryAfter  int    `json:"retry_after"`

	cause  error
	status int
}

// Error implements the error interface.
func (e *HTTPError) Error() string {
	if e == nil {
		return "<nil>"
	}
	return e.Message
}

// Unwrap enables errors.Is / errors.As to reach the underlying cause.
func (e *HTTPError) Unwrap() error { return e.cause }

// WithCause associates the underlying cause
//
// Examples:
//
//	return models.Internal(ctx, "").WithCause(err) // just wrap cause
//
// Notes:
//   - The first argument remains the underlying cause used by errors.Is/As.
func (e *HTTPError) WithCause(cause error) *HTTPError {
	// Keep the original cause for errors.Is/As
	if cause != nil {
		e.cause = cause
	}
	return e
}

// Status returns the HTTP status code for this error.
// If no explicit status is set, it falls back to defaults by Code.
func (e *HTTPError) Status() int {
	if e.status != 0 {
		return e.status
	}
	switch e.Code {
	case errC.BadRequest, errC.ValidationFailed:
		return http.StatusBadRequest
	case errC.Unauthorized:
		return http.StatusUnauthorized
	case errC.Forbidden:
		return http.StatusForbidden
	case errC.NotFound:
		return http.StatusNotFound
	case errC.Conflict:
		return http.StatusConflict
	case errC.TooManyRequests:
		return http.StatusTooManyRequests
	case errC.BadGateway:
		return http.StatusBadGateway
	default:
		return http.StatusInternalServerError
	}
}

// Response returns (statusCode, *HttpError).
//
//	status, httpErr := err.Response(); return c.JSON(status, httpErr)
func (e *HTTPError) Response() (int, http.Header, *HTTPError) {
	headers := make(http.Header)
	if e.RetryAfter > 0 {
		headers.Set("Retry-After", strconv.Itoa(e.RetryAfter))
	}
	return e.Status(), headers, e
}

// -----------------------------------------------------------------------------
// Constructors
// -----------------------------------------------------------------------------

// NewHTTPError creates a new error with the given code, message, and status.
// Normally you should use one of the convenience helpers (Internal, NotFound...).
func NewHTTPError(requestID string, code string, message string, status int) *HTTPError {
	return &HTTPError{
		Code:      code,
		Message:   message,
		RequestID: requestID,
		status:    status,
	}
}

// Internal returns a 500 Internal Server Error.
func Internal(requestID string, msg string) *HTTPError {
	if msg == "" {
		msg = "internal server error"
	}
	return NewHTTPError(requestID, errC.Internal, msg, http.StatusInternalServerError)
}

// BadGateway returns a 502 Bad Gateway Error.
func BadGateway(requestID string, msg string) *HTTPError {
	if msg == "" {
		msg = errC.HumanMessage(errC.BadGateway)
	}
	return NewHTTPError(requestID, errC.BadGateway, msg, http.StatusBadGateway)
}

// Unauthorized returns a 401 Unauthorized error.
func Unauthorized(requestID string, msg string) *HTTPError {
	if msg == "" {
		msg = "unauthorized"
	}
	return NewHTTPError(requestID, errC.Unauthorized, msg, http.StatusUnauthorized)
}

// Forbidden returns a 403 Forbidden error.
func Forbidden(requestID string, msg string) *HTTPError {
	if msg == "" {
		msg = "forbidden"
	}
	return NewHTTPError(requestID, errC.Forbidden, msg, http.StatusForbidden)
}

// NotFound returns a 404 Not Found error.
func NotFound(requestID string, msg string) *HTTPError {
	if msg == "" {
		msg = "not found"
	}
	return NewHTTPError(requestID, errC.NotFound, msg, http.StatusNotFound)
}

// Conflict returns a 409 Conflict error.
func Conflict(requestID string, msg string) *HTTPError {
	if msg == "" {
		msg = "conflict"
	}
	return NewHTTPError(requestID, errC.Conflict, msg, http.StatusConflict)
}

// TooManyRequests returns a 429 Too Many Requests error.
func TooManyRequests(requestID string, msg string) *HTTPError {
	if msg == "" {
		msg = "too many requests"
	}
	return NewHTTPError(requestID, errC.TooManyRequests, msg, http.StatusTooManyRequests).
		WithRetry(60) // Defaults to recommending retry after 60 seconds
}

// ServiceUnavailable returns a 503 Service Unavailable error.
func ServiceUnavailable(requestID string, msg string) *HTTPError {
	if msg == "" {
		msg = "service temporarily unavailable"
	}
	return NewHTTPError(requestID, errC.Internal, msg, http.StatusServiceUnavailable).
		WithRetry(30) // Defaults to recommending retry after 30 seconds
}

// BadRequest returns a 400 Bad Request error.
func BadRequest(requestID string, msg string) *HTTPError {
	if msg == "" {
		msg = "bad request"
	}
	return NewHTTPError(requestID, errC.BadRequest, msg, http.StatusBadRequest)
}

// FromError converts any error into an *HTTPError.
//   - If the error is already *HTTPError, it is returned as-is.
//   - Otherwise, it is wrapped as Internal with the original cause.
func FromError(requestID string, err error) *HTTPError {
	if err == nil {
		return nil
	}
	var he *HTTPError
	if errors.As(err, &he) {
		return he
	}
	return Internal(requestID, "").WithCause(err)
}

// WithRetry sets the RetryAfter field and returns the modified error.
// If retryAfterSeconds > 0, it also sets Recoverable to true, otherwise false.
func (e *HTTPError) WithRetry(retryAfterSeconds int) *HTTPError {
	e.RetryAfter = retryAfterSeconds

	if retryAfterSeconds > 0 {
		e.Recoverable = true
	} else {
		e.Recoverable = false
	}
	return e
}
