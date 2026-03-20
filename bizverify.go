package bizverify

import (
	"net/http"
	"strings"
	"time"
)

// Option configures a Client.
type Option func(*Client)

// WithAPIKey sets the API key for authentication.
func WithAPIKey(key string) Option {
	return func(c *Client) { c.client.apiKey = key }
}

// WithBaseURL sets the base URL for API requests.
func WithBaseURL(url string) Option {
	return func(c *Client) { c.client.baseURL = strings.TrimRight(url, "/") }
}

// WithMaxRetries sets the maximum number of retries for failed requests.
func WithMaxRetries(n int) Option {
	return func(c *Client) { c.client.maxRetries = n }
}

// WithTimeout sets the request timeout.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) { c.client.client.Timeout = d }
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) { c.client.client = hc }
}

// Client is the main BizVerify API client.
type Client struct {
	client       *httpClient
	Auth         *AuthService
	Verification *VerificationService
	Entities     *EntitiesService
	Search       *SearchService
	Account      *AccountService
	Billing      *BillingService
	Checker      *CheckerService
	Config       *ConfigService
}

// LastResponseMeta returns metadata parsed from the most recent API response headers,
// including credit balance and rate limit information.
func (c *Client) LastResponseMeta() *ResponseMeta {
	return c.client.LastResponseMeta()
}

// New creates a new BizVerify client with the given options.
func New(opts ...Option) *Client {
	c := &Client{
		client: &httpClient{
			baseURL:    defaultBaseURL,
			client:     &http.Client{Timeout: defaultTimeout},
			maxRetries: defaultMaxRetries,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	c.Auth = &AuthService{client: c.client}
	c.Verification = &VerificationService{client: c.client}
	c.Entities = &EntitiesService{client: c.client}
	c.Search = &SearchService{client: c.client}
	c.Account = &AccountService{client: c.client}
	c.Billing = &BillingService{client: c.client}
	c.Checker = &CheckerService{client: c.client}
	c.Config = &ConfigService{client: c.client}

	return c
}
