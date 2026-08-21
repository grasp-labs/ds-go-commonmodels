package dto_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grasp-labs/ds-go-commonmodels/v3/commonmodels/dto"
)

func TestNullable_UnmarshalJSON_ValidValue(t *testing.T) {
	var n dto.Nullable[time.Time]
	err := json.Unmarshal([]byte(`"2026-01-15T00:00:00Z"`), &n)
	require.NoError(t, err)
	assert.True(t, n.IsSet())
	assert.True(t, n.HasValue())
	assert.False(t, n.IsNull())
	assert.False(t, n.IsZero())
	expected := time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC)
	require.NotNil(t, n.Ptr())
	assert.True(t, expected.Equal(*n.Ptr()))
}

func TestNullable_UnmarshalJSON_Null(t *testing.T) {
	n := dto.Some(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
	err := json.Unmarshal([]byte(`null`), &n)
	require.NoError(t, err)
	assert.True(t, n.IsSet())
	assert.False(t, n.HasValue())
	assert.True(t, n.IsNull())
	assert.Nil(t, n.Ptr())
}

func TestNullable_NotSetByDefault(t *testing.T) {
	var n dto.Nullable[time.Time]
	assert.False(t, n.IsSet())
	assert.False(t, n.HasValue())
	assert.False(t, n.IsNull())
	assert.True(t, n.IsZero())
}

func TestNullable_MarshalJSON_Valid(t *testing.T) {
	tm := time.Date(2026, 3, 15, 12, 0, 0, 0, time.UTC)
	n := dto.Some(tm)

	data, err := json.Marshal(n)
	require.NoError(t, err)
	assert.Equal(t, `"2026-03-15T12:00:00Z"`, string(data))
}

func TestNullable_MarshalJSON_Unset(t *testing.T) {
	var n dto.Nullable[time.Time]

	data, err := json.Marshal(n)
	require.NoError(t, err)
	assert.Equal(t, `null`, string(data))
}

func TestNullable_MarshalJSON_Null(t *testing.T) {
	n := dto.Null[time.Time]()

	data, err := json.Marshal(n)
	require.NoError(t, err)
	assert.Equal(t, `null`, string(data))
}

func TestNullable_FromPtr(t *testing.T) {
	tm := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	n := dto.FromPtr(&tm)
	assert.True(t, n.IsSet())
	assert.True(t, n.HasValue())
	assert.False(t, n.IsNull())
	require.NotNil(t, n.Ptr())
	assert.True(t, tm.Equal(*n.Ptr()))
}

func TestNullable_FromPtr_Nil(t *testing.T) {
	n := dto.FromPtr[time.Time](nil)
	assert.True(t, n.IsSet())
	assert.False(t, n.HasValue())
	assert.True(t, n.IsNull())
	assert.False(t, n.IsZero())
	assert.Nil(t, n.Ptr())
}

func TestNullable_Ptr(t *testing.T) {
	tm := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	n := dto.Some(tm)
	p := n.Ptr()
	require.NotNil(t, p)
	assert.True(t, tm.Equal(*p))
}

func TestNullable_Ptr_Nil(t *testing.T) {
	var n dto.Nullable[time.Time]
	assert.Nil(t, n.Ptr())
	assert.Nil(t, dto.Null[time.Time]().Ptr())
}

func TestNullable_UnmarshalJSON_ErrorDoesNotMutate(t *testing.T) {
	tm := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	n := dto.Some(tm)
	err := json.Unmarshal([]byte(`"not-a-time"`), &n)
	require.Error(t, err)
	assert.True(t, n.IsSet())
	assert.True(t, n.HasValue())
	require.NotNil(t, n.Ptr())
	assert.True(t, tm.Equal(*n.Ptr()))
}

func TestNullable_MarshalOmitzero(t *testing.T) {
	type wrap struct {
		Date dto.Nullable[time.Time] `json:"date,omitzero"`
	}

	unset, err := json.Marshal(wrap{})
	require.NoError(t, err)
	assert.Equal(t, `{}`, string(unset))

	explicitNull, err := json.Marshal(wrap{Date: dto.Null[time.Time]()})
	require.NoError(t, err)
	assert.Equal(t, `{"date":null}`, string(explicitNull))

	tm := time.Date(2026, 3, 15, 12, 0, 0, 0, time.UTC)
	value, err := json.Marshal(wrap{Date: dto.Some(tm)})
	require.NoError(t, err)
	assert.Equal(t, `{"date":"2026-03-15T12:00:00Z"}`, string(value))
}

func TestNullable_UnmarshalOmitzeroRoundTrip(t *testing.T) {
	type wrap struct {
		Date dto.Nullable[time.Time] `json:"date,omitzero"`
	}

	var omitted wrap
	require.NoError(t, json.Unmarshal([]byte(`{}`), &omitted))
	assert.False(t, omitted.Date.IsSet())
	assert.True(t, omitted.Date.IsZero())

	var explicitNull wrap
	require.NoError(t, json.Unmarshal([]byte(`{"date":null}`), &explicitNull))
	assert.True(t, explicitNull.Date.IsSet())
	assert.True(t, explicitNull.Date.IsNull())
}
