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
	"strings"
)

// -----------------------------------------------------------------------------
// Machine-readable error codes
// -----------------------------------------------------------------------------

const (
	Internal                      = "internal_error"
	Unauthorized                  = "unauthorized"
	Forbidden                     = "forbidden"
	NotFound                      = "not_found"
	Conflict                      = "conflict"
	BadRequest                    = "bad_request"
	ValidationFailed              = "validation_failed"
	TooManyRequests               = "too_many_requests"
	Required                      = "required"
	InvalidEmailFormat            = "invalid_email_format"
	InvalidJSONFormat             = "invalid_json_format"
	InvalidStatus                 = "invalid_status"
	Invalid                       = "invalid"
	InvalidDataType               = "invalid_data_type"
	BadGateway                    = "bad_gateway"
	OK                            = "ok"
	Created                       = "created"
	Accepted                      = "accepted"
	NonAuthoritativeInfo          = "non_authoritative_info"
	NoContent                     = "no_content"
	ResetContent                  = "reset_content"
	PartialContent                = "partial_content"
	MultiStatus                   = "multi_status"
	AlreadyReported               = "already_reported"
	IMUsed                        = "im_used"
	MultipleChoices               = "multiple_choices"
	MovedPermanently              = "moved_permanently"
	Found                         = "found"
	SeeOther                      = "see_other"
	NotModified                   = "not_modified"
	UseProxy                      = "use_proxy"
	Unused                        = "unused"
	TemporaryRedirect             = "temporary_redirect"
	PermanentRedirect             = "permanent_redirect"
	PaymentRequired               = "payment_required"
	MethodNotAllowed              = "method_not_allowed"
	NotAcceptable                 = "not_acceptable"
	ProxyAuthRequired             = "proxy_auth_required"
	RequestTimeout                = "request_timeout"
	RequestCancelled              = "request_cancelled"
	Gone                          = "gone"
	LengthRequired                = "length_required"
	PreconditionFailed            = "precondition_failed"
	ContentTooLarge               = "content_too_large"
	URITooLong                    = "uri_too_long"
	UnsupportedMediaType          = "unsupported_media_type"
	RangeNotSatisfiable           = "range_not_satisfiable"
	ExpectationFailed             = "expectation_failed"
	ImATeapot                     = "im_a_teapot"
	MisdirectedRequest            = "misdirected_request"
	UnprocessableContent          = "unprocessable_content"
	Locked                        = "locked"
	FailedDependency              = "failed_dependency"
	TooEarly                      = "too_early"
	UpgradeRequired               = "upgrade_required"
	PreconditionRequired          = "precondition_required"
	RequestHeaderFieldsTooLarge   = "request_header_fields_too_large"
	UnavailableForLegalReasons    = "unavailable_for_legal_reasons"
	NotImplemented                = "not_implemented"
	ServiceUnavailable            = "service_unavailable"
	GatewayTimeout                = "gateway_timeout"
	HTTPVersionNotSupported       = "http_version_not_supported"
	VariantAlsoNegotiates         = "variant_also_negotiates"
	InsufficientStorage           = "insufficient_storage"
	LoopDetected                  = "loop_detected"
	NotExtended                   = "not_extended"
	NetworkAuthenticationRequired = "network_auth_required"
	RequirePositiveInt            = "require_positive_int"
)

// -----------------------------------------------------------------------------
// Human readable error messages EN
// -----------------------------------------------------------------------------

