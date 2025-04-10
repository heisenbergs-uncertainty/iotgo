package controllers

import (
	"html/template"
	"iotgo/models"
	"strconv"
	"time"

	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/core/validation"
	admin "github.com/beego/beego/v2/server/web"
)

// WebIntegrationController handles integration-related operations
type WebIntegrationController struct {
	BaseController
}

// Prepare sets up the controller
func (c *WebIntegrationController) Prepare() {
	c.BaseController.Prepare()
}

// List shows the list of all integrations
func (c *WebIntegrationController) List() {
	start := time.Now()

	integrations, err := models.GetAllIntegrations()
	if err != nil {
		userId, _ := c.GetSession("user_id").(int)
		models.AddErrorLog("Error in List: "+err.Error(), userId)
		c.Data["Error"] = "Error fetching integrations: " + err.Error()
		c.Data["Redirect"] = "/"
		c.Data["Title"] = "Error"
		c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
		c.Data["RetryURL"] = "/integrations"
		logs.Error("Error in List:", err)
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("GET", "/integrations", "&WebIntegrationController.List", time.Since(start))
		return
	}

	c.Data["Integrations"] = integrations
	c.Data["Title"] = "Integrations List"
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.TplName = "integrations/list.html"
	admin.StatisticsMap.AddStatistics("GET", "/integrations", "&WebIntegrationController.List", time.Since(start))
}

// Show displays an integration's details
func (c *WebIntegrationController) Show() {
	start := time.Now()

	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		userId, _ := c.GetSession("user_id").(int)
		models.AddErrorLog("Error in Show: "+err.Error(), userId)
		c.Data["Error"] = "Invalid integration ID"
		c.Data["Redirect"] = "/integrations"
		c.Data["Title"] = "Error"
		c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
		c.Data["RetryURL"] = "/integrations/" + c.Ctx.Input.Param(":id")
		logs.Error("Error in Show:", err)
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("GET", "/integrations/"+strconv.Itoa(id), "&WebIntegrationController.Show", time.Since(start))
		return
	}

	integration, err := models.GetIntegrationById(id)
	if err != nil || integration == nil {
		userId, _ := c.GetSession("user_id").(int)
		models.AddErrorLog("Error in Show: "+err.Error(), userId)
		c.Data["Error"] = "Integration not found"
		c.Data["Redirect"] = "/integrations"
		c.Data["Title"] = "Error"
		c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
		c.Data["RetryURL"] = "/integrations/" + c.Ctx.Input.Param(":id")
		logs.Error("Error in Show:", err)
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("GET", "/integrations/"+strconv.Itoa(id), "&WebIntegrationController.Show", time.Since(start))
		return
	}

	c.Data["Integration"] = integration
	c.Data["Title"] = "Integration: " + integration.IntegrationType
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.TplName = "integrations/show.html"
	admin.StatisticsMap.AddStatistics("GET", "/integrations/"+strconv.Itoa(id), "&WebIntegrationController.Show", time.Since(start))
}

// New shows the form to add a new integration
func (c *WebIntegrationController) New() {
	start := time.Now()

	deviceId, err := strconv.Atoi(c.Ctx.Input.Param(":device_id"))
	if err != nil {
		userId, _ := c.GetSession("user_id").(int)
		models.AddErrorLog("Error in New: "+err.Error(), userId)
		c.Data["Error"] = "Invalid device ID"
		c.Data["Redirect"] = "/devices"
		c.Data["Title"] = "Error"
		c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
		c.Data["RetryURL"] = "/devices/" + c.Ctx.Input.Param(":device_id") + "/integrations/new"
		logs.Error("Error in New:", err)
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("GET", "/devices/"+strconv.Itoa(deviceId)+"/integrations/new", "&WebIntegrationController.New", time.Since(start))
		return
	}

	c.Data["DeviceId"] = deviceId
	c.Data["Title"] = "New Integration"
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.TplName = "integrations/new.html"
	admin.StatisticsMap.AddStatistics("GET", "/devices/"+strconv.Itoa(deviceId)+"/integrations/new", "&WebIntegrationController.New", time.Since(start))
}

