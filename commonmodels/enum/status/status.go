package status

type Status string
type ProcessStatus string

const (
	Active    Status = "active"
	Deleted   Status = "deleted"
	Suspended Status = "suspended"
	Rejected  Status = "rejected"
	Draft     Status = "draft"
	Closed    Status = "closed"

	Queued    ProcessStatus = "queued"
	Running   ProcessStatus = "running"
	Completed ProcessStatus = "completed"
	Failed    ProcessStatus = "failed"
	Cancelled ProcessStatus = "cancelled"
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

// Map of statuses that can be used to check if a status is valid,
//
// Example:
//
// s := "hello"
// _, ok := ValidProcessStatus[s]
// if !ok {...}
var ValidProcessStatus = map[ProcessStatus]struct{}{
	Queued:    {},
	Running:   {},
	Completed: {},
	Failed:    {},
	Cancelled: {},
}
