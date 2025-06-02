package drivers

import (
	"app/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

// InfluxDBDriver implements the PlatformDriver interface for InfluxDB platforms.
type InfluxDBDriver struct {
	client influxdb2.Client
	org    string
}

// NewInfluxDBDriver creates a new InfluxDBDriver instance from platform metadata.
func NewInfluxDBDriver(metadata string) (*InfluxDBDriver, error) {
	var config model.InfluxDBMetadata
	if err := json.Unmarshal([]byte(metadata), &config); err != nil {
		return nil, fmt.Errorf("invalid metadata: %w", err)
	}
	if config.URL == "" {
		return nil, errors.New("missing url in metadata")
	}
	if config.Token == "" {
		return nil, errors.New("missing token in metadata")
	}
	if config.Org == "" {
		return nil, errors.New("missing org in metadata")
	}
	if config.Bucket == "" {
		return nil, errors.New("missing bucket in metadata")
	}
	if config.Timeout == 0 {
		config.Timeout = 10
	}

	// Create InfluxDB client with timeout
	client := influxdb2.NewClientWithOptions(config.URL, config.Token, influxdb2.DefaultOptions().
		SetHTTPRequestTimeout(uint(config.Timeout)))

	return &InfluxDBDriver{
		client: client,
		org:    config.Org,
	}, nil
}

// ValidateConfig checks if the InfluxDB configuration is valid by performing a health check.
func (d *InfluxDBDriver) ValidateConfig(ctx context.Context) error {
	health, err := d.client.Health(ctx)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}
	if health.Status != "pass" {
		return fmt.Errorf("influxDB health check failed: %s", *health.Message)
	}
	return nil
}

// Connect is a no-op for InfluxDB since the client manages connections.
func (d *InfluxDBDriver) Connect(ctx context.Context) error {
	return nil
}

// FetchData sends a Flux query to InfluxDB to retrieve device data.
func (d *InfluxDBDriver) FetchData(ctx context.Context, resourceDetails string) (interface{}, error) {
	var details model.InfluxDBResourceDetails
	if err := json.Unmarshal([]byte(resourceDetails), &details); err != nil {
		return nil, fmt.Errorf("invalid resource details: %w", err)
	}

	if details.Bucket == "" {
		return nil, errors.New("bucket is required in resource details")
	}
	if details.Measurement == "" || details.Field == "" || details.TimeRange == "" {
		return nil, errors.New("measurement, field, and time_range are required in resource details")
	}

	// Construct Flux query
	fluxQuery := fmt.Sprintf(
		`from(bucket: "%s")
		|> range(start: %s)
		|> filter(fn: (r) => r._measurement == "%s")
		|> filter(fn: (r) => r._field == "%s")`,
		details.Bucket, details.TimeRange, details.Measurement, details.Field)

	// Get query API
	queryAPI := d.client.QueryAPI(d.org)

	// Execute query
	result, err := queryAPI.Query(ctx, fluxQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer result.Close()

	// Parse query results
	data := []map[string]interface{}{}
	for result.Next() {
		record := result.Record()
		data = append(data, map[string]interface{}{
			"time":        record.Time(),
			"value":       record.Value(),
			"field":       record.Field(),
			"measurement": record.Measurement(),
		})
	}
	if result.Err() != nil {
		return nil, fmt.Errorf("query error: %w", result.Err())
	}

	return data, nil
}

// TestResource tests a specific InfluxDB resource by executing the query.
func (d *InfluxDBDriver) TestResource(ctx context.Context, resourceDetails string) (interface{}, error) {
	return d.FetchData(ctx, resourceDetails)
}

// Disconnect closes the InfluxDB client.
func (d *InfluxDBDriver) Disconnect(ctx context.Context) error {
	d.client.Close()
	return nil
}
