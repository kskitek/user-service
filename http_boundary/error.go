package http_boundary

import "fmt"

type Handler interface {
	Routes() []*Route
}
type HttpError struct {
	Href *Link `json:"self"`
	*ApiError
}

type ApiError struct {
	Message    string `json:"event,omitempty"`
	StatusCode int    `json:"statusCode,omitempty"`
}

func (a ApiError) Error() string {
	return fmt.Sprintf("HTTPErr: %d=%s", a.StatusCode, a.Message)
}
