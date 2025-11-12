package httperror_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/uuid"
	he "github.com/grasp-labs/ds-go-commonmodels/v3/commonmodels/http_error"
)

func TestHttpError_NotFound(t *testing.T) {
	requestID := uuid.MustParse("31ac4e2a-10a1-471d-ac7c-fd6ee13a526d").String()
	httpErr := he.NotFound(requestID, "")
	status := httpErr.Status()
	if status != http.StatusNotFound {
		t.Fatalf("expected %d, got %d", http.StatusNotFound, status)
	}
	jsonBytes, err := json.Marshal(httpErr)
	if err != nil {
		t.Fatalf("did not expect json marshal error, err: %v", err)
	}
	jsonStr := string(jsonBytes)
	expected := `{"code":"not_found","message":"The requested resource was not found.","request_id":"31ac4e2a-10a1-471d-ac7c-fd6ee13a526d","recoverable":false,"retry_after":0}`
	if jsonStr != expected {
		t.Fatalf("failed to marshal struct, got %s, expected %s", jsonStr, expected)
	}

	statusCode, headers, hErr := httpErr.Response()
	if statusCode != http.StatusNotFound {
		t.Fatalf("expected %d, got %d", http.StatusNotFound, statusCode)
	}
	if retryHeader := headers.Get("Retry-After"); retryHeader != "" {
		t.Fatalf("expected empty Retry-After header, got %s", retryHeader)
	}
	if recoverableHeader := headers.Get("Recoverable"); recoverableHeader != "" {
		t.Fatalf("expected empty Recoverable header, got %s", recoverableHeader)
	}

	expect := "not_found"
	if hErr.Code != expect {
		t.Fatalf("expected %s, got %s", expect, hErr.Code)
	}
	expect = "The requested resource was not found."
	if hErr.Message != expect {
		t.Fatalf("expected %s, got %s", expect, hErr.Code)
	}
	if _, err := uuid.Parse(hErr.RequestID); err != nil {
		t.Fatalf("failed to parse, err: %v", err)
	}

	jsonBytes, err = json.Marshal(hErr)
	if err != nil {
		t.Fatalf("did not expect json marshal error, err: %v", err)
	}
	jsonStr = string(jsonBytes)
	expected = `{"code":"not_found","message":"The requested resource was not found.","request_id":"31ac4e2a-10a1-471d-ac7c-fd6ee13a526d","recoverable":false,"retry_after":0}`
	if jsonStr != expected {
		t.Fatalf("failed to marshal struct, got %s, expected %s", jsonStr, expected)
	}
}

func TestBadRequest(t *testing.T) {
	requestID := uuid.MustParse("31ac4e2a-10a1-471d-ac7c-fd6ee13a526d").String()
	httpErr := he.BadRequest(requestID, "invalid payload")
	status := httpErr.Status()
	if status != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, status)
	}
	jsonBytes, err := json.Marshal(httpErr)
	if err != nil {
		t.Fatalf("did not expect json marshal error, err: %v", err)
	}
	jsonStr := string(jsonBytes)
	expected := `{"code":"bad_request","message":"invalid payload","request_id":"31ac4e2a-10a1-471d-ac7c-fd6ee13a526d","recoverable":false,"retry_after":0}`
	if jsonStr != expected {
		t.Fatalf("failed to marshal struct, got %s, expected %s", jsonStr, expected)
	}

	statusCode, headers, hErr := httpErr.Response()
	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, statusCode)
	}
	if retryHeader := headers.Get("Retry-After"); retryHeader != "" {
		t.Fatalf("expected empty Retry-After header, got %s", retryHeader)
	}
	if recoverableHeader := headers.Get("Recoverable"); recoverableHeader != "" {
		t.Fatalf("expected empty Recoverable header, got %s", recoverableHeader)
	}

	expect := "bad_request"
	if hErr.Code != expect {
		t.Fatalf("expected %s, got %s", expect, hErr.Code)
	}
	expect = "invalid payload"
	if hErr.Message != expect {
		t.Fatalf("expected %s, got %s", expect, hErr.Message)
	}
	if _, err := uuid.Parse(hErr.RequestID); err != nil {
		t.Fatalf("failed to parse, err: %v", err)
	}

	jsonBytes, err = json.Marshal(hErr)
	if err != nil {
		t.Fatalf("did not expect json marshal error, err: %v", err)
	}
	jsonStr = string(jsonBytes)
	expected = `{"code":"bad_request","message":"invalid payload","request_id":"31ac4e2a-10a1-471d-ac7c-fd6ee13a526d","recoverable":false,"retry_after":0}`
	if jsonStr != expected {
		t.Fatalf("failed to marshal struct, got %s, expected %s", jsonStr, expected)
	}
	recoverable := hErr.Recoverable
	if recoverable {
		t.Fatalf("expected recoverable to be false, got true")
	}
	retryAfter := hErr.RetryAfter
	if retryAfter != 0 {
		t.Fatalf("expected retryAfter to be 0, got %d", retryAfter)
	}
}

