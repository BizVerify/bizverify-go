package bizverify

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestClientSuccessfulGet(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}
		writeJSON(w, 200, map[string]interface{}{"id": "123"})
	})
	ctx := context.Background()
	var result map[string]interface{}
	err := client.client.request(ctx, requestOptions{method: "GET", path: "/test", auth: authAPIKey}, &result)
	if err != nil {
		t.Fatal(err)
	}
	if result["id"] != "123" {
		t.Errorf("unexpected result: %v", result)
	}
}

func TestClient204ReturnsNil(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(204)
	})
	ctx := context.Background()
	err := client.client.request(ctx, requestOptions{method: "DELETE", path: "/test", auth: authJWT}, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestClientErrorResponse(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 400, errBody("BAD", "bad request"))
	})
	ctx := context.Background()
	err := client.client.request(ctx, requestOptions{method: "GET", path: "/test", auth: authNone}, nil)
	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected ValidationError, got %T: %v", err, err)
	}
}

func TestClientRetryOn500(t *testing.T) {
	var count int64
	client, _ := setupTestClientRetries(t, func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt64(&count, 1)
		if n < 3 {
			writeJSON(w, 500, errBody("INTERNAL", "fail"))
			return
		}
		writeJSON(w, 200, map[string]interface{}{"ok": true})
	}, 2)

	ctx := context.Background()
	var result map[string]interface{}
	err := client.client.request(ctx, requestOptions{method: "GET", path: "/test", auth: authNone}, &result)
	if err != nil {
		t.Fatal(err)
	}
	if atomic.LoadInt64(&count) != 3 {
		t.Errorf("expected 3 attempts, got %d", atomic.LoadInt64(&count))
	}
}

func TestClientNoRetryOn400(t *testing.T) {
	var count int64
	client, _ := setupTestClientRetries(t, func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&count, 1)
		writeJSON(w, 400, errBody("BAD", "bad"))
	}, 2)

	ctx := context.Background()
	err := client.client.request(ctx, requestOptions{method: "GET", path: "/test", auth: authNone}, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if atomic.LoadInt64(&count) != 1 {
		t.Errorf("expected 1 attempt, got %d", atomic.LoadInt64(&count))
	}
}

func TestClientAPIKeyHeader(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("X-API-Key"); got != "test-key" {
			t.Errorf("expected X-API-Key 'test-key', got '%s'", got)
		}
		writeJSON(w, 200, map[string]interface{}{})
	})
	ctx := context.Background()
	client.client.request(ctx, requestOptions{method: "GET", path: "/test", auth: authAPIKey}, nil)
}

func TestClientJWTHeader(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Errorf("expected 'Bearer test-token', got '%s'", got)
		}
		writeJSON(w, 200, map[string]interface{}{})
	})
	ctx := context.Background()
	client.client.request(ctx, requestOptions{method: "GET", path: "/test", auth: authJWT}, nil)
}

func TestClientNoAuthHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API-Key") != "" {
			t.Error("expected no X-API-Key header")
		}
		if r.Header.Get("Authorization") != "" {
			t.Error("expected no Authorization header")
		}
		writeJSON(w, 200, map[string]interface{}{})
	}))
	t.Cleanup(server.Close)
	c := New(WithBaseURL(server.URL), WithMaxRetries(0))
	ctx := context.Background()
	c.client.request(ctx, requestOptions{method: "GET", path: "/test", auth: authNone}, nil)
}

func TestClientQueryParams(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("days") != "30" {
			t.Errorf("expected days=30, got %s", r.URL.Query().Get("days"))
		}
		writeJSON(w, 200, map[string]interface{}{})
	})
	ctx := context.Background()
	client.client.request(ctx, requestOptions{
		method: "GET", path: "/test",
		query: map[string]string{"days": "30"},
		auth:  authNone,
	}, nil)
}