var messagesEN = map[string]string{
	Internal:                      "Something went wrong on our side. Please try again.",
	Unauthorized:                  "You need to sign in to continue.",
	Forbidden:                     "You don't have permission to perform this action.",
	NotFound:                      "The requested resource was not found.",
	Conflict:                      "The resource is in a conflicting state.",
	BadRequest:                    "Your request could not be understood.",
	ValidationFailed:              "One or more fields failed validation.",
	TooManyRequests:               "Too many requests. Please slow down.",
	Required:                      "%s is required.",
	InvalidEmailFormat:            "%s must be a valid email address.",
	InvalidJSONFormat:             "Request body must be valid JSON.",
	InvalidStatus:                 "The provided status %s is invalid. Allowed values include active, deleted, suspended, rejected, draft",
	Invalid:                       "Invalid",
	InvalidDataType:               "The provided data type %s is invalid. Allowed values include string, int64, float64, decimal, time, datetime, bytes, uuid, map",
	BadGateway:                    "The server received an invalid response from an upstream service. Please try again later.",
	OK:                            "Success.",
	Created:                       "Resource created successfully.",
	Accepted:                      "Request accepted for processing.",
	NonAuthoritativeInfo:          "Non-authoritative information.",
	NoContent:                     "No content.",
	ResetContent:                  "Reset content.",
	PartialContent:                "Partial content delivered.",
	MultiStatus:                   "Multiple status responses.",
	AlreadyReported:               "Already reported.",
	IMUsed:                        "IM used.",
	MultipleChoices:               "Multiple choices available.",
	MovedPermanently:              "Resource moved permanently.",
	Found:                         "Resource found.",
	SeeOther:                      "See other resource.",
	NotModified:                   "Resource not modified.",
	UseProxy:                      "Use proxy.",
	Unused:                        "Unused status code.",
	TemporaryRedirect:             "Temporary redirect.",
	PermanentRedirect:             "Permanent redirect.",
	PaymentRequired:               "Payment required.",
	MethodNotAllowed:              "Method not allowed.",
	NotAcceptable:                 "Not acceptable.",
	ProxyAuthRequired:             "Proxy authentication required.",
	RequestTimeout:                "Request timeout.",
	RequestCancelled:              "Request cancelled.",
	Gone:                          "Resource gone.",
	LengthRequired:                "Content length required.",
	PreconditionFailed:            "Precondition failed.",
	ContentTooLarge:               "Content too large.",
	URITooLong:                    "URI too long.",
	UnsupportedMediaType:          "Unsupported media type.",
	RangeNotSatisfiable:           "Range not satisfiable.",
	ExpectationFailed:             "Expectation failed.",
	ImATeapot:                     "I'm a teapot.",
	MisdirectedRequest:            "Misdirected request.",
	UnprocessableContent:          "Unprocessable content.",
	Locked:                        "Resource locked.",
	FailedDependency:              "Failed dependency.",
	TooEarly:                      "Too early.",
	UpgradeRequired:               "Upgrade required.",
	PreconditionRequired:          "Precondition required.",
	RequestHeaderFieldsTooLarge:   "Request header fields too large.",
	UnavailableForLegalReasons:    "Unavailable for legal reasons.",
	NotImplemented:                "Not implemented.",
	ServiceUnavailable:            "Service unavailable.",
	GatewayTimeout:                "Gateway timeout.",
	HTTPVersionNotSupported:       "HTTP version not supported.",
	VariantAlsoNegotiates:         "Variant also negotiates.",
	InsufficientStorage:           "Insufficient storage.",
	LoopDetected:                  "Loop detected.",
	NotExtended:                   "Not extended.",
	NetworkAuthenticationRequired: "Network authentication required.",
	RequirePositiveInt:            "Integer must be positive.",
}

// -----------------------------------------------------------------------------
// Human readable error messages NB
// -----------------------------------------------------------------------------

