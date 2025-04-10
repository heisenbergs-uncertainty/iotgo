package controllers

import (
	"context"
	"html/template"
	"iotgo/models"
	"strconv"
	"time"

	"github.com/beego/beego/v2/client/orm"
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
	admin "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/pagination"
	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
)

type WebDeviceController struct {
	BaseController
}

func (c *WebDeviceController) Prepare() {
	c.BaseController.Prepare()
}

func (c *WebDeviceController) List() {
	start := time.Now()
	devices, err := models.GetAllDevices()
	if err != nil {
		c.Ctx.WriteString("Error: " + err.Error())
		return
	}
	c.Data["Devices"] = devices
	c.Data["Title"] = "Device List"
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.TplName = "devices/list.html"
	admin.StatisticsMap.AddStatistics("GET", "/devices", "&WebDeviceController.List", time.Since(start))
}

func (c *WebDeviceController) New() {
	start := time.Now()
	c.Data["Title"] = "New Device"
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.TplName = "devices/new.html"
	admin.StatisticsMap.AddStatistics("GET", "/devices/new", "&WebDeviceController.New", time.Since(start))
}

func (c *WebDeviceController) Create() {
	start := time.Now()
	name := c.GetString("name")
	manufacturer := c.GetString("manufacturer")
	deviceType := c.GetString("type")

	v := validation.Validation{}
	v.Required(name, "name").Message("Device name is required")
	v.MaxSize(name, 100, "name").Message("Name must be 100 characters or less")
	v.Required(manufacturer, "manufacturer").Message("Manufacturer is required")
	v.MaxSize(manufacturer, 100, "manufacturer").Message("Manufacturer must be 100 characters or less")
	v.Required(deviceType, "type").Message("Type is required")
	v.MaxSize(deviceType, 50, "type").Message("Type must be 50 characters or less")

	if v.HasErrors() {
		c.Data["Errors"] = v.Errors
		c.Data["Title"] = "New Device"
		c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
		c.TplName = "devices/new.html"
		admin.StatisticsMap.AddStatistics("POST", "/devices", "&WebDeviceController.Create", time.Since(start))
		return
	}

	device := &models.Device{Name: name, Manufacturer: manufacturer, Type: deviceType}
	if err := models.AddDevice(device); err != nil {
		logs.Error("Failed to add device:", err)
		c.Data["Error"] = "Failed to add device"
		c.Data["Redirect"] = "/devices"
		c.Data["Title"] = "Error"
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("POST", "/devices", "&WebDeviceController.Create", time.Since(start))
		return
	}
	logs.Info("Device added:", device.Name)
	c.Redirect("/devices", 302)
	admin.StatisticsMap.AddStatistics("POST", "/devices", "&WebDeviceController.Create", time.Since(start))
}

// Edit shows the form to edit an existing device
func (c *WebDeviceController) Edit() {
	start := time.Now()

	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		userId, _ := c.GetSession("user_id").(int)
		models.AddErrorLog("Error in Edit: "+err.Error(), userId)
		c.Data["Error"] = "Invalid device ID"
		c.Data["Redirect"] = "/devices"
		c.Data["Title"] = "Error"
		c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
		c.Data["RetryURL"] = "/devices/" + c.Ctx.Input.Param(":id") + "/edit"
		logs.Error("Error in Edit:", err)
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("GET", "/devices/"+strconv.Itoa(id)+"/edit", "&WebDeviceController.Edit", time.Since(start))
		return
	}

	device, err := models.GetDeviceById(id)
	if err != nil || device == nil {
		userId, _ := c.GetSession("user_id").(int)
		models.AddErrorLog("Error in Edit: "+err.Error(), userId)
		c.Data["Error"] = "Device not found"
		c.Data["Redirect"] = "/devices"
		c.Data["Title"] = "Error"
		c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
		c.Data["RetryURL"] = "/devices/" + c.Ctx.Input.Param(":id") + "/edit"
		logs.Error("Error in Edit:", err)
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("GET", "/devices/"+strconv.Itoa(id)+"/edit", "&WebDeviceController.Edit", time.Since(start))
		return
	}

	c.Data["Device"] = device
	c.Data["Title"] = "Edit Device"
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.TplName = "devices/edit.html"
	admin.StatisticsMap.AddStatistics("GET", "/devices/"+strconv.Itoa(id)+"/edit", "&WebDeviceController.Edit", time.Since(start))
}

