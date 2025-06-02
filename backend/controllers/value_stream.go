package controllers

import (
	"app/dal"
	"app/model"
	"errors"
	"strconv"
)

type ValueStreamController struct {
	BaseController
}

// GetAll retrieves all value streams with pagination, filtering, and sorting (API)
func (c *ValueStreamController) GetAll() {
	limit, _ := c.GetInt("limit", 10)
	offset, _ := c.GetInt("offset", 0)
	nameFilter := c.GetString("name")
	nameSort := c.GetString("sort", "name")

	q := dal.Q
	query := q.ValueStream
	if nameFilter != "" {
		query.Where(q.ValueStream.Name.Like("%" + nameFilter + "%"))
	}

	switch nameSort {
	case "name":
		query.Order(q.ValueStream.Name.Asc())
	case "-name":
		query.Order(q.ValueStream.Name.Desc())
	}

	valueStreams, err := query.Offset(offset).Limit(limit).Find()
	if err != nil {
		c.PaginatedResponse([]model.ValueStream{}, 0, limit, offset, err)
		return
	}

	total, err := query.Count()
	c.PaginatedResponse(valueStreams, total, limit, offset, err)
}

// Get retrieves a single value stream by ID (API)
func (c *ValueStreamController) Get() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}

	q := dal.Q
	valueStream, err := q.ValueStream.Where(q.ValueStream.ID.Eq(uint(id))).First()
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}

	c.JSONResponse(valueStream, nil)
}

// Post creates a new value stream with validation (API)
func (c *ValueStreamController) Post() {
	var valueStream model.ValueStream
	if err := c.BindJSON(&valueStream); err != nil {
		c.JSONResponse(nil, err)
		return
	}

	if valueStream.Name == "" || valueStream.Type == "" {
		c.JSONResponse(nil, errors.New("name and type are required"))
		return
	}

	q := dal.Q
	if err := q.ValueStream.Create(&valueStream); err != nil {
		c.JSONResponse(nil, err)
		return
	}

	c.JSONResponse(valueStream, nil)
}

// Put updates an existing value stream with validation (API)
func (c *ValueStreamController) Put() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}

	var valueStream model.ValueStream
	if err := c.BindJSON(&valueStream); err != nil {
		c.JSONResponse(nil, err)
		return
	}

	if valueStream.Name == "" || valueStream.Type == "" {
		c.JSONResponse(nil, errors.New("name and type are required"))
		return
	}

	valueStream.ID = uint(id)

	q := dal.Q

	info, err := q.ValueStream.Where(q.ValueStream.ID.Eq(uint(id))).Updates(&valueStream)
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}
	if info.RowsAffected == 0 {
		c.JSONResponse(nil, errors.New("no rows affected"))
		return
	}

	c.JSONResponse(valueStream, nil)
}

// Delete removes a value stream by ID (API)
func (c *ValueStreamController) Delete() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}

	q := dal.Q

	info, err := q.ValueStream.Where(q.ValueStream.ID.Eq(uint(id))).Delete()
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}
	if info.RowsAffected == 0 {
		c.JSONResponse(nil, errors.New("no rows affected"))
		return
	}

	c.JSONResponse(map[string]string{"message": "Value Stream deleted successfully"}, nil)
}

// ListValueStreams renders the value stream management page (Web)
func (c *ValueStreamController) ListValueStreams() {
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

	c.TplName = "value_streams/index.html"
}

// NewValueStream renders the create value stream form (Web)
func (c *ValueStreamController) NewValueStream() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Redirect("/login", 302)
		return
	}

	c.TplName = "value_streams/new.html"
}

// EditValueStream renders the edit value stream form (Web)
func (c *ValueStreamController) EditValueStream() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Redirect("/login", 302)
		return
	}

	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Data["Error"] = "Invalid value stream ID"
		c.TplName = "value_streams/edit.html"
		return
	}

	valueStream, err := dal.ValueStream.Where(dal.ValueStream.ID.Eq(uint(id))).First()
	if err != nil {
		c.Data["Error"] = "Value Stream not found"
		c.TplName = "value_streams/edit.html"
		return
	}

	c.Data["ValueStream"] = valueStream
	c.TplName = "value_streams/edit.html"
}

// ViewValueStream renders the view value stream page (Web)
func (c *ValueStreamController) ViewValueStream() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Redirect("/login", 302)
		return
	}

	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Data["Error"] = "Invalid value stream ID"
		c.TplName = "value_streams/view.html"
		return
	}

	valueStream, err := dal.ValueStream.Where(dal.ValueStream.ID.Eq(uint(id))).First()
	if err != nil {
		c.Data["Error"] = "Value Stream not found"
		c.TplName = "value_streams/view.html"
		return
	}

	c.Data["ValueStream"] = valueStream
	c.TplName = "value_streams/view.html"
}
