package trigger

type TriggerType string

const (
	Manual   TriggerType = "manual"
	Webhook  TriggerType = "webhook"
	Schedule TriggerType = "schedule"
)