// Create handles the creation of a new integration
func (c *WebIntegrationController) Create() {
	start := time.Now()

	deviceId, err := strconv.Atoi(c.Ctx.Input.Param(":device_id"))
	if err != nil {
		userId, _ := c.GetSession("user_id").(int)
		models.AddErrorLog("Error in Create: "+err.Error(), userId)
		c.Data["Error"] = "Invalid device ID"
		c.Data["Redirect"] = "/devices"
		c.Data["Title"] = "Error"
		c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
		c.Data["RetryURL"] = "/devices/" + c.Ctx.Input.Param(":device_id") + "/integrations/new"
		logs.Error("Error in Create:", err)
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("POST", "/devices/"+strconv.Itoa(deviceId)+"/integrations", "&WebIntegrationController.Create", time.Since(start))
		return
	}

	integrationType := c.GetString("integration_type")
	identifier := c.GetString("identifier")
	host := c.GetString("host")
	port := c.GetString("port")
	protocol := c.GetString("protocol")

	// Validate input
	valid := validation.Validation{}
	valid.Required(integrationType, "integration_type").Message("Integration type is required")
	valid.Required(identifier, "identifier").Message("Identifier is required")
	valid.Required(host, "host").Message("Host is required")

	// Validate port if provided
	if port != "" {
		portNum, err := strconv.Atoi(port)
		if err != nil {
			c.Data["Errors"] = []string{"Port must be a valid number"}
			c.Data["DeviceId"] = deviceId
			c.Data["Title"] = "New Integration"
			c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
			c.TplName = "integrations/new.html"
			admin.StatisticsMap.AddStatistics("POST", "/devices/"+strconv.Itoa(deviceId)+"/integrations", "&WebIntegrationController.Create", time.Since(start))
			return
		}
		valid.Range(portNum, 1, 65535, "port").Message("Port must be between 1 and 65535")
	}

	if valid.HasErrors() {
		var errors []string
		for _, err := range valid.Errors {
			errors = append(errors, err.Message)
		}
		c.Data["Errors"] = errors
		c.Data["DeviceId"] = deviceId
		c.Data["Title"] = "New Integration"
		c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
		c.TplName = "integrations/new.html"
		admin.StatisticsMap.AddStatistics("POST", "/devices/"+strconv.Itoa(deviceId)+"/integrations", "&WebIntegrationController.Create", time.Since(start))
		return
	}

	device, err := models.GetDeviceById(deviceId)
	if err != nil || device == nil {
		userId, _ := c.GetSession("user_id").(int)
		models.AddErrorLog("Error in Create: "+err.Error(), userId)
		c.Data["Error"] = "Device not found"
		c.Data["Redirect"] = "/devices"
		c.Data["Title"] = "Error"
		c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
		c.Data["RetryURL"] = "/devices/" + c.Ctx.Input.Param(":device_id") + "/integrations/new"
		logs.Error("Error in Create:", err)
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("POST", "/devices/"+strconv.Itoa(deviceId)+"/integrations", "&WebIntegrationController.Create", time.Since(start))
		return
	}

	integration := &models.Integration{
		Device:          device,
		IntegrationType: integrationType,
		Identifier:      identifier,
		Host:            host,
		Port:            port,
		Protocol:        protocol,
	}

	if err := models.AddIntegration(integration); err != nil {
		userId, _ := c.GetSession("user_id").(int)
		models.AddErrorLog("Error in Create: "+err.Error(), userId)
		c.Data["Error"] = "Failed to add integration: " + err.Error()
		c.Data["Redirect"] = "/devices/" + strconv.Itoa(deviceId)
		c.Data["Title"] = "Error"
		c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
		c.Data["RetryURL"] = "/devices/" + c.Ctx.Input.Param(":device_id") + "/integrations/new"
		logs.Error("Error in Create:", err)
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("POST", "/devices/"+strconv.Itoa(deviceId)+"/integrations", "&WebIntegrationController.Create", time.Since(start))
		return
	}

	c.SetFlash("Integration created successfully")
	c.Redirect("/devices/"+strconv.Itoa(deviceId), 302)
	admin.StatisticsMap.AddStatistics("POST", "/devices/"+strconv.Itoa(deviceId)+"/integrations", "&WebIntegrationController.Create", time.Since(start))
}