func (c *WebDeviceController) Show() {
	start := time.Now()
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Data["Error"] = "Invalid device ID"
		c.Data["Redirect"] = "/devices"
		c.Data["Title"] = "Error"
		c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
		c.Data["RetryURL"] = "/devices/" + c.Ctx.Input.Param(":id")
		logs.Error("Error in Show:", err)
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("GET", "/devices/"+strconv.Itoa(id), "&WebDeviceController.Show", time.Since(start))
		return
	}

	device, err := models.GetDeviceById(id)
	if err != nil || device == nil {
		c.Data["Error"] = "Device not found"
		c.Data["Redirect"] = "/devices"
		c.Data["Title"] = "Error"
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("GET", "/devices/"+strconv.Itoa(id), "&WebDeviceController.Show", time.Since(start))
		return
	}

	integrations, err := models.GetIntegrationsByDeviceId(id)
	if err != nil {
		c.Data["Error"] = "Failed to fetch integrations: " + err.Error()
		c.Data["Redirect"] = "/devices"
		c.Data["Title"] = "Error"
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("GET", "/devices/"+strconv.Itoa(id), "&WebDeviceController.Show", time.Since(start))
		return
	}

	// Fetch recordings and snapshots
	recordings, err := models.GetRecordingsByDeviceId(id)
	if err != nil {
		logs.Error("Failed to fetch recordings:", err)
	}
	snapshots, err := models.GetSnapshotsByDeviceId(id)
	if err != nil {
		logs.Error("Failed to fetch snapshots:", err)
	}

	c.Data["Device"] = device
	c.Data["Integrations"] = integrations
	c.Data["Recordings"] = recordings
	c.Data["Snapshots"] = snapshots
	c.Data["Title"] = "Device: " + device.Name
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.TplName = "devices/show.html"
	admin.StatisticsMap.AddStatistics("GET", "/devices/"+strconv.Itoa(id), "&WebDeviceController.Show", time.Since(start))
}

func (c *WebDeviceController) Update() {
	start := time.Now()

	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Data["Error"] = "Invalid device ID"
		c.Data["Redirect"] = "/devices"
		c.Data["Title"] = "Error"
		c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
		c.Data["RetryURL"] = "/devices/" + c.Ctx.Input.Param(":id") + "/edit"
		logs.Error("Error in Update:", err)
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("POST", "/devices/"+strconv.Itoa(id)+"/edit", "&WebDeviceController.Update", time.Since(start))
		return
	}

	device, err := models.GetDeviceById(id)
	if err != nil || device == nil {
		c.Data["Error"] = "Device not found"
		c.Data["Redirect"] = "/devices"
		c.Data["Title"] = "Error"
		c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
		c.Data["RetryURL"] = "/devices/" + c.Ctx.Input.Param(":id") + "/edit"
		logs.Error("Error in Update:", err)
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("POST", "/devices/"+strconv.Itoa(id)+"/edit", "&WebDeviceController.Update", time.Since(start))
		return
	}

	name := c.GetString("name")
	manufacturer := c.GetString("manufacturer")
	deviceType := c.GetString("type")

	device.Name = name
	device.Manufacturer = manufacturer
	device.Type = deviceType

	if err := models.UpdateDevice(device); err != nil {
		logs.Error("Failed to update device:", err)
		c.Data["Error"] = "Failed to update device: " + err.Error()
		c.Data["Redirect"] = "/devices/" + strconv.Itoa(id)
		c.Data["Title"] = "Error"
		c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
		c.Data["RetryURL"] = "/devices/" + c.Ctx.Input.Param(":id") + "/edit"
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("POST", "/devices/"+strconv.Itoa(id)+"/edit", "&WebDeviceController.Update", time.Since(start))
		return
	}
	logs.Info("Device updated:", device.Name)
	c.Redirect("/devices/"+strconv.Itoa(id), 302)
	admin.StatisticsMap.AddStatistics("POST", "/devices/"+strconv.Itoa(id)+"/edit", "&WebDeviceController.Update", time.Since(start))
}

func (c *WebDeviceController) GetAllDeviceWithOffsetLimit(offset int, limit int) ([]*models.Device, error) {
	o := orm.NewOrm()
	var devices []*models.Device
	_, err := o.QueryTable("device").Limit(limit).Offset(int64(offset)).All(&devices)
	return devices, err
}

