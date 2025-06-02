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

// OPCUADriver implements the PlatformDriver interface for OPCUA platforms.
type OPCUADriver struct {
	client *opcua.Client
	config OPCUAMetadata
}

// OPCUAMetadata defines the configuration for OPCUA platforms.
type OPCUAMetadata struct {
	Endpoint string `json:"endpoint"`
	Username string `json:"username"`
	Password string `json:"password"`
	Timeout  int    `json:"timeout"`
}

// OPCUAResourceDetails defines the details for an OPCUA resource.
type OPCUAResourceDetails struct {
	NodeID string `json:"node_id"`
}

// NewOPCUADriver creates a new OPCUADriver instance from platform metadata.
func NewOPCUADriver(metadata string) (*OPCUADriver, error) {
	var config OPCUAMetadata
	if err := json.Unmarshal([]byte(metadata), &config); err != nil {
		return nil, fmt.Errorf("invalid metadata: %w", err)
	}
	if config.Endpoint == "" {
		return nil, errors.New("missing endpoint in metadata")
	}
	if config.Timeout == 0 {
		config.Timeout = 10
	}

	opts := []opcua.Option{
		opcua.SecurityMode(ua.MessageSecurityModeNone),
		opcua.AutoReconnect(true),
		opcua.ReconnectInterval(time.Second * time.Duration(config.Timeout)),
	}
	if config.Username != "" && config.Password != "" {
		opts = append(opts, opcua.AuthUsername(config.Username, config.Password))
	}

	client, err := opcua.NewClient(config.Endpoint, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create OPCUA client: %w", err)
	}

	return &OPCUADriver{
		client: client,
		config: config,
	}, nil
}

// Connect establishes a connection to the OPCUA server.
func (d *OPCUADriver) Connect(ctx context.Context) error {
	if err := d.client.Connect(ctx); err != nil {
		return fmt.Errorf("failed to connect to OPCUA server: %w", err)
	}
	return nil
}

// FetchData retrieves data from an OPCUA node.
func (d *OPCUADriver) FetchData(ctx context.Context, resourceDetails string) (interface{}, error) {
	var details OPCUAResourceDetails
	if err := json.Unmarshal([]byte(resourceDetails), &details); err != nil {
		return nil, fmt.Errorf("invalid resource details: %w", err)
	}

	if details.NodeID == "" {
		return nil, errors.New("node_id is required in resource details")
	}

	nodeID, err := ua.ParseNodeID(details.NodeID)
	if err != nil {
		return nil, fmt.Errorf("invalid node ID: %w", err)
	}

	req := &ua.ReadRequest{
		NodesToRead: []*ua.ReadValueID{
			{NodeID: nodeID, AttributeID: ua.AttributeIDValue},
		},
	}
	resp, err := d.client.Read(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to read node: %w", err)
	}

	if len(resp.Results) == 0 {
		return nil, errors.New("no data returned")
	}

	return map[string]interface{}{
		"node_id": details.NodeID,
		"value":   resp.Results[0].Value,
		"status":  resp.Results[0].Status,
	}, nil
}

// Disconnect closes the OPCUA connection.
func (d *OPCUADriver) Disconnect(ctx context.Context) error {
	return d.client.Close(ctx)
}

// ValidateConfig checks if the OPCUA configuration is valid.
func (d *OPCUADriver) ValidateConfig(ctx context.Context) error {
	if err := d.Connect(ctx); err != nil {
		return err
	}
	defer d.Disconnect(ctx)
	return nil
}

// TestResource tests an OPCUA resource by reading the node.
func (d *OPCUADriver) TestResource(ctx context.Context, resourceDetails string) (interface{}, error) {
	return d.FetchData(ctx, resourceDetails)
}