func TestBadRequest_CustomTime(t *testing.T) {
	requestID := uuid.MustParse("31ac4e2a-10a1-471d-ac7c-fd6ee13a526d").String()
	httpErr := he.BadRequest(requestID, "invalid payload").WithRetry(120)
	status := httpErr.Status()
	if status != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, status)
	}
	jsonBytes, err := json.Marshal(httpErr)
	if err != nil {
		t.Fatalf("did not expect json marshal error, err: %v", err)
	}
	jsonStr := string(jsonBytes)
	expected := `{"code":"bad_request","message":"invalid payload","request_id":"31ac4e2a-10a1-471d-ac7c-fd6ee13a526d","recoverable":true,"retry_after":120}`
	if jsonStr != expected {
		t.Fatalf("failed to marshal struct, got %s, expected %s", jsonStr, expected)
	}

	statusCode, headers, hErr := httpErr.Response()
	if statusCode != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, statusCode)
	}
	if retryHeader := headers.Get("Retry-After"); retryHeader != "120" {
		t.Fatalf("expected Retry-After header to be 120, got %s", retryHeader)
	}

	expect := "bad_request"
	if hErr.Code != expect {
		t.Fatalf("expected %s, got %s", expect, hErr.Code)
	}
	expect = "invalid payload"
	if hErr.Message != expect {
		t.Fatalf("expected %s, got %s", expect, hErr.Message)
	}
	if _, err := uuid.Parse(hErr.RequestID); err != nil {
		t.Fatalf("failed to parse, err: %v", err)
	}

	jsonBytes, err = json.Marshal(hErr)
	if err != nil {
		t.Fatalf("did not expect json marshal error, err: %v", err)
	}
	jsonStr = string(jsonBytes)
	expected = `{"code":"bad_request","message":"invalid payload","request_id":"31ac4e2a-10a1-471d-ac7c-fd6ee13a526d","recoverable":true,"retry_after":120}`
	if jsonStr != expected {
		t.Fatalf("failed to marshal struct, got %s, expected %s", jsonStr, expected)
	}
	recoverable := hErr.Recoverable
	if !recoverable {
		t.Fatalf("expected recoverable to be true, got false")
	}
	retryAfter := hErr.RetryAfter
	if retryAfter != 120 {
		t.Fatalf("expected retryAfter to be 120, got %d", retryAfter)
	}
}

