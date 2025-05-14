package controllers

import (
	"app/dal"
	"app/drivers"
	"app/model"
	"context"
	"errors"
	"strconv"
)

type PlatformController struct {
	BaseController
}

// GetAll retrieves all platforms with pagination, filtering, and sorting (API)
func (c *PlatformController) GetAll() {
	limit, _ := c.GetInt("limit", 10)
	offset, _ := c.GetInt("offset", 0)
	nameFilter := c.GetString("name")
	nameSort := c.GetString("sort", "name")

	q := dal.Q
	query := q.Platform
	if nameFilter != "" {
		query.Where(q.Platform.Name.Like("%" + nameFilter + "%"))
	}

	switch nameSort {
	case "name":
		query.Order(q.Platform.Name.Asc())
	case "-name":
		query.Order(q.Platform.Name.Desc())
	}

	platforms, err := query.Offset(offset).Limit(limit).Find()
	if err != nil {
		c.PaginatedResponse([]model.Platform{}, 0, limit, offset, err)
		return
	}

	total, err := query.Count()
	c.PaginatedResponse(platforms, total, limit, offset, err)
}

// Get retrieves a single platform by ID (API)
func (c *PlatformController) Get() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}

	q := dal.Q
	platform, err := q.Platform.Where(q.Platform.ID.Eq(uint(id))).First()
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}

	c.JSONResponse(platform, nil)
}

// Post creates a new platform with validation (API)
func (c *PlatformController) Post() {
	var platform model.Platform
	if err := c.BindJSON(&platform); err != nil {
		c.JSONResponse(nil, err)
		return
	}

	if platform.Name == "" || platform.Type == "" {
		c.JSONResponse(nil, errors.New("name and type are required"))
		return
	}

	q := dal.Q
	if err := q.Platform.Create(&platform); err != nil {
		c.JSONResponse(nil, err)
		return
	}

	c.JSONResponse(platform, nil)
}

// Put updates an existing platform with validation (API)
func (c *PlatformController) Put() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}

	var platform model.Platform
	if err := c.BindJSON(&platform); err != nil {
		c.JSONResponse(nil, err)
		return
	}

	if platform.Name == "" || platform.Type == "" {
		c.JSONResponse(nil, errors.New("name and type are required"))
		return
	}

	platform.ID = uint(id)

	q := dal.Q
	info, err := q.Platform.Where(q.Platform.ID.Eq(uint(id))).Updates(&platform)
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}
	if info.RowsAffected == 0 {
		c.JSONResponse(nil, errors.New("no rows affected"))
		return
	}

	c.JSONResponse(platform, info.Error)
}

// Delete removes a platform by ID (API)
func (c *PlatformController) Delete() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}

	q := dal.Q
	info, err := q.Platform.Where(q.Platform.ID.Eq(uint(id))).Delete()
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}
	if info.RowsAffected == 0 {
		c.JSONResponse(nil, errors.New("no rows affected"))
		return
	}

	c.JSONResponse(map[string]string{"message": "Platform deleted successfully"}, info.Error)
}

// FetchDeviceData fetches data for a specific device on a specific platform (API)
func (c *PlatformController) FetchDeviceData() {
	platformID, _ := strconv.Atoi(c.Ctx.Input.Param(":platform_id"))
	deviceID, _ := strconv.Atoi(c.Ctx.Input.Param(":device_id"))

	q := dal.Q

	// Find the DevicePlatform entry
	dp, err := q.DevicePlatform.Where(
		q.DevicePlatform.PlatformID.Eq(uint(platformID)),
		q.DevicePlatform.DeviceID.Eq(uint(deviceID)),
	).First()
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}

	// Get the platform
	platform, err := q.Platform.Where(q.Platform.ID.Eq(uint(platformID))).First()
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}

	// Get the driver
	driver, err := drivers.GetDriver(platform.Type, platform.Metadata)
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}

	ctx := context.Background()
	err = driver.Connect(ctx)
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}
	defer driver.Disconnect(ctx)

	data, err := driver.FetchData(ctx, dp.DeviceAlias)
	c.JSONResponse(data, err)
}

// ListPlatforms renders the platform management page (Web)
func (c *PlatformController) ListPlatforms() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Redirect("/login", 302)
		return
	}

	q := dal.Q

	// Fetch user's active API key
	apiKey, err := q.ApiKey.Where(
		q.ApiKey.UserID.Eq(userID.(uint)),
		q.ApiKey.IsActive.Is(true),
	).First()
	if err != nil || apiKey.ID == 0 {
		c.Data["ApiToken"] = ""
		c.Data["TokenError"] = "No active API key found. Please generate one."
	} else {
		c.Data["ApiToken"] = apiKey.Token
	}

	c.TplName = "platforms.html"
}
