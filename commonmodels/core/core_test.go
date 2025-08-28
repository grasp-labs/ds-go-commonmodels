package core_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/grasp-labs/ds-go-commonmodels/commonmodels/core"
	"github.com/grasp-labs/ds-go-commonmodels/commonmodels/types"
	"github.com/stretchr/testify/assert"
)

func TestCoreModel_Validate_empty(t *testing.T) {
	coreModel := core.CoreModel{}
	validationErrors := coreModel.Validate()
	if len(validationErrors) < 1 {
		t.Fatalf("Expected more than 1 validation error.")
	}

}

func TestCoreModel_Validate_ID_is_required_and_uuid(t *testing.T) {
	jsonStr := `
	{
	  "id": null
    }
	`

	var coreModel core.CoreModel
	err := json.Unmarshal([]byte(jsonStr), &coreModel)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	validationErrors := coreModel.Validate()
	exist := false
	var message string
	for _, e := range validationErrors {
		if e.Field == "id" {
			exist = true
			message = e.Message
		}
	}
	assert.True(t, exist)
	assert.Equal(t, message, "required")

	coreModel.ID = uuid.New()
}

func TestCoreModel_Validate_TenantID_is_required_and_uuid(t *testing.T) {
	jsonStr := `
	{
	  "tenant_id": null
    }
	`

	var coreModel core.CoreModel
	err := json.Unmarshal([]byte(jsonStr), &coreModel)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	validationErrors := coreModel.Validate()
	exist := false
	var message string
	for _, e := range validationErrors {
		if e.Field == "tenant_id" {
			exist = true
			message = e.Message
		}
	}
	assert.True(t, exist)
	assert.Equal(t, message, "required")

	coreModel.TenantID = uuid.New()
}

func TestCoreModel_Validate_Name_is_requiredand_string(t *testing.T) {
	jsonStr := `
	{
	  "name": null
    }
	`

	var coreModel core.CoreModel
	err := json.Unmarshal([]byte(jsonStr), &coreModel)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	validationErrors := coreModel.Validate()
	exist := false
	var message string
	for _, e := range validationErrors {
		if e.Field == "name" {
			exist = true
			message = e.Message
		}
	}
	assert.True(t, exist)
	assert.Equal(t, message, "required")

	coreModel.Name = "test"
}

func TestCoreModel_Validate_CreatedBy_is_requiredand_string(t *testing.T) {
	jsonStr := `
	{
	  "created_by": null
    }
	`

	var coreModel core.CoreModel
	err := json.Unmarshal([]byte(jsonStr), &coreModel)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	validationErrors := coreModel.Validate()
	exist := false
	var message string
	for _, e := range validationErrors {
		if e.Field == "created_by" {
			exist = true
			message = e.Message
		}
	}
	assert.True(t, exist)
	allowed := map[string]struct{}{
		"required": {}, "not a valid email format": {},
	}
	_, ok := allowed[message]

	assert.True(t, ok)

	coreModel.CreatedBy = "user@domain.com"
}

func TestCoreModel_Validate_ModifiedBy_is_required_string(t *testing.T) {
	jsonStr := `
	{
	  "modified_by": null
    }
	`

	var coreModel core.CoreModel
	err := json.Unmarshal([]byte(jsonStr), &coreModel)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	validationErrors := coreModel.Validate()
	exist := false
	var message string
	for _, e := range validationErrors {
		if e.Field == "modified_by" {
			exist = true
			message = e.Message
		}
	}
	assert.True(t, exist)
	allowed := map[string]struct{}{
		"required": {}, "not a valid email format": {},
	}
	_, ok := allowed[message]

	assert.True(t, ok)

	coreModel.ModifiedBy = "user@domain.com"
}

func TestCoreModel_Validate_CreatedAt_is_required_timestamp(t *testing.T) {
	jsonStr := `
	{
	  "created_at": null
    }
	`

	var coreModel core.CoreModel
	err := json.Unmarshal([]byte(jsonStr), &coreModel)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	validationErrors := coreModel.Validate()
	exist := false
	var message string
	for _, e := range validationErrors {
		if e.Field == "created_at" {
			exist = true
			message = e.Message
		}
	}
	assert.True(t, exist)
	assert.Equal(t, message, "required")

	coreModel.CreatedAt = time.Now().UTC()
}

func TestCoreModel_Validate_ModifiedAt_is_required_timestamp(t *testing.T) {
	jsonStr := `
	{
	  "modified_at": null
    }
	`

	var coreModel core.CoreModel
	err := json.Unmarshal([]byte(jsonStr), &coreModel)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	validationErrors := coreModel.Validate()
	exist := false
	var message string
	for _, e := range validationErrors {
		if e.Field == "modified_at" {
			exist = true
			message = e.Message
		}
	}
	assert.True(t, exist)
	assert.Equal(t, message, "required")

	coreModel.CreatedAt = time.Now().UTC()
}

func TestManualLifeCycle_create(t *testing.T) {
	subject := "user@domain.com"
	issuer := "grasp-labs"
	tenantId := uuid.MustParse("25948ccc-cf16-491e-9cd4-44d5ebb7bc54")

	meta := map[string]string{"owner_id": "xyz123", "retention": "365"}

	tags := map[string]string{
		"env":       "prod",
		"tenant_id": tenantId.String(),
	}

	c := core.CoreModel{
		Name:     "Webhook config for X",
		Metadata: types.JSONB[map[string]string]{Data: meta},
		Tags:     types.JSONB[map[string]string]{Data: tags},
		Status:   "active",
	}

	// Apply create-time defaults & audit
	c.Create(subject, issuer, tenantId)
	validationErrors := c.Validate()
	if len(validationErrors) > 0 {
		t.Fatalf("Expected 0 validationErrors, got %v", validationErrors)
	}
}
