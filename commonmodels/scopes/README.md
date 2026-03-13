# Scopes

The **scopes** package provides composable query scopes for GORM, enabling flexible filtering, ordering and pagination of database queries. It includes helpers for common query patterns, JSONB field operations and dynamic query construction.

## Features

- Composable Scopes: Build GORM queries using reusable functions
- Dynamic Filtering: Create scopes from query parameters, including UUID, strings and timestamp fields
- Comparison Operators: Support for equality and comparison (e.g., `>=`, `<=`) on fields
- Ordering and Pagination: Safely order results and pagi ss snate within limit/offset
- JSONB Helpers: Query JSONB fields for key value and key existence

---

## Usage

The usage here is quite flexible, you can utilize the different Scope functions directly from the library if you need single off values, and if you need complex combination of scopes it's easy to build an adapter to combine whatever you need for the Scopes in your different handlers.

### Direct usage:

Handler:

```go
import (
	cmScopes "github.com/grasp-labs/ds-go-commonmodels/v3/commonmodels/scopes"
)

func (h *Handler) List(c echo.Context) error {
	// pagination
	limit, _ := strconv.ParseInt(c.QueryParam("limit"), 10, 64)
	offset, _ := strconv.ParseInt(c.QueryParam("offset"), 10, 64)
	if limit < 1 {
		limit = 500
	}
	if offset < 0 {
		offset = 0
	}

	// scopes
	scopes := cmScopes.QueryParamScopes((map[string]string{
		"id":                  c.QueryParam("id"),
		"owner_id":            c.QueryParam("owner_id"),
		"status":              c.QueryParam("status"),
		"scheduled_minute_eq": c.QueryParam("scheduled_minute_eq"),
		"created_at_gte":      c.QueryParam("created_at_gte"),
		"created_at_lte":      c.QueryParam("created_at_lte"),
		"modified_at_gte":     c.QueryParam("modified_at_gte"),
		"modified_at_lte":     c.QueryParam("modified_at_lte"),
	}))
	scopes = append(
		scopes,
		cmScopes.Field("tenant_id", rc.TenantID),
		cmScopes.OrderBy(c.QueryParam("order_by"), orderByWhitelist),
	)
	scopes = append(
		scopes,
		cmScopes.QueryParamScopes(map[string]string{
			"topic_id":            c.QueryParam("topic_id"),
			"session_id":          c.QueryParam("session_id"),
			"trigger_type":        c.QueryParam("trigger_type"),
			"trigger_id":          c.QueryParam("trigger_id"),
			"cb_id":               c.QueryParam("cb_id"),
			"cb_type":             c.QueryParam("cb_type"),
			"owner_id":            c.QueryParam("owner_id"),
			"status":              c.QueryParam("status"),
			"scheduled_minute_eq": c.QueryParam("scheduled_minute_eq"),
		})...,
	)

	items, err := h.Svc.List(rc.Ctx, append(scopes, cmScopes.Paginate(int(limit), int(offset)))...)

	------------------- rest of the handler ------------------>

}

```

### Indirect usage via adapter:

Adapter:

```go
package repo

import	"github.com/grasp-labs/ds-go-commonmodels/v3/commonmodels/scopes"

// TenantID scopes to rows with the given tenant_id.
func TenantID(t uuid.UUID) Scope {
	return cmScopes.Field("tenant_id", t)
}

func OrderBy(orderBy string, whitelist map[string]string) Scope {
	return cmScopes.OrderBy(orderBy, whitelist)
}

func TankScope(c echo.Context) []Scope {
	params := map[string]string{
		"id":                  c.QueryParam("id"),
		"owner_id":            c.QueryParam("owner_id"),
		"status":              c.QueryParam("status"),
		"scheduled_minute_eq": c.QueryParam("scheduled_minute_eq"),
		"created_at_gte":      c.QueryParam("created_at_gte"),
		"created_at_lte":      c.QueryParam("created_at_lte"),
		"modified_at_gte":     c.QueryParam("modified_at_gte"),
		"modified_at_lte":     c.QueryParam("modified_at_lte"),
	}
	return cmScopes.TankScope(params)
}

func EnqueuedJobScopes(c echo.Context) []Scope {
	params := map[string]string{
		"topic_id":     c.QueryParam("topic_id"),
		"session_id":   c.QueryParam("session_id"),
		"trigger_type": c.QueryParam("trigger_type"),
		"trigger_id":   c.QueryParam("trigger_id"),
		"cb_id":        c.QueryParam("cb_id"),
		"cb_type":      c.QueryParam("cb_type"),
	}
	return cmScopes.QueryParamScopes(params)
}
```

Handler:

```go

func (h *Handler) List(c echo.Context) error {
    limit, _ := strconv.Atoi(c.QueryParam("limit"))
    offset, _ := strconv.Atoi(c.QueryParam("offset"))
    orderBy := c.QueryParam("order_by")

    scopes := repo.TankScope(c)
    scopes = append(scopes,
        repo.TenantID(rc.TenantID),
        repo.OrderBy(orderBy, orderByWhitelist),
    )
    scopes = append(
	scopes,
	repo.EnqueuedJobScopes(c)...,
    )

    items, err := h.Svc.List(rc.Ctx, append(scopes, repo.Paginate(limit, offset))...)

    ------------------- rest of the handler ------------------>
}
```

**Note**: It's split in two operations since Go does not allow for unpacking in the same operation as the others
