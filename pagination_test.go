package bizverify

import (
	"context"
	"net/http"
	"testing"
)

func TestPaginationSinglePage(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, fixtureSearchResp)
	})
	iter := newSearchIterator(context.Background(), client.client, SearchParams{EntityName: "Acme"})
	count := 0
	for iter.Next() {
		_ = iter.Value()
		count++
	}
	if iter.Err() != nil {
		t.Fatal(iter.Err())
	}
	if count != 1 {
		t.Errorf("expected 1, got %d", count)
	}
	if iter.TotalResults() != 1 {
		t.Errorf("expected total 1, got %d", iter.TotalResults())
	}
}

func TestPaginationEmpty(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]interface{}{
			"results": []interface{}{}, "total": float64(0), "limit": float64(50), "offset": float64(0),
			"jurisdictions_searched": []interface{}{}, "jurisdictions_failed": []interface{}{},
			"credits_charged": float64(0),
		})
	})
	iter := newSearchIterator(context.Background(), client.client, SearchParams{EntityName: "Nothing"})
	if iter.Next() {
		t.Error("expected no results")
	}
	if iter.Err() != nil {
		t.Fatal(iter.Err())
	}
}
