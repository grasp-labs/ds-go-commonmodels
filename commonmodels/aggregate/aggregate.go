// Package aggregate is the shared contract for record statistics.
//
// A resource that offers them serves, under /<resource>/aggregates/:
//
//	GET /<resource>/aggregates/          Index: the levels it offers
//	GET /<resource>/aggregates/total/    Total
//	GET /<resource>/aggregates/<level>/  one page of Group in the page.Page envelope
//
// Every endpoint counts what the resource's list endpoint would return for the
// same filters, so one filter reads the same at every level. A level groups the
// records by one field; each group's Key is a value its Filter accepts, which
// drills down from the group to its records or to the level below.
package aggregate

import "time"

// Pagination and order of a level. Groups are ordered by OrderBy fields,
// prefixed "-" for descending or "+" for ascending; Key breaks ties.
const (
	DefaultLimit   = 500
	MaxLimit       = 10000
	DefaultOrderBy = "-count"
)

// OrderBy are the fields a level orders by; resources without bytes drop
// "bytes" from their Index.
var OrderBy = []string{"key", "count", "bytes", "modified_at"}

// Intervals of a time level, which keys each group by the UTC start of its
// interval in RFC 3339.
const (
	IntervalHour = "hour"
	IntervalDay  = "day"
)

// Total counts the records. Bytes is omitted by resources that store no
// payload; ModifiedAt is the latest change, nil without records.
type Total struct {
	Count      int64      `json:"count"`
	Bytes      *int64     `json:"bytes,omitempty"`
	ModifiedAt *time.Time `json:"modified_at"`
}

// Group counts the records whose level field is Key. Key is "" for records
// without a value, which no Filter narrows to.
type Group struct {
	Key        string     `json:"key"`
	Count      int64      `json:"count"`
	Bytes      *int64     `json:"bytes,omitempty"`
	ModifiedAt *time.Time `json:"modified_at"`
}

// Index lists the levels a resource offers.
type Index struct {
	Levels  []Level  `json:"levels"`
	OrderBy []string `json:"order_by"`
}

// Level is one way a resource groups its records.
type Level struct {
	// Name is the path segment: /<resource>/aggregates/<name>/.
	Name string `json:"name"`
	// Field is the record field grouped, e.g. "status".
	Field string `json:"field"`
	// Filter is the list query parameter a group's Key narrows to, e.g.
	// "status", or "key_prefix" for a folder.
	Filter string `json:"filter"`
	// Bytes reports whether groups carry bytes.
	Bytes bool `json:"bytes"`
	// Intervals are the interval values a time level accepts; the first is
	// the default. Empty for other levels.
	Intervals []string `json:"intervals,omitempty"`
}