func TestTooManyRequests(t *testing.T) {
	requestID := uuid.MustParse("31ac4e2a-10a1-471d-ac7c-fd6ee13a526d").String()
	httpErr := he.TooManyRequests(requestID, "")
	status := httpErr.Status()
	if status != http.StatusTooManyRequests {
		t.Fatalf("expected %d, got %d", http.StatusTooManyRequests, status)
	}
	jsonBytes, err := json.Marshal(httpErr)
	if err != nil {
		t.Fatalf("did not expect json marshal error, err: %v", err)
	}
	jsonStr := string(jsonBytes)
	expected := `{"code":"too_many_requests","message":"Too many requests. Please slow down.","request_id":"31ac4e2a-10a1-471d-ac7c-fd6ee13a526d","recoverable":true,"retry_after":60}`
	if jsonStr != expected {
		t.Fatalf("failed to marshal struct, got %s, expected %s", jsonStr, expected)
	}

	statusCode, headers, hErr := httpErr.Response()
	if statusCode != http.StatusTooManyRequests {
		t.Fatalf("expected %d, got %d", http.StatusTooManyRequests, statusCode)
	}
	if retryHeader := headers.Get("Retry-After"); retryHeader != "60" {
		t.Fatalf("expected Retry-After header to be 60, got %s", retryHeader)
	}
	expect := "too_many_requests"
	if hErr.Code != expect {
		t.Fatalf("expected %s, got %s", expect, hErr.Code)
	}
	expect = "Too many requests. Please slow down."
	if hErr.Message != expect {
		t.Fatalf("expected %s, got %s", expect, hErr.Message)
	}
	if _, err := uuid.Parse(hErr.RequestID); err != nil {
		t.Fatalf("failed to parse, err: %v", err)
	}

	jsonBytes, err = json.Marshal(hErr)
	if err != nil {
		t.Fatalf("did not expect json marshal error, err: %v", err)
	}
	jsonStr = string(jsonBytes)
	expected = `{"code":"too_many_requests","message":"Too many requests. Please slow down.","request_id":"31ac4e2a-10a1-471d-ac7c-fd6ee13a526d","recoverable":true,"retry_after":60}`
	if jsonStr != expected {
		t.Fatalf("failed to marshal struct, got %s, expected %s", jsonStr, expected)
	}
	recoverable := hErr.Recoverable
	if !recoverable {
		t.Fatalf("expected recoverable to be true, got false")
	}
	retryAfter := hErr.RetryAfter
	if retryAfter != 60 {
		t.Fatalf("expected retryAfter to be 60, got %d", retryAfter)
	}
}

func TestWithRetry(t *testing.T) {
	requestID := uuid.MustParse("31ac4e2a-10a1-471d-ac7c-fd6ee13a526d").String()
	httpErr := he.TooManyRequests(requestID, "").WithRetry(90)
	status := httpErr.Status()
	if status != http.StatusTooManyRequests {
		t.Fatalf("expected %d, got %d", http.StatusTooManyRequests, status)
	}
	jsonBytes, err := json.Marshal(httpErr)
	if err != nil {
		t.Fatalf("did not expect json marshal error, err: %v", err)
	}
	jsonStr := string(jsonBytes)
	expected := `{"code":"too_many_requests","message":"Too many requests. Please slow down.","request_id":"31ac4e2a-10a1-471d-ac7c-fd6ee13a526d","recoverable":true,"retry_after":90}`
	if jsonStr != expected {
		t.Fatalf("failed to marshal struct, got %s, expected %s", jsonStr, expected)
	}

	statusCode, headers, hErr := httpErr.Response()
	if statusCode != http.StatusTooManyRequests {
		t.Fatalf("expected %d, got %d", http.StatusTooManyRequests, statusCode)
	}
	if retryHeader := headers.Get("Retry-After"); retryHeader != "90" {
		t.Fatalf("expected Retry-After header to be 90, got %s", retryHeader)
	}
	expect := "too_many_requests"
	if hErr.Code != expect {
		t.Fatalf("expected %s, got %s", expect, hErr.Code)
	}
	expect = "Too many requests. Please slow down."
	if hErr.Message != expect {
		t.Fatalf("expected %s, got %s", expect, hErr.Message)
	}
	if _, err := uuid.Parse(hErr.RequestID); err != nil {
		t.Fatalf("failed to parse, err: %v", err)
	}

	jsonBytes, err = json.Marshal(hErr)
	if err != nil {
		t.Fatalf("did not expect json marshal error, err: %v", err)
	}
	jsonStr = string(jsonBytes)
	expected = `{"code":"too_many_requests","message":"Too many requests. Please slow down.","request_id":"31ac4e2a-10a1-471d-ac7c-fd6ee13a526d","recoverable":true,"retry_after":90}`
	if jsonStr != expected {
		t.Fatalf("failed to marshal struct, got %s, expected %s", jsonStr, expected)
	}
	recoverable := hErr.Recoverable
	if !recoverable {
		t.Fatalf("expected recoverable to be true, got false")
	}
	retryAfter := hErr.RetryAfter
	if retryAfter != 90 {
		t.Fatalf("expected retryAfter to be 90, got %d", retryAfter)
	}
}
func TestBadGateway_DefaultMessage(t *testing.T) {
	requestID := uuid.MustParse("31ac4e2a-10a1-471d-ac7c-fd6ee13a526d").String()
	httpErr := he.BadGateway(requestID, "")
	status := httpErr.Status()
	if status != http.StatusBadGateway {
		t.Fatalf("expected %d, got %d", http.StatusBadGateway, status)
	}
	jsonBytes, err := json.Marshal(httpErr)
	if err != nil {
		t.Fatalf("did not expect json marshal error, err: %v", err)
	}
	expected := `{"code":"bad_gateway","message":"The server received an invalid response from an upstream service. Please try again later.","request_id":"31ac4e2a-10a1-471d-ac7c-fd6ee13a526d","recoverable":false,"retry_after":0}`
	if string(jsonBytes) != expected {
		t.Fatalf("failed to marshal struct, got %s, expected %s", string(jsonBytes), expected)
	}
	statusCode, headers, hErr := httpErr.Response()
	if statusCode != http.StatusBadGateway {
		t.Fatalf("expected %d, got %d", http.StatusBadGateway, statusCode)
	}
	if retryHeader := headers.Get("Retry-After"); retryHeader != "" {
		t.Fatalf("expected empty Retry-After header, got %s", retryHeader)
	}
	if hErr.Code != "bad_gateway" {
		t.Fatalf("expected bad_gateway, got %s", hErr.Code)
	}
	if hErr.Message != "The server received an invalid response from an upstream service. Please try again later." {
		t.Fatalf("expected The server received an invalid response from an upstream service. Please try again later., got %s", hErr.Message)
	}
	if _, err := uuid.Parse(hErr.RequestID); err != nil {
		t.Fatalf("failed to parse, err: %v", err)
	}
}

