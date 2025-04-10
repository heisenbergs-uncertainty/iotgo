package controllers

import (
	"context"
	"encoding/json"
	"iotgo/models"
	"strconv"
	"time"

	"github.com/beego/beego/v2/core/logs"
	admin "github.com/beego/beego/v2/server/web"
	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
)

// SnapshotController handles snapshot-related operations for OPC UA integrations
type SnapshotController struct {
	BaseController
}

// Prepare controller
func (c *SnapshotController) Prepare() {
	c.BaseController.Prepare()
}

// Create Snapshot with current values of selected nodes
func (c *SnapshotController) TakeSnapshot() {
	start := time.Now()

	deviceId, err := strconv.Atoi(c.Ctx.Input.Param(":device_id"))
	if err != nil {
		c.Data["Error"] = "Invalid device ID"
		c.Data["Redirect"] = "/devices"
		c.Data["Title"] = "Error"
		c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
		c.Data["RetryURL"] = "/devices/" + c.Ctx.Input.Param(":device_id") + "/snapshot"
		logs.Error("Error in TakeSnapshot:", err)
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("POST", "/devices/"+strconv.Itoa(deviceId)+"/integrations/"+c.Ctx.Input.Param(":integration_id")+"/snapshot", "&SnapshotController.TakeSnapshot", time.Since(start))
		return
	}

	integrationId, err := strconv.Atoi(c.Ctx.Input.Param(":integration_id"))
	if err != nil {
		c.Data["Error"] = "Invalid integration ID"
		c.Data["Redirect"] = "/devices/" + strconv.Itoa(deviceId)
		c.Data["Title"] = "Error"
		c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
		c.Data["RetryURL"] = "/devices/" + c.Ctx.Input.Param(":device_id") + "/integrations/" + c.Ctx.Input.Param(":integration_id") + "/snapshot"
		logs.Error("Error in TakeSnapshot:", err)
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("POST", "/devices/"+strconv.Itoa(deviceId)+"/integrations/"+c.Ctx.Input.Param(":integration_id")+"/snapshot", "&SnapshotController.TakeSnapshot", time.Since(start))
		return
	}

	device, err := models.GetDeviceById(deviceId)
	if err != nil || device == nil {
		c.Data["Error"] = "Device not found"
		c.Data["Redirect"] = "/devices"
		c.Data["Title"] = "Error"
		c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
		c.Data["RetryURL"] = "/devices/" + c.Ctx.Input.Param(":device_id") + "/integrations/" + c.Ctx.Input.Param(":integration_id") + "/snapshot"
		logs.Error("Error in TakeSnapshot:", err)
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("POST", "/devices/"+strconv.Itoa(deviceId)+"/integrations/"+c.Ctx.Input.Param(":integration_id")+"/snapshot", "&SnapshotController.TakeSnapshot", time.Since(start))
		return
	}

	integration, err := models.GetIntegrationById(integrationId)
	if err != nil || integration == nil {
		c.Data["Error"] = "Integration not found"
		c.Data["Redirect"] = "/devices/" + strconv.Itoa(deviceId)
		c.Data["Title"] = "Error"
		c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
		c.Data["RetryURL"] = "/devices/" + c.Ctx.Input.Param(":device_id") + "/integrations/" + c.Ctx.Input.Param(":integration_id") + "/snapshot"
		logs.Error("Error in TakeSnapshot:", err)
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("POST", "/devices/"+strconv.Itoa(deviceId)+"/integrations/"+c.Ctx.Input.Param(":integration_id")+"/snapshot", "&SnapshotController.TakeSnapshot", time.Since(start))
		return
	}

	if integration.IntegrationType != "OPC UA" {
		c.Data["Error"] = "Integration is not of type OPC UA"
		c.Data["Redirect"] = "/devices/" + strconv.Itoa(deviceId)
		c.Data["Title"] = "Error"
		c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
		c.Data["RetryURL"] = "/devices/" + c.Ctx.Input.Param(":device_id") + "/integrations/" + c.Ctx.Input.Param(":integration_id") + "/snapshot"
		logs.Error("Error in TakeSnapshot: Integration is not OPC UA")
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("POST", "/devices/"+strconv.Itoa(deviceId)+"/integrations/"+c.Ctx.Input.Param(":integration_id")+"/snapshot", "&SnapshotController.TakeSnapshot", time.Since(start))
		return
	}

	// Get selected nodes
	nodeIds := c.GetStrings("node_ids")
	nodeNames := c.GetStrings("node_names")
	if len(nodeIds) == 0 || len(nodeNames) == 0 || len(nodeIds) != len(nodeNames) {
		c.Data["Error"] = "Please select at least one node for the snapshot"
		c.Data["Redirect"] = "/devices/" + strconv.Itoa(deviceId)
		c.Data["Title"] = "Error"
		c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
		c.Data["RetryURL"] = "/devices/" + c.Ctx.Input.Param(":device_id") + "/integrations/" + c.Ctx.Input.Param(":integration_id") + "/snapshot"
		logs.Error("Error in TakeSnapshot: No nodes selected or node IDs/names mismatch")
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("POST", "/devices/"+strconv.Itoa(deviceId)+"/integrations/"+c.Ctx.Input.Param(":integration_id")+"/snapshot", "&SnapshotController.TakeSnapshot", time.Since(start))
		return
	}

	ctx := c.Ctx.Request.Context()

	// Connect to the OPC UA server with security
	endpoint := "opc.tcp://" + integration.Host + ":" + integration.Port
	opts := []opcua.Option{
		opcua.SecurityMode(ua.MessageSecurityModeSignAndEncrypt),
		opcua.CertificateFile("path/to/cert.pem"), // Replace with actual certificate path
		opcua.PrivateKeyFile("path/to/key.pem"),   // Replace with actual private key path
	}
	client, err := opcua.NewClient(endpoint, opts...)
	if err := client.Connect(context.Background()); err != nil {
		c.Data["Error"] = "Failed to connect to OPC UA server: " + err.Error()
		c.Data["Redirect"] = "/devices/" + strconv.Itoa(deviceId)
		c.Data["Title"] = "Error"
		c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
		c.Data["RetryURL"] = "/devices/" + c.Ctx.Input.Param(":device_id") + "/integrations/" + c.Ctx.Input.Param(":integration_id") + "/snapshot"
		logs.Error("Error in TakeSnapshot:", err)
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("POST", "/devices/"+strconv.Itoa(deviceId)+"/integrations/"+c.Ctx.Input.Param(":integration_id")+"/snapshot", "&SnapshotController.TakeSnapshot", time.Since(start))
		return
	}
	defer client.Close(ctx)

	// Read values of the selected nodes
	nodesData := make(map[string]string)
	for i, nodeId := range nodeIds {
		node := ua.MustParseNodeID(nodeId)
		req := &ua.ReadRequest{
			NodesToRead: []*ua.ReadValueID{
				{NodeID: node},
			},
		}
		resp, err := client.Read(ctx, req)
		if err != nil {
			c.Data["Error"] = "Failed to read node values: " + err.Error()
			c.Data["Redirect"] = "/devices/" + strconv.Itoa(deviceId)
			c.Data["Title"] = "Error"
			c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
			c.Data["RetryURL"] = "/devices/" + c.Ctx.Input.Param(":device_id") + "/integrations/" + c.Ctx.Input.Param(":integration_id") + "/snapshot"
			logs.Error("Error in TakeSnapshot:", err)
			c.TplName = "error.html"
			admin.StatisticsMap.AddStatistics("POST", "/devices/"+strconv.Itoa(deviceId)+"/integrations/"+c.Ctx.Input.Param(":integration_id")+"/snapshot", "&SnapshotController.TakeSnapshot", time.Since(start))
			return
		}
		for _, result := range resp.Results {
			if result.Status != ua.StatusOK {
				logs.Warn("Failed to read node", nodeId, ":", result.Status)
				continue
			}
			nodesData[nodeNames[i]] = result.Value.String()
		}
	}

	// Store the snapshot
	nodesDataJSON, err := json.Marshal(nodesData)
	if err != nil {
		c.Data["Error"] = "Failed to marshal snapshot data: " + err.Error()
		c.Data["Redirect"] = "/devices/" + strconv.Itoa(deviceId)
		c.Data["Title"] = "Error"
		c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
		c.Data["RetryURL"] = "/devices/" + c.Ctx.Input.Param(":device_id") + "/integrations/" + c.Ctx.Input.Param(":integration_id") + "/snapshot"
		logs.Error("Error in TakeSnapshot:", err)
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("POST", "/devices/"+strconv.Itoa(deviceId)+"/integrations/"+c.Ctx.Input.Param(":integration_id")+"/snapshot", "&SnapshotController.TakeSnapshot", time.Since(start))
		return
	}

	snapshot := &models.Snapshot{
		DeviceId:      deviceId,
		IntegrationId: integrationId,
		Nodes:         string(nodesDataJSON),
	}
	if err := models.AddSnapshot(snapshot); err != nil {
		c.Data["Error"] = "Failed to store snapshot: " + err.Error()
		c.Data["Redirect"] = "/devices/" + strconv.Itoa(deviceId)
		c.Data["Title"] = "Error"
		c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
		c.Data["RetryURL"] = "/devices/" + c.Ctx.Input.Param(":device_id") + "/integrations/" + c.Ctx.Input.Param(":integration_id") + "/snapshot"
		logs.Error("Error in TakeSnapshot:", err)
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("POST", "/devices/"+strconv.Itoa(deviceId)+"/integrations/"+c.Ctx.Input.Param(":integration_id")+"/snapshot", "&SnapshotController.TakeSnapshot", time.Since(start))
		return
	}

	c.SetFlash("Snapshot taken successfully")
	c.Redirect("/devices/"+strconv.Itoa(deviceId), 302)
	admin.StatisticsMap.AddStatistics("POST", "/devices/"+strconv.Itoa(deviceId)+"/integrations/"+c.Ctx.Input.Param(":integration_id")+"/snapshot", "&SnapshotController.TakeSnapshot", time.Since(start))
}
