package bizverify

import (
	"testing"
	"time"
)

func TestNewCreatesAllServices(t *testing.T) {
	c := New(WithAPIKey("key"))
	if c.Auth == nil {
		t.Error("Auth service is nil")
	}
	if c.Verification == nil {
		t.Error("Verification service is nil")
	}
	if c.Entities == nil {
		t.Error("Entities service is nil")
	}
	if c.Search == nil {
		t.Error("Search service is nil")
	}
	if c.Account == nil {
		t.Error("Account service is nil")
	}
	if c.Billing == nil {
		t.Error("Billing service is nil")
	}
	if c.Checker == nil {
		t.Error("Checker service is nil")
	}
	if c.Config == nil {
		t.Error("Config service is nil")
	}
}

func TestWithAPIKey(t *testing.T) {
	c := New(WithAPIKey("my-key"))
	if c.client.apiKey != "my-key" {
		t.Errorf("expected api key 'my-key', got '%s'", c.client.apiKey)
	}
}

func TestWithBaseURL(t *testing.T) {
	c := New(WithBaseURL("https://custom.api.com/"))
	if c.client.baseURL != "https://custom.api.com" {
		t.Errorf("expected trimmed base URL, got '%s'", c.client.baseURL)
	}
}

func TestWithMaxRetries(t *testing.T) {
	c := New(WithMaxRetries(5))
	if c.client.maxRetries != 5 {
		t.Errorf("expected 5 retries, got %d", c.client.maxRetries)
	}
}

func TestWithTimeout(t *testing.T) {
	c := New(WithTimeout(10 * time.Second))
	if c.client.client.Timeout != 10*time.Second {
		t.Errorf("expected 10s timeout, got %v", c.client.client.Timeout)
	}
}

func TestDefaultValues(t *testing.T) {
	c := New()
	if c.client.baseURL != defaultBaseURL {
		t.Errorf("expected default base URL, got '%s'", c.client.baseURL)
	}
	if c.client.maxRetries != defaultMaxRetries {
		t.Errorf("expected default retries, got %d", c.client.maxRetries)
	}
}

func TestLastResponseMetaDelegatesToHTTPClient(t *testing.T) {
	c := New()
	if c.LastResponseMeta() != nil {
		t.Error("expected nil meta before any request")
	}
}