// BrowseNodes connects to the OPC UA server and browses its nodes
func (c *WebDeviceController) BrowseNodes() {
	start := time.Now()

	deviceId, err := strconv.Atoi(c.Ctx.Input.Param(":device_id"))
	if err != nil {
		c.Data["Error"] = "Invalid device ID"
		c.Data["Redirect"] = "/devices"
		c.Data["Title"] = "Error"
		c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
		c.Data["RetryURL"] = "/devices/" + c.Ctx.Input.Param(":device_id") + "/browse"
		logs.Error("Error in BrowseNodes:", err)
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("GET", "/devices/"+strconv.Itoa(deviceId)+"/browse", "&WebDeviceController.BrowseNodes", time.Since(start))
		return
	}

	integrationId, err := strconv.Atoi(c.Ctx.Input.Param(":integration_id"))
	if err != nil {
		c.Data["Error"] = "Invalid integration ID"
		c.Data["Redirect"] = "/devices/" + strconv.Itoa(deviceId)
		c.Data["Title"] = "Error"
		c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
		c.Data["RetryURL"] = "/devices/" + c.Ctx.Input.Param(":device_id") + "/browse"
		logs.Error("Error in BrowseNodes:", err)
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("GET", "/devices/"+strconv.Itoa(deviceId)+"/browse", "&WebDeviceController.BrowseNodes", time.Since(start))
		return
	}

	device, err := models.GetDeviceById(deviceId)
	if err != nil || device == nil {
		c.Data["Error"] = "Device not found"
		c.Data["Redirect"] = "/devices"
		c.Data["Title"] = "Error"
		c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
		c.Data["RetryURL"] = "/devices/" + c.Ctx.Input.Param(":device_id") + "/browse"
		logs.Error("Error in BrowseNodes:", err)
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("GET", "/devices/"+strconv.Itoa(deviceId)+"/browse", "&WebDeviceController.BrowseNodes", time.Since(start))
		return
	}

	integration, err := models.GetIntegrationById(integrationId)
	if err != nil || integration == nil {
		c.Data["Error"] = "Integration not found"
		c.Data["Redirect"] = "/devices/" + strconv.Itoa(deviceId)
		c.Data["Title"] = "Error"
		c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
		c.Data["RetryURL"] = "/devices/" + c.Ctx.Input.Param(":device_id") + "/browse"
		logs.Error("Error in BrowseNodes:", err)
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("GET", "/devices/"+strconv.Itoa(deviceId)+"/browse", "&WebDeviceController.BrowseNodes", time.Since(start))
		return
	}

	if integration.IntegrationType != "OPC UA" {
		c.Data["Error"] = "Integration is not of type OPC UA"
		c.Data["Redirect"] = "/devices/" + strconv.Itoa(deviceId)
		c.Data["Title"] = "Error"
		c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
		c.Data["RetryURL"] = "/devices/" + c.Ctx.Input.Param(":device_id") + "/browse"
		logs.Error("Error in BrowseNodes: Integration is not OPC UA")
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("GET", "/devices/"+strconv.Itoa(deviceId)+"/browse", "&WebDeviceController.BrowseNodes", time.Since(start))
		return
	}

	// Connect to the OPC UA server with security
	ctx := context.Background()
	endpoint := "opc.tcp://" + integration.Host + ":" + integration.Port
	opts := []opcua.Option{}
	client, err := opcua.NewClient(endpoint, opts...)
	if err := client.Connect(ctx); err != nil {
		c.Data["Error"] = "Failed to connect to OPC UA server: " + err.Error()
		c.Data["Redirect"] = "/devices/" + strconv.Itoa(deviceId)
		c.Data["Title"] = "Error"
		c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
		c.Data["RetryURL"] = "/devices/" + c.Ctx.Input.Param(":device_id") + "/browse"
		logs.Error("Error in BrowseNodes:", err)
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("GET", "/devices/"+strconv.Itoa(deviceId)+"/browse", "&WebDeviceController.BrowseNodes", time.Since(start))
		return
	}
	defer client.Close(ctx)

	// Pagination parameters
	page, _ := c.GetInt("page", 1)
	perPage := 50
	offset := (page - 1) * perPage

	// Browse the Objects folder (default starting point)
	nodeId := ua.MustParseNodeID("ns=0;i=85") // Objects folder (Root/Objects)
	allNodes, err := browseNodes(client, nodeId, 0)
	if err != nil {
		c.Data["Error"] = "Failed to browse OPC UA nodes: " + err.Error()
		c.Data["Redirect"] = "/devices/" + strconv.Itoa(deviceId)
		c.Data["Title"] = "Error"
		c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
		c.Data["RetryURL"] = "/devices/" + c.Ctx.Input.Param(":device_id") + "/browse"
		logs.Error("Error in BrowseNodes:", err)
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("GET", "/devices/"+strconv.Itoa(deviceId)+"/browse", "&WebDeviceController.BrowseNodes", time.Since(start))
		return
	}

	// Flatten the node tree for pagination
	var flatNodes []*Node
	flattenNodes(allNodes, &flatNodes)

	// Apply pagination
	total := len(flatNodes)
	paginator := pagination.SetPaginator(c.Ctx, perPage, int64(total))
	end := offset + perPage
	if end > total {
		end = total
	}
	if offset > total {
		offset = total
	}
	nodes := flatNodes[offset:end]

	c.Data["Device"] = device
	c.Data["Integration"] = integration
	c.Data["Nodes"] = nodes
	c.Data["TotalNodes"] = total
	c.Data["Paginator"] = paginator
	c.Data["Title"] = "Browse OPC UA Nodes - " + device.Name
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.TplName = "devices/browse_nodes.html"
	admin.StatisticsMap.AddStatistics("GET", "/devices/"+strconv.Itoa(deviceId)+"/browse", "&WebDeviceController.BrowseNodes", time.Since(start))
}

