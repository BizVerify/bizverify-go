package bizverify

import (
	"context"
	"strconv"
)

// BillingService handles billing and credit purchases.
type BillingService struct {
	client *httpClient
}

// Get retrieves billing information including balance, packages, and transactions.
func (s *BillingService) Get(ctx context.Context, params *BillingParams) (*BillingInfo, error) {
	query := map[string]string{}
	if params != nil {
		if params.Limit != nil {
			query["limit"] = strconv.Itoa(*params.Limit)
		}
		if params.Offset != nil {
			query["offset"] = strconv.Itoa(*params.Offset)
		}
	}
	var resp BillingInfo
	err := s.client.request(ctx, requestOptions{
		method: "GET",
		path:   "/v1/billing",
		query:  query,
		auth:   authAPIKey,
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// Purchase initiates a credit package purchase via Stripe.
func (s *BillingService) Purchase(ctx context.Context, packageID string) (*PurchaseResponse, error) {
	var resp PurchaseResponse
	err := s.client.request(ctx, requestOptions{
		method: "POST",
		path:   "/v1/billing/purchase",
		body:   map[string]string{"package_id": packageID},
		auth:   authAPIKey,
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
