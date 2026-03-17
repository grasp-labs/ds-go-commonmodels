package scopes

import (
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Scope = func(*gorm.DB) *gorm.DB

// isValidIdentifier checks if a string is a valid SQL identifier.
// Valid identifiers must start with a letter or underscore, followed by letters, digits, or underscores.
// This enforces SQL naming standards and prevents injection attacks on field names.
var validIdentifierRegex = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

func isValidIdentifier(s string) bool {
	return validIdentifierRegex.MatchString(strings.TrimSpace(s))
}

// TankScope is a helper to create scopes for queries based on common query parameters.
func TankScope(params map[string]string) []Scope {
	scopes := QueryParamScopes(map[string]string{
		"id":       params["id"],
		"owner_id": params["owner_id"],
		"status":   params["status"],
	})
	// Handle timestamp comparisons
	if s := strings.TrimSpace(params["created_at_gte"]); s != "" {
		if ts, ok := ParseTimeFlexible(s); ok {
			scopes = append(scopes, FieldCmp("created_at", ">=", ts))
		}
	}
	if s := strings.TrimSpace(params["created_at_lte"]); s != "" {
		if ts, ok := ParseTimeFlexible(s); ok {
			scopes = append(scopes, FieldCmp("created_at", "<=", ts))
		}
	}
	if s := strings.TrimSpace(params["modified_at_gte"]); s != "" {
		if ts, ok := ParseTimeFlexible(s); ok {
			scopes = append(scopes, FieldCmp("modified_at", ">=", ts))
		}
	}
	if s := strings.TrimSpace(params["modified_at_lte"]); s != "" {
		if ts, ok := ParseTimeFlexible(s); ok {
			scopes = append(scopes, FieldCmp("modified_at", "<=", ts))
		}
	}
	return scopes
}

// QueryParamScopes creates scopes from a map of field names to values.
// For each field, it tries to parse as UUID first, then falls back to string.
// Empty values are skipped (no-op). Fields are processed in sorted order for deterministic output.
func QueryParamScopes(params map[string]string) []Scope {
	var scopes []Scope

	// Sort keys for deterministic scope ordering and stable SQL generation
	keys := make([]string, 0, len(params))
	for field := range params {
		keys = append(keys, field)
	}
	sort.Strings(keys)

	for _, field := range keys {
		val := strings.TrimSpace(params[field])
		if val == "" {
			continue
		}
		// Try UUID first
		if id, err := uuid.Parse(val); err == nil {
			scopes = append(scopes, Field(field, id))
			continue
		}
		// Fall back to string
		scopes = append(scopes, Field(field, val))
	}
	return scopes
}

// Field creates an equality scope. Zero values are no-ops.
// Validates fieldName to prevent SQL injection attacks on identifier names.
func Field(fieldName string, value interface{}) Scope {
	if !isValidIdentifier(fieldName) || IsZero(value) {
		return func(tx *gorm.DB) *gorm.DB { return tx }
	}
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Where("? = ?", clause.Column{Table: clause.CurrentTable, Name: fieldName}, value)
	}
}

// FieldCmp creates a comparison scope (>=, <=, !=, etc.). Zero values are no-ops.
// Only allows a whitelist of operators: =, !=, >, >=, <, <= to prevent SQL injection.
// Validates fieldName to prevent SQL injection attacks on identifier names.
func FieldCmp(fieldName string, operator string, value interface{}) Scope {
	// Validate fieldName and value
	if !isValidIdentifier(fieldName) || IsZero(value) {
		return func(tx *gorm.DB) *gorm.DB { return tx }
	}

	// Validate operator against whitelist to prevent SQL injection
	allowedOperators := map[string]bool{
		"=":  true,
		"!=": true,
		">":  true,
		">=": true,
		"<":  true,
		"<=": true,
	}
	operator = strings.TrimSpace(operator)
	if !allowedOperators[operator] {
		return func(tx *gorm.DB) *gorm.DB { return tx }
	}

	return func(tx *gorm.DB) *gorm.DB {
		return tx.Where("? "+operator+" ?", clause.Column{Table: clause.CurrentTable, Name: fieldName}, value)
	}
}

// Helper: isZero checks if a value is considered empty/zero.
func IsZero(v interface{}) bool {
	if v == nil {
		return true
	}
	switch val := v.(type) {
	case uuid.UUID:
		return val == uuid.Nil
	case time.Time:
		return val.IsZero()
	case string:
		return val == ""
	default:
		return reflect.ValueOf(v).IsZero()
	}
}

