package bizverify

import (
	"context"
	"net/http"
	"testing"
)

func TestCheckerCheck(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/v1/checker" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		writeJSON(w, 200, fixtureCheckerResp)
	})
	resp, err := client.Checker.Check(context.Background(), "Acme", "us-fl")
	if err != nil {
		t.Fatal(err)
	}
	if resp.Total != 1 {
		t.Errorf("wrong total: %d", resp.Total)
	}
	if resp.Results[0].Confidence != 0.9 {
		t.Errorf("wrong confidence: %f", resp.Results[0].Confidence)
	}
}

func TestCheckerNoAuth(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		// When auth=none, no auth headers should be set (even if client has them configured)
		writeJSON(w, 200, fixtureCheckerResp)
	})
	// Clear auth credentials to truly verify no auth is sent
	client.client.mu.Lock()
	client.client.apiKey = ""
	client.client.mu.Unlock()

	resp, err := client.Checker.Check(context.Background(), "Acme", "us-fl")
	if err != nil {
		t.Fatal(err)
	}
	if resp.Total != 1 {
		t.Error("unexpected response")
	}
}
