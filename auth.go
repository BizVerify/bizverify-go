package bizverify

import "context"

// AuthService handles authentication operations.
type AuthService struct {
	client *httpClient
}

// Register creates a new user account.
func (s *AuthService) Register(ctx context.Context, email, password string, acceptTerms bool) (*RegisterResponse, error) {
	var resp RegisterResponse
	err := s.client.request(ctx, requestOptions{
		method: "POST",
		path:   "/v1/auth/register",
		body:   map[string]interface{}{"email": email, "password": password, "accept_terms": acceptTerms},
		auth:   authNone,
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// Login authenticates a user and stores the JWT token.
func (s *AuthService) Login(ctx context.Context, email, password string) (*LoginResponse, error) {
	var resp LoginResponse
	err := s.client.request(ctx, requestOptions{
		method: "POST",
		path:   "/v1/auth/login",
		body:   map[string]interface{}{"email": email, "password": password},
		auth:   authNone,
	}, &resp)
	if err != nil {
		return nil, err
	}
	s.client.setToken(resp.Token)
	return &resp, nil
}

// VerifyEmail confirms a user's email address with a verification code.
func (s *AuthService) VerifyEmail(ctx context.Context, email, code string) (*MessageResponse, error) {
	var resp MessageResponse
	err := s.client.request(ctx, requestOptions{
		method: "POST",
		path:   "/v1/auth/verify-email",
		body:   map[string]interface{}{"email": email, "code": code},
		auth:   authNone,
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// ResendVerification resends the email verification code.
func (s *AuthService) ResendVerification(ctx context.Context, email string) (*MessageResponse, error) {
	var resp MessageResponse
	err := s.client.request(ctx, requestOptions{
		method: "POST",
		path:   "/v1/auth/resend-verification",
		body:   map[string]interface{}{"email": email},
		auth:   authNone,
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// ForgotPassword initiates the password reset flow.
func (s *AuthService) ForgotPassword(ctx context.Context, email string) (*MessageResponse, error) {
	var resp MessageResponse
	err := s.client.request(ctx, requestOptions{
		method: "POST",
		path:   "/v1/auth/forgot-password",
		body:   map[string]interface{}{"email": email},
		auth:   authNone,
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// ResetPassword sets a new password using a reset code.
func (s *AuthService) ResetPassword(ctx context.Context, email, code, newPassword string) (*MessageResponse, error) {
	var resp MessageResponse
	err := s.client.request(ctx, requestOptions{
		method: "POST",
		path:   "/v1/auth/reset-password",
		body:   map[string]interface{}{"email": email, "code": code, "new_password": newPassword},
		auth:   authNone,
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