func TestBadGateway_CustomMessage(t *testing.T) {
	requestID := uuid.MustParse("31ac4e2a-10a1-471d-ac7c-fd6ee13a526d").String()
	customMsg := "custom bad gateway error"
	httpErr := he.BadGateway(requestID, customMsg)
	status := httpErr.Status()
	if status != http.StatusBadGateway {
		t.Fatalf("expected %d, got %d", http.StatusBadGateway, status)
	}
	jsonBytes, err := json.Marshal(httpErr)
	if err != nil {
		t.Fatalf("did not expect json marshal error, err: %v", err)
	}
	expected := `{"code":"bad_gateway","message":"custom bad gateway error","request_id":"31ac4e2a-10a1-471d-ac7c-fd6ee13a526d","recoverable":false,"retry_after":0}`
	if string(jsonBytes) != expected {
		t.Fatalf("failed to marshal struct, got %s, expected %s", string(jsonBytes), expected)
	}
	statusCode, headers, hErr := httpErr.Response()
	if statusCode != http.StatusBadGateway {
		t.Fatalf("expected %d, got %d", http.StatusBadGateway, statusCode)
	}
	if retryHeader := headers.Get("Retry-After"); retryHeader != "" {
		t.Fatalf("expected empty Retry-After header, got %s", retryHeader)
	}
	if hErr.Code != "bad_gateway" {
		t.Fatalf("expected bad_gateway, got %s", hErr.Code)
	}
	if hErr.Message != customMsg {
		t.Fatalf("expected %s, got %s", customMsg, hErr.Message)
	}
	if _, err := uuid.Parse(hErr.RequestID); err != nil {
		t.Fatalf("failed to parse, err: %v", err)
	}
}
func TestGetLocale(t *testing.T) {
	tests := []struct {
		name   string
		input  []string
		expect string
	}{
		{"no locale", []string{}, "en"},
		{"empty locale", []string{""}, "en"},
		{"single locale", []string{"fr"}, "fr"},
		{"multiple locales", []string{"es", "de"}, "es"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := he.GetLocale(tt.input...)
			if got != tt.expect {
				t.Fatalf("expected %s, got %s", tt.expect, got)
			}
		})
	}
}