// Edit shows the form to edit an existing integration
func (c *WebIntegrationController) Edit() {
	start := time.Now()

	deviceId, err := strconv.Atoi(c.Ctx.Input.Param(":device_id"))
	if err != nil {
		userId, _ := c.GetSession("user_id").(int)
		models.AddErrorLog("Error in Edit: "+err.Error(), userId)
		c.Data["Error"] = "Invalid device ID"
		c.Data["Redirect"] = "/devices"
		c.Data["Title"] = "Error"
		c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
		c.Data["RetryURL"] = "/devices/" + c.Ctx.Input.Param(":device_id") + "/integrations/" + c.Ctx.Input.Param(":id") + "/edit"
		logs.Error("Error in Edit:", err)
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("GET", "/devices/"+strconv.Itoa(deviceId)+"/integrations/"+c.Ctx.Input.Param(":id")+"/edit", "&WebIntegrationController.Edit", time.Since(start))
		return
	}

	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		userId, _ := c.GetSession("user_id").(int)
		models.AddErrorLog("Error in Edit: "+err.Error(), userId)
		c.Data["Error"] = "Invalid integration ID"
		c.Data["Redirect"] = "/devices/" + strconv.Itoa(deviceId)
		c.Data["Title"] = "Error"
		c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
		c.Data["RetryURL"] = "/devices/" + c.Ctx.Input.Param(":device_id") + "/integrations/" + c.Ctx.Input.Param(":id") + "/edit"
		logs.Error("Error in Edit:", err)
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("GET", "/devices/"+strconv.Itoa(deviceId)+"/integrations/"+strconv.Itoa(id)+"/edit", "&WebIntegrationController.Edit", time.Since(start))
		return
	}

	integration, err := models.GetIntegrationById(id)
	if err != nil || integration == nil {
		userId, _ := c.GetSession("user_id").(int)
		models.AddErrorLog("Error in Edit: "+err.Error(), userId)
		c.Data["Error"] = "Integration not found"
		c.Data["Redirect"] = "/devices/" + strconv.Itoa(deviceId)
		c.Data["Title"] = "Error"
		c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
		c.Data["RetryURL"] = "/devices/" + c.Ctx.Input.Param(":device_id") + "/integrations/" + c.Ctx.Input.Param(":id") + "/edit"
		logs.Error("Error in Edit:", err)
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("GET", "/devices/"+strconv.Itoa(deviceId)+"/integrations/"+strconv.Itoa(id)+"/edit", "&WebIntegrationController.Edit", time.Since(start))
		return
	}

	c.Data["DeviceId"] = deviceId
	c.Data["Integration"] = integration
	c.Data["Title"] = "Edit Integration"
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	c.TplName = "integrations/edit.html"
	admin.StatisticsMap.AddStatistics("GET", "/devices/"+strconv.Itoa(deviceId)+"/integrations/"+strconv.Itoa(id)+"/edit", "&WebIntegrationController.Edit", time.Since(start))
}

