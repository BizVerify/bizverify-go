package bizverify

import "context"

// AuthService handles authentication operations.
type AuthService struct {
	client *httpClient
}

// RequestAccess initiates the passwordless auth flow by sending an OTP to the given email.
func (s *AuthService) RequestAccess(ctx context.Context, email string, acceptTerms bool) (*RequestAccessResponse, error) {
	var resp RequestAccessResponse
	err := s.client.request(ctx, requestOptions{
		method: "POST",
		path:   "/v1/auth/request-access",
		body:   RequestAccessParams{Email: email, AcceptTerms: acceptTerms},
		auth:   authNone,
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// VerifyAccess completes the passwordless auth flow by verifying the OTP code.
// On success, the client is automatically configured with the returned API key.
func (s *AuthService) VerifyAccess(ctx context.Context, email string, code string, label *string) (*VerifyAccessResponse, error) {
	var resp VerifyAccessResponse
	err := s.client.request(ctx, requestOptions{
		method: "POST",
		path:   "/v1/auth/verify-access",
		body:   VerifyAccessParams{Email: email, Code: code, Label: label},
		auth:   authNone,
	}, &resp)
	if err != nil {
		return nil, err
	}
	s.client.setAPIKey(resp.APIKey)
	return &resp, nil
}
