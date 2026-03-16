package bizverify

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

func TestParseError400(t *testing.T) {
	body, _ := json.Marshal(errBody("VALIDATION_ERROR", "Bad input"))
	err := parseErrorResponse(400, body, http.Header{})
	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected ValidationError, got %T", err)
	}
	if ve.Code != "VALIDATION_ERROR" {
		t.Errorf("expected code VALIDATION_ERROR, got %s", ve.Code)
	}
}

func TestParseError401(t *testing.T) {
	body, _ := json.Marshal(errBody("INVALID_API_KEY", "Bad key"))
	err := parseErrorResponse(401, body, http.Header{})
	var ae *AuthenticationError
	if !errors.As(err, &ae) {
		t.Fatalf("expected AuthenticationError, got %T", err)
	}
}

func TestParseError402(t *testing.T) {
	body, _ := json.Marshal(errBody("INSUFFICIENT_CREDITS", "No credits"))
	err := parseErrorResponse(402, body, http.Header{})
	var ie *InsufficientCreditsError
	if !errors.As(err, &ie) {
		t.Fatalf("expected InsufficientCreditsError, got %T", err)
	}
}

func TestParseError403(t *testing.T) {
	body, _ := json.Marshal(errBody("FORBIDDEN", "Not allowed"))
	err := parseErrorResponse(403, body, http.Header{})
	var ae *AuthorizationError
	if !errors.As(err, &ae) {
		t.Fatalf("expected AuthorizationError, got %T", err)
	}
}

func TestParseError404(t *testing.T) {
	body, _ := json.Marshal(errBody("NOT_FOUND", "Missing"))
	err := parseErrorResponse(404, body, http.Header{})
	var ne *NotFoundError
	if !errors.As(err, &ne) {
		t.Fatalf("expected NotFoundError, got %T", err)
	}
}

func TestParseError409(t *testing.T) {
	body, _ := json.Marshal(errBody("ALREADY_REFUNDED", "Dup"))
	err := parseErrorResponse(409, body, http.Header{})
	var ce *ConflictError
	if !errors.As(err, &ce) {
		t.Fatalf("expected ConflictError, got %T", err)
	}
}

func TestParseError422(t *testing.T) {
	body, _ := json.Marshal(errBody("VALIDATION_ERROR", "Bad data"))
	err := parseErrorResponse(422, body, http.Header{})
	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected ValidationError, got %T", err)
	}
}

func TestParseError429WithRetryAfter(t *testing.T) {
	body, _ := json.Marshal(errBody("RATE_LIMIT", "Too fast"))
	headers := http.Header{}
	headers.Set("Retry-After", "30")
	err := parseErrorResponse(429, body, headers)
	var re *RateLimitError
	if !errors.As(err, &re) {
		t.Fatalf("expected RateLimitError, got %T", err)
	}
	if re.RetryAfter != 30 {
		t.Errorf("expected RetryAfter 30, got %d", re.RetryAfter)
	}
}

func TestParseError429WithoutRetryAfter(t *testing.T) {
	body, _ := json.Marshal(errBody("RATE_LIMIT", "Too fast"))
	err := parseErrorResponse(429, body, http.Header{})
	var re *RateLimitError
	if !errors.As(err, &re) {
		t.Fatalf("expected RateLimitError, got %T", err)
	}
	if re.RetryAfter != 0 {
		t.Errorf("expected RetryAfter 0, got %d", re.RetryAfter)
	}
}

func TestParseError500(t *testing.T) {
	body, _ := json.Marshal(errBody("INTERNAL", "Server error"))
	err := parseErrorResponse(500, body, http.Header{})
	var ie *InternalError
	if !errors.As(err, &ie) {
		t.Fatalf("expected InternalError, got %T", err)
	}
}

func TestParseError503(t *testing.T) {
	body, _ := json.Marshal(errBody("UNAVAILABLE", "Down"))
	err := parseErrorResponse(503, body, http.Header{})
	var ie *InternalError
	if !errors.As(err, &ie) {
		t.Fatalf("expected InternalError, got %T", err)
	}
}

func TestParseErrorUnknownStatus(t *testing.T) {
	body, _ := json.Marshal(errBody("TEAPOT", "I'm a teapot"))
	err := parseErrorResponse(418, body, http.Header{})
	var be *APIError
	if !errors.As(err, &be) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if be.StatusCode != 418 {
		t.Errorf("expected status 418, got %d", be.StatusCode)
	}
}

func TestErrorsAsUnwrap(t *testing.T) {
	body, _ := json.Marshal(errBody("TEST", "test"))
	err := parseErrorResponse(401, body, http.Header{})
	var base *APIError
	if !errors.As(err, &base) {
		t.Error("expected errors.As to work for base *Error")
	}
}
