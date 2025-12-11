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
	return errC.StatusFor(e.Code)
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

// GetLocale returns the first locale if provided, otherwise "en".
func GetLocale(locale ...string) string {
	if len(locale) > 0 && locale[0] != "" {
		return locale[0]
	}
	return "en"
}

// Internal returns a 500 Internal Server Error.
// Locale is optional; defaults to "en".
func Internal(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.Internal)
	}
	return NewHTTPError(requestID, errC.Internal, msg, http.StatusInternalServerError)
}

// BadGateway returns a 502 Bad Gateway Error.
// Locale is optional; defaults to "en".
func BadGateway(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.BadGateway)
	}
	return NewHTTPError(requestID, errC.BadGateway, msg, http.StatusBadGateway)
}

// PartialContent returns a 206 Partial Content Error.
// Locale is optional; defaults to "en".
func PartialContent(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.PartialContent)
	}
	return NewHTTPError(requestID, errC.PartialContent, msg, http.StatusPartialContent)
}

// Unauthorized returns a 401 Unauthorized error.
// Locale is optional; defaults to "en".
func Unauthorized(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.Unauthorized)
	}
	return NewHTTPError(requestID, errC.Unauthorized, msg, http.StatusUnauthorized)
}

// Forbidden returns a 403 Forbidden error.
// Locale is optional; defaults to "en".
func Forbidden(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.Forbidden)
	}
	return NewHTTPError(requestID, errC.Forbidden, msg, http.StatusForbidden)
}

// NotFound returns a 404 Not Found error.
// Locale is optional; defaults to "en".
func NotFound(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.NotFound)
	}
	return NewHTTPError(requestID, errC.NotFound, msg, http.StatusNotFound)
}

// Conflict returns a 409 Conflict error.
// Locale is optional; defaults to "en".
func Conflict(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.Conflict)
	}
	return NewHTTPError(requestID, errC.Conflict, msg, http.StatusConflict)
}

// TooManyRequests returns a 429 Too Many Requests error.
// Locale is optional; defaults to "en".
func TooManyRequests(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.TooManyRequests)
	}
	return NewHTTPError(requestID, errC.TooManyRequests, msg, http.StatusTooManyRequests).
		WithRetry(60)
}

// ServiceUnavailable returns a 503 Service Unavailable error.
// Locale is optional; defaults to "en".
func ServiceUnavailable(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.ServiceUnavailable)
	}
	return NewHTTPError(requestID, errC.ServiceUnavailable, msg, http.StatusServiceUnavailable).
		WithRetry(30)
}

// BadRequest returns a 400 Bad Request error.
// Locale is optional; defaults to "en".
func BadRequest(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.BadRequest)
	}
	return NewHTTPError(requestID, errC.BadRequest, msg, http.StatusBadRequest)
}

// OK returns a 200 OK response.
// Locale is optional; defaults to "en".
func OK(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.OK)
	}
	return NewHTTPError(requestID, errC.OK, msg, http.StatusOK)
}

// Created returns a 201 Created response.
// Locale is optional; defaults to "en".
func Created(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.Created)
	}
	return NewHTTPError(requestID, errC.Created, msg, http.StatusCreated)
}

// Accepted returns a 202 Accepted response.
// Locale is optional; defaults to "en".
func Accepted(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.Accepted)
	}
	return NewHTTPError(requestID, errC.Accepted, msg, http.StatusAccepted)
}

// NonAuthoritativeInfo returns a 203 Non-Authoritative Information response.
// Locale is optional; defaults to "en".
func NonAuthoritativeInfo(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.NonAuthoritativeInfo)
	}
	return NewHTTPError(requestID, errC.NonAuthoritativeInfo, msg, http.StatusNonAuthoritativeInfo)
}

// NoContent returns a 204 No Content response.
// Locale is optional; defaults to "en".
func NoContent(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.NoContent)
	}
	return NewHTTPError(requestID, errC.NoContent, msg, http.StatusNoContent)
}

