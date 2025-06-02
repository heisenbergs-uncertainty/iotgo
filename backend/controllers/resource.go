package controllers

import (
	"app/dal"
	"app/drivers"
	"app/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"path"
	"strconv"
	"strings"

	"github.com/beego/beego/logs"
)

type ResourceController struct {
	BaseController
}

func (c *ResourceController) GetResources(platformId uint, nameSort string, limit *int, offset *int, nameFilter *string) ([]*model.Resource, error) {
	q := dal.Q
	query := q.Resource.Where(q.Resource.PlatformID.Eq(uint(platformId)))

	if nameFilter != nil && *nameFilter != "" {
		query.Where(q.Resource.Name.Like("%" + *nameFilter + "%"))
	}

	switch nameSort {
	case "name":
		query.Order(q.Resource.Name.Asc())
	case "-name":
		query.Order(q.Resource.Name.Desc())
	default:
		query.Order(q.Resource.Name.Asc())
	}

	if limit != nil && *limit > 0 {
		query.Limit(*limit)
	}
	if offset != nil && *offset > 0 {
		query.Offset(*offset)
	}

	resources, err := query.Find()
	if err != nil {
		return nil, err
	}

	return resources, nil
}

func (c *ResourceController) GetAll() {
	platformID, err := strconv.Atoi(c.Ctx.Input.Param(":platform_id"))
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}

	limit, _ := c.GetInt("limit", 10)
	offset, _ := c.GetInt("offset", 0)
	nameFilter := c.GetString("name")
	nameSort := c.GetString("sort", "name")

	resources, err := c.GetResources(uint(platformID), nameSort, &limit, &offset, &nameFilter)
	if err != nil {
		c.PaginatedResponse([]model.Resource{}, 0, limit, offset, err)
		return
	}
	logs.Info("Fetched resources for platform:", platformID)

	total, err := dal.Q.Resource.Count()
	if err != nil {
		c.PaginatedResponse(resources, 0, limit, offset, err)
		return
	}

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

// GetResourceForEdit retrieves a resource for editing, parsing details for form pre-filling (API)
func (c *ResourceController) GetResourceForEdit() {
	logs.Info("Received GET request to fetch resource %s for editing", c.Ctx.Input.Param(":id"))
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		logs.Error("Invalid resource ID:", err)
		c.JSONResponse(nil, err)
		return
	}

	resource, err := dal.Resource.Where(dal.Resource.ID.Eq(uint(id))).First()
	if err != nil {
		logs.Error("Failed to find resource:", err)
		c.JSONResponse(nil, err)
		return
	}

	// Prepare response with parsed details
	response := map[string]interface{}{
		"id":   resource.ID,
		"name": resource.Name,
		"type": resource.Type,
	}

	if resource.Type == "rest_endpoint" {
		var details model.RESTResourceDetails
		if err := json.Unmarshal([]byte(resource.Details), &details); err != nil {
			logs.Error("Failed to parse REST resource details:", err)
			c.JSONResponse(nil, err)
			return
		}
		response["details"] = details
	} else if resource.Type == "influxdb_query" {
		var details model.InfluxDBResourceDetails
		if err := json.Unmarshal([]byte(resource.Details), &details); err != nil {
			logs.Error("Failed to parse InfluxDB resource details:", err)
			c.JSONResponse(nil, err)
			return
		}
		response["details"] = details
	} else {
		err := errors.New("unsupported resource type")
		logs.Error("Validation failed:", err)
		c.JSONResponse(nil, err)
		return
	}

	logs.Info("Resource %d fetched for editing", resource.ID)
	c.JSONResponse(response, nil)
}

