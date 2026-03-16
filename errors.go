package bizverify

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// APIError is the base error type for all BizVerify API errors.
type APIError struct {
	Message    string      `json:"message"`
	Code       string      `json:"code"`
	StatusCode int         `json:"status_code"`
	Details    interface{} `json:"details,omitempty"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("bizverify: %s (code=%s, status=%d)", e.Message, e.Code, e.StatusCode)
}

type AuthenticationError struct{ *APIError }
type AuthorizationError struct{ *APIError }
type ValidationError struct{ *APIError }
type NotFoundError struct{ *APIError }
type InsufficientCreditsError struct{ *APIError }
type ConflictError struct{ *APIError }
type InternalError struct{ *APIError }
type TimeoutError struct{ *APIError }

type RateLimitError struct {
	*APIError
	RetryAfter int
}

type JobFailedError struct {
	*APIError
	JobID string
}

func (e *AuthenticationError) Unwrap() error      { return e.APIError }
func (e *AuthorizationError) Unwrap() error       { return e.APIError }
func (e *ValidationError) Unwrap() error          { return e.APIError }
func (e *NotFoundError) Unwrap() error            { return e.APIError }
func (e *InsufficientCreditsError) Unwrap() error { return e.APIError }
func (e *ConflictError) Unwrap() error            { return e.APIError }
func (e *InternalError) Unwrap() error            { return e.APIError }
func (e *TimeoutError) Unwrap() error             { return e.APIError }
func (e *RateLimitError) Unwrap() error           { return e.APIError }
func (e *JobFailedError) Unwrap() error           { return e.APIError }

type errorResponseBody struct {
	ErrorData struct {
		Code    string      `json:"code"`
		Message string      `json:"message"`
		Details interface{} `json:"details,omitempty"`
	} `json:"error"`
}

func parseErrorResponse(statusCode int, body []byte, headers http.Header) error {
	var resp errorResponseBody
	if err := json.Unmarshal(body, &resp); err != nil {
		return &APIError{Message: "Unknown error", Code: "UNKNOWN", StatusCode: statusCode}
	}

	base := &APIError{
		Message:    resp.ErrorData.Message,
		Code:       resp.ErrorData.Code,
		StatusCode: statusCode,
		Details:    resp.ErrorData.Details,
	}

	switch statusCode {
	case 400, 422:
		return &ValidationError{APIError: base}
	case 401:
		return &AuthenticationError{APIError: base}
	case 402:
		return &InsufficientCreditsError{APIError: base}
	case 403:
		return &AuthorizationError{APIError: base}
	case 404:
		return &NotFoundError{APIError: base}
	case 409:
		return &ConflictError{APIError: base}
	case 429:
		retryAfter := 0
		if ra := headers.Get("Retry-After"); ra != "" {
			if v, err := strconv.Atoi(ra); err == nil {
				retryAfter = v
			}
		}
		return &RateLimitError{APIError: base, RetryAfter: retryAfter}
	default:
		if statusCode >= 500 {
			return &InternalError{APIError: base}
		}
		return base
	}
}
