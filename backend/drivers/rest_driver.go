package drivers

import (
	"app/model"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/beego/beego/logs"
)

// RESTDriver implements the PlatformDriver interface for REST-based platforms.
type RESTDriver struct {
	BaseURL string
	Auth    model.RESTAuth
	Timeout time.Duration
	client  *http.Client
}

// NewRESTDriver creates a new RESTDriver instance from platform metadata.
func NewRESTDriver(metadata string) (*RESTDriver, error) {
	var config model.RESTMetadata
	if err := json.Unmarshal([]byte(metadata), &config); err != nil {
		return nil, fmt.Errorf("invalid metadata: %w", err)
	}
	if config.BaseEndpoint == "" {
		return nil, errors.New("missing base_endpoint in metadata")
	}
	if config.Timeout == 0 {
		config.Timeout = 10 // Default timeout
	}
	return &RESTDriver{
		BaseURL: config.BaseEndpoint,
		Auth:    config.Auth,
		Timeout: time.Duration(config.Timeout) * time.Second,
		client:  &http.Client{Timeout: time.Duration(config.Timeout) * time.Second},
	}, nil
}

// Connect is a no-op for REST since it's stateless.
func (d *RESTDriver) Connect(ctx context.Context) error {
	return nil
}

// FetchData sends an HTTP request to the REST endpoint and returns the response.
func (d *RESTDriver) FetchData(ctx context.Context, resourceDetails string) (interface{}, error) {
	var details model.RESTResourceDetails
	if err := json.Unmarshal([]byte(resourceDetails), &details); err != nil {
		return nil, fmt.Errorf("invalid resource details: %w", err)
	}

	if details.Method == "" || details.Path == "" {
		return nil, errors.New("method and path are required in resource details")
	}

	// Construct URL with query parameters
	baseURL, err := url.Parse(d.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL %s: %w", d.BaseURL, err)
	}

	// Construct the full path by joining base URL path and resource path
	// Use path.Join to ensure proper path concatenation
	resourcePath := strings.TrimLeft(details.Path, "/") // Remove leading slash for clean joining
	fullPath := path.Join(baseURL.Path, resourcePath)
	if !strings.HasPrefix(fullPath, "/") {
		fullPath = "/" + fullPath // Ensure path starts with /
	}

	// Create a new URL with the combined path
	fullURL := *baseURL
	fullURL.Path = fullPath

	// Add query parameters if any
	query := fullURL.Query()
	for key, value := range details.QueryParams {
		query.Set(key, value)
	}
	fullURL.RawQuery = query.Encode()

	// Create request
	req, err := http.NewRequestWithContext(ctx, details.Method, fullURL.String(), bytes.NewBufferString(details.Body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set authentication headers
	switch d.Auth.Type {
	case "api_key":
		if d.Auth.APIKey == nil {
			return nil, errors.New("API key missing for api_key authentication")
		}
		req.Header.Set("X-API-Key", *d.Auth.APIKey)
		logs.Debug("Setting API key header")
	case "bearer":
		if d.Auth.BearerToken == nil {
			return nil, errors.New("bearer token missing for bearer authentication")
		}
		req.Header.Set("Authorization", "Bearer "+*d.Auth.BearerToken)
		logs.Debug("Setting Bearer token header")
	case "basic":
		if d.Auth.BasicAuth == nil {
			return nil, errors.New("basic auth configuration missing")
		}
		auth := base64.StdEncoding.EncodeToString([]byte(d.Auth.BasicAuth.Username + ":" + d.Auth.BasicAuth.Password))
		req.Header.Set("Authorization", "Basic "+auth)
		logs.Debug("Setting Basic auth header")
	}

	// Set custom headers
	for key, value := range details.Headers {
		req.Header.Set(key, value)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Log request details
	logs.Info("REST request: method=%s, url=%s, headers=%v, body=%s", details.Method, fullURL.String(), req.Header, details.Body)
	resp, err := d.client.Do(req)
	if err != nil {
		logs.Error("REST request failed: %v", err)
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logs.Error("Failed to read response body: %v", err)
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Prepare structured response
	responseData := map[string]interface{}{
		"status_code": resp.StatusCode,
		"status":      resp.Status,
		"headers":     resp.Header,
	}

	// Parse body based on Content-Type
	contentType := strings.ToLower(resp.Header.Get("Content-Type"))
	if strings.Contains(contentType, "application/json") {
		var bodyData interface{}
		if len(bodyBytes) > 0 {
			if err := json.Unmarshal(bodyBytes, &bodyData); err != nil {
				logs.Error("Failed to decode JSON response: %v", err)
				return nil, fmt.Errorf("failed to decode JSON response: %w", err)
			}
		}
		responseData["body"] = bodyData
	} else {
		responseData["body"] = string(bodyBytes)
	}

	if resp.StatusCode >= 400 {
		logs.Error("REST request returned error: status=%s, body=%s", resp.Status, string(bodyBytes))
		return responseData, fmt.Errorf("HTTP error: %s", resp.Status)
	}

	logs.Info("REST request successful: status=%s", resp.Status)
	return responseData, nil
}

// ValidateConfig checks if the REST configuration is valid by sending a HEAD request.
func (d *RESTDriver) ValidateConfig(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "HEAD", d.BaseURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create HEAD request: %w", err)
	}

	// Set authentication headers
	switch d.Auth.Type {
	case "api_key":
		req.Header.Set("X-API-Key", *d.Auth.APIKey)
	case "bearer":
		req.Header.Set("Authorization", "Bearer "+*d.Auth.BearerToken)
	case "basic":
		if d.Auth.BasicAuth == nil {
			return errors.New("basic auth configuration missing")
		}
		auth := base64.StdEncoding.EncodeToString([]byte(d.Auth.BasicAuth.Username + ":" + d.Auth.BasicAuth.Password))
		req.Header.Set("Authorization", "Basic "+auth)
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("validation failed with status: %s", resp.Status)
	}

	return nil
}

// TestResource tests a specific REST resource by executing the request.
func (d *RESTDriver) TestResource(ctx context.Context, resourceDetails string) (interface{}, error) {
	return d.FetchData(ctx, resourceDetails)
}

// Disconnect is a no-op for REST since it's stateless.
func (d *RESTDriver) Disconnect(ctx context.Context) error {
	return nil
}
