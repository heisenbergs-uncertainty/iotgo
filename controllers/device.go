package controllers

import (
	"iotgo/models"
	"strconv"

	"github.com/beego/beego/v2/server/web"
)

type DeviceController struct {
	web.Controller
}

func (c *DeviceController) GetAll() {
	devices, err := models.GetAllDevices()
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Ctx.Output.JSON(map[string]string{"error": err.Error()}, false, false)
		return
	}
	c.Ctx.Output.JSON(devices, false, false)
}

func (c *DeviceController) Post() {
	var device models.Device
	if err := c.BindJSON(&device); err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Ctx.Output.JSON(map[string]string{"error": "Invalid JSON"}, false, false)
		return
	}
	if err := models.AddDevice(&device); err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Ctx.Output.JSON(map[string]string{"error": err.Error()}, false, false)
		return
	}
	c.Ctx.Output.SetStatus(201)
	c.Ctx.Output.JSON(device, false, false)
}

func (c *DeviceController) GetOne() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Ctx.Output.JSON(map[string]string{"error": "Invalid ID"}, false, false)
		return
	}
	device, err := models.GetDeviceById(id)
	if err != nil || device == nil {
		c.Ctx.Output.SetStatus(404)
		c.Ctx.Output.JSON(map[string]string{"error": "Device not found"}, false, false)
		return
	}
	c.Ctx.Output.JSON(device, false, false)
}

func (c *DeviceController) Put() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Ctx.Output.JSON(map[string]string{"error": "Invalid ID"}, false, false)
		return
	}
	var device models.Device
	if err := c.BindJSON(&device); err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Ctx.Output.JSON(map[string]string{"error": "Invalid JSON"}, false, false)
		return
	}
	device.Id = id
	if err := models.UpdateDevice(&device); err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Ctx.Output.JSON(map[string]string{"error": err.Error()}, false, false)
		return
	}
	c.Ctx.Output.JSON(device, false, false)
}

func (c *DeviceController) Delete() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Ctx.Output.JSON(map[string]string{"error": "Invalid ID"}, false, false)
		return
	}
	if err := models.DeleteDevice(id); err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Ctx.Output.JSON(map[string]string{"error": err.Error()}, false, false)
		return
	}
	c.Ctx.Output.SetStatus(204)
}
