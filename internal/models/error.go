package models

// ErrorResponse represents an error response following RFC 7807
type ErrorResponse struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
}

// NewErrorResponse creates a new error response
func NewErrorResponse(title string, status int, detail string, instance string) *ErrorResponse {
	return &ErrorResponse{
		Type:     "about:blank", // As specified in requirements
		Title:    title,
		Status:   status,
		Detail:   detail,
		Instance: instance,
	}
}

// Common error responses
var (
	ErrBadRequest = func(detail, instance string) *ErrorResponse {
		return NewErrorResponse("Bad Request", 400, detail, instance)
	}

	ErrNotFound = func(resource string, id interface{}, instance string) *ErrorResponse {
		return NewErrorResponse(
			resource+" Not Found",
			404,
			"The "+resource+" with ID "+string(rune(id.(int)))+" was not found.",
			instance,
		)
	}

	ErrInternalServer = func(detail, instance string) *ErrorResponse {
		return NewErrorResponse("Internal Server Error", 500, detail, instance)
	}
)
