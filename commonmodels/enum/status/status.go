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
