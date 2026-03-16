package bizverify

import (
	"context"
	"net/http"
	"testing"
)

func TestEntitiesGet(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/entity/ent_123" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		writeJSON(w, 200, fixtureEntityResp)
	})
	entity, err := client.Entities.Get(context.Background(), "ent_123")
	if err != nil {
		t.Fatal(err)
	}
	if entity.EntityName != "Acme Inc" {
		t.Errorf("wrong name: %s", entity.EntityName)
	}
	if entity.RegisteredAgent == nil || entity.RegisteredAgent.Name != "Agent Corp" {
		t.Error("wrong registered agent")
	}
}

func TestEntitiesHistoryNoParams(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/entity/ent_123/history" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		writeJSON(w, 200, fixtureHistoryResp)
	})
	resp, err := client.Entities.History(context.Background(), "ent_123", nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Total != 1 {
		t.Errorf("wrong total: %d", resp.Total)
	}
}

func TestEntitiesHistoryWithParams(t *testing.T) {
	client, _ := setupTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("limit") != "10" {
			t.Errorf("expected limit=10, got %s", r.URL.Query().Get("limit"))
		}
		if r.URL.Query().Get("offset") != "5" {
			t.Errorf("expected offset=5, got %s", r.URL.Query().Get("offset"))
		}
		writeJSON(w, 200, fixtureHistoryResp)
	})
	client.Entities.History(context.Background(), "ent_123", &HistoryParams{
		Limit: ptr(10), Offset: ptr(5),
	})
}
