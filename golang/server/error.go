package server

import (
	"fmt"
)

const (
	ErrProcessingRequestCode = iota + 1000
	ErrNotFoundCode
	ErrNotAllowedCode
	ErrServiceUnavailableCode
	ErrRequiredParamCode
	ErrUnauthorizedCode
	ErrFailedConversionCode
)

// Error represents a error response.
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error formats an error response as a string.
func (e Error) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}