// Update handles the updating of an integration
func (c *WebIntegrationController) Update() {
	start := time.Now()

	deviceId, err := strconv.Atoi(c.Ctx.Input.Param(":device_id"))
	if err != nil {
		userId, _ := c.GetSession("user_id").(int)
		models.AddErrorLog("Error in Update: "+err.Error(), userId)
		c.Data["Error"] = "Invalid device ID"
		c.Data["Redirect"] = "/devices"
		c.Data["Title"] = "Error"
		c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
		c.Data["RetryURL"] = "/devices/" + c.Ctx.Input.Param(":device_id") + "/integrations/" + c.Ctx.Input.Param(":id") + "/edit"
		logs.Error("Error in Update:", err)
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("POST", "/devices/"+strconv.Itoa(deviceId)+"/integrations/"+c.Ctx.Input.Param(":id")+"/edit", "&WebIntegrationController.Update", time.Since(start))
		return
	}

	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		userId, _ := c.GetSession("user_id").(int)
		models.AddErrorLog("Error in Update: "+err.Error(), userId)
		c.Data["Error"] = "Invalid integration ID"
		c.Data["Redirect"] = "/devices/" + strconv.Itoa(deviceId)
		c.Data["Title"] = "Error"
		c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
		c.Data["RetryURL"] = "/devices/" + c.Ctx.Input.Param(":device_id") + "/integrations/" + c.Ctx.Input.Param(":id") + "/edit"
		logs.Error("Error in Update:", err)
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("POST", "/devices/"+strconv.Itoa(deviceId)+"/integrations/"+strconv.Itoa(id)+"/edit", "&WebIntegrationController.Update", time.Since(start))
		return
	}

	integration, err := models.GetIntegrationById(id)
	if err != nil || integration == nil {
		userId, _ := c.GetSession("user_id").(int)
		models.AddErrorLog("Error in Update: "+err.Error(), userId)
		c.Data["Error"] = "Integration not found"
		c.Data["Redirect"] = "/devices/" + strconv.Itoa(deviceId)
		c.Data["Title"] = "Error"
		c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
		c.Data["RetryURL"] = "/devices/" + c.Ctx.Input.Param(":device_id") + "/integrations/" + c.Ctx.Input.Param(":id") + "/edit"
		logs.Error("Error in Update:", err)
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("POST", "/devices/"+strconv.Itoa(deviceId)+"/integrations/"+strconv.Itoa(id)+"/edit", "&WebIntegrationController.Update", time.Since(start))
		return
	}

	integrationType := c.GetString("integration_type")
	identifier := c.GetString("identifier")
	host := c.GetString("host")
	port := c.GetString("port")
	protocol := c.GetString("protocol")

	// Validate input
	valid := validation.Validation{}
	valid.Required(integrationType, "integration_type").Message("Integration type is required")
	valid.Required(identifier, "identifier").Message("Identifier is required")
	valid.Required(host, "host").Message("Host is required")

	// Validate port if provided
	if port != "" {
		portNum, err := strconv.Atoi(port)
		if err != nil {
			c.Data["Errors"] = []string{"Port must be a valid number"}
			c.Data["DeviceId"] = deviceId
			c.Data["Integration"] = integration
			c.Data["Title"] = "Edit Integration"
			c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
			c.TplName = "integrations/edit.html"
			admin.StatisticsMap.AddStatistics("POST", "/devices/"+strconv.Itoa(deviceId)+"/integrations/"+strconv.Itoa(id)+"/edit", "&WebIntegrationController.Update", time.Since(start))
			return
		}
		valid.Range(portNum, 1, 65535, "port").Message("Port must be between 1 and 65535")
	}

	if valid.HasErrors() {
		var errors []string
		for _, err := range valid.Errors {
			errors = append(errors, err.Message)
		}
		c.Data["Errors"] = errors
		c.Data["DeviceId"] = deviceId
		c.Data["Integration"] = integration
		c.Data["Title"] = "Edit Integration"
		c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
		c.TplName = "integrations/edit.html"
		admin.StatisticsMap.AddStatistics("POST", "/devices/"+strconv.Itoa(deviceId)+"/integrations/"+strconv.Itoa(id)+"/edit", "&WebIntegrationController.Update", time.Since(start))
		return
	}

	// Update fields
	integration.IntegrationType = integrationType
	integration.Identifier = identifier
	integration.Host = host
	integration.Port = port
	integration.Protocol = protocol

	if err := models.UpdateIntegration(integration); err != nil {
		userId, _ := c.GetSession("user_id").(int)
		models.AddErrorLog("Error in Update: "+err.Error(), userId)
		c.Data["Error"] = "Failed to update integration: " + err.Error()
		c.Data["Redirect"] = "/devices/" + strconv.Itoa(deviceId)
		c.Data["Title"] = "Error"
		c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
		c.Data["RetryURL"] = "/devices/" + c.Ctx.Input.Param(":device_id") + "/integrations/" + c.Ctx.Input.Param(":id") + "/edit"
		logs.Error("Error in Update:", err)
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("POST", "/devices/"+strconv.Itoa(deviceId)+"/integrations/"+strconv.Itoa(id)+"/edit", "&WebIntegrationController.Update", time.Since(start))
		return
	}

	c.SetFlash("Integration updated successfully")
	c.Redirect("/devices/"+strconv.Itoa(deviceId), 302)
	admin.StatisticsMap.AddStatistics("POST", "/devices/"+strconv.Itoa(deviceId)+"/integrations/"+strconv.Itoa(id)+"/edit", "&WebIntegrationController.Update", time.Since(start))
}

