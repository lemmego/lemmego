package framework

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// HTTPClient represents an HTTP API client.
type HTTPClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewHTTPClient creates a new instance of HTTPClient.
func NewHTTPClient(baseURL string) *HTTPClient {
	return &HTTPClient{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{
			// You can customize the HTTP client settings here.
		},
	}
}

// request performs an HTTP request and returns the response.
func (c *HTTPClient) request(ctx context.Context, method, path string, body io.Reader, headers map[string]string) (*http.Response, error) {
	url := c.BaseURL + path

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Perform the request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Get performs an HTTP GET request.
func (c *HTTPClient) Get(ctx context.Context, path string, headers map[string]string) (*http.Response, error) {
	return c.request(ctx, http.MethodGet, path, nil, headers)
}

// Post performs an HTTP POST request.
func (c *HTTPClient) Post(ctx context.Context, path string, body io.Reader, headers map[string]string) (*http.Response, error) {
	return c.request(ctx, http.MethodPost, path, body, headers)
}

// JSONRequest performs an HTTP request with JSON payload and parses the JSON response.
func (c *HTTPClient) JSONRequest(ctx context.Context, method, path string, requestBody interface{}, responseBody interface{}, headers map[string]string) error {
	var body io.Reader

	if requestBody != nil {
		jsonBody, err := json.Marshal(requestBody)
		if err != nil {
			return err
		}
		body = strings.NewReader(string(jsonBody))
	}

	resp, err := c.request(ctx, method, path, body, headers)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	if responseBody != nil {
		err := json.NewDecoder(resp.Body).Decode(responseBody)
		if err != nil {
			return err
		}
	}

	return nil
}
