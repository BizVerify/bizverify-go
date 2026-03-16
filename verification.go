package bizverify

import "context"

// VerificationService handles business entity verification.
type VerificationService struct {
	client *httpClient
}

// Verify initiates a business entity verification.
func (s *VerificationService) Verify(ctx context.Context, params VerifyParams) (*VerifyResponse, error) {
	var resp VerifyResponse
	err := s.client.request(ctx, requestOptions{
		method: "POST",
		path:   "/v1/verify",
		body:   params,
		auth:   authAPIKey,
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// VerifyAndWait initiates verification and polls until completion.
func (s *VerificationService) VerifyAndWait(ctx context.Context, params VerifyParams, opts *PollOptions) (*JobStatusResponse, error) {
	resp, err := s.Verify(ctx, params)
	if err != nil {
		return nil, err
	}
	return pollUntilComplete(ctx, s.client, resp, opts)
}

// GetStatus retrieves the status of a verification job.
func (s *VerificationService) GetStatus(ctx context.Context, jobID string) (*JobStatusResponse, error) {
	var resp JobStatusResponse
	err := s.client.request(ctx, requestOptions{
		method: "GET",
		path:   "/v1/verify/status/" + jobID,
		auth:   authAPIKey,
	}, &resp)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
