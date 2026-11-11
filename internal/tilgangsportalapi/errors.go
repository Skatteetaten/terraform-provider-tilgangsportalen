package tilgangsportalapi

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// APIError is returned when the API responds with a non-2xx status.
// The API always responds with HTTP 500; the actual error code is the prefix of Message, e.g. "605: ...".
type APIError struct {
	Status     string
	StatusCode int
	Code       string
	Message    string
}

func (e *APIError) Error() string {
	if e.Message == "" {
		return fmt.Sprintf("response: %s, HTTP request failed with status code: %d", e.Status, e.StatusCode)
	}
	return fmt.Sprintf("response: %s, HTTP request failed with status code: %d, message: %s", e.Status, e.StatusCode, e.Message)
}

func newAPIError(status string, statusCode int, message string) *APIError {
	apiErr := &APIError{Status: status, StatusCode: statusCode, Message: message}
	if code, _, found := strings.Cut(message, ":"); found {
		code = strings.TrimSpace(code)
		if _, err := strconv.Atoi(code); err == nil {
			apiErr.Code = code
		}
	}
	return apiErr
}

func apiErrorCode(err error) string {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.Code
	}
	return ""
}

// IsErrorDoesNotExistOrUnauthorized reports API error code 605.
func IsErrorDoesNotExistOrUnauthorized(err error) bool {
	return apiErrorCode(err) == "605"
}

// IsErrorDoesNotExist reports API error code 606.
func IsErrorDoesNotExist(err error) bool {
	return apiErrorCode(err) == "606"
}

// IsErrorAlreadyExists reports API error code 609.
func IsErrorAlreadyExists(err error) bool {
	return apiErrorCode(err) == "609"
}
