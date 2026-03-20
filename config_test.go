package bizverify

import (
	"context"
	"net/http"
	"testing"
)

func TestConfigGet(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" || r.URL.Path != "/v1/config" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		if r.Header.Get("X-API-Key") != "" {
			t.Error("config endpoint should not send auth headers")
		}
		writeJSON(w, 200, fixtureConfigResp)
	})
	// Clear API key to verify authNone is used
	client.client.mu.Lock()
	client.client.apiKey = ""
	client.client.mu.Unlock()

	resp, err := client.Config.Get(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if resp.Jurisdictions == nil {
		t.Error("expected jurisdictions to be set")
	}
}

func TestConfigJurisdictions(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" || r.URL.Path != "/v1/jurisdictions" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		writeJSON(w, 200, fixtureJurisdictionsResp)
	})
	resp, err := client.Config.Jurisdictions(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.Jurisdictions) != 1 {
		t.Errorf("expected 1 jurisdiction, got %d", len(resp.Jurisdictions))
	}
	if resp.Jurisdictions[0].Code != "us-fl" {
		t.Errorf("wrong jurisdiction code: %s", resp.Jurisdictions[0].Code)
	}
	if resp.Jurisdictions[0].Name != "Florida" {
		t.Errorf("wrong jurisdiction name: %s", resp.Jurisdictions[0].Name)
	}
}
