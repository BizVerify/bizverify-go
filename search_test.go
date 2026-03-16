package bizverify

import (
	"context"
	"net/http"
	"sync/atomic"
	"testing"
)

func TestSearchFind(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/v1/search" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		writeJSON(w, 200, fixtureSearchResp)
	})
	resp, err := client.Search.Find(context.Background(), SearchParams{EntityName: "Acme"})
	if err != nil {
		t.Fatal(err)
	}
	if resp.Total != 1 {
		t.Errorf("wrong total: %d", resp.Total)
	}
	if resp.Results[0].EntityName != "Acme Inc" {
		t.Errorf("wrong name: %s", resp.Results[0].EntityName)
	}
}

func TestSearchFindAllSinglePage(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, fixtureSearchResp)
	})
	iter := client.Search.FindAll(context.Background(), SearchParams{EntityName: "Acme"})
	var results []SearchResult
	for iter.Next() {
		results = append(results, iter.Value())
	}
	if iter.Err() != nil {
		t.Fatal(iter.Err())
	}
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
}

func TestSearchFindAllMultiPage(t *testing.T) {
	var count int64
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt64(&count, 1)
		if n == 1 {
			writeJSON(w, 200, map[string]interface{}{
				"results": []interface{}{map[string]interface{}{
					"entity_name": "A", "jurisdiction": "us-fl", "entity_type": "llc",
					"status": "active", "jurisdiction_id": nil, "confidence": 0.9,
				}},
				"total": float64(2), "limit": float64(50), "offset": float64(0),
				"jurisdictions_searched": []interface{}{"us-fl"}, "jurisdictions_failed": []interface{}{},
				"credits_charged": float64(2),
			})
			return
		}
		writeJSON(w, 200, map[string]interface{}{
			"results": []interface{}{map[string]interface{}{
				"entity_name": "B", "jurisdiction": "us-fl", "entity_type": "llc",
				"status": "active", "jurisdiction_id": nil, "confidence": 0.8,
			}},
			"total": float64(2), "limit": float64(50), "offset": float64(1),
			"jurisdictions_searched": []interface{}{"us-fl"}, "jurisdictions_failed": []interface{}{},
			"credits_charged": float64(0),
		})
	})
	iter := client.Search.FindAll(context.Background(), SearchParams{EntityName: "Acme"})
	var results []SearchResult
	for iter.Next() {
		results = append(results, iter.Value())
	}
	if iter.Err() != nil {
		t.Fatal(iter.Err())
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
	if results[0].EntityName != "A" || results[1].EntityName != "B" {
		t.Error("wrong result order")
	}
}

func TestSearchFindAllEmpty(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]interface{}{
			"results": []interface{}{}, "total": float64(0), "limit": float64(50), "offset": float64(0),
			"jurisdictions_searched": []interface{}{}, "jurisdictions_failed": []interface{}{},
			"credits_charged": float64(0),
		})
	})
	iter := client.Search.FindAll(context.Background(), SearchParams{EntityName: "Nothing"})
	var results []SearchResult
	for iter.Next() {
		results = append(results, iter.Value())
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}