// validateRESTResourceDetails validates and sanitizes REST endpoint details
func (c *ResourceController) validateRESTResourceDetails(detailsJSON string) error {
	var details model.RESTResourceDetails
	if err := json.Unmarshal([]byte(detailsJSON), &details); err != nil {
		return errors.New("invalid details JSON")
	}

	// Validate method
	validMethods := map[string]bool{"GET": true, "POST": true, "PUT": true, "DELETE": true}
	if details.Method == "" {
		return errors.New("method is required for rest_endpoint")
	}
	if !validMethods[strings.ToUpper(details.Method)] {
		return errors.New("method must be GET, POST, PUT, or DELETE")
	}
	details.Method = strings.ToUpper(details.Method)

	// Validate and sanitize path
	if details.Path == "" {
		return errors.New("path is required for rest_endpoint")
	}
	cleanPath := path.Clean("/" + strings.TrimPrefix(details.Path, "/"))
	if cleanPath == "/" {
		return errors.New("path must not be empty")
	}
	details.Path = cleanPath

	// Validate headers and query parameters
	if details.Headers == nil {
		details.Headers = make(map[string]string)
	}
	if details.QueryParams == nil {
		details.QueryParams = make(map[string]string)
	}
	for key, value := range details.Headers {
		if strings.TrimSpace(key) == "" {
			return errors.New("header keys must not be empty")
		}
		details.Headers[strings.TrimSpace(key)] = strings.TrimSpace(value)
	}
	for key, value := range details.QueryParams {
		if strings.TrimSpace(key) == "" {
			return errors.New("query parameter keys must not be empty")
		}
		details.QueryParams[strings.TrimSpace(key)] = strings.TrimSpace(value)
	}

	// Validate body (optional, must be valid JSON if provided)
	if details.Body != "" {
		var body interface{}
		if err := json.Unmarshal([]byte(details.Body), &body); err != nil {
			return errors.New("body must be valid JSON")
		}
	}

	// Re-serialize sanitized details
	serialized, err := json.Marshal(details)
	if err != nil {
		return errors.New("failed to serialize sanitized details")
	}
	c.Ctx.Input.SetData("sanitized_details", string(serialized))
	return nil
}

// validateInfluxDBResourceDetails validates and sanitizes InfluxDB query details
func (c *ResourceController) validateInfluxDBResourceDetails(detailsJSON string) error {
	var details model.InfluxDBResourceDetails
	if err := json.Unmarshal([]byte(detailsJSON), &details); err != nil {
		return errors.New("invalid details JSON")
	}

	// Validate required fields
	if details.Bucket == "" {
		return errors.New("bucket is required for influxdb_query")
	}
	if details.Measurement == "" {
		return errors.New("measurement is required for influxdb_query")
	}
	if details.Field == "" {
		return errors.New("field is required for influxdb_query")
	}
	if details.TimeRange == "" {
		return errors.New("time_range is required for influxdb_query")
	}

	// Validate time_range format (e.g., "-1h", "-30m", "-1d")
	if !strings.HasPrefix(details.TimeRange, "-") || !strings.ContainsAny(details.TimeRange, "smhdwy") {
		return errors.New("time_range must be a negative duration (e.g., '-1h', '-30m', '-1d')")
	}

	// Sanitize fields
	details.Bucket = strings.TrimSpace(details.Bucket)
	details.Measurement = strings.TrimSpace(details.Measurement)
	details.Field = strings.TrimSpace(details.Field)
	details.TimeRange = strings.TrimSpace(details.TimeRange)

	// Re-serialize sanitized details
	serialized, err := json.Marshal(details)
	if err != nil {
		return errors.New("failed to serialize sanitized details")
	}
	c.Ctx.Input.SetData("sanitized_details", string(serialized))
	return nil
}

func (c *ResourceController) Post() {
	logs.Info("Received POST request to create resource for platform %s", c.Ctx.Input.Param(":platform_id"))

	platformID, err := strconv.Atoi(c.Ctx.Input.Param(":platform_id"))
	if err != nil {
		logs.Error("Invalid platform ID:", err)
		c.JSONResponse(nil, err)
		return
	}

	var resource model.Resource
	if err := c.BindJSON(&resource); err != nil {
		logs.Error("Failed to bind JSON:", err)
		c.JSONResponse(nil, err)
		return
	}

	logs.Debug("Resource input:", resource)

	if resource.Name == "" || resource.Type == "" || resource.Details == "" {
		err := errors.New("name, type, and details are required")
		logs.Error("Validation failed:", err)
		c.JSONResponse(nil, err)
		return
	}

	// Validate resource details based on type
	switch resource.Type {
	case "rest_endpoint":
		if err := c.validateRESTResourceDetails(resource.Details); err != nil {
			logs.Error("REST resource details validation failed:", err)
			c.JSONResponse(nil, err)
			return
		}
		resource.Details = c.Ctx.Input.GetData("sanitized_details").(string)
	case "influxdb_query":
		if err := c.validateInfluxDBResourceDetails(resource.Details); err != nil {
			logs.Error("InfluxDB resource details validation failed:", err)
			c.JSONResponse(nil, err)
			return
		}
		resource.Details = c.Ctx.Input.GetData("sanitized_details").(string)
	default:
		err := errors.New("unsupported resource type")
		logs.Error("Validation failed:", err)
		c.JSONResponse(nil, err)
		return
	}

	resource.PlatformID = uint(platformID)

	q := dal.Q
	if err := q.Resource.Create(&resource); err != nil {
		logs.Error("Failed to create resource:", err)
		c.JSONResponse(nil, err)
		return
	}

	logs.Info("Resource created successfully:", resource.ID)
	c.JSONResponse(resource, nil)
}