// Delete handles the deletion of an integration
func (c *WebIntegrationController) Delete() {
	start := time.Now()

	deviceId, err := strconv.Atoi(c.Ctx.Input.Param(":device_id"))
	if err != nil {
		userId, _ := c.GetSession("user_id").(int)
		models.AddErrorLog("Error in Delete: "+err.Error(), userId)
		c.Data["Error"] = "Invalid device ID"
		c.Data["Redirect"] = "/devices"
		c.Data["Title"] = "Error"
		c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
		c.Data["RetryURL"] = "/devices/" + c.Ctx.Input.Param(":device_id") + "/integrations/" + c.Ctx.Input.Param(":id") + "/delete"
		logs.Error("Error in Delete:", err)
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("POST", "/devices/"+strconv.Itoa(deviceId)+"/integrations/"+c.Ctx.Input.Param(":id")+"/delete", "&WebIntegrationController.Delete", time.Since(start))
		return
	}

	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		userId, _ := c.GetSession("user_id").(int)
		models.AddErrorLog("Error in Delete: "+err.Error(), userId)
		c.Data["Error"] = "Invalid integration ID"
		c.Data["Redirect"] = "/devices/" + strconv.Itoa(deviceId)
		c.Data["Title"] = "Error"
		c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
		c.Data["RetryURL"] = "/devices/" + c.Ctx.Input.Param(":device_id") + "/integrations/" + c.Ctx.Input.Param(":id") + "/delete"
		logs.Error("Error in Delete:", err)
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("POST", "/devices/"+strconv.Itoa(deviceId)+"/integrations/"+strconv.Itoa(id)+"/delete", "&WebIntegrationController.Delete", time.Since(start))
		return
	}

	if err := models.DeleteIntegration(id); err != nil {
		userId, _ := c.GetSession("user_id").(int)
		models.AddErrorLog("Error in Delete: "+err.Error(), userId)
		c.Data["Error"] = "Failed to delete integration: " + err.Error()
		c.Data["Redirect"] = "/devices/" + strconv.Itoa(deviceId)
		c.Data["Title"] = "Error"
		c.Data["Timestamp"] = time.Now().Format(time.RFC3339)
		c.Data["RetryURL"] = "/devices/" + c.Ctx.Input.Param(":device_id") + "/integrations/" + c.Ctx.Input.Param(":id") + "/delete"
		logs.Error("Error in Delete:", err)
		c.TplName = "error.html"
		admin.StatisticsMap.AddStatistics("POST", "/devices/"+strconv.Itoa(deviceId)+"/integrations/"+strconv.Itoa(id)+"/delete", "&WebIntegrationController.Delete", time.Since(start))
		return
	}

	c.SetFlash("Integration deleted successfully")
	c.Redirect("/devices/"+strconv.Itoa(deviceId), 302)
	admin.StatisticsMap.AddStatistics("POST", "/devices/"+strconv.Itoa(deviceId)+"/integrations/"+strconv.Itoa(id)+"/delete", "&WebIntegrationController.Delete", time.Since(start))
}
