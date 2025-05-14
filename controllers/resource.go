package controllers

import (
	"app/dal"
	"app/model"
	"encoding/json"
	"errors"
	"strconv"
)

type ResourceController struct {
	BaseController
}

func (c *ResourceController) GetAll() {
	limit, _ := c.GetInt("limit", 10)
	offset, _ := c.GetInt("offset", 0)
	nameSort := c.GetString("sort", "name")

	q := dal.Q
	query := q.Resource

	if nameSort == "name" {
		query.Order(q.Resource.Name.Asc())
	} else {
		query.Order(q.Resource.Name.Desc())
	}

	resources, err := query.Offset(offset).Limit(limit).Find()
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}

	total, err := query.Count()

	c.PaginatedResponse(resources, total, limit, offset, err)
}

func (c *ResourceController) Get() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}

	q := dal.Q
	resource, err := q.Resource.Where(q.Resource.ID.Eq(uint(id))).First()
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}

	c.JSONResponse(resource, err)
}

func (c *ResourceController) Post() {
	var resource model.Resource
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &resource); err != nil {
		c.JSONResponse(nil, err)
		return
	}

	if resource.Name == "" {
		c.JSONResponse(nil, errors.New("name is required"))
		return
	}

	q := dal.Q
	if err := q.Resource.Create(&resource); err != nil {
		c.JSONResponse(nil, err)
		return
	}

	c.JSONResponse(resource, nil)
}

func (c *ResourceController) Put() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}

	var resource model.Resource
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &resource); err != nil {
		c.JSONResponse(nil, err)
		return
	}

	if resource.Name == "" {
		c.JSONResponse(nil, errors.New("name is required"))
		return
	}

	resource.ID = uint(id)

	q := dal.Q
	info, err := q.Resource.Where(q.Resource.ID.Eq(uint(id))).Updates(resource)
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}
	if info.RowsAffected == 0 {
		c.JSONResponse(nil, errors.New("resource not found"))
		return
	}

	c.JSONResponse(resource, info.Error)
}

func (c *ResourceController) Delete() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}

	q := dal.Q
	info, err := q.Resource.Where(q.Resource.ID.Eq(uint(id))).Delete()
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}
	if info.RowsAffected == 0 {
		c.JSONResponse(nil, errors.New("no rows affected"))
		return
	}

	c.JSONResponse(nil, info.Error)
}
