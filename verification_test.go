package bizverify

import (
	"context"
	"net/http"
	"testing"
)

func TestVerify(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/v1/verify" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		writeJSON(w, 200, fixtureVerifySyncResp)
	})
	resp, err := client.Verification.Verify(context.Background(), VerifyParams{
		EntityName: "Acme Inc", Jurisdiction: "us-fl",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Status != "completed" {
		t.Errorf("wrong status: %s", resp.Status)
	}
	if !resp.Cached {
		t.Error("expected cached=true")
	}
}

func TestVerifyAsync(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, fixtureVerifyAsyncResp)
	})
	resp, err := client.Verification.Verify(context.Background(), VerifyParams{
		EntityName: "Acme Inc", Jurisdiction: "us-fl",
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Status != "pending" {
		t.Errorf("wrong status: %s", resp.Status)
	}
	if resp.JobID == nil || *resp.JobID != "job_456" {
		t.Error("wrong job_id")
	}
}

func TestGetStatus(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/verify/status/job_456" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		writeJSON(w, 200, fixtureJobCompletedResp)
	})
	resp, err := client.Verification.GetStatus(context.Background(), "job_456")
	if err != nil {
		t.Fatal(err)
	}
	if resp.Status != "completed" {
		t.Errorf("wrong status: %s", resp.Status)
	}
}

func TestVerifyUsesAPIKey(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API-Key") != "test-key" {
			t.Error("expected API key header")
		}
		writeJSON(w, 200, fixtureVerifySyncResp)
	})
	client.Verification.Verify(context.Background(), VerifyParams{
		EntityName: "Acme", Jurisdiction: "us-fl",
	})
}
