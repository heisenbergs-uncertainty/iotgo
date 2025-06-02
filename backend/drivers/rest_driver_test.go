package drivers

import (
	"app/model"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
)

func TestRESTDriver_FetchData_URLConstruction(t *testing.T) {
	metadata := model.RESTMetadata{
		BaseEndpoint: "https://api.limblecmms.com:443/v2",
		Timeout:      5,
		Auth:         model.RESTAuth{Type: "none"},
	}
	metadataJSON, _ := json.Marshal(metadata)

	driver, err := NewRESTDriver(string(metadataJSON))
	if err != nil {
		t.Fatalf("Failed to create RESTDriver: %v", err)
	}

	resourceDetails := model.RESTResourceDetails{
		Method: "GET",
		Path:   "/assets",
	}
	detailsJSON, _ := json.Marshal(resourceDetails)

	// Mock HTTP client to capture the request URL
	var capturedURL string
	originalClient := driver.client
	driver.client = &http.Client{
		Transport: &mockTransport{
			RoundTripFunc: func(req *http.Request) (*http.Response, error) {
				capturedURL = req.URL.String()
				return nil, errors.New("mock error")
			},
		},
	}

	ctx := context.Background()
	_, err = driver.FetchData(ctx, string(detailsJSON))
	if err == nil {
		t.Error("Expected error due to mock transport, got nil")
	}

	expectedURL := "https://api.limblecmms.com:443/v2/assets"
	if capturedURL != expectedURL {
		t.Errorf("Expected URL %s, got %s", expectedURL, capturedURL)
	}

	// Restore original client
	driver.client = originalClient
}

type mockTransport struct {
	RoundTripFunc func(*http.Request) (*http.Response, error)
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req)
}
