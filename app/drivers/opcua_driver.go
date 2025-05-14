package drivers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
)

// OPCUADriver implements the PlatformDriver interface for OPC UA platforms.
type OPCUADriver struct {
	Endpoint string
	Timeout  time.Duration
	client   *opcua.Client
}

// NewOPCUADriver creates a new OPCUADriver instance from platform metadata.
func NewOPCUADriver(metadata string) (*OPCUADriver, error) {
	var config struct {
		Endpoint string `json:"endpoint"`
		Timeout  int    `json:"timeout"`
	}
	if err := json.Unmarshal([]byte(metadata), &config); err != nil {
		return nil, fmt.Errorf("invalid metadata: %w", err)
	}
	if config.Endpoint == "" {
		return nil, errors.New("missing endpoint in metadata")
	}
	return &OPCUADriver{
		Endpoint: config.Endpoint,
		Timeout:  time.Duration(config.Timeout) * time.Second,
	}, nil
}

// Connect establishes a connection to the OPC UA server.
func (d *OPCUADriver) Connect(ctx context.Context) error {
	opts := []opcua.Option{
		opcua.SecurityMode(ua.MessageSecurityModeNone),
	}
	client, err := opcua.NewClient(d.Endpoint, opts...)
	if err != nil {
		return fmt.Errorf("failed to create OPC UA Client: %w", err)
	}
	d.client = client
	if err := d.client.Connect(ctx); err != nil {
		return fmt.Errorf("failed to connect to OPC UA server: %w", err)
	}
	return nil
}

// FetchData reads a value from the specified OPC UA node.
func (d *OPCUADriver) FetchData(ctx context.Context, nodeID string) (interface{}, error) {
	if d.client == nil {
		return nil, errors.New("not connected")
	}
	id, err := ua.ParseNodeID(nodeID)
	if err != nil {
		return nil, fmt.Errorf("invalid node ID: %w", err)
	}
	req := &ua.ReadRequest{
		NodesToRead: []*ua.ReadValueID{
			{NodeID: id},
		},
	}
	resp, err := d.client.Read(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to read node: %w", err)
	}
	if len(resp.Results) == 0 {
		return nil, errors.New("no results returned")
	}
	return resp.Results[0].Value, nil
}

// Disconnect closes the connection to the OPC UA server.
func (d *OPCUADriver) Disconnect(ctx context.Context) error {
	if d.client != nil {
		return d.client.Close(context.Background())
	}
	return nil
}
