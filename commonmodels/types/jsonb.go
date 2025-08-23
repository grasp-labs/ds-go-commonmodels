package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// JSONB wraps any Go type T and tells GORM to persist it as Postgres JSONB.
type JSONB[T any] struct {
	Data T
}

// Validate returns an error if Data can’t round-trip to JSON.
func (j JSONB[T]) Validate() error {
	v := j.Data

	// Allow nil or empty values
	switch t := any(v).(type) {
	case map[string]string:
		if len(t) == 0 {
			return nil
		}
	case []map[string]string:
		if len(t) == 0 {
			return nil
		}
	}

	// Marshal to check JSON round-trip
	_, err := json.Marshal(v)
	return err
}

// GormDataType makes GORM map this to a JSONB column.
func (JSONB[T]) GormDataType() string {
	return "jsonb"
}

// Value is called by database/sql to turn our Data into JSON bytes.
func (j JSONB[T]) Value() (driver.Value, error) {
	return json.Marshal(j.Data)
}

// Scan is called by database/sql to populate j.Data from JSON bytes.
func (j *JSONB[T]) Scan(src any) error {
	if src == nil {
		// Accept null JSONB column
		var zero T
		j.Data = zero
		return nil
	}

	switch v := src.(type) {
	case []byte:
		return json.Unmarshal(v, &j.Data)
	default:
		return fmt.Errorf("JSONB.Scan: Unsupported type %T", v)
	}
}

// MarshalJSON/UnmarshalJSON make sure encoding/json
// treats JSONB[T] as if it were T directly.
func (j JSONB[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.Data)
}
func (j *JSONB[T]) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &j.Data)
}
