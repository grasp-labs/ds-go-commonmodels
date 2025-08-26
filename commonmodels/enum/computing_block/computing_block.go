package computingblock

type ComputingBlock string

const (
	Workflow ComputingBlock = "workflow"
	Pipeline ComputingBlock = "pipeline"
	Clone    ComputingBlock = "clone"
)
