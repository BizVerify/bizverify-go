package bizverify

import (
	"context"
	"strconv"
)

// AccountService manages user account operations.
type AccountService struct {
	client *httpClient
}

// Get retrieves the current user's account information.
func (s *AccountService) Get(ctx context.Context) (*Account, error) {
	var resp Account
	err := s.client.request(ctx, requestOptions{
		method: "GET",
		path:   "/v1/account",
		auth:   authJWT,
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// Usage retrieves usage statistics. Pass nil for default period.
func (s *AccountService) Usage(ctx context.Context, days *int) (*UsageStats, error) {
	query := map[string]string{}
	if days != nil {
		query["days"] = strconv.Itoa(*days)
	}
	var resp UsageStats
	err := s.client.request(ctx, requestOptions{
		method: "GET",
		path:   "/v1/account/usage",
		query:  query,
		auth:   authJWT,
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// DataExport exports all user data (GDPR).
func (s *AccountService) DataExport(ctx context.Context) (*DataExport, error) {
	var resp DataExport
	err := s.client.request(ctx, requestOptions{
		method: "GET",
		path:   "/v1/account/data-export",
		auth:   authJWT,
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateEmail updates the user's email address.
func (s *AccountService) UpdateEmail(ctx context.Context, email string) (*MessageResponse, error) {
	var resp MessageResponse
	err := s.client.request(ctx, requestOptions{
		method: "PATCH",
		path:   "/v1/account",
		body:   map[string]string{"email": email},
		auth:   authJWT,
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdatePassword changes the user's password.
func (s *AccountService) UpdatePassword(ctx context.Context, currentPassword, newPassword string) (*MessageResponse, error) {
	var resp MessageResponse
	err := s.client.request(ctx, requestOptions{
		method: "PUT",
		path:   "/v1/account/password",
		body:   map[string]interface{}{"current_password": currentPassword, "new_password": newPassword},
		auth:   authJWT,
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// Delete permanently deletes the user's account.
func (s *AccountService) Delete(ctx context.Context, password string) error {
	return s.client.request(ctx, requestOptions{
		method: "DELETE",
		path:   "/v1/account",
		body:   map[string]string{"password": password},
		auth:   authJWT,
	}, nil)
}

// CreateKey creates a new API key.
func (s *AccountService) CreateKey(ctx context.Context, label string) (*CreateKeyResponse, error) {
	var resp CreateKeyResponse
	err := s.client.request(ctx, requestOptions{
		method: "POST",
		path:   "/v1/account/keys",
		body:   map[string]string{"label": label},
		auth:   authJWT,
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// RevokeKey revokes an existing API key.
func (s *AccountService) RevokeKey(ctx context.Context, keyID string) error {
	return s.client.request(ctx, requestOptions{
		method: "DELETE",
		path:   "/v1/account/keys/" + keyID,
		auth:   authJWT,
	}, nil)
}