var messagesNB = map[string]string{
	Internal:                      "Noe gikk galt hos oss. Prøv igjen.",
	Unauthorized:                  "Du må være innlogget for å fortsette.",
	Forbidden:                     "Du har ikke tilgang til denne handlingen.",
	NotFound:                      "Forespurt ressurs ble ikke funnet.",
	Conflict:                      "Ressursen er i konflikt.",
	BadRequest:                    "Forespørselen kunne ikke forstås.",
	ValidationFailed:              "Ett eller flere felt feilet validering.",
	TooManyRequests:               "For mange forespørsler. Vent litt.",
	Required:                      "%s er påkrevd.",
	InvalidEmailFormat:            "%s må være en gyldig e-postadresse.",
	InvalidJSONFormat:             "Kroppen i forespørselen må være gyldig JSON.",
	InvalidStatus:                 "Oppgitt statusverdi %s er ugyldig. Gyldige verdier inkluderer active, deleted, suspended, rejected, draft",
	Invalid:                       "Ugyldig",
	InvalidDataType:               "Oppgitt datatypeverdi %s er ugyldig. Gyldige verdier inkluderer string, int64, float64, decimal, time, datetime, bytes, uuid, map",
	BadGateway:                    "Serveren mottok et ugyldig svar fra en ekstern tjeneste. Prøv igjen senere.",
	OK:                            "Vellykket.",
	Created:                       "Ressurs opprettet.",
	Accepted:                      "Forespørsel akseptert for behandling.",
	NonAuthoritativeInfo:          "Ikke-autoritativ informasjon.",
	NoContent:                     "Ingen innhold.",
	ResetContent:                  "Tilbakestill innhold.",
	PartialContent:                "Delvis innhold levert.",
	MultiStatus:                   "Flere statusresponser.",
	AlreadyReported:               "Allerede rapportert.",
	IMUsed:                        "IM brukt.",
	MultipleChoices:               "Flere valg tilgjengelig.",
	MovedPermanently:              "Ressursen er permanent flyttet.",
	Found:                         "Ressurs funnet.",
	SeeOther:                      "Se annen ressurs.",
	NotModified:                   "Ressurs ikke endret.",
	UseProxy:                      "Bruk proxy.",
	Unused:                        "Ubrukt statuskode.",
	TemporaryRedirect:             "Midlertidig omdirigering.",
	PermanentRedirect:             "Permanent omdirigering.",
	PaymentRequired:               "Betaling kreves.",
	MethodNotAllowed:              "Metode ikke tillatt.",
	NotAcceptable:                 "Ikke akseptabelt.",
	ProxyAuthRequired:             "Proxy-autentisering kreves.",
	RequestTimeout:                "Forespørselen har tidsavbrudd.",
	RequestCancelled:              "Forespørselen ble avbrutt.",
	Gone:                          "Ressursen er fjernet.",
	LengthRequired:                "Content-Length kreves.",
	PreconditionFailed:            "Forutsetning feilet.",
	ContentTooLarge:               "Innholdet er for stort.",
	URITooLong:                    "URI er for lang.",
	UnsupportedMediaType:          "Medietype ikke støttet.",
	RangeNotSatisfiable:           "Område ikke tilfredsstillende.",
	ExpectationFailed:             "Forventning feilet.",
	ImATeapot:                     "Jeg er en tekanne.",
	MisdirectedRequest:            "Feilrettet forespørsel.",
	UnprocessableContent:          "Kan ikke behandle innhold.",
	Locked:                        "Ressursen er låst.",
	FailedDependency:              "Avhengighet feilet.",
	TooEarly:                      "For tidlig.",
	UpgradeRequired:               "Oppgradering kreves.",
	PreconditionRequired:          "Forutsetning kreves.",
	RequestHeaderFieldsTooLarge:   "Forespørselens header-felt er for store.",
	UnavailableForLegalReasons:    "Utilgjengelig av juridiske årsaker.",
	NotImplemented:                "Ikke implementert.",
	ServiceUnavailable:            "Tjenesten er utilgjengelig.",
	GatewayTimeout:                "Gateway tidsavbrudd.",
	HTTPVersionNotSupported:       "HTTP-versjon ikke støttet.",
	VariantAlsoNegotiates:         "Variant forhandler også.",
	InsufficientStorage:           "Utilstrekkelig lagringsplass.",
	LoopDetected:                  "Sløyfe oppdaget.",
	NotExtended:                   "Ikke utvidet.",
	NetworkAuthenticationRequired: "Nettverksautentisering kreves.",
	RequirePositiveInt:            "Heltallet må være positivt.",
}

// -----------------------------------------------------------------------------
// Message catalog
// -----------------------------------------------------------------------------

