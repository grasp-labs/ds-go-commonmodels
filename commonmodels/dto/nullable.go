package dto

import (
	"bytes"
	"encoding/json"
)

type Nullable[T any] struct {
	value T
	valid bool
	set   bool
}

func Some[T any](v T) Nullable[T] {
	return Nullable[T]{value: v, valid: true, set: true}
}

func Null[T any]() Nullable[T] {
	return Nullable[T]{set: true}
}

func FromPtr[T any](p *T) Nullable[T] {
	if p == nil {
		return Null[T]()
	}
	return Some(*p)
}

func (n Nullable[T]) Ptr() *T {
	if !n.valid {
		return nil
	}
	v := n.value
	return &v
}

func (n Nullable[T]) IsSet() bool {
	return n.set
}

func (n Nullable[T]) HasValue() bool {
	return n.valid
}

func (n Nullable[T]) IsNull() bool {
	return n.set && !n.valid
}

// On a DTO we can use keyword: omitzero in order to omit the field if it is zero. Useful for PATCH.
func (n Nullable[T]) IsZero() bool {
	return !n.set
}

func (n *Nullable[T]) UnmarshalJSON(data []byte) error {
	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		*n = Null[T]()
		return nil
	}
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*n = Some(v)
	return nil
}

func (n Nullable[T]) MarshalJSON() ([]byte, error) {
	if n.valid {
		return json.Marshal(n.value)
	}
	return []byte(`null`), nil
}
