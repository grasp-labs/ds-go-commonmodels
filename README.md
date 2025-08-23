# ds-go-commonmodels — Common Models for the DS Platform in Go

![Build](https://github.com/grasp-labs/ds-go-commonmodels/actions/workflows/ci.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/grasp-labs/ds-go-commonmodels)](https://goreportcard.com/report/github.com/grasp-labs/ds-go-commonmodels)
[![codecov](https://codecov.io/gh/grasp-labs/ds-go-commonmodels/branch/main/graph/badge.svg)](https://codecov.io/gh/grasp-labs/ds-go-commonmodels)
[![Latest tag](https://img.shields.io/github/v/tag/grasp-labs/ds-go-commonmodels?sort=semver)](https://github.com/grasp-labs/ds-go-commonmodels/tags)
![License](https://img.shields.io/github/license/grasp-labs/ds-go-commonmodels?cacheSeconds=60)

A small, focused set of Go models and helpers intended for API services on Grasp Data platform (GraspDP).

## Features

### Core - our base attributes

The Core model define shared attributes (e.g. `ID`, `Name`) and patterns for flexible fields like `Metadata` and `Tags`. Model stricktly enforce types, managed required fields and content, and provide helpers for defaults and updates.

- `Create`: Applies create-time defaults and audit fields to Core model
- `Touch`: Updates modification audit fields.

Gorm
Model has Gorm support and implement the following Gorm hooks:

- `BeforeCreate`: Applied safe defaults and validate/normalize metadata.
- `BeforeUpdate`: Refresh modification audit fields.

## Install

Latest

```bash
go get github.com/grasp-labs/ds-go-commonmodels@latest
```

Tag

```bash
go get github.com/grasp-labs/ds-go-commonmodels@v1.1.0-rc.1
```

## Examples

### JSON shape

```json
{
  "id": "c1c9f1e8-2a7a-4f7f-93f1-7e8b9e4c1234",
  "name": "webhook",
  "tags": ["erp", "25948ccc-cf16-491e-9cd4-44d5ebb7bc54"],
  "metadata": {
    "owner_id": "xyz123",
    "retention_days": "365"
  },
  "created_at": "2025-08-23T10:15:30Z",
  "created_by": "user@domain.com",
  "updated_at": "2025-08-23T10:15:30Z",
  "updated_by": "user@domain.com"
}

```

### Basic

```go
package example

import (
    "time"

    cm "github.com/grasp-labs/ds-go-commonmodels"
)

func manualLifecycle() {
	subject := "user@domain.com"
	issuer := "grasp-labs"
	tenantId := uuid.MustParse("25948ccc-cf16-491e-9cd4-44d5ebb7bc54")

	meta := []map[string]string{
		{"owner_id": "xyz123"},
		{"retention": "365"},
	}
	tags := map[string]string{
		"env":       "prod",
		"tenant_id": tenantId.String(),
	}

	c := core.CoreModel{
		Name:     "Webhook config for X",
		Metadata: types.JSONB[[]map[string]string]{Data: meta},
		Tags:     types.JSONB[map[string]string]{Data: tags},
	}

	// Apply create-time defaults & audit
	c.Create(subject, issuer, tenantId)
}

```

### Gorm - autoflatten

GORM v2 automatically flattens anonymous (embedded) structs into the same table.

```go
package example

import (
    "time"

    cm "github.com/grasp-labs/ds-go-commonmodels"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

type Dataset struct {
    cm.Core  // ID, Name, Tags, Metadata, CreatedAt, UpdatedAt, etc.

    // Local fields
    FieldX string `json:"field_x"`
}

func setup() (*gorm.DB, error) {
    db, err := gorm.Open(postgres.Open("postgres://user:pass@host:5432/db?sslmode=enable"), &gorm.Config{})
    if err != nil { return nil, err }
    return db, db.AutoMigrate(&Dataset{})
}

func createExample(db *gorm.DB) error {
    subject := "user@domain.com"
	issuer := "grasp-labs"
	tenantId := uuid.MustParse("25948ccc-cf16-491e-9cd4-44d5ebb7bc54")

	meta := []map[string]string{
		{"owner_id": "xyz123"},
		{"retention": "365"},
	}
	tags := map[string]string{
		"env":       "prod",
		"tenant_id": tenantId.String(),
	}
    ds := &Dataset{
        Core: cm.CoreModel{
            Name:     "Webhook config for X",
		    Metadata: types.JSONB[[]map[string]string]{Data: meta},
		    Tags:     types.JSONB[map[string]string]{Data: tags},
		    Status:   "active",
        },
        FieldX: "awesome new field",
    }

    ds.CoreModel.Create(subject, issuer, tenantId)

    // GORM will trigger Core.BeforeCreate to apply defaults/audit fields.
    return db.Create(ds).Error
}

func updateExample(db *gorm.DB, id string) error {
    var ds Dataset
    if err := db.First(&ds, "id = ?", id).Error; err != nil { return err }

    ds.FieldX = "new awesome field data"
    // GORM will trigger Core.BeforeUpdate to refresh audit fields.
    return db.Save(&ds).Error
}

```