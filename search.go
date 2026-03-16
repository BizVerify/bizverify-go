package bizverify

import "context"

// SearchService searches for business entities.
type SearchService struct {
	client *httpClient
}

// Find performs a search with pagination support.
func (s *SearchService) Find(ctx context.Context, params SearchParams) (*SearchResponse, error) {
	var resp SearchResponse
	err := s.client.request(ctx, requestOptions{
		method: "POST",
		path:   "/v1/search",
		body:   params,
		auth:   authAPIKey,
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// FindAll returns an iterator that auto-paginates through all search results.
func (s *SearchService) FindAll(ctx context.Context, params SearchParams) *SearchIterator {
	return newSearchIterator(ctx, s.client, params)
}