// BulkPost creates multiple resources for a platform (API)
func (c *ResourceController) BulkPost() {
	logs.Info("Received POST request to create multiple resources for platform %s", c.Ctx.Input.Param(":platform_id"))

	platformID, err := strconv.Atoi(c.Ctx.Input.Param(":platform_id"))
	if err != nil {
		logs.Error("Invalid platform ID:", err)
		c.JSONResponse(nil, err)
		return
	}

	var resources []model.Resource
	if err := c.BindJSON(&resources); err != nil {
		logs.Error("Failed to bind JSON:", err)
		c.JSONResponse(nil, err)
		return
	}

	var createdResources []model.Resource
	var errorsList []string

	for i, resource := range resources {
		if resource.Name == "" || resource.Type == "" || resource.Details == "" {
			err := fmt.Errorf("resource %d: name, type, and details are required", i)
			logs.Error("Validation failed:", err)
			errorsList = append(errorsList, err.Error())
			continue
		}

		// Validate resource details based on type
		switch resource.Type {
		case "rest_endpoint":
			if err := c.validateRESTResourceDetails(resource.Details); err != nil {
				logs.Error("Resource %d REST details validation failed: %v", i, err)
				errorsList = append(errorsList, fmt.Sprintf("resource %d: %v", i, err))
				continue
			}
			resource.Details = c.Ctx.Input.GetData("sanitized_details").(string)
		case "influxdb_query":
			if err := c.validateInfluxDBResourceDetails(resource.Details); err != nil {
				logs.Error("Resource %d InfluxDB details validation failed: %v", i, err)
				errorsList = append(errorsList, fmt.Sprintf("resource %d: %v", i, err))
				continue
			}
			resource.Details = c.Ctx.Input.GetData("sanitized_details").(string)
		default:
			err := fmt.Errorf("resource %d: unsupported resource type", i)
			logs.Error("Validation failed:", err)
			errorsList = append(errorsList, err.Error())
			continue
		}

		resource.PlatformID = uint(platformID)

		q := dal.Q
		if err := q.Resource.Create(&resource); err != nil {
			logs.Error("Failed to create resource %d: %v", i, err)
			errorsList = append(errorsList, fmt.Sprintf("resource %d: %v", i, err))
			continue
		}

		createdResources = append(createdResources, resource)
	}

	response := map[string]interface{}{
		"created": createdResources,
	}
	if len(errorsList) > 0 {
		response["errors"] = errorsList
	}

	logs.Info("Bulk resource creation completed: %d created, %d errors", len(createdResources), len(errorsList))
	c.JSONResponse(response, nil)
}

func (c *ResourceController) Put() {
	logs.Info("Received PUT request to update resource %s", c.Ctx.Input.Param(":id"))
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		logs.Error("Invalid resource ID:", err)
		c.JSONResponse(nil, err)
		return
	}

	var resource model.Resource
	if err := c.BindJSON(&resource); err != nil {
		logs.Error("Failed to bind JSON:", err)
		c.JSONResponse(nil, err)
		return
	}

	logs.Debug("Resource input:", resource)

	if resource.Name == "" || resource.Type == "" || resource.Details == "" {
		err := errors.New("name, type, and details are required")
		logs.Error("Validation failed:", err)
		c.JSONResponse(nil, err)
		return
	}

	// Validate resource details based on type
	switch resource.Type {
	case "rest_endpoint":
		if err := c.validateRESTResourceDetails(resource.Details); err != nil {
			logs.Error("REST resource details validation failed:", err)
			c.JSONResponse(nil, err)
			return
		}
		resource.Details = c.Ctx.Input.GetData("sanitized_details").(string)
	case "influxdb_query":
		if err := c.validateInfluxDBResourceDetails(resource.Details); err != nil {
			logs.Error("InfluxDB resource details validation failed:", err)
			c.JSONResponse(nil, err)
			return
		}
		resource.Details = c.Ctx.Input.GetData("sanitized_details").(string)
	default:
		err := errors.New("unsupported resource type")
		logs.Error("Validation failed:", err)
		c.JSONResponse(nil, err)
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
		logs.Error("Failed to update resource:", err)
		c.JSONResponse(nil, errors.New("resource not found"))
		return
	}

	logs.Info("Resource updated successfully:", resource.ID)
	c.JSONResponse(resource, info.Error)
}

