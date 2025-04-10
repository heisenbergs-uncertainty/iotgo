package controllers

import (
	"iotgo/models"
	"strconv"

	"github.com/beego/beego/v2/server/web"
)

type IntegrationController struct {
	web.Controller
}

func (c *IntegrationController) GetAllByDevice() {
	deviceId, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Ctx.Output.JSON(map[string]string{"error": "Invalid device ID"}, false, false)
		return
	}
	integrations, err := models.GetIntegrationsByDeviceId(deviceId)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Ctx.Output.JSON(map[string]string{"error": err.Error()}, false, false)
		return
	}
	c.Ctx.Output.JSON(integrations, false, false)
}

func (c *IntegrationController) Post() {
	deviceId, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Ctx.Output.JSON(map[string]string{"error": "Invalid device ID"}, false, false)
		return
	}
	var integration models.Integration
	if err := c.BindJSON(&integration); err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Ctx.Output.JSON(map[string]string{"error": "Invalid JSON"}, false, false)
		return
	}
	device, err := models.GetDeviceById(deviceId)
	if err != nil || device == nil {
		c.Ctx.Output.SetStatus(404)
		c.Ctx.Output.JSON(map[string]string{"error": "Device not found"}, false, false)
		return
	}
	integration.Device = device
	if err := models.AddIntegration(&integration); err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Ctx.Output.JSON(map[string]string{"error": err.Error()}, false, false)
		return
	}
	c.Ctx.Output.SetStatus(201)
	c.Ctx.Output.JSON(integration, false, false)
}

func (c *IntegrationController) GetOne() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Ctx.Output.JSON(map[string]string{"error": "Invalid ID"}, false, false)
		return
	}
	integration, err := models.GetIntegrationById(id)
	if err != nil || integration == nil {
		c.Ctx.Output.SetStatus(404)
		c.Ctx.Output.JSON(map[string]string{"error": "Integration not found"}, false, false)
		return
	}
	c.Ctx.Output.JSON(integration, false, false)
}

func (c *IntegrationController) Put() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Ctx.Output.JSON(map[string]string{"error": "Invalid ID"}, false, false)
		return
	}
	var integration models.Integration
	if err := c.BindJSON(&integration); err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Ctx.Output.JSON(map[string]string{"error": "Invalid JSON"}, false, false)
		return
	}
	integration.Id = id
	if err := models.UpdateIntegration(&integration); err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Ctx.Output.JSON(map[string]string{"error": err.Error()}, false, false)
		return
	}
	c.Ctx.Output.JSON(integration, false, false)
}

func (c *IntegrationController) Delete() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Ctx.Output.JSON(map[string]string{"error": "Invalid ID"}, false, false)
		return
	}
	if err := models.DeleteIntegration(id); err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Ctx.Output.JSON(map[string]string{"error": err.Error()}, false, false)
		return
	}
	c.Ctx.Output.SetStatus(204)
}