const (
	LocaleEN = "en"
	LocaleNB = "nb"
)

var catalogs = map[string]map[string]string{
	LocaleEN: messagesEN,
	LocaleNB: messagesNB,
}

// genericSubject stands in for a missing Sprintf argument, so a caller that
// omits the field name gets a sentence instead of a raw "%s".
var genericSubject = map[string]string{
	LocaleEN: "This field",
	LocaleNB: "Dette feltet",
}

// normalizeLocale maps a request locale onto a catalog key. Callers hand us
// whatever Accept-Language negotiation produced, which is a language tag and not
// a catalog key: "no", "nb", "nb-NO" and "nn" all mean the Norwegian catalog,
// "en-US" means the English one. Returns "" when nothing matches, which leaves
// the fallback decision to the caller.
//
// Getting this wrong is not a small mistake: a service defaulting to "no" used
// to miss the "nb" catalog entirely and answer in English, while
// CustomHumanMessageLocale accepted "no" and answered in Norwegian — so a single
// validation response mixed the two languages.
func normalizeLocale(locale string) string {
	l := strings.ToLower(strings.TrimSpace(locale))
	if i := strings.IndexAny(l, "-_"); i > 0 {
		l = l[:i]
	}
	switch l {
	case "no", "nb", "nn", "nob", "nno":
		return LocaleNB
	case "en", "eng":
		return LocaleEN
	default:
		return ""
	}
}

// countVerbs counts the fmt verbs in a template, ignoring escaped "%%".
func countVerbs(msg string) int {
	n := 0
	for i := 0; i < len(msg); i++ {
		if msg[i] != '%' {
			continue
		}
		if i+1 < len(msg) && msg[i+1] == '%' {
			i++
			continue
		}
		n++
	}
	return n
}