// ResetContent returns a 205 Reset Content response.
// Locale is optional; defaults to "en".
func ResetContent(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.ResetContent)
	}
	return NewHTTPError(requestID, errC.ResetContent, msg, http.StatusResetContent)
}

// MultiStatus returns a 207 Multi-Status response (WebDAV).
// Locale is optional; defaults to "en".
func MultiStatus(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.MultiStatus)
	}
	return NewHTTPError(requestID, errC.MultiStatus, msg, 207)
}

// AlreadyReported returns a 208 Already Reported response (WebDAV).
// Locale is optional; defaults to "en".
func AlreadyReported(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.AlreadyReported)
	}
	return NewHTTPError(requestID, errC.AlreadyReported, msg, 208)
}

// IMUsed returns a 226 IM Used response.
// Locale is optional; defaults to "en".
func IMUsed(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.IMUsed)
	}
	return NewHTTPError(requestID, errC.IMUsed, msg, 226)
}

// MultipleChoices returns a 300 Multiple Choices response.
// Locale is optional; defaults to "en".
func MultipleChoices(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.MultipleChoices)
	}
	return NewHTTPError(requestID, errC.MultipleChoices, msg, http.StatusMultipleChoices)
}

// MovedPermanently returns a 301 Moved Permanently response.
// Locale is optional; defaults to "en".
func MovedPermanently(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.MovedPermanently)
	}
	return NewHTTPError(requestID, errC.MovedPermanently, msg, http.StatusMovedPermanently)
}

// Found returns a 302 Found response.
// Locale is optional; defaults to "en".
func Found(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.Found)
	}
	return NewHTTPError(requestID, errC.Found, msg, http.StatusFound)
}

// SeeOther returns a 303 See Other response.
// Locale is optional; defaults to "en".
func SeeOther(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.SeeOther)
	}
	return NewHTTPError(requestID, errC.SeeOther, msg, http.StatusSeeOther)
}

// NotModified returns a 304 Not Modified response.
// Locale is optional; defaults to "en".
func NotModified(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.NotModified)
	}
	return NewHTTPError(requestID, errC.NotModified, msg, http.StatusNotModified)
}

// UseProxy returns a 305 Use Proxy response.
// Locale is optional; defaults to "en".
func UseProxy(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.UseProxy)
	}
	return NewHTTPError(requestID, errC.UseProxy, msg, http.StatusUseProxy)
}

// Unused returns a 306 Unused response.
// Locale is optional; defaults to "en".
func Unused(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.Unused)
	}
	return NewHTTPError(requestID, errC.Unused, msg, 306)
}

// TemporaryRedirect returns a 307 Temporary Redirect response.
// Locale is optional; defaults to "en".
func TemporaryRedirect(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.TemporaryRedirect)
	}
	return NewHTTPError(requestID, errC.TemporaryRedirect, msg, http.StatusTemporaryRedirect)
}

// PermanentRedirect returns a 308 Permanent Redirect response.
// Locale is optional; defaults to "en".
func PermanentRedirect(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.PermanentRedirect)
	}
	return NewHTTPError(requestID, errC.PermanentRedirect, msg, http.StatusPermanentRedirect)
}

// PaymentRequired returns a 402 Payment Required response.
// Locale is optional; defaults to "en".
func PaymentRequired(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.PaymentRequired)
	}
	return NewHTTPError(requestID, errC.PaymentRequired, msg, http.StatusPaymentRequired)
}

// MethodNotAllowed returns a 405 Method Not Allowed response.
// Locale is optional; defaults to "en".
func MethodNotAllowed(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.MethodNotAllowed)
	}
	return NewHTTPError(requestID, errC.MethodNotAllowed, msg, http.StatusMethodNotAllowed)
}

// NotAcceptable returns a 406 Not Acceptable response.
// Locale is optional; defaults to "en".
func NotAcceptable(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.NotAcceptable)
	}
	return NewHTTPError(requestID, errC.NotAcceptable, msg, http.StatusNotAcceptable)
}

