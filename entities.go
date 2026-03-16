package bizverify

import (
	"context"
	"strconv"
)

// EntitiesService retrieves entity information.
type EntitiesService struct {
	client *httpClient
}

// Get retrieves a single entity by ID.
func (s *EntitiesService) Get(ctx context.Context, entityID string) (*Entity, error) {
	var resp Entity
	err := s.client.request(ctx, requestOptions{
		method: "GET",
		path:   "/v1/entity/" + entityID,
		auth:   authAPIKey,
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// History retrieves historical snapshots for an entity.
func (s *EntitiesService) History(ctx context.Context, entityID string, params *HistoryParams) (*PaginatedSnapshots, error) {
	query := map[string]string{}
	if params != nil {
		if params.Limit != nil {
			query["limit"] = strconv.Itoa(*params.Limit)
		}
		if params.Offset != nil {
			query["offset"] = strconv.Itoa(*params.Offset)
		}
	}
	var resp PaginatedSnapshots
	err := s.client.request(ctx, requestOptions{
		method: "GET",
		path:   "/v1/entity/" + entityID + "/history",
		query:  query,
		auth:   authAPIKey,
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
