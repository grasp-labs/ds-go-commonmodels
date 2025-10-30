package usage_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	models "github.com/grasp-labs/ds-go-commonmodels/v2/commonmodels/usage"
)

func TestUsageEntry_JSON_MarshalBasic(t *testing.T) {
	id := uuid.New()
	tenantID := uuid.New()
	productID := uuid.New()
	ownerID := "fancy-owner-id"

	entry := models.UsageEntry{
		ID:             id,
		TenantID:       tenantID,
		OwnerID:        &ownerID,
		ProductID:      productID,
		MemoryMB:       128,
		StartTimestamp: time.Date(2025, 7, 9, 14, 30, 0, 0, time.UTC),
		EndTimestamp:   time.Date(2025, 7, 9, 14, 32, 0, 0, time.UTC),
		Duration:       120,
		Status:         "active",
		Metadata: []map[string]string{
			{"key1": "value1"},
			{"key2": "value2"},
		},
		Tags: map[string]string{
			"env": "production",
			"app": "webserver",
		},
		CreatedAt: time.Date(2025, 7, 9, 14, 29, 0, 0, time.UTC),
		CreatedBy: "super admin 42",
	}

	data, err := json.Marshal(entry)
	assert.NoError(t, err, "JSON Marshalling should succeed")

	jsonStr := string(data)
	assert.Contains(t, jsonStr, `"id"`)
	assert.Contains(t, jsonStr, `"tenant_id"`)
	assert.Contains(t, jsonStr, `"owner_id"`)
	assert.Contains(t, jsonStr, `"product_id"`)
	assert.Contains(t, jsonStr, `"memory_mb"`)
	assert.Contains(t, jsonStr, `"start_timestamp"`)
	assert.Contains(t, jsonStr, `"end_timestamp"`)
	assert.Contains(t, jsonStr, `"duration"`)
	assert.Contains(t, jsonStr, `"status"`)
	assert.Contains(t, jsonStr, `"metadata"`)
	assert.Contains(t, jsonStr, `"tags"`)
	assert.Contains(t, jsonStr, `"created_at"`)
	assert.Contains(t, jsonStr, `"created_by"`)
}

func TestFieldValidation(t *testing.T) {
	entry := models.UsageEntry{}

	validationErrors := entry.Validate("en")
	assert.Len(t, validationErrors, 11, "Expected 11 validation errors for missing required fields")

	expectedFields := []string{
		"id",
		"tenant_id",
		"product_id",
		"memory_mb",
		"start_timestamp",
		"end_timestamp",
		"duration",
		"status",
		"tags",
		"created_at",
		"created_by",
	}

	for _, field := range expectedFields {
		found := false
		for _, ve := range validationErrors {
			if ve.Field == field {
				found = true
				break
			}
		}
		assert.True(t, found, "Expected validation error for field: %s", field)
	}
}

func TestUsageEntry_Validate_EndBeforeStart(t *testing.T) {
	entry := models.UsageEntry{
		ID:             uuid.New(),
		TenantID:       uuid.New(),
		ProductID:      uuid.New(),
		MemoryMB:       128,
		StartTimestamp: time.Date(2025, 7, 9, 14, 32, 0, 0, time.UTC),
		EndTimestamp:   time.Date(2025, 7, 9, 14, 30, 0, 0, time.UTC), // Before start!
		Duration:       120.0,
		Status:         "active",
		Tags:           map[string]string{"test": "value"},
		CreatedAt:      time.Now(),
		CreatedBy:      "system",
	}

	errors := entry.Validate("en")
	assert.NotEmpty(t, errors)

	found := false
	for _, err := range errors {
		if err.Field == "end_timestamp" {
			found = true
			break
		}
	}
	assert.True(t, found, "Should have validation error for end_timestamp")
}

func TestUsageEntry_JSON_WithNullOwnerID(t *testing.T) {
	entry := models.UsageEntry{
		ID:             uuid.New(),
		TenantID:       uuid.New(),
		OwnerID:        nil,
		ProductID:      uuid.New(),
		MemoryMB:       128,
		StartTimestamp: time.Now().UTC(),
		EndTimestamp:   time.Now().UTC().Add(time.Minute),
		Duration:       60.0,
		Status:         "active",
		Tags:           map[string]string{"test": "value"},
		CreatedAt:      time.Now().UTC(),
		CreatedBy:      "system",
	}

	data, err := json.Marshal(entry)
	assert.NoError(t, err)

	var unmarshaled models.UsageEntry
	err = json.Unmarshal(data, &unmarshaled)
	assert.NoError(t, err)
	assert.Nil(t, unmarshaled.OwnerID)
}

func TestUsageEntry_Validate_NegativeValues(t *testing.T) {
	entry := models.UsageEntry{
		ID:             uuid.New(),
		TenantID:       uuid.New(),
		ProductID:      uuid.New(),
		MemoryMB:       -128,
		StartTimestamp: time.Now(),
		EndTimestamp:   time.Now().Add(time.Minute),
		Duration:       -60.0,
		Status:         "active",
		Tags:           map[string]string{"test": "value"},
		CreatedAt:      time.Now(),
		CreatedBy:      "system",
	}

	errors := entry.Validate("en")

	// Should have errors for both memory_mb and duration
	memoryError := false
	durationError := false

	for _, err := range errors {
		if err.Field == "memory_mb" {
			memoryError = true
		}
		if err.Field == "duration" {
			durationError = true
		}
	}

	assert.True(t, memoryError, "Should have validation error for negative memory_mb")
	assert.True(t, durationError, "Should have validation error for negative duration")
}
