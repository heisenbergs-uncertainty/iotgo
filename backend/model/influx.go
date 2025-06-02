package model

// InfluxDBMetadata defines the structure for InfluxDB platform metadata
type InfluxDBMetadata struct {
	URL     string `json:"url"`               // e.g., "http://localhost:8086"
	Token   string `json:"token"`             // InfluxDB API token
	Org     string `json:"org"`               // Organization name
	Bucket  string `json:"bucket"`            // Default bucket
	Timeout int    `json:"timeout,omitempty"` // Timeout in seconds
}

// InfluxDBResourceDetails defines the structure for InfluxDB query details
type InfluxDBResourceDetails struct {
	Bucket      string `json:"bucket"`      // Bucket to query
	Measurement string `json:"measurement"` // Measurement name
	Field       string `json:"field"`       // Field to retrieve
	TimeRange   string `json:"time_range"`  // e.g., "-1h" for last hour
}
