package drivers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// RESTDriver implements the PlatformDriver interface for REST-based platforms.
type RESTDriver struct {
	BaseURL   string
	AuthType  string
	AuthToken string
	Timeout   time.Duration
	client    *http.Client
}

// NewRESTDriver creates a new RESTDriver instance from platform metadata.
func NewRESTDriver(metadata string) (*RESTDriver, error) {
	var config struct {
		BaseURL   string `json:"base_url"`
		AuthType  string `json:"auth_type"` // e.g., "none", "api_key", "bearer"
		AuthToken string `json:"auth_token"`
		Timeout   int    `json:"timeout"`
	}
	if err := json.Unmarshal([]byte(metadata), &config); err != nil {
		return nil, fmt.Errorf("invalid metadata: %w", err)
	}
	if config.BaseURL == "" {
		return nil, errors.New("missing base_url in metadata")
	}
	return &RESTDriver{
		BaseURL:   config.BaseURL,
		AuthType:  config.AuthType,
		AuthToken: config.AuthToken,
		Timeout:   time.Duration(config.Timeout) * time.Second,
		client:    &http.Client{Timeout: time.Duration(config.Timeout) * time.Second},
	}, nil
}

// Connect is a no-op for REST since it's stateless.
func (d *RESTDriver) Connect(ctx context.Context) error {
	return nil
}

// FetchData sends an HTTP request to the REST endpoint and returns the response.
func (d *RESTDriver) FetchData(ctx context.Context, resourceDetails string) (interface{}, error) {
	var details struct {
		Path   string `json:"path"`
		Method string `json:"method"`
		Body   string `json:"body"`
	}
	if err := json.Unmarshal([]byte(resourceDetails), &details); err != nil {
		return nil, fmt.Errorf("invalid resource details: %w", err)
	}

	url := d.BaseURL + details.Path
	req, err := http.NewRequestWithContext(ctx, details.Method, url, bytes.NewBufferString(details.Body))
	if err != nil {
		return nil, err
	}

	// Set authentication headers
	switch d.AuthType {
	case "api_key":
		req.Header.Set("X-API-Key", d.AuthToken)
	case "bearer":
		req.Header.Set("Authorization", "Bearer "+d.AuthToken)
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP error: %s", resp.Status)
	}

	var data interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return data, nil
}

// Disconnect is a no-op for REST since it's stateless.
func (d *RESTDriver) Disconnect(ctx context.Context) error {
	return nil
}
