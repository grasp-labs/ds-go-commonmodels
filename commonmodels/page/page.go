package page

import (
	"math"
)

const (
	DefaultPageSize = 20
	DefaultOffset   = 0
)

type Page struct {
	// Page number (1-based)
	Page int `json:"page" validate:"required"`
	// Items per page
	PageSize int `json:"page_size" validate:"required"`
	// Total items
	Total int64 `json:"total" validate:"required"`
	// Total pages (derived)
	TotalPages int64 `json:"total_pages" validate:"required"`
	// Is there a previous page
	HasPrev bool `json:"has_prev" validate:"required"`
	// Is there a next page
	HasNext bool `json:"has_next" validate:"required"`
}

// Sanitize normalizes inputs and enforces minimal bounds.
func Sanitize(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = DefaultPageSize
	}
	return page, pageSize
}

// Compute builds the Page metadata from total + inputs.
// If total == 0 -> page=1, totalPages=0, no prev/next.
// If total > 0 -> clamps page to [1..totalPages].
func Compute(total int64, page, pageSize int) Page {
	page, pageSize = Sanitize(page, pageSize)

	var totalPages int64
	if total > 0 {
		totalPages = (total + int64(pageSize) - 1) / int64(pageSize)

		// Clamp page into range while guarding against casting overflow.
		switch {
		case totalPages == 0:
			page = 1
		case int64(page) > totalPages:
			// If totalPages exceeds int range (extremely large datasets),
			// clamp to MaxInt to avoid overflow.
			if totalPages > int64(math.MaxInt) {
				page = math.MaxInt
			} else {
				page = int(totalPages)
			}
		case page < 1:
			page = 1
		}
	}

	return Page{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
		HasPrev:    totalPages > 0 && page > 1,
		HasNext:    totalPages > 0 && int64(page) < totalPages,
	}
}

func CalculatePage(limit, offset, total int64) Page {
	// Set pageSize from limit, if limit is zero or negative, use default
	pageSize := int(limit)
	if pageSize <= 0 {
		pageSize = DefaultPageSize
	}

	// Ensure offset is non-negative
	if offset < 0 {
		offset = DefaultOffset
	}

	page := int(offset)/pageSize + 1

	return Compute(total, page, pageSize)
}

// Generic response envelope: works for any element type T,
// e.g. should support any datamodel we choose to implement.
type Response[T any] struct {
	Data []T  `json:"data"`
	Page Page `json:"page"`
}

// Convenience constructor.
//
// # Example:
// r := page.NewResponse(data, total, pageNr, pageSize)
//
//	    {
//		  "data":[
//		    {"id":1,"foo":"bar"},
//		    {"id":2,"foo":"bar"}],
//		  "page":{
//		    "page":1,
//		    "page_size":20,
//		    "total":2,
//		    "total_pages":1,
//		    "has_prev":false,
//		    "has_next":false
//		  }
//		}
func NewResponse[T any](items []T, total int64, page, pageSize int) Response[T] {
	return Response[T]{Data: items, Page: Compute(total, page, pageSize)}
}