func TestInternal_DefaultMessage(t *testing.T) {
	requestID := uuid.New().String()
	httpErr := he.Internal(requestID, "")
	if httpErr.Code != "internal_error" {
		t.Fatalf("expected code internal_error, got %s", httpErr.Code)
	}
	if httpErr.Message == "" {
		t.Fatalf("expected non-empty message")
	}
	if httpErr.RequestID != requestID {
		t.Fatalf("expected requestID %s, got %s", requestID, httpErr.RequestID)
	}
	if httpErr.Status() != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, httpErr.Status())
	}
}

func TestInternal_CustomMessage(t *testing.T) {
	requestID := uuid.New().String()
	msg := "custom internal error"
	httpErr := he.Internal(requestID, msg)
	if httpErr.Message != msg {
		t.Fatalf("expected message %s, got %s", msg, httpErr.Message)
	}
}

func TestPartialContent_DefaultMessage(t *testing.T) {
	requestID := uuid.New().String()
	httpErr := he.PartialContent(requestID, "")
	if httpErr.Code != "partial_content" {
		t.Fatalf("expected code partial_content, got %s", httpErr.Code)
	}
	if httpErr.Status() != http.StatusPartialContent {
		t.Fatalf("expected status %d, got %d", http.StatusPartialContent, httpErr.Status())
	}
}

func TestUnauthorized_DefaultMessage(t *testing.T) {
	requestID := uuid.New().String()
	httpErr := he.Unauthorized(requestID, "")
	if httpErr.Code != "unauthorized" {
		t.Fatalf("expected code unauthorized, got %s", httpErr.Code)
	}
	if httpErr.Status() != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, httpErr.Status())
	}
}

func TestForbidden_DefaultMessage(t *testing.T) {
	requestID := uuid.New().String()
	httpErr := he.Forbidden(requestID, "")
	if httpErr.Code != "forbidden" {
		t.Fatalf("expected code forbidden, got %s", httpErr.Code)
	}
	if httpErr.Status() != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, httpErr.Status())
	}
}

func TestConflict_DefaultMessage(t *testing.T) {
	requestID := uuid.New().String()
	httpErr := he.Conflict(requestID, "")
	if httpErr.Code != "conflict" {
		t.Fatalf("expected code conflict, got %s", httpErr.Code)
	}
	if httpErr.Status() != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, httpErr.Status())
	}
}

func TestServiceUnavailable_DefaultMessage(t *testing.T) {
	requestID := uuid.New().String()
	httpErr := he.ServiceUnavailable(requestID, "")
	if httpErr.Code != "service_unavailable" {
		t.Fatalf("expected code service_unavailable, got %s", httpErr.Code)
	}
	if httpErr.Status() != http.StatusServiceUnavailable {
		t.Fatalf("expected status %d, got %d", http.StatusServiceUnavailable, httpErr.Status())
	}
	if !httpErr.Recoverable {
		t.Fatalf("expected recoverable true")
	}
	if httpErr.RetryAfter != 30 {
		t.Fatalf("expected retry_after 30, got %d", httpErr.RetryAfter)
	}
}

func TestCreated_DefaultMessage(t *testing.T) {
	requestID := uuid.New().String()
	httpErr := he.Created(requestID, "")
	if httpErr.Code != "created" {
		t.Fatalf("expected code created, got %s", httpErr.Code)
	}
	if httpErr.Status() != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, httpErr.Status())
	}
}

func TestLocaleOverride(t *testing.T) {
	requestID := uuid.New().String()
	httpErr := he.NotFound(requestID, "", "nb")
	if httpErr.Code != "not_found" {
		t.Fatalf("expected code not_found, got %s", httpErr.Code)
	}
	expected := "Forespurt ressurs ble ikke funnet."
	if httpErr.Message != expected {
		t.Fatalf("expected message %q, got %q", expected, httpErr.Message)
	}
}
