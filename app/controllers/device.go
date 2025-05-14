package controllers

import (
	"app/dal"
	"app/model"
	"errors"
	"log"
	"strconv"
)

type DeviceController struct {
	BaseController
}

func (c *DeviceController) GetAll() {
	limit, _ := c.GetInt("limit", 10)
	offset, _ := c.GetInt("offset", 0)
	nameFilter := c.GetString("name")
	nameSort := c.GetString("sort", "name")

	q := dal.Q
	query := dal.Q.Device.Preload(q.Device.Site, q.Device.ValueStream)

	if nameFilter != "" {
		query.Where(q.Device.Name.Like("%" + nameFilter + "%"))
	}

	switch nameSort {
	case "name":
		query.Order(q.Device.Name.Asc())
	case "-name":
		query.Order(q.Device.Name.Desc())
	}

	devices, err := query.Offset(offset).Limit(limit).Find()
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}

	total, err := query.Count()
	c.PaginatedResponse(devices, total, limit, offset, err)
}

func (c *DeviceController) Get() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}

	q := dal.Q
	device, err := q.Device.Preload(q.Device.Site, q.Device.ValueStream).Where(q.Device.ID.Eq(uint(id))).First()
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}

	c.JSONResponse(device, err)
}

func (c *DeviceController) Post() {
	var device model.Device
	if err := c.BindJSON(&device); err != nil {
		c.JSONResponse(nil, err)
		return
	}

	// Validate required fields
	if device.Name == "" {
		c.JSONResponse(nil, errors.New("name is required"))
		return
	}

	q := dal.Q
	// Validate SiteID if provided
	if device.SiteID != nil {
		site, err := q.Site.Where(q.Site.ID.Eq(*device.SiteID)).First()
		if err != nil {
			c.JSONResponse(nil, errors.New("invalid site id"))
			return
		}
		device.Site = site
	}

	// Validate ValueStreamID if provided
	if device.ValueStreamID != nil {
		vs, err := q.ValueStream.Where(q.ValueStream.ID.Eq(*device.ValueStreamID)).First()
		if err != nil {
			c.JSONResponse(nil, errors.New("invalid value stream id"))
			return
		}
		device.ValueStream = vs
	}

	if err := q.Device.Create(&device); err != nil {
		c.JSONResponse(nil, err)
		return
	}

	c.JSONResponse(device, nil)
}

func (c *DeviceController) Put() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}

	var device model.Device
	if err := c.BindJSON(&device); err != nil {
		c.JSONResponse(nil, err)
		return
	}

	if device.Name == "" {
		c.JSONResponse(nil, errors.New("name is required"))
		return
	}

	device.ID = uint(id)
	q := dal.Q
	info, err := q.Device.Where(q.Device.ID.Eq(uint(id))).Updates(&device)
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}

	if info.RowsAffected == 0 {
		c.JSONResponse(nil, errors.New("no rows affected"))
		return
	}

	c.JSONResponse(device, err)
}

// Delete removes a device by ID
func (c *DeviceController) Delete() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	log.Print(id)
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}

	q := dal.Q
	info, err := q.Device.WithContext(c.Ctx.Request.Context()).Where(q.Device.ID.Eq(uint(id))).Delete()
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}
	if info.RowsAffected == 0 {
		c.JSONResponse(nil, errors.New("no rows affected"))
		return
	}

	c.JSONResponse(map[string]string{"message": "Device deleted successfully"}, nil)
}

// ListDevices renders the device management page (Web)
func (c *DeviceController) ListDevices() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Redirect("/login", 302)
		return
	}

	q := dal.Q
	sites, err := q.Site.Find()
	if err != nil {
		c.Data["Error"] = "Failed to load sites"
	}

	valueStreams, err := q.ValueStream.Find()
	if err != nil {
		c.Data["Error"] = "Failed to load value streams"
	}

	apiKey, err := q.ApiKey.Where(q.ApiKey.UserID.Eq(userID.(uint)), q.ApiKey.IsActive.Is(true)).First()
	if err != nil || apiKey.ID == 0 {
		c.Data["ApiToken"] = ""
		c.Data["TokenError"] = "No active API key found. Please generate one."
	} else {
		c.Data["ApiToken"] = apiKey.Token
	}

	devices, err := q.Device.Find()
	if err != nil {
		c.Data["Error"] = "Failed to load devices"
	}

	for _, device := range devices {
		log.Printf("%+v\n", device)
	}
	c.Data["Devices"] = devices
	c.Data["Sites"] = sites
	c.Data["ValueStreams"] = valueStreams

	c.TplName = "devices.html"
}
