package controllers

import (
	"app/dal"
	"app/model"
	"errors"
	"strconv"
)

type DevicePlatformController struct {
	BaseController
}

// GetAll retrieves all device-platform associations for a device (API)
func (c *DevicePlatformController) GetAll() {
	deviceID, err := strconv.Atoi(c.Ctx.Input.Param(":device_id"))
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}

	limit, _ := c.GetInt("limit", 10)
	offset, _ := c.GetInt("offset", 0)

	q := dal.Q
	query := q.DevicePlatform.Where(q.DevicePlatform.DeviceID.Eq(uint(deviceID)))
	associations, err := query.Offset(offset).Limit(limit).Find()
	if err != nil {
		c.PaginatedResponse([]model.DevicePlatform{}, 0, limit, offset, err)
		return
	}

	total, err := query.Count()
	c.PaginatedResponse(associations, total, limit, offset, err)
}

// Post creates a new device-platform association (API)
func (c *DevicePlatformController) Post() {
	deviceID, err := strconv.Atoi(c.Ctx.Input.Param(":device_id"))
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}

	var association model.DevicePlatform
	if err := c.BindJSON(&association); err != nil {
		c.JSONResponse(nil, err)
		return
	}

	if association.PlatformID == 0 || association.DeviceAlias == "" {
		c.JSONResponse(nil, errors.New("platform_id and device_alias are required"))
		return
	}

	association.DeviceID = uint(deviceID)

	q := dal.Q
	// Validate device and platform existence
	_, err = q.Device.Where(dal.Device.ID.Eq(uint(deviceID))).First()
	if err != nil {
		c.JSONResponse(nil, errors.New("device not found"))
		return
	}

	_, err = q.Platform.Where(dal.Platform.ID.Eq(association.PlatformID)).First()
	if err != nil {
		c.JSONResponse(nil, errors.New("platform not found"))
		return
	}

	// Check for duplicate DeviceAlias for this PlatformID
	count, err := q.DevicePlatform.Where(
		q.DevicePlatform.PlatformID.Eq(association.PlatformID),
		q.DevicePlatform.DeviceAlias.Eq(association.DeviceAlias),
	).Count()
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}
	if count > 0 {
		c.JSONResponse(nil, errors.New("device_alias already exists for this platform"))
		return
	}

	if err := q.DevicePlatform.Create(&association); err != nil {
		c.JSONResponse(nil, err)
		return
	}

	c.JSONResponse(association, nil)
}

// Delete removes a device-platform association (API)
func (c *DevicePlatformController) Delete() {
	deviceID, err := strconv.Atoi(c.Ctx.Input.Param(":device_id"))
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}

	platformID, err := strconv.Atoi(c.Ctx.Input.Param(":platform_id"))
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}

	q := dal.Q
	_, err = q.DevicePlatform.Where(
		q.DevicePlatform.DeviceID.Eq(uint(deviceID)),
		q.DevicePlatform.PlatformID.Eq(uint(platformID)),
	).Delete()
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}

	c.JSONResponse(map[string]string{"message": "Association deleted successfully"}, nil)
}

// ListAssociations renders the device-platform association management page (Web)
func (c *DevicePlatformController) ListAssociations() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Redirect("/login", 302)
		return
	}

	deviceID, err := strconv.Atoi(c.Ctx.Input.Param(":device_id"))
	if err != nil {
		c.Data["Error"] = "Invalid device ID"
		c.TplName = "device_platforms.html"
		c.Layout = "layout.html"
		return
	}

	q := dal.Q

	// Fetch device
	device, err := q.Device.Where(dal.Device.ID.Eq(uint(deviceID))).First()
	if err != nil {
		c.Data["Error"] = "Device not found"
		c.TplName = "device_platforms.html"
		c.Layout = "layout.html"
		return
	}

	// Fetch all platforms for dropdown
	platforms, err := q.Platform.Find()
	if err != nil {
		c.Data["Error"] = "Failed to load platforms"
	}

	// Fetch user's active API key
	apiKey, err := q.ApiKey.Where(
		dal.ApiKey.UserID.Eq(userID.(uint)),
		dal.ApiKey.IsActive.Is(true),
	).First()
	if err != nil || apiKey.ID == 0 {
		c.Data["ApiToken"] = ""
		c.Data["TokenError"] = "No active API key found. Please generate one."
	} else {
		c.Data["ApiToken"] = apiKey.Token
	}

	c.Data["Device"] = device
	c.Data["Platforms"] = platforms
	c.TplName = "device_platforms.html"
	c.Layout = "layout.html"
}
