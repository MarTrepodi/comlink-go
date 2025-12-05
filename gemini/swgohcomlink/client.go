package swgohcomlink

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// DefaultBaseURL is a placeholder. You should replace this with the actual base URL of the API.
const DefaultBaseURL = "https://api.example.com/swgoh-comlink"

// Client manages communication with the swgoh-comlink API.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a new API client with a default HTTP client and base URL.
func NewClient(baseURL string) *Client {
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

// request performs an API request and unmarshals the response into the provided interface.
func (c *Client) request(method, path string, requestBody interface{}, responseBody interface{}) error {
	var body io.Reader
	if requestBody != nil {
		buf := new(bytes.Buffer)
		if err := json.NewEncoder(buf).Encode(requestBody); err != nil {
			return fmt.Errorf("failed to encode request body: %w", err)
		}
		body = buf
	}

	url := c.BaseURL + path
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// The API uses POST for most data endpoints and expects JSON
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return c.handleErrorResponse(resp)
	}

	if responseBody != nil {
		if err := json.NewDecoder(resp.Body).Decode(responseBody); err != nil {
			return fmt.Errorf("failed to decode response body: %w", err)
		}
	}

	return nil
}

// handleErrorResponse handles non-2xx API responses based on the OpenAPI spec's ErrorResponse schema.
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("API Error (%s): %s", e.Code, e.Message)
}

func (c *Client) handleErrorResponse(resp *http.Response) error {
	apiError := ErrorResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&apiError); err != nil {
		// Fallback to a simple error if the body can't be decoded into ErrorResponse
		return fmt.Errorf("API returned status %d: failed to decode error response", resp.StatusCode)
	}
	return fmt.Errorf("API error: %s (Status: %d)", apiError.Message, resp.StatusCode)
}
