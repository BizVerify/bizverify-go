package bizverify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const (
	defaultBaseURL    = "https://api.bizverify.co"
	defaultTimeout    = 30 * time.Second
	defaultMaxRetries = 2
)

type authMode int

const (
	authAPIKey authMode = iota
	authJWT
	authNone
)

type httpClient struct {
	baseURL    string
	client     *http.Client
	maxRetries int
	apiKey     string
	token      string
	mu         sync.RWMutex
}

type requestOptions struct {
	method string
	path   string
	body   interface{}
	query  map[string]string
	auth   authMode
}

func (c *httpClient) setToken(token string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.token = token
}

func (c *httpClient) setAPIKey(apiKey string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.apiKey = apiKey
}

func (c *httpClient) request(ctx context.Context, opts requestOptions, result interface{}) error {
	reqURL := c.buildURL(opts.path, opts.query)

	var bodyBytes []byte
	if opts.body != nil {
		var err error
		bodyBytes, err = json.Marshal(opts.body)
		if err != nil {
			return fmt.Errorf("bizverify: failed to marshal request body: %w", err)
		}
	}

	var lastErr error

	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			delay := time.Duration(1000*(1<<(attempt-1))) * time.Millisecond
			if delay > 4*time.Second {
				delay = 4 * time.Second
			}
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return &TimeoutError{APIError: &APIError{Message: ctx.Err().Error(), Code: "TIMEOUT", StatusCode: 0}}
			}
		}

		var bodyReader io.Reader
		if bodyBytes != nil {
			bodyReader = bytes.NewReader(bodyBytes)
		}

		req, err := http.NewRequestWithContext(ctx, opts.method, reqURL, bodyReader)
		if err != nil {
			return fmt.Errorf("bizverify: failed to create request: %w", err)
		}

		c.setHeaders(req, opts.auth, bodyBytes != nil)

		resp, err := c.client.Do(req)
		if err != nil {
			if ctx.Err() != nil {
				return &TimeoutError{APIError: &APIError{Message: "Request timed out", Code: "TIMEOUT", StatusCode: 0}}
			}
			lastErr = &APIError{Message: err.Error(), Code: "NETWORK_ERROR", StatusCode: 0}
			if attempt < c.maxRetries {
				continue
			}
			return lastErr
		}

		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = &APIError{Message: "failed to read response body", Code: "NETWORK_ERROR", StatusCode: 0}
			if attempt < c.maxRetries {
				continue
			}
			return lastErr
		}

		if resp.StatusCode == 204 {
			return nil
		}

		if resp.StatusCode >= 400 {
			apiErr := parseErrorResponse(resp.StatusCode, respBody, resp.Header)
			if resp.StatusCode >= 500 && attempt < c.maxRetries {
				lastErr = apiErr
				continue
			}
			return apiErr
		}

		if result != nil {
			if err := json.Unmarshal(respBody, result); err != nil {
				return fmt.Errorf("bizverify: failed to decode response: %w", err)
			}
		}
		return nil
	}

	if lastErr != nil {
		return lastErr
	}
	return &APIError{Message: "Request failed", Code: "UNKNOWN", StatusCode: 0}
}

func (c *httpClient) buildURL(path string, query map[string]string) string {
	u := c.baseURL + path
	if len(query) > 0 {
		params := url.Values{}
		for k, v := range query {
			params.Set(k, v)
		}
		u += "?" + params.Encode()
	}
	return u
}

func (c *httpClient) setHeaders(req *http.Request, auth authMode, hasBody bool) {
	req.Header.Set("Accept", "application/json")
	if hasBody {
		req.Header.Set("Content-Type", "application/json")
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	switch auth {
	case authAPIKey:
		if c.apiKey != "" {
			req.Header.Set("X-API-Key", c.apiKey)
		}
	case authJWT:
		if c.token != "" {
			req.Header.Set("Authorization", "Bearer "+c.token)
		}
	}
}