// ProxyAuthRequired returns a 407 Proxy Authentication Required response.
// Locale is optional; defaults to "en".
func ProxyAuthRequired(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.ProxyAuthRequired)
	}
	return NewHTTPError(requestID, errC.ProxyAuthRequired, msg, http.StatusProxyAuthRequired)
}

// RequestTimeout returns a 408 Request Timeout response.
// Locale is optional; defaults to "en".
func RequestTimeout(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.RequestTimeout)
	}
	return NewHTTPError(requestID, errC.RequestTimeout, msg, http.StatusRequestTimeout)
}

// Gone returns a 410 Gone response.
// Locale is optional; defaults to "en".
func Gone(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.Gone)
	}
	return NewHTTPError(requestID, errC.Gone, msg, http.StatusGone)
}

// LengthRequired returns a 411 Length Required response.
// Locale is optional; defaults to "en".
func LengthRequired(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.LengthRequired)
	}
	return NewHTTPError(requestID, errC.LengthRequired, msg, http.StatusLengthRequired)
}

// PreconditionFailed returns a 412 Precondition Failed response.
// Locale is optional; defaults to "en".
func PreconditionFailed(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.PreconditionFailed)
	}
	return NewHTTPError(requestID, errC.PreconditionFailed, msg, http.StatusPreconditionFailed)
}

// ContentTooLarge returns a 413 Content Too Large response.
// Locale is optional; defaults to "en".
func ContentTooLarge(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.ContentTooLarge)
	}
	return NewHTTPError(requestID, errC.ContentTooLarge, msg, http.StatusRequestEntityTooLarge)
}

// URITooLong returns a 414 URI Too Long response.
// Locale is optional; defaults to "en".
func URITooLong(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.URITooLong)
	}
	return NewHTTPError(requestID, errC.URITooLong, msg, http.StatusRequestURITooLong)
}

// UnsupportedMediaType returns a 415 Unsupported Media Type response.
// Locale is optional; defaults to "en".
func UnsupportedMediaType(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.UnsupportedMediaType)
	}
	return NewHTTPError(requestID, errC.UnsupportedMediaType, msg, http.StatusUnsupportedMediaType)
}

// RangeNotSatisfiable returns a 416 Range Not Satisfiable response.
// Locale is optional; defaults to "en".
func RangeNotSatisfiable(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.RangeNotSatisfiable)
	}
	return NewHTTPError(requestID, errC.RangeNotSatisfiable, msg, http.StatusRequestedRangeNotSatisfiable)
}

// ExpectationFailed returns a 417 Expectation Failed response.
// Locale is optional; defaults to "en".
func ExpectationFailed(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.ExpectationFailed)
	}
	return NewHTTPError(requestID, errC.ExpectationFailed, msg, http.StatusExpectationFailed)
}

// ImATeapot returns a 418 I'm a teapot response.
// Locale is optional; defaults to "en".
func ImATeapot(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.ImATeapot)
	}
	return NewHTTPError(requestID, errC.ImATeapot, msg, http.StatusTeapot)
}

// MisdirectedRequest returns a 421 Misdirected Request response.
// Locale is optional; defaults to "en".
func MisdirectedRequest(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.MisdirectedRequest)
	}
	return NewHTTPError(requestID, errC.MisdirectedRequest, msg, 421)
}

// UnprocessableContent returns a 422 Unprocessable Content response.
// Locale is optional; defaults to "en".
func UnprocessableContent(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.UnprocessableContent)
	}
	return NewHTTPError(requestID, errC.UnprocessableContent, msg, 422)
}

// Locked returns a 423 Locked response.
// Locale is optional; defaults to "en".
func Locked(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.Locked)
	}
	return NewHTTPError(requestID, errC.Locked, msg, 423)
}

