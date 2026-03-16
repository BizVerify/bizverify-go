package bizverify

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestAuthRegister(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/v1/auth/register" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		var body map[string]interface{}
		json.NewDecoder(r.Body).Decode(&body)
		if body["email"] != "test@example.com" {
			t.Errorf("wrong email: %v", body["email"])
		}
		writeJSON(w, 200, fixtureRegisterResp)
	})
	resp, err := client.Auth.Register(context.Background(), "test@example.com", "pass123", true)
	if err != nil {
		t.Fatal(err)
	}
	if resp.APIKey != "bv_live_abc123def456" {
		t.Errorf("wrong api key: %s", resp.APIKey)
	}
	if resp.User.Email != "test@example.com" {
		t.Errorf("wrong email: %s", resp.User.Email)
	}
}

func TestAuthLoginStoresToken(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, fixtureLoginResp)
	})
	resp, err := client.Auth.Login(context.Background(), "test@example.com", "pass123")
	if err != nil {
		t.Fatal(err)
	}
	if resp.Token != "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test" {
		t.Errorf("wrong token: %s", resp.Token)
	}
	client.client.mu.RLock()
	if client.client.token != resp.Token {
		t.Error("token not stored on client")
	}
	client.client.mu.RUnlock()
}

func TestAuthVerifyEmail(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/auth/verify-email" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		writeJSON(w, 200, fixtureMsgResp)
	})
	resp, err := client.Auth.VerifyEmail(context.Background(), "test@example.com", "123456")
	if err != nil {
		t.Fatal(err)
	}
	if resp.Message != "Success" {
		t.Errorf("wrong message: %s", resp.Message)
	}
}

func TestAuthResendVerification(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, fixtureMsgResp)
	})
	resp, err := client.Auth.ResendVerification(context.Background(), "test@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if resp.Message != "Success" {
		t.Errorf("wrong message: %s", resp.Message)
	}
}

func TestAuthForgotPassword(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, fixtureMsgResp)
	})
	resp, err := client.Auth.ForgotPassword(context.Background(), "test@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if resp.Message != "Success" {
		t.Errorf("wrong message: %s", resp.Message)
	}
}

func TestAuthResetPassword(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/auth/reset-password" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		writeJSON(w, 200, fixtureMsgResp)
	})
	resp, err := client.Auth.ResetPassword(context.Background(), "test@example.com", "123456", "newpass")
	if err != nil {
		t.Fatal(err)
	}
	if resp.Message != "Success" {
		t.Errorf("wrong message: %s", resp.Message)
	}
}

func TestAuthNoAuthHeaders(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API-Key") != "" || r.Header.Get("Authorization") != "" {
			t.Error("auth headers should not be set for auth endpoints")
		}
		writeJSON(w, 200, fixtureRegisterResp)
	})
	// Clear auth to verify none is sent by the endpoint itself
	client.client.mu.Lock()
	client.client.apiKey = ""
	client.client.token = ""
	client.client.mu.Unlock()
	client.Auth.Register(context.Background(), "t@t.com", "p", true)
}
