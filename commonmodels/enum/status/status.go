package status

type Status string

const (
	Active    Status = "active"
	Deleted   Status = "deleted"
	Suspended Status = "suspended"
	Rejected  Status = "rejected"
	Draft     Status = "draft"
	Closed    Status = "closed"
)

// Map of statuses that can be used to check if a status is valid,
//
// Example:
//
// s := "hello"
// _, ok := ValidStatus[s]
// if !ok {...}
var ValidStatus = map[Status]struct{}{
	Active:    {},
	Deleted:   {},
	Suspended: {},
	Rejected:  {},
	Draft:     {},
	Closed:    {},
}
