package controllers

import (
	"app/dal"
	"app/model"
	"errors"
	"strconv"
)

type SiteController struct {
	BaseController
}

// GetAll retrieves all sites with pagination, filtering, and sorting (API)
func (c *SiteController) GetAll() {
	limit, _ := c.GetInt("limit", 10)
	offset, _ := c.GetInt("offset", 0)
	nameFilter := c.GetString("name")
	nameSort := c.GetString("sort", "name")

	q := dal.Q
	query := q.Site
	if nameFilter != "" {
		query.Where(q.Site.Name.Like("%" + nameFilter + "%"))
	}

	switch nameSort {
	case "name":
		query.Order(q.Site.Name.Asc())
	case "-name":
		query.Order(q.Site.Name.Desc())
	}

	sites, err := query.Offset(offset).Limit(limit).Find()
	if err != nil {
		c.PaginatedResponse([]model.Site{}, 0, limit, offset, err)
		return
	}

	total, err := query.Count()
	c.PaginatedResponse(sites, total, limit, offset, err)
}

// Get retrieves a single site by ID (API)
func (c *SiteController) Get() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}

	q := dal.Q
	site, err := q.Site.Where(q.Site.ID.Eq(uint(id))).First()
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}

	c.JSONResponse(site, nil)
}

// Post creates a new site with validation (API)
func (c *SiteController) Post() {
	var site model.Site
	if err := c.BindJSON(&site); err != nil {
		c.JSONResponse(nil, err)
		return
	}

	if site.Name == "" || site.Address == "" || site.City == "" || site.State == "" || site.Country == "" {
		c.JSONResponse(nil, errors.New("name, address, city, state, and country are required"))
		return
	}

	q := dal.Q
	if err := q.Site.Create(&site); err != nil {
		c.JSONResponse(nil, err)
		return
	}

	c.JSONResponse(site, nil)
}

// Put updates an existing site with validation (API)
func (c *SiteController) Put() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}

	var site model.Site
	if err := c.BindJSON(&site); err != nil {
		c.JSONResponse(nil, err)
		return
	}

	if site.Name == "" || site.Address == "" || site.City == "" || site.State == "" || site.Country == "" {
		c.JSONResponse(nil, errors.New("name, address, city, state, and country are required"))
		return
	}

	site.ID = uint(id)
	q := dal.Q
	info, err := q.Site.Where(q.Site.ID.Eq(uint(id))).Updates(&site)
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}
	if info.RowsAffected == 0 {
		c.JSONResponse(nil, errors.New("no rows affected"))
		return
	}

	c.JSONResponse(site, info.Error)
}

// Delete removes a site by ID (API)
func (c *SiteController) Delete() {
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}

	q := dal.Q
	info, err := q.Site.Where(q.Site.ID.Eq(uint(id))).Delete()
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}
	if info.RowsAffected == 0 {
		c.JSONResponse(nil, errors.New("no rows affected"))
		return
	}

	c.JSONResponse(map[string]string{"message": "Site deleted successfully"}, info.Error)
}

// ListSites renders the site management page (Web)
func (c *SiteController) ListSites() {
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

	c.TplName = "sites/index.html"
}

// NewSite renders the create site form (Web)
func (c *SiteController) NewSite() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Redirect("/login", 302)
		return
	}

	c.TplName = "sites/new.html"
}

// EditSite renders the edit site form (Web)
func (c *SiteController) EditSite() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Redirect("/login", 302)
		return
	}

	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Data["Error"] = "Invalid site ID"
		c.TplName = "sites/edit.html"
		return
	}

	site, err := dal.Site.Where(dal.Site.ID.Eq(uint(id))).First()
	if err != nil {
		c.Data["Error"] = "Site not found"
		c.TplName = "sites/edit.html"
		return
	}

	c.Data["Site"] = site
	c.TplName = "sites/edit.html"
}

// ViewSite renders the view site page (Web)
func (c *SiteController) ViewSite() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Redirect("/login", 302)
		return
	}

	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		c.Data["Error"] = "Invalid site ID"
		c.TplName = "sites/view.html"
		return
	}

	site, err := dal.Site.Where(dal.Site.ID.Eq(uint(id))).First()
	if err != nil {
		c.Data["Error"] = "Site not found"
		c.TplName = "sites/view.html"
		return
	}

	c.Data["Site"] = site
	c.TplName = "sites/view.html"
}