func ParseTimeFlexible(s string) (time.Time, bool) {
	// Try common formats
	for _, layout := range []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02", // date only
	} {
		if t, err := time.Parse(layout, s); err == nil {
			return t, true
		}
	}
	// Try unix seconds
	if sec, err := strconv.ParseInt(s, 10, 64); err == nil && sec > 0 {
		return time.Unix(sec, 0).UTC(), true
	}
	return time.Time{}, false
}

func OrderBy(param string, whitelist map[string]string) Scope {
	return func(tx *gorm.DB) *gorm.DB {
		param = strings.TrimSpace(param)
		if param == "" {
			return tx
		}
		// Parse
		var field string
		var desc bool
		switch {
		case strings.HasPrefix(param, "-"):
			field, desc = param[1:], true
		case strings.HasPrefix(param, "+"):
			field, desc = param[1:], false
		default:
			field, desc = param, false
		}
		// Validate against whitelist
		col, ok := whitelist[strings.ToLower(strings.TrimSpace(field))]
		if !ok || col == "" {
			return tx
		}
		return tx.Order(clause.OrderByColumn{
			Column: clause.Column{Name: col},
			Desc:   desc,
		})
	}
}

func Paginate(limit int, offset int) Scope {
	return func(tx *gorm.DB) *gorm.DB {
		if limit <= 0 {
			limit = 100 // default limit
		}
		if offset < 0 {
			offset = 0
		}
		return tx.Offset(offset).Limit(limit)
	}
}

// JSONBContains creates a scope to search for a key-value pair within a JSONB field.
// For example: JSONBContains("metadata", "env", "production") searches where metadata->>'env' = 'production'
// Validates fieldName to prevent SQL injection attacks on identifier names.
func JSONBContains(fieldName, key, value string) Scope {
	if key == "" || value == "" || !isValidIdentifier(fieldName) {
		return func(tx *gorm.DB) *gorm.DB { return tx }
	}
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Where(clause.Expr{
			SQL: "? ->> ? = ?",
			Vars: []interface{}{
				clause.Column{Table: clause.CurrentTable, Name: fieldName},
				key,
				value,
			},
		})
	}
}

// JSONBHasKey creates a scope to check if a JSONB field has a specific key.
// For example: JSONBHasKey("metadata", "env") searches where metadata ? 'env'
// Validates fieldName to prevent SQL injection attacks on identifier names.
func JSONBHasKey(fieldName, key string) Scope {
	if key == "" || !isValidIdentifier(fieldName) {
		return func(tx *gorm.DB) *gorm.DB { return tx }
	}
	return func(tx *gorm.DB) *gorm.DB {
		// Build the SQL with the ? operator as a literal.
		// fieldName is validated as a valid identifier, then quoted by GORM.
		// Only the key is parameterized to prevent injection.
		quotedField := tx.Statement.Quote(fieldName)
		return tx.Where(quotedField+" ? ?", key)
	}
}

// JSONBContainsValue creates a scope to search all keys in a JSONB field for a specific value.
// For example: JSONBContainsValue("metadata", "production") searches where any value in metadata equals 'production'
// Validates fieldName to prevent SQL injection attacks on identifier names.
func JSONBContainsValue(fieldName, value string) Scope {
	if value == "" || !isValidIdentifier(fieldName) {
		return func(tx *gorm.DB) *gorm.DB { return tx }
	}
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Where(clause.Expr{
			SQL: "EXISTS (SELECT 1 FROM jsonb_each_text(?) AS kv WHERE kv.value = ?)",
			Vars: []interface{}{
				clause.Column{Table: clause.CurrentTable, Name: fieldName},
				value,
			},
		})
	}
}

// JSONBContainsAll creates a scope to search for multiple key-value pairs in a JSONB field.
// All pairs must match (AND logic). Keys are sorted for deterministic SQL generation.
// Validates fieldName to prevent SQL injection attacks on identifier names.
func JSONBContainsAll(fieldName string, pairs map[string]string) Scope {
	if len(pairs) == 0 || !isValidIdentifier(fieldName) {
		return func(tx *gorm.DB) *gorm.DB { return tx }
	}
	return func(tx *gorm.DB) *gorm.DB {
		// Sort keys for deterministic SQL generation and better prepared statement reuse
		keys := make([]string, 0, len(pairs))
		for k := range pairs {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			v := pairs[k]
			if v == "" {
				continue
			}
			tx = tx.Where(clause.Expr{
				SQL: "? ->> ? = ?",
				Vars: []interface{}{
					clause.Column{Table: clause.CurrentTable, Name: fieldName},
					k,
					v,
				},
			})
		}
		return tx
	}
}
