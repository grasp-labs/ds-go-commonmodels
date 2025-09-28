package page_test

import (
	"encoding/json"
	"testing"

	page "github.com/grasp-labs/ds-go-commonmodels/v2/commonmodels/page"
	"github.com/stretchr/testify/assert"
)

type AnyModel struct {
	ID  int    `json:"id"`
	Foo string `json:"foo"`
}

func TestPage_NewResponse_small(t *testing.T) {
	data := []AnyModel{
		{
			ID:  1,
			Foo: "bar",
		},
		{
			ID:  2,
			Foo: "bar",
		},
	}
	total := int64(2)
	pageNr := 1
	pageSize := 20
	r := page.NewResponse(data, total, pageNr, pageSize)
	assert.Equal(t, r.Page.Total, total)
	assert.Equal(t, r.Page.HasNext, false)
	assert.Equal(t, r.Page.HasPrev, false)
	assert.Equal(t, r.Page.PageSize, page.DefaultPageSize)
	assert.Equal(t, r.Page.TotalPages, int64(1))

	jsonBytes, err := json.Marshal(r)
	jsonStr := string(jsonBytes)
	if err != nil {
		t.Fatalf("failed to marshal, err: %v", err)
	}
	expected := `{"data":[{"id":1,"foo":"bar"},{"id":2,"foo":"bar"}],"page":{"page":1,"page_size":20,"total":2,"total_pages":1,"has_prev":false,"has_next":false}}`
	assert.Equal(t, expected, jsonStr)
}

func TestPage_NewResponseXL(t *testing.T) {
	type Foo struct {
		ID  int    `json:"id"`
		Bar string `json:"bar"`
	}

	limit := page.DefaultPageSize
	offset := 0
	var data []Foo
	size := 500
	for n := range size {
		data = append(data, Foo{ID: n, Bar: "yea"})
	}

	totalCount := len(data)
	s := data[offset:limit]
	r := page.NewResponse(s, int64(totalCount), 2, page.DefaultPageSize)
	assert.Equal(t, r.Page.Total, int64(500))
	assert.Equal(t, len(r.Data), r.Page.PageSize)
}
