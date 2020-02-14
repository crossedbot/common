package simpleserver

import (
	"fmt"
)

// list of error response status codes
const (
	_ = iota
	ErrorProcessingStatus
	ErrorNotFoundStatus
	ErrorNotAllowedStatus
	ErrorServiceUnavailableStatus
)

// list of error response messages
const (
	ErrorProcessingText         = "failed to process request"
	ErrorNotFoundText           = "failed to find resource"
	ErrorNotAllowedText         = "method is not allowed"
	ErrorServiceUnavailableText = "service is unavailable"
)

// Error represents a error response.
type Error struct {
	Status int    `json:"status"`
	Text   string `json:"text"`
}

// String formats an error response as a string.
func (e Error) String() string {
	return fmt.Sprintf("[%d] %s", e.Status, e.Text)
}
