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
	expected := `{"code":"not_found","message":"not found","request_id":"31ac4e2a-10a1-471d-ac7c-fd6ee13a526d"}`
	if jsonStr != expected {
		t.Fatalf("failed to marshal struct, got %s, expected %s", jsonStr, expected)
	}

	statusCode, hErr := httpErr.Response()
	if statusCode != http.StatusNotFound {
		t.Fatalf("expected %d, got %d", http.StatusNotFound, status)
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
	expected = `{"code":"not_found","message":"not found","request_id":"31ac4e2a-10a1-471d-ac7c-fd6ee13a526d"}`
	if jsonStr != expected {
		t.Fatalf("failed to marshal struct, got %s, expected %s", jsonStr, expected)
	}
}
