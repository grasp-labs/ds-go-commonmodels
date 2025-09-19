package validation_error

type Location string

const (
	Query Location = "query"
	Body  Location = "body"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Loc     string `json:"loc"`
	Code    string `json:"code"`
}

var ValidLocations = map[Location]struct{}{
	Query: {},
	Body:  {},
}
