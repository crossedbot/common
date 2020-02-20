package server

import (
	"fmt"
)

// list of common error response objects
var (
	ErrProcessingRequest    = Error{Code: 1000, Title: "Error Processing Request", Text: "failed to process request"}
	ErrNotFound             = Error{Code: 1001, Title: "Not Found", Text: "failed to find resource"}
	ErrNotAllowed           = Error{Code: 1002, Title: "Not Allowed", Text: "method not allowed"}
	ErrorServiceUnavailable = Error{Code: 1003, Title: "Service Unavailable", Text: "service is unavailable"}
)

// Error represents a error response.
type Error struct {
	Code  int    `json:"code"`
	Title string `json:"title"`
	Text  string `json:"text"`
}

// String formats an error response as a string.
func (e Error) String() string {
	return fmt.Sprintf("[%d] %s: %s", e.Code, e.Title, e.Text)
}