func (c *ResourceController) Delete() {
	logs.Info("Received DELETE request to delete resource %s", c.Ctx.Input.Param(":id"))
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		logs.Error("Invalid resource ID:", err)
		c.JSONResponse(nil, err)
		return
	}

	q := dal.Q
	info, err := q.Resource.Where(q.Resource.ID.Eq(uint(id))).Delete()
	if err != nil {
		logs.Error("Failed to delete resource:", err)
		c.JSONResponse(nil, err)
		return
	}

	if info.RowsAffected == 0 {
		c.JSONResponse(nil, errors.New("no rows affected"))
		return
	}

	logs.Info("Resource deleted successfully:", id)
	c.JSONResponse(map[string]string{"message": "Resource deleted sucessfully"}, info.Error)
}

// TestResource tests a resource's query or endpoint (API)
func (c *ResourceController) TestResource() {
	logs.Info("Received POST request to test resource %s", c.Ctx.Input.Param(":id"))
	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		logs.Error("Invalid resource ID:", err)
		c.JSONResponse(nil, err)
		return
	}

	q := dal.Q
	// Fetch the resource
	resource, err := q.Resource.Where(q.Resource.ID.Eq(uint(id))).First()
	if err != nil {
		logs.Error("Failed to find resource:", err)
		c.JSONResponse(nil, err)
		return
	}

	// Fetch the associated platform
	platform, err := q.Platform.Where(q.Platform.ID.Eq(uint(resource.PlatformID))).First()
	if err != nil {
		logs.Error("Failed to find platform:", err)
		c.JSONResponse(nil, err)
		return
	}

	// Get the driver
	driver, err := drivers.GetDriver(platform.Type, platform.Metadata)
	if err != nil {
		logs.Error("Failed to get driver:", err)
		c.JSONResponse(nil, err)
		return
	}

	ctx := context.Background()
	err = driver.Connect(ctx)
	if err != nil {
		logs.Error("Failed to connect driver:", err)
		c.JSONResponse(nil, err)
		return
	}
	defer driver.Disconnect(ctx)

	// Test the resource
	var result interface{}
	if (platform.Type == "REST" && resource.Type == "rest_endpoint") || (platform.Type == "InfluxDB" && resource.Type == "influxdb_query") {
		// For REST, log the constructed URL
		if platform.Type == "REST" {
			var details model.RESTResourceDetails
			if err := json.Unmarshal([]byte(resource.Details), &details); err != nil {
				logs.Error("Failed to parse REST resource details: %v", err)
				c.JSONResponse(nil, fmt.Errorf("invalid resource details: %w", err))
				return
			}
			var config model.RESTMetadata
			if err := json.Unmarshal([]byte(platform.Metadata), &config); err != nil {
				logs.Error("Failed to parse REST platform metadata: %v", err)
				c.JSONResponse(nil, fmt.Errorf("invalid platform metadata: %w", err))
				return
			}
			baseURL, err := url.Parse(config.BaseEndpoint)
			if err != nil {
				logs.Error("Invalid base URL %s: %v", config.BaseEndpoint, err)
				c.JSONResponse(nil, fmt.Errorf("invalid base URL: %w", err))
				return
			}
			pathURL, err := url.Parse(details.Path)
			if err != nil {
				logs.Error("Invalid path %s: %v", details.Path, err)
				c.JSONResponse(nil, fmt.Errorf("invalid path: %w", err))
				return
			}
			fullURL := baseURL.ResolveReference(pathURL)
			query := fullURL.Query()
			for key, value := range details.QueryParams {
				query.Set(key, value)
			}
			fullURL.RawQuery = query.Encode()
			logs.Info("Testing REST resource %s: %s %s", resource.Name, details.Method, fullURL.String())
		}

		result, err = driver.FetchData(ctx, resource.Details)
		if err != nil {
			logs.Error("Failed to test resource %s: %v", resource.Name, err)
			c.JSONResponse(map[string]interface{}{
				"resource_id": resource.ID,
				"error":       err.Error(),
			}, fmt.Errorf("test failed: %w", err))
			return
		}
		logs.Debug("Test result for resource %s: %v", resource.Name, result)
	} else {
		err = fmt.Errorf("incompatible resource type %s for platform type %s", resource.Type, platform.Type)
		logs.Error("Validation failed: %v", err)
		c.JSONResponse(nil, err)
		return
	}

	logs.Info("Resource %d tested successfully", resource.ID)
	c.JSONResponse(map[string]interface{}{
		"resource_id":   resource.ID,
		"name":          resource.Name,
		"type":          resource.Type,
		"platform_id":   platform.ID,
		"platform_type": platform.Type,
		"result":        result,
	}, nil)
}
