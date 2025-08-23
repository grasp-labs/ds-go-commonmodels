package health_check

type HealthCheck struct {
	Message string `json:"message" example:"Server is running."`
	Time    string `json:"time" example:"2006-01-02T15:04:05Z07:00"`
}