// Node represents an OPC UA node for display purposes
type Node struct {
	NodeId      string
	BrowseName  string
	DisplayName string
	NodeClass   string
	Level       int
	Children    []*Node
}

// browseNodes recursively browses OPC UA nodes with optimization
func browseNodes(client *opcua.Client, nodeId *ua.NodeID, level int) ([]*Node, error) {
	if level > 10 { // Prevent infinite recursion in cyclic graphs
		return nil, nil
	}

	browseReq := &ua.BrowseRequest{
		NodesToBrowse: []*ua.BrowseDescription{
			{
				NodeID:          nodeId,
				BrowseDirection: ua.BrowseDirectionForward,
				ReferenceTypeID: ua.MustParseNodeID("i=35"), // HierarchicalReferences
				IncludeSubtypes: true,
				NodeClassMask:   uint32(ua.NodeClassObject | ua.NodeClassVariable),
				ResultMask:      uint32(ua.BrowseResultMaskAll),
			},
		},
	}

	ctx := context.Background()
	resp, err := client.Browse(ctx, browseReq)
	if err != nil {
		return nil, err
	}

	var nodes []*Node
	for _, result := range resp.Results {
		for _, ref := range result.References {
			node := &Node{
				NodeId:      ref.NodeID.String(),
				BrowseName:  ref.BrowseName.Name,
				DisplayName: ref.DisplayName.Text,
				NodeClass:   ref.NodeClass.String(),
				Level:       level,
			}

			// Recursively browse children
			children, err := browseNodes(client, ref.NodeID.NodeID, level+1)
			if err != nil {
				logs.Warn("Error browsing children for node", ref.NodeID.String(), ":", err)
				continue
			}
			node.Children = children
			nodes = append(nodes, node)
		}
	}
	return nodes, nil
}

// flattenNodes flattens the hierarchical node tree into a single list for pagination
func flattenNodes(nodes []*Node, flatNodes *[]*Node) {
	for _, node := range nodes {
		*flatNodes = append(*flatNodes, node)
		if len(node.Children) > 0 {
			flattenNodes(node.Children, flatNodes)
		}
	}
}

func (c *WebDeviceController) Delete() {
	start := time.Now()

	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Data["Error"] = "Invalid device ID"
		c.Data["Redirect"] = "/devices"
		c.Data["Title"] = "Error"
		c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
		c.Data["RetryURL"] = "/devices/" + c.Ctx.Input.Param(":id") + "/delete"
		logs.Error("Error in Delete:", err)
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("POST", "/devices/"+strconv.Itoa(id)+"/delete", "&WebDeviceController.Delete", time.Since(start))
		return
	}

	if err := models.DeleteDevice(id); err != nil {
		logs.Error("Failed to delete device:", err)
		c.Data["Error"] = "Failed to delete device: " + err.Error()
		c.Data["Redirect"] = "/devices"
		c.Data["Title"] = "Error"
		c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
		c.Data["RetryURL"] = "/devices/" + c.Ctx.Input.Param(":id") + "/delete"
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("POST", "/devices/"+strconv.Itoa(id)+"/delete", "&WebDeviceController.Delete", time.Since(start))
		return
	}
	logs.Info("Device deleted, ID:", id)
	c.Redirect("/devices", 302)
	admin.StatisticsMap.AddStatistics("POST", "/devices/"+strconv.Itoa(id)+"/delete", "&WebDeviceController.Delete", time.Since(start))
}