// FailedDependency returns a 424 Failed Dependency response.
// Locale is optional; defaults to "en".
func FailedDependency(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.FailedDependency)
	}
	return NewHTTPError(requestID, errC.FailedDependency, msg, 424)
}

// TooEarly returns a 425 Too Early response.
// Locale is optional; defaults to "en".
func TooEarly(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.TooEarly)
	}
	return NewHTTPError(requestID, errC.TooEarly, msg, 425)
}

// UpgradeRequired returns a 426 Upgrade Required response.
// Locale is optional; defaults to "en".
func UpgradeRequired(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.UpgradeRequired)
	}
	return NewHTTPError(requestID, errC.UpgradeRequired, msg, 426)
}

// PreconditionRequired returns a 428 Precondition Required response.
// Locale is optional; defaults to "en".
func PreconditionRequired(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.PreconditionRequired)
	}
	return NewHTTPError(requestID, errC.PreconditionRequired, msg, 428)
}

// RequestHeaderFieldsTooLarge returns a 431 Request Header Fields Too Large response.
// Locale is optional; defaults to "en".
func RequestHeaderFieldsTooLarge(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.RequestHeaderFieldsTooLarge)
	}
	return NewHTTPError(requestID, errC.RequestHeaderFieldsTooLarge, msg, 431)
}

// UnavailableForLegalReasons returns a 451 Unavailable For Legal Reasons response.
// Locale is optional; defaults to "en".
func UnavailableForLegalReasons(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.UnavailableForLegalReasons)
	}
	return NewHTTPError(requestID, errC.UnavailableForLegalReasons, msg, 451)
}

// NotImplemented returns a 501 Not Implemented response.
// Locale is optional; defaults to "en".
func NotImplemented(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.NotImplemented)
	}
	return NewHTTPError(requestID, errC.NotImplemented, msg, http.StatusNotImplemented)
}

// GatewayTimeout returns a 504 Gateway Timeout response.
// Locale is optional; defaults to "en".
func GatewayTimeout(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.GatewayTimeout)
	}
	return NewHTTPError(requestID, errC.GatewayTimeout, msg, http.StatusGatewayTimeout)
}

// HTTPVersionNotSupported returns a 505 HTTP Version Not Supported response.
// Locale is optional; defaults to "en".
func HTTPVersionNotSupported(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.HTTPVersionNotSupported)
	}
	return NewHTTPError(requestID, errC.HTTPVersionNotSupported, msg, http.StatusHTTPVersionNotSupported)
}

// VariantAlsoNegotiates returns a 506 Variant Also Negotiates response.
// Locale is optional; defaults to "en".
func VariantAlsoNegotiates(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.VariantAlsoNegotiates)
	}
	return NewHTTPError(requestID, errC.VariantAlsoNegotiates, msg, 506)
}

// InsufficientStorage returns a 507 Insufficient Storage response.
// Locale is optional; defaults to "en".
func InsufficientStorage(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.InsufficientStorage)
	}
	return NewHTTPError(requestID, errC.InsufficientStorage, msg, 507)
}

// LoopDetected returns a 508 Loop Detected response.
// Locale is optional; defaults to "en".
func LoopDetected(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.LoopDetected)
	}
	return NewHTTPError(requestID, errC.LoopDetected, msg, 508)
}

// NotExtended returns a 510 Not Extended response.
// Locale is optional; defaults to "en".
func NotExtended(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.NotExtended)
	}
	return NewHTTPError(requestID, errC.NotExtended, msg, 510)
}

// NetworkAuthenticationRequired returns a 511 Network Authentication Required response.
// Locale is optional; defaults to "en".
func NetworkAuthenticationRequired(requestID, msg string, locale ...string) *HTTPError {
	loc := GetLocale(locale...)
	if msg == "" {
		msg = errC.HumanMessageLocale(loc, errC.NetworkAuthenticationRequired)
	}
	return NewHTTPError(requestID, errC.NetworkAuthenticationRequired, msg, 511)
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
