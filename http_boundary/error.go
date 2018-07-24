package http_boundary

import "fmt"

type HttpError struct {
	Href *Link `json:"self"`
	*ApiError
}

type ApiError struct {
	Message    string `json:"message,omitempty"`
	StatusCode int    `json:"statusCode,omitempty"`
}

func (a ApiError) Error() string {
	return fmt.Sprintf("HTTPErr: %d=%s", a.StatusCode, a.Message)
}
