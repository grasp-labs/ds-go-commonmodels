package validation_error

// ErrorEnvelope is used in order to wrap errors for response.
// Example:
//
//	{
//		"details": [
//			{
//				"field": "",
//				"message": "",
//				"loc": "",
//				"code": "",
//			}
//		]
//	}
type ErrorEnvelope struct {
	Details []ValidationError `json:"details"`
}
