package bizverify

import (
	"context"
	"net/http"
	"testing"
)

func TestBillingGet(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/billing" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		writeJSON(w, 200, fixtureBillingResp)
	})
	resp, err := client.Billing.Get(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Balance != 100 {
		t.Errorf("wrong balance: %d", resp.Balance)
	}
}

func TestBillingGetWithParams(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("limit") != "10" {
			t.Errorf("expected limit=10, got %s", r.URL.Query().Get("limit"))
		}
		writeJSON(w, 200, fixtureBillingResp)
	})
	client.Billing.Get(context.Background(), &BillingParams{Limit: ptr(10)})
}

func TestBillingPurchase(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/v1/billing/purchase" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		writeJSON(w, 200, fixturePurchaseResp)
	})
	resp, err := client.Billing.Purchase(context.Background(), "pkg_100")
	if err != nil {
		t.Fatal(err)
	}
	if resp.SessionID != "cs_test_123" {
		t.Errorf("wrong session ID: %s", resp.SessionID)
	}
}
