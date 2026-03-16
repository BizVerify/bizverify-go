package bizverify

import (
	"context"
	"net/http"
	"testing"
)

func TestAccountGet(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/account" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Error("expected JWT auth")
		}
		writeJSON(w, 200, fixtureAccountResp)
	})
	account, err := client.Account.Get(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if account.CreditBalance != 100 {
		t.Errorf("wrong balance: %d", account.CreditBalance)
	}
}

func TestAccountUsageWithDays(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("days") != "7" {
			t.Errorf("expected days=7, got %s", r.URL.Query().Get("days"))
		}
		writeJSON(w, 200, fixtureUsageResp)
	})
	days := 7
	usage, err := client.Account.Usage(context.Background(), &days)
	if err != nil {
		t.Fatal(err)
	}
	if usage.PeriodDays != 30 {
		t.Errorf("wrong period: %d", usage.PeriodDays)
	}
}

func TestAccountUsageNoDays(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, fixtureUsageResp)
	})
	_, err := client.Account.Usage(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAccountDataExport(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, fixtureDataExportResp)
	})
	export, err := client.Account.DataExport(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if export.Profile.Email != "test@example.com" {
		t.Errorf("wrong email: %s", export.Profile.Email)
	}
}

func TestAccountUpdateEmail(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		writeJSON(w, 200, fixtureMsgResp)
	})
	resp, err := client.Account.UpdateEmail(context.Background(), "new@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if resp.Message != "Success" {
		t.Errorf("wrong message: %s", resp.Message)
	}
}

func TestAccountUpdatePassword(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		writeJSON(w, 200, fixtureMsgResp)
	})
	resp, err := client.Account.UpdatePassword(context.Background(), "old", "new")
	if err != nil {
		t.Fatal(err)
	}
	if resp.Message != "Success" {
		t.Errorf("wrong message: %s", resp.Message)
	}
}

func TestAccountDelete(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(204)
	})
	err := client.Account.Delete(context.Background(), "mypassword")
	if err != nil {
		t.Fatal(err)
	}
}

func TestAccountCreateKey(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, fixtureCreateKeyResp)
	})
	resp, err := client.Account.CreateKey(context.Background(), "My Key")
	if err != nil {
		t.Fatal(err)
	}
	if resp.Key != "bv_live_newkey" {
		t.Errorf("wrong key: %s", resp.Key)
	}
}

func TestAccountRevokeKey(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/account/keys/key_123" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.WriteHeader(204)
	})
	err := client.Account.RevokeKey(context.Background(), "key_123")
	if err != nil {
		t.Fatal(err)
	}
}
