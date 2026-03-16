package bizverify

import "context"

// CheckerService provides a simple entity checker (no auth required).
type CheckerService struct {
	client *httpClient
}

// Check performs a quick entity check without authentication.
func (s *CheckerService) Check(ctx context.Context, entityName, jurisdiction string) (*CheckerResponse, error) {
	var resp CheckerResponse
	err := s.client.request(ctx, requestOptions{
		method: "POST",
		path:   "/v1/checker",
		body:   map[string]string{"entity_name": entityName, "jurisdiction": jurisdiction},
		auth:   authNone,
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
