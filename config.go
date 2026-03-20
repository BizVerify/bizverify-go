package bizverify

import "context"

// ConfigService provides access to public API configuration.
type ConfigService struct {
	client *httpClient
}

// Get retrieves the full API configuration including jurisdictions, pricing, and features.
func (s *ConfigService) Get(ctx context.Context) (*ConfigResponse, error) {
	var resp ConfigResponse
	if err := s.client.request(ctx, requestOptions{
		method: "GET",
		path:   "/v1/config",
		auth:   authNone,
	}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Jurisdictions retrieves the list of available jurisdictions.
func (s *ConfigService) Jurisdictions(ctx context.Context) (*JurisdictionsResponse, error) {
	var resp JurisdictionsResponse
	if err := s.client.request(ctx, requestOptions{
		method: "GET",
		path:   "/v1/jurisdictions",
		auth:   authNone,
	}, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