// fill applies args to a template without ever leaking formatting noise into a
// client-facing message. A template and its call sites drift — a translation
// loses a placeholder, a caller passes a field name to a message that has none —
// and the cost of that drift should be a slightly generic sentence, never
// "%!(EXTRA string=memory_mb)" or a literal "%s" in front of a user.
func fill(locale, msg string, args ...any) string {
	n := countVerbs(msg)
	if n == 0 {
		return msg
	}
	if len(args) > n {
		args = args[:n]
	}
	for len(args) < n {
		subject, ok := genericSubject[locale]
		if !ok {
			subject = genericSubject[LocaleEN]
		}
		args = append(args, subject)
	}
	return fmt.Sprintf(msg, args...)
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

// HumanMessageLocale lets you choose a locale. It accepts any language tag the
// Accept-Language negotiation produces ("en", "en-US", "no", "nb", "nb-NO",
// "nn"); anything it does not recognise falls back to English.
func HumanMessageLocale(locale, code string, args ...any) string {
	key := normalizeLocale(locale)
	cat, ok := catalogs[key]
	if !ok {
		key, cat = LocaleEN, messagesEN
	}
	msg, ok := cat[code]
	if !ok {
		// An untranslated code is still better answered in English than as a
		// generic internal error.
		if msg, ok = messagesEN[code]; !ok {
			msg = messagesEN[Internal]
		}
	}
	return fill(key, msg, args...)
}

type CustomMessage struct {
	En string
	No string
}

// CustomHumanMessageLocale picks the translation for a locale, accepting the same
// language tags as HumanMessageLocale so the two never disagree about what "nb"
// or "no" means. Anything unrecognised, or a missing translation, falls back to
// English.
func CustomHumanMessageLocale(locale string, c CustomMessage) string {
	if normalizeLocale(locale) == LocaleNB && c.No != "" {
		return c.No
	}
	return c.En
}

// (Optional) map codes to HTTP status for consistent responses.
var statusByCode = map[string]int{
	Internal:                      http.StatusInternalServerError,
	Unauthorized:                  http.StatusUnauthorized,
	Forbidden:                     http.StatusForbidden,
	NotFound:                      http.StatusNotFound,
	Conflict:                      http.StatusConflict,
	BadRequest:                    http.StatusBadRequest,
	ValidationFailed:              http.StatusUnprocessableEntity,
	TooManyRequests:               http.StatusTooManyRequests,
	Required:                      http.StatusBadRequest,
	InvalidEmailFormat:            http.StatusBadRequest,
	InvalidJSONFormat:             http.StatusBadRequest,
	InvalidStatus:                 http.StatusBadRequest,
	InvalidDataType:               http.StatusBadRequest,
	Invalid:                       http.StatusBadRequest,
	BadGateway:                    http.StatusBadGateway,
	OK:                            http.StatusOK,
	Created:                       http.StatusCreated,
	Accepted:                      http.StatusAccepted,
	NonAuthoritativeInfo:          http.StatusNonAuthoritativeInfo,
	NoContent:                     http.StatusNoContent,
	ResetContent:                  http.StatusResetContent,
	PartialContent:                http.StatusPartialContent,
	MultiStatus:                   http.StatusMultiStatus,
	AlreadyReported:               http.StatusAlreadyReported,
	IMUsed:                        http.StatusIMUsed,
	MultipleChoices:               http.StatusMultipleChoices,
	MovedPermanently:              http.StatusMovedPermanently,
	Found:                         http.StatusFound,
	SeeOther:                      http.StatusSeeOther,
	NotModified:                   http.StatusNotModified,
	UseProxy:                      http.StatusUseProxy,
	Unused:                        306,
	TemporaryRedirect:             http.StatusTemporaryRedirect,
	PermanentRedirect:             http.StatusPermanentRedirect,
	PaymentRequired:               http.StatusPaymentRequired,
	MethodNotAllowed:              http.StatusMethodNotAllowed,
	NotAcceptable:                 http.StatusNotAcceptable,
	ProxyAuthRequired:             http.StatusProxyAuthRequired,
	RequestTimeout:                http.StatusRequestTimeout,
	RequestCancelled:              499, // Client Closed Request (non-standard)
	Gone:                          http.StatusGone,
	LengthRequired:                http.StatusLengthRequired,
	PreconditionFailed:            http.StatusPreconditionFailed,
	ContentTooLarge:               http.StatusRequestEntityTooLarge,
	URITooLong:                    http.StatusRequestURITooLong,
	UnsupportedMediaType:          http.StatusUnsupportedMediaType,
	RangeNotSatisfiable:           http.StatusRequestedRangeNotSatisfiable,
	ExpectationFailed:             http.StatusExpectationFailed,
	ImATeapot:                     http.StatusTeapot,
	MisdirectedRequest:            http.StatusMisdirectedRequest,
	UnprocessableContent:          http.StatusUnprocessableEntity,
	Locked:                        http.StatusLocked,
	FailedDependency:              http.StatusFailedDependency,
	TooEarly:                      http.StatusTooEarly,
	UpgradeRequired:               http.StatusUpgradeRequired,
	PreconditionRequired:          http.StatusPreconditionRequired,
	RequestHeaderFieldsTooLarge:   http.StatusRequestHeaderFieldsTooLarge,
	UnavailableForLegalReasons:    http.StatusUnavailableForLegalReasons,
	NotImplemented:                http.StatusNotImplemented,
	ServiceUnavailable:            http.StatusServiceUnavailable,
	GatewayTimeout:                http.StatusGatewayTimeout,
	HTTPVersionNotSupported:       http.StatusHTTPVersionNotSupported,
	VariantAlsoNegotiates:         http.StatusVariantAlsoNegotiates,
	InsufficientStorage:           http.StatusInsufficientStorage,
	LoopDetected:                  http.StatusLoopDetected,
	NotExtended:                   http.StatusNotExtended,
	NetworkAuthenticationRequired: http.StatusNetworkAuthenticationRequired,
	RequirePositiveInt:            http.StatusBadRequest,
}

func StatusFor(code string) int {
	if s, ok := statusByCode[code]; ok {
		return s
	}
	return http.StatusInternalServerError
}
