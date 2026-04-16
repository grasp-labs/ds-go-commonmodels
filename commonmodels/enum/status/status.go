package status

type Status string

const (
	Inactive  Status = "inactive"
	Active    Status = "active"
	Deleted   Status = "deleted"
	Suspended Status = "suspended"
	Rejected  Status = "rejected"
	Draft     Status = "draft"
	Closed    Status = "closed"
)

type JobStatus string

const (
	JobStatusNew       JobStatus = "new"
	JobStatusQueued    JobStatus = "queued"
	JobStatusRunning   JobStatus = "running"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"
	JobStatusCancelled JobStatus = "cancelled"
)

// Map of statuses that can be used to check if a status is valid,
//
// Example:
//
// s := "hello"
// _, ok := ValidStatus[s]
// if !ok {...}
var ValidStatus = map[Status]struct{}{
	Inactive:  {},
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
var ValidJobStatus = map[JobStatus]struct{}{
	JobStatusNew:       {},
	JobStatusQueued:    {},
	JobStatusRunning:   {},
	JobStatusCompleted: {},
	JobStatusFailed:    {},
	JobStatusCancelled: {},
}
