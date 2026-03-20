package bizverify

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestAuthRequestAccess(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/v1/auth/request-access" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		if body["email"] != "test@example.com" {
			t.Errorf("wrong email: %v", body["email"])
		}
		if body["accept_terms"] != true {
			t.Errorf("wrong accept_terms: %v", body["accept_terms"])
		}
		writeJSON(w, 200, fixtureRequestAccessResp)
	})
	resp, err := client.Auth.RequestAccess(context.Background(), "test@example.com", true)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Message != "Verification code sent to your email" {
		t.Errorf("wrong message: %s", resp.Message)
	}
}

func TestAuthVerifyAccess(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/v1/auth/verify-access" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		if body["email"] != "test@example.com" {
			t.Errorf("wrong email: %v", body["email"])
		}
		if body["code"] != "123456" {
			t.Errorf("wrong code: %v", body["code"])
		}
		writeJSON(w, 200, fixtureVerifyAccessResp)
	})
	resp, err := client.Auth.VerifyAccess(context.Background(), "test@example.com", "123456", nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.APIKey != "bv_live_abc123def456" {
		t.Errorf("wrong api key: %s", resp.APIKey)
	}
	if resp.KeyID != "key_123" {
		t.Errorf("wrong key id: %s", resp.KeyID)
	}
	if resp.Label != "my-key" {
		t.Errorf("wrong label: %s", resp.Label)
	}
}

func TestAuthVerifyAccessAutoConfiguresAPIKey(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, fixtureVerifyAccessResp)
	})
	client.Auth.VerifyAccess(context.Background(), "test@example.com", "123456", nil)
	client.client.mu.RLock()
	if client.client.apiKey != "bv_live_abc123def456" {
		t.Errorf("expected API key to be set, got '%s'", client.client.apiKey)
	}
	client.client.mu.RUnlock()
}

func TestAuthVerifyAccessWithLabel(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		if body["label"] != "my-agent" {
			t.Errorf("wrong label: %v", body["label"])
		}
		writeJSON(w, 200, fixtureVerifyAccessResp)
	})
	label := "my-agent"
	_, err := client.Auth.VerifyAccess(context.Background(), "test@example.com", "123456", &label)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAuthNoAuthHeaders(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API-Key") != "" || r.Header.Get("Authorization") != "" {
			t.Error("auth headers should not be set for auth endpoints")
		}
		writeJSON(w, 200, fixtureRequestAccessResp)
	})
	// Clear auth to verify none is sent by the endpoint itself
	client.client.mu.Lock()
	client.client.apiKey = ""
	client.client.mu.Unlock()
	client.Auth.RequestAccess(context.Background(), "t@t.com", true)
}
