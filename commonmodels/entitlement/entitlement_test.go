package entitlement_test

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	models "github.com/grasp-labs/ds-go-commonmodels/v2/commonmodels/entitlement"
)

func TestEntitlement_JSON_MarshalBasic(t *testing.T) {
	id := uuid.New().String()
	tenantID := uuid.New().String()

	entitlement := models.Entitlement{
		ID:       id,
		Name:     "Premium Access",
		TenantId: tenantID,
	}

	data, err := json.Marshal(entitlement)
	assert.NoError(t, err, "JSON Marshalling should succeed")

	jsonStr := string(data)
	assert.Contains(t, jsonStr, `"id"`)
	assert.Contains(t, jsonStr, `"name"`)
	assert.Contains(t, jsonStr, `"tenant_id"`)
}

func TestEntitlement_Validate(t *testing.T) {
	entitlement := models.Entitlement{
		ID:       "",
		Name:     "",
		TenantId: "",
	}

	validationErrors := entitlement.Validate()
	assert.Len(t, validationErrors, 3, "Should have 3 validation errors for missing fields")

	expectedFields := map[string]bool{
		"id":        false,
		"name":      false,
		"tenant_id": false,
	}

	for _, ve := range validationErrors {
		if _, exists := expectedFields[ve.Field]; exists {
			expectedFields[ve.Field] = true
		}
	}

	for field, found := range expectedFields {
		assert.True(t, found, "Expected validation error for field: %s", field)
	}
}
