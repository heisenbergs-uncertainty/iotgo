package drivers

import (
	"context"
	"errors"
	"fmt"
)

// PlatformDriver defines the interface for all platform drivers.
type PlatformDriver interface {
	// Connect establishes a connection to the platform.
	Connect(ctx context.Context) error

	// FetchData retrieves data from the platform for a specific resource.
	FetchData(ctx context.Context, resourceDetails string) (interface{}, error)

	// Disconnect closes the connection to the platform.
	Disconnect(ctx context.Context) error

	// ValidateConfig checks if the platform configuration is valid.
	ValidateConfig(ctx context.Context) error

	// TestResource tests a specific resource configuration.
	TestResource(ctx context.Context, resourceDetails string) (interface{}, error)
}

// ErrNotImplemented is returned when a driver does not implement a method.
var ErrNotImplemented = errors.New("method not implemented for this platform type")

// GetDriver returns the appropriate PlatformDriver based on the platform type and metadata.
func GetDriver(platformType string, metadata string) (PlatformDriver, error) {
	switch platformType {
	case "REST":
		return NewRESTDriver(metadata)
	// case "OPCUA":
	// 	return NewOPCUADriver(metadata)
	case "InfluxDB":
		return NewInfluxDBDriver(metadata)
	default:
		return nil, fmt.Errorf("unsupported platform type: %s", platformType)
	}
}
