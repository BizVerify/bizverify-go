package bizverify

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestClient(t *testing.T, handler http.HandlerFunc) (*Client, *httptest.Server) {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	client := New(
		WithBaseURL(server.URL),
		WithAPIKey("test-key"),
		WithMaxRetries(0),
	)
	return client, server
}

func setupTestClientRetries(t *testing.T, handler http.HandlerFunc, retries int) (*Client, *httptest.Server) {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)
	client := New(
		WithBaseURL(server.URL),
		WithAPIKey("test-key"),
		WithMaxRetries(retries),
	)
	return client, server
}

func writeJSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func writeJSONWithMeta(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Credits-Remaining", "95")
	w.Header().Set("X-Credits-Charged", "5")
	w.Header().Set("X-Ratelimit-Limit", "100")
	w.Header().Set("X-Ratelimit-Remaining", "99")
	w.Header().Set("X-Ratelimit-Reset", "1700000000")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func errBody(code, message string) map[string]interface{} {
	return map[string]interface{}{
		"error": map[string]interface{}{
			"code":    code,
			"message": message,
		},
	}
}

func ptr[T any](v T) *T { return &v }

var fixtureRequestAccessResp = map[string]interface{}{
	"message": "Verification code sent to your email",
}

var fixtureVerifyAccessResp = map[string]interface{}{
	"api_key": "bv_live_abc123def456",
	"key_id":  "key_123",
	"label":   "my-key",
}

var fixtureMsgResp = map[string]interface{}{"message": "Success"}

var fixtureVerifySyncResp = map[string]interface{}{
	"status": "completed", "data": map[string]interface{}{"exists": true},
	"entity_id": "ent_123", "cached": true, "credits_charged": float64(1),
}

var fixtureVerifyAsyncResp = map[string]interface{}{
	"status": "pending", "job_id": "job_456", "cached": false, "credits_charged": float64(15),
}

var fixtureJobPendingResp = map[string]interface{}{
	"id": "job_456", "status": "pending", "jurisdiction": "us-fl", "query": "Acme Inc",
	"verification_level": "full", "credits_charged": float64(15),
	"created_at": "2026-01-01T00:00:00.000Z",
}

var fixtureJobCompletedResp = map[string]interface{}{
	"id": "job_456", "status": "completed", "jurisdiction": "us-fl", "query": "Acme Inc",
	"verification_level": "full", "credits_charged": float64(15),
	"result": map[string]interface{}{"exists": true},
	"created_at": "2026-01-01T00:00:00.000Z", "completed_at": "2026-01-01T00:01:00.000Z",
}

var fixtureJobFailedResp = map[string]interface{}{
	"id": "job_456", "status": "failed", "jurisdiction": "us-fl", "query": "Acme Inc",
	"verification_level": "full", "credits_charged": float64(15), "error": "Scraper timeout",
	"created_at": "2026-01-01T00:00:00.000Z", "completed_at": "2026-01-01T00:01:00.000Z",
}

var fixtureEntityResp = map[string]interface{}{
	"id": "ent_123", "entity_name": "Acme Inc", "jurisdiction": "us-fl",
	"entity_type": "corporation", "status": "active", "jurisdiction_id": "FL12345",
	"good_standing": true, "formation_date": "2020-01-15",
	"registered_agent": map[string]interface{}{"name": "Agent Corp", "address": nil},
	"officers": []interface{}{}, "principal_address": nil, "filing_history_summary": []interface{}{},
	"created_at": "2026-01-01T00:00:00.000Z", "updated_at": "2026-01-01T00:00:00.000Z",
}

var fixtureHistoryResp = map[string]interface{}{
	"snapshots": []interface{}{map[string]interface{}{"id": "snap_1"}},
	"total": float64(1), "limit": float64(50), "offset": float64(0),
}

var fixtureSearchResp = map[string]interface{}{
	"results": []interface{}{map[string]interface{}{
		"entity_name": "Acme Inc", "jurisdiction": "us-fl", "entity_type": "corporation",
		"status": "active", "jurisdiction_id": "FL12345", "confidence": 0.95,
	}},
	"total": float64(1), "limit": float64(50), "offset": float64(0),
	"jurisdictions_searched": []interface{}{"us-fl"}, "jurisdictions_failed": []interface{}{},
	"credits_charged": float64(2),
}

var fixtureAccountResp = map[string]interface{}{
	"id": "550e8400-e29b-41d4-a716-446655440000", "email": "test@example.com",
	"email_verified": true, "plan": "free", "credit_balance": float64(100),
	"api_keys": []interface{}{}, "created_at": "2026-01-01T00:00:00.000Z",
}

var fixtureUsageResp = map[string]interface{}{
	"period_days": float64(30), "daily": []interface{}{},
	"by_endpoint": []interface{}{}, "by_jurisdiction": []interface{}{},
}

var fixtureDataExportResp = map[string]interface{}{
	"profile": map[string]interface{}{
		"id": "550e8400-e29b-41d4-a716-446655440000", "email": "test@example.com",
		"email_verified": true, "plan": "free", "credit_balance": float64(100),
		"terms_accepted_at": nil, "terms_version": nil, "created_at": "2026-01-01T00:00:00.000Z",
	},
	"api_keys": []interface{}{}, "credit_transactions": []interface{}{},
	"verification_jobs": []interface{}{}, "usage_stats": []interface{}{},
}

var fixtureCreateKeyResp = map[string]interface{}{
	"id": "key_123", "key": "bv_live_newkey", "prefix": "bv_live_new",
	"label": "My Key", "message": "Store this API key securely. It will not be shown again.",
}

var fixtureBillingResp = map[string]interface{}{
	"balance": float64(100), "packages": []interface{}{}, "transactions": []interface{}{},
}

var fixturePurchaseResp = map[string]interface{}{
	"session_id": "cs_test_123", "url": "https://checkout.stripe.com/session/cs_test_123",
}

var fixtureCheckerResp = map[string]interface{}{
	"results": []interface{}{map[string]interface{}{
		"entity_name": "Acme Inc", "entity_type": "corporation", "status": "active",
		"jurisdiction": "us-fl", "confidence": 0.9,
	}},
	"query": "Acme", "jurisdiction": "us-fl", "total": float64(1),
}

var fixtureConfigResp = map[string]interface{}{
	"jurisdictions": map[string]interface{}{"us-fl": map[string]interface{}{"name": "Florida"}},
	"checker":       map[string]interface{}{},
	"pricing":       map[string]interface{}{},
	"features":      map[string]interface{}{},
	"rateLimits":    map[string]interface{}{},
	"status":        map[string]interface{}{},
	"legal":         map[string]interface{}{},
	"docs":          map[string]interface{}{},
}

var fixtureJurisdictionsResp = map[string]interface{}{
	"jurisdictions": []interface{}{
		map[string]interface{}{
			"code": "us-fl", "name": "Florida",
			"features": map[string]interface{}{"search": true, "details": true},
		},
	},
}
