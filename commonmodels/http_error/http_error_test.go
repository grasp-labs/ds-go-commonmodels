package httperror_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/uuid"
	he "github.com/grasp-labs/ds-go-commonmodels/v2/commonmodels/http_error"
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
	expected := `{"code":"not_found","message":"not found","request_id":"31ac4e2a-10a1-471d-ac7c-fd6ee13a526d","recoverable":false,"retry_after":0}`
	if jsonStr != expected {
		t.Fatalf("failed to marshal struct, got %s, expected %s", jsonStr, expected)
	}

	statusCode, headers, hErr := httpErr.Response()
	if statusCode != http.StatusNotFound {
		t.Fatalf("expected %d, got %d", http.StatusNotFound, status)
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
	expect = "not found"
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
	expected = `{"code":"not_found","message":"not found","request_id":"31ac4e2a-10a1-471d-ac7c-fd6ee13a526d","recoverable":false,"retry_after":0}`
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
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, status)
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
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, status)
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
	expected := `{"code":"too_many_requests","message":"too many requests","request_id":"31ac4e2a-10a1-471d-ac7c-fd6ee13a526d","recoverable":true,"retry_after":60}`
	if jsonStr != expected {
		t.Fatalf("failed to marshal struct, got %s, expected %s", jsonStr, expected)
	}

	statusCode, headers, hErr := httpErr.Response()
	if statusCode != http.StatusTooManyRequests {
		t.Fatalf("expected %d, got %d", http.StatusTooManyRequests, status)
	}
	if retryHeader := headers.Get("Retry-After"); retryHeader != "60" {
		t.Fatalf("expected Retry-After header to be 60, got %s", retryHeader)
	}
	expect := "too_many_requests"
	if hErr.Code != expect {
		t.Fatalf("expected %s, got %s", expect, hErr.Code)
	}
	expect = "too many requests"
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
	expected = `{"code":"too_many_requests","message":"too many requests","request_id":"31ac4e2a-10a1-471d-ac7c-fd6ee13a526d","recoverable":true,"retry_after":60}`
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
	expected := `{"code":"too_many_requests","message":"too many requests","request_id":"31ac4e2a-10a1-471d-ac7c-fd6ee13a526d","recoverable":true,"retry_after":90}`
	if jsonStr != expected {
		t.Fatalf("failed to marshal struct, got %s, expected %s", jsonStr, expected)
	}

	statusCode, headers, hErr := httpErr.Response()
	if statusCode != http.StatusTooManyRequests {
		t.Fatalf("expected %d, got %d", http.StatusTooManyRequests, status)
	}
	if retryHeader := headers.Get("Retry-After"); retryHeader != "90" {
		t.Fatalf("expected Retry-After header to be 90, got %s", retryHeader)
	}
	expect := "too_many_requests"
	if hErr.Code != expect {
		t.Fatalf("expected %s, got %s", expect, hErr.Code)
	}
	expect = "too many requests"
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
	expected = `{"code":"too_many_requests","message":"too many requests","request_id":"31ac4e2a-10a1-471d-ac7c-fd6ee13a526d","recoverable":true,"retry_after":90}`
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
