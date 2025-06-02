package controllers

import (
	"app/dal"
	"app/drivers"
	"app/model"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/logs"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type PlatformController struct {
	BaseController
}

func (c *PlatformController) GetPlatforms(limit *int, offset *int, nameFilter *string, nameSort *string) ([]*model.Platform, error) {
	q := dal.Q
	query := q.Platform
	if nameFilter != nil {
		query.Where(q.Platform.Name.Like("%" + *nameFilter + "%"))
	}
	if nameSort != nil {
		switch *nameSort {
		case "name":
			query.Order(q.Platform.Name.Asc())
		case "-name":
			query.Order(q.Platform.Name.Desc())
		}
	}
	if limit != nil {
		query.Limit(*limit)
	}
	if offset != nil {
		query.Offset(*offset)
	}
	platforms, err := query.Find()

	if err != nil {
		return nil, err
	}
	return platforms, nil
}

// GetAll retrieves all platforms with pagination, filtering, and sorting (API)
func (c *PlatformController) GetAll() {
	limit, _ := c.GetInt("limit", 10)
	offset, _ := c.GetInt("offset", 0)
	nameFilter := c.GetString("name")
	nameSort := c.GetString("sort", "name")

	//@deprecated
	// q := dal.Q
	// query := q.Platform
	// if nameFilter != "" {
	// 	query.Where(q.Platform.Name.Like("%" + nameFilter + "%"))
	// }

	// switch nameSort {
	// case "name":
	// 	query.Order(q.Platform.Name.Asc())
	// case "-name":
	// 	query.Order(q.Platform.Name.Desc())
	// }

	// platforms, err := query.Offset(offset).Limit(limit).Find()
	// if err != nil {
	// 	c.PaginatedResponse([]model.Platform{}, 0, limit, offset, err)
	// 	return
	// }

	platforms, err := c.GetPlatforms(&limit, &offset, &nameFilter, &nameSort)
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}

	total, err := dal.Q.Platform.Count()
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}

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

// validateRESTMetadata validates and sanitizes REST platform metadata
func (c *PlatformController) validateRESTMetadata(metadataJSON string) error {
	var metadata model.RESTMetadata
	if err := json.Unmarshal([]byte(metadataJSON), &metadata); err != nil {
		return errors.New("invalid metadata JSON")
	}

	// Validate and sanitize base_endpoint
	if metadata.BaseEndpoint == "" {
		return errors.New("base_endpoint is required for REST platforms")
	}
	parsedURL, err := url.Parse(metadata.BaseEndpoint)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") || parsedURL.Host == "" {
		return errors.New("base_endpoint must be a valid HTTP/HTTPS URL")
	}
	// Sanitize by reconstructing the URL to prevent malicious input
	metadata.BaseEndpoint = parsedURL.String()

	// Validate authentication
	if metadata.Auth.Type == "" {
		metadata.Auth.Type = "none"
	}
	validAuthTypes := map[string]bool{"none": true, "api_key": true, "bearer": true, "basic": true}
	if !validAuthTypes[metadata.Auth.Type] {
		return errors.New("auth.type must be none, api_key, bearer, or basic")
	}

	switch metadata.Auth.Type {
	case "api_key":
		if metadata.Auth.APIKey == nil || *metadata.Auth.APIKey == "" {
			return errors.New("auth.api_key is required for api_key authentication")
		}
		*metadata.Auth.APIKey = strings.TrimSpace(*metadata.Auth.APIKey)
	case "bearer":
		if metadata.Auth.APIKey == nil || *metadata.Auth.BearerToken == "" {
			return errors.New("auth.bearer_token is required for bearer authentication")
		}
		*metadata.Auth.BearerToken = strings.TrimSpace(*metadata.Auth.BearerToken)
	case "basic":
		if metadata.Auth.BasicAuth == nil || metadata.Auth.BasicAuth.Username == "" || metadata.Auth.BasicAuth.Password == "" {
			return errors.New("auth.basic_auth.username and auth.basic_auth.password are required for basic authentication")
		}
		metadata.Auth.BasicAuth.Username = strings.TrimSpace(metadata.Auth.BasicAuth.Username)
		metadata.Auth.BasicAuth.Password = strings.TrimSpace(metadata.Auth.BasicAuth.Password)
	case "none":
		metadata.Auth = model.RESTAuth{Type: "none"} // Clear other auth fields
	}

	// Re-serialize sanitized metadata
	serialized, err := json.Marshal(metadata)
	if err != nil {
		return errors.New("failed to serialize sanitized metadata")
	}
	c.Ctx.Input.SetData("sanitized_metadata", string(serialized))
	return nil
}

// validateInfluxDBMetadata validates and sanitizes InfluxDB platform metadata
func (c *PlatformController) validateInfluxDBMetadata(metadataJSON string) error {
	var metadata model.InfluxDBMetadata
	if err := json.Unmarshal([]byte(metadataJSON), &metadata); err != nil {
		return errors.New("invalid metadata JSON")
	}

	if metadata.URL == "" {
		return errors.New("url is required for InfluxDB platforms")
	}
	parsedURL, err := url.Parse(metadata.URL)
	if err != nil || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") || parsedURL.Host == "" {
		return errors.New("url must be a valid HTTP/HTTPS URL")
	}
	metadata.URL = parsedURL.String()

	if metadata.Token == "" {
		return errors.New("token is required for InfluxDB platforms")
	}
	if metadata.Org == "" {
		return errors.New("org is required for InfluxDB platforms")
	}
	if metadata.Bucket == "" {
		return errors.New("bucket is required for InfluxDB platforms")
	}
	if metadata.Timeout == 0 {
		metadata.Timeout = 10
	}

	metadata.Token = strings.TrimSpace(metadata.Token)
	metadata.Org = strings.TrimSpace(metadata.Org)
	metadata.Bucket = strings.TrimSpace(metadata.Bucket)

	serialized, err := json.Marshal(metadata)
	if err != nil {
		return errors.New("failed to serialize sanitized metadata")
	}
	c.Ctx.Input.SetData("sanitized_metadata", string(serialized))
	return nil
}

// Post creates a new platform with validation (API)
func (c *PlatformController) Post() {
	logs.Info("Received POST request to /api/platforms")

	var platform model.Platform
	if err := c.BindJSON(&platform); err != nil {
		logs.Error("Failed to bind JSON:", err)
		c.JSONResponse(nil, err)
		return
	}

	logs.Debug("Platform input:", platform)

	if platform.Name == "" || platform.Type == "" {
		err := errors.New("name and type are required")
		logs.Error("Validation failed:", err)
		c.JSONResponse(nil, err)
		return
	}

	switch platform.Type {
	case "REST":
		if err := c.validateRESTMetadata(platform.Metadata); err != nil {
			logs.Error("REST metadata validation failed:", err)
			c.JSONResponse(nil, err)
			return
		}
		// Retrieve sanitized metadata from context
		sanitizedMetadata, ok := c.Ctx.Input.GetData("sanitized_metadata").(string)
		if !ok {
			c.JSONResponse(nil, errors.New("failed to retrieve sanitized metadata"))
			return
		}
		platform.Metadata = sanitizedMetadata
	case "InfluxDB":
		if err := c.validateInfluxDBMetadata(platform.Metadata); err != nil {
			logs.Error("InfluxDB metadata validation failed:", err)
			c.JSONResponse(nil, err)
			return
		}
		// Retrieve sanitized metadata from context
		sanitizedMetadata, ok := c.Ctx.Input.GetData("sanitized_metadata").(string)
		if !ok {
			logs.Error("Failed to retrieve sanitized metadata")
			c.JSONResponse(nil, errors.New("failed to retrieve sanitized metadata"))
			return
		}
		platform.Metadata = sanitizedMetadata
	}

	q := dal.Q
	if err := q.Platform.Create(&platform); err != nil {
		logs.Error("Failed to create platform:", err)
		c.JSONResponse(nil, err)
		return
	}

	logs.Info("Platform created successfully:", platform.ID)
	c.JSONResponse(platform, nil)
}

// Put updates an existing platform with validation (API)
func (c *PlatformController) Put() {
	logs.Info("Received PUT request to /api/platforms/%s", c.Ctx.Input.Param(":id"))

	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		logs.Error("Invalid platform ID:", err)
		c.JSONResponse(nil, err)
		return
	}

	var platform model.Platform
	if err := c.BindJSON(&platform); err != nil {
		logs.Error("Failed to bind JSON:", err)
		c.JSONResponse(nil, err)
		return
	}

	logs.Debug("Platform input:", platform)

	if platform.Name == "" || platform.Type == "" {
		err := errors.New("name and type are required")
		logs.Error("Validation failed:", err)
		c.JSONResponse(nil, err)
		return
	}

	// Validate REST platform metadata if Type is REST
	switch platform.Type {
	case "REST":
		if err := c.validateRESTMetadata(platform.Metadata); err != nil {
			logs.Error("REST metadata validation failed:", err)
			c.JSONResponse(nil, err)
			return
		}
		// Retrieve sanitized metadata from context
		sanitizedMetadata, ok := c.Ctx.Input.GetData("sanitized_metadata").(string)
		if !ok {
			c.JSONResponse(nil, errors.New("failed to retrieve sanitized metadata"))
			return
		}
		platform.Metadata = sanitizedMetadata
	case "InfluxDB":
		if err := c.validateInfluxDBMetadata(platform.Metadata); err != nil {
			logs.Error("InfluxDB metadata validation failed:", err)
			c.JSONResponse(nil, err)
			return
		}
		// Retrieve sanitized metadata from context
		sanitizedMetadata, ok := c.Ctx.Input.GetData("sanitized_metadata").(string)
		if !ok {
			logs.Error("Failed to retrieve sanitized metadata")
			c.JSONResponse(nil, errors.New("failed to retrieve sanitized metadata"))
			return
		}
		platform.Metadata = sanitizedMetadata
	}

	platform.ID = uint(id)

	q := dal.Q
	info, err := q.Platform.Where(q.Platform.ID.Eq(uint(id))).Updates(&platform)
	if err != nil {
		logs.Error("Failed to update platform:", err)
		c.JSONResponse(nil, err)
		return
	}
	if info.RowsAffected == 0 {
		logs.Error("Failed to update platform:", err)
		c.JSONResponse(nil, errors.New("no rows affected"))
		return
	}

	logs.Info("Platform updated successfully:", platform.ID)
	c.JSONResponse(platform, info.Error)
}

// Delete removes a platform by ID (API)
func (c *PlatformController) Delete() {
	logs.Info("Received DELETE request to /api/platforms/%s", c.Ctx.Input.Param(":id"))

	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		logs.Error("Invalid platform ID:", err)
		c.JSONResponse(nil, err)
		return
	}

	q := dal.Q
	info, err := q.Platform.Where(q.Platform.ID.Eq(uint(id))).Delete()
	if err != nil {
		logs.Error("Failed to delete platform:", err)
		c.JSONResponse(nil, err)
		return
	}
	if info.RowsAffected == 0 {
		logs.Error("Failed to delete platform:", err)
		c.JSONResponse(nil, errors.New("no rows affected"))
		return
	}

	c.JSONResponse(map[string]string{"message": "Platform deleted successfully"}, info.Error)
}

// FetchDeviceData fetches data for a specific device on a specific platform (API)
func (c *PlatformController) FetchDeviceData() {
	logs.Info("Received request to fetch device data for platform %s, device %s", c.Ctx.Input.Param(":platform_id"), c.Ctx.Input.Param(":device_id"))
	platformID, err := strconv.Atoi(c.Ctx.Input.Param(":platform_id"))
	if err != nil {
		logs.Error("Invalid platform ID:", err)
		c.JSONResponse(nil, err)
		return
	}
	deviceID, err := strconv.Atoi(c.Ctx.Input.Param(":device_id"))
	if err != nil {
		logs.Error("Invalid device ID:", err)
		c.JSONResponse(nil, err)
		return
	}

	q := dal.Q
	// Find the DevicePlatform entry
	dp, err := q.DevicePlatform.Where(
		q.DevicePlatform.PlatformID.Eq(uint(platformID)),
		q.DevicePlatform.DeviceID.Eq(uint(deviceID)),
	).First()
	if err != nil {
		logs.Error("Failed to find DevicePlatform:", err)
		c.JSONResponse(nil, err)
		return
	}

	// Get the platform
	platform, err := dal.Platform.Where(dal.Platform.ID.Eq(uint(platformID))).First()
	if err != nil {
		logs.Error("Failed to find platform:", err)
		c.JSONResponse(nil, err)
		return
	}

	// Find resources for the platform
	resources, err := dal.Resource.Where(dal.Resource.PlatformID.Eq(uint(platformID))).Find()
	if err != nil {
		logs.Error("Failed to find resources:", err)
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

	// Get query parameters for customization
	queryParams := c.Ctx.Request.URL.Query()
	// Fetch data for each resource
	results := make(map[string]interface{})
	for _, resource := range resources {
		if (platform.Type == "REST" && resource.Type == "rest_endpoint") || (platform.Type == "InfluxDB" && resource.Type == "influxdb_query") {
			// Prepare resource details with query parameter overrides
			var modifiedDetails string
			if platform.Type == "InfluxDB" {
				var details model.InfluxDBResourceDetails
				if err := json.Unmarshal([]byte(resource.Details), &details); err != nil {
					logs.Error("Failed to parse InfluxDB resource details:", err)
					results[resource.Name] = map[string]interface{}{"error": err.Error()}
					continue
				}
				// Override time_range or field if provided in query params
				if timeRange := queryParams.Get("time_range"); timeRange != "" {
					if strings.HasPrefix(timeRange, "-") && strings.ContainsAny(timeRange, "smhdwy") {
						details.TimeRange = timeRange
					} else {
						logs.Error("Invalid time_range parameter:", timeRange)
						results[resource.Name] = map[string]interface{}{"error": "invalid time_range, must be negative duration (e.g., '-1h')"}
						continue
					}
				}
				if field := queryParams.Get("field"); field != "" {
					details.Field = strings.TrimSpace(field)
				}
				// Use DeviceAlias as measurement if specified
				if dp.DeviceAlias != "" {
					details.Measurement = dp.DeviceAlias
				}
				modifiedDetailsBytes, err := json.Marshal(details)
				if err != nil {
					logs.Error("Failed to serialize modified InfluxDB details:", err)
					results[resource.Name] = map[string]interface{}{"error": err.Error()}
					continue
				}
				modifiedDetails = string(modifiedDetailsBytes)
			} else if platform.Type == "REST" {
				var details model.RESTResourceDetails
				if err := json.Unmarshal([]byte(resource.Details), &details); err != nil {
					logs.Error("Failed to parse REST resource details:", err)
					results[resource.Name] = map[string]interface{}{"error": err.Error()}
					continue
				}
				// Merge query parameters from request
				if len(queryParams) > 0 {
					if details.QueryParams == nil {
						details.QueryParams = make(map[string]string)
					}
					for key, values := range queryParams {
						if key != "time_range" && key != "field" { // Ignore InfluxDB-specific params
							details.QueryParams[key] = values[0]
						}
					}
				}
				// Append DeviceAlias to query params or path
				if dp.DeviceAlias != "" {
					if strings.Contains(details.Path, ":device_alias") {
						details.Path = strings.ReplaceAll(details.Path, ":device_alias", url.PathEscape(dp.DeviceAlias))
					} else {
						details.QueryParams["device_alias"] = dp.DeviceAlias
					}
				}
				modifiedDetailsBytes, err := json.Marshal(details)
				if err != nil {
					logs.Error("Failed to serialize modified REST details:", err)
					results[resource.Name] = map[string]interface{}{"error": err.Error()}
					continue
				}
				modifiedDetails = string(modifiedDetailsBytes)
			}

			// Fetch data with modified details
			data, err := driver.FetchData(ctx, modifiedDetails)
			if err != nil {
				logs.Error("Failed to fetch data for resource %s: %v", resource.Name, err)
				results[resource.Name] = map[string]interface{}{"error": err.Error()}
				continue
			}
			results[resource.Name] = data
		}
	}

	if len(results) == 0 {
		err := errors.New("no compatible resources found for platform")
		logs.Error(err.Error())
		c.JSONResponse(nil, err)
		return
	}

	logs.Info("Data fetched successfully for platform %d, device %d", platformID, deviceID)
	c.JSONResponse(map[string]interface{}{
		"device_id":   deviceID,
		"platform_id": platformID,
		"alias":       dp.DeviceAlias,
		"data":        results,
	}, nil)
}

// TestConnection tests connectivity to a platform's endpoint (API)
func (c *PlatformController) TestConnection() {
	logs.Info("Received POST request to /api/platforms/test")
	var input struct {
		Type     string `json:"type"`
		Metadata string `json:"metadata"`
	}
	if err := c.BindJSON(&input); err != nil {
		logs.Error("Failed to bind JSON:", err)
		c.JSONResponse(nil, err)
		return
	}

	logs.Debug("Test connection input: type=%s, metadata=%s", input.Type, input.Metadata)

	if input.Type == "REST" {
		if err := c.validateRESTMetadata(input.Metadata); err != nil {
			logs.Error("REST metadata validation failed:", err)
			c.JSONResponse(nil, err)
			return
		}
		var metadata model.RESTMetadata
		if err := json.Unmarshal([]byte(input.Metadata), &metadata); err != nil {
			logs.Error("Invalid metadata JSON:", err)
			c.JSONResponse(nil, errors.New("invalid metadata JSON"))
			return
		}

		client := &http.Client{Timeout: time.Duration(metadata.Timeout) * time.Second}
		req, err := http.NewRequest("HEAD", metadata.BaseEndpoint, nil)
		if err != nil {
			logs.Error("Failed to create HEAD request:", err)
			c.JSONResponse(nil, err)
			return
		}

		switch metadata.Auth.Type {
		case "api_key":
			req.Header.Set("X-API-Key", *metadata.Auth.APIKey)
		case "bearer":
			req.Header.Set("Authorization", "Bearer "+*metadata.Auth.BearerToken)
		case "basic":
			if metadata.Auth.BasicAuth == nil {
				err := errors.New("basic auth configuration missing")
				logs.Error(err.Error())
				c.JSONResponse(nil, err)
				return
			}
			auth := base64.StdEncoding.EncodeToString([]byte(metadata.Auth.BasicAuth.Username + ":" + metadata.Auth.BasicAuth.Password))
			req.Header.Set("Authorization", "Basic "+auth)
		}

		resp, err := client.Do(req)
		if err != nil {
			logs.Error("Connection failed:", err)
			c.JSONResponse(nil, fmt.Errorf("connection failed: %w", err))
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 400 {
			err := fmt.Errorf("connection failed with status: %s", resp.Status)
			logs.Error(err.Error())
			c.JSONResponse(nil, err)
			return
		}

		logs.Info("Connection test successful for base_endpoint:", metadata.BaseEndpoint)
	} else if input.Type == "InfluxDB" {
		if err := c.validateInfluxDBMetadata(input.Metadata); err != nil {
			logs.Error("InfluxDB metadata validation failed:", err)
			c.JSONResponse(nil, err)
			return
		}
		var metadata model.InfluxDBMetadata
		if err := json.Unmarshal([]byte(input.Metadata), &metadata); err != nil {
			logs.Error("Invalid metadata JSON:", err)
			c.JSONResponse(nil, errors.New("invalid metadata JSON"))
			return
		}

		client := influxdb2.NewClientWithOptions(metadata.URL, metadata.Token, influxdb2.DefaultOptions().
			SetHTTPRequestTimeout(uint(metadata.Timeout)))
		defer client.Close()

		ctx := context.Background()
		health, err := client.Health(ctx)
		if err != nil {
			logs.Error("InfluxDB health check failed:", err)
			c.JSONResponse(nil, fmt.Errorf("connection failed: %w", err))
			return
		}
		if health.Status != "pass" {
			err := fmt.Errorf("influxDB health check failed: %s", *health.Message)
			logs.Error(err.Error())
			c.JSONResponse(nil, err)
			return
		}

		logs.Info("Connection test successful for InfluxDB URL:", metadata.URL)
	} else {
		err := fmt.Errorf("unsupported platform type: %s", input.Type)
		logs.Error(err.Error())
		c.JSONResponse(nil, err)
		return
	}

	c.JSONResponse(map[string]string{"message": "Connection successful"}, nil)
}

// ListPlatforms renders the platform management page (Web)
func (c *PlatformController) PlatformIndex() {
	logs.Info("Rendering platform management page")
	userID := c.GetSession("user_id")
	if userID == nil {
		logs.Warn("Unauthenticated access to /platforms, redirecting to /login")
		c.Redirect("/login", 302)
		return
	}
	page, offset := 10, 0

	platforms, err := c.GetPlatforms(&page, &offset, nil, nil)
	if err != nil {
		logs.Error("Failed to get platforms:", err)
		c.JSONResponse(nil, err)
		return
	}

	c.Data["Platforms"] = platforms
	c.TplName = "platforms/index.html"
}

// NewPlatform renders the create platform form (Web)
func (c *PlatformController) NewPlatform() {
	logs.Info("Rendering create platform page")
	userID := c.GetSession("user_id")
	if userID == nil {
		logs.Warn("Unauthenticated access to /platforms/new, redirecting to /login")
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
		logs.Warn("No active API key found for user:", userID)
		c.Data["ApiToken"] = ""
		c.Data["TokenError"] = "No active API key found. Please generate one."
	} else {
		logs.Info("API key found for user:", userID, "Token:", apiKey.Token)
		c.Data["ApiToken"] = apiKey.Token
	}

	c.TplName = "platforms/new.html"
}

// EditPlatform renders the edit platform form (Web)
func (c *PlatformController) EditPlatform() {
	logs.Info("Rendering edit platform page for ID:", c.Ctx.Input.Param(":id"))
	userID := c.GetSession("user_id")
	if userID == nil {
		logs.Warn("Unauthenticated access to /platforms/%s/edit, redirecting to /login", c.Ctx.Input.Param(":id"))
		c.Redirect("/login", 302)
		return
	}

	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		logs.Error("Invalid platform ID:", err)
		c.Data["Error"] = "Invalid platform ID"
		c.TplName = "platforms/edit.html"
		return
	}

	platform, err := dal.Platform.Where(dal.Platform.ID.Eq(uint(id))).First()
	if err != nil {
		logs.Error("Platform not found:", err)
		c.Data["Error"] = "Platform not found"
		c.TplName = "platforms/edit.html"
		return
	}

	// Initialize default values for template data
	c.Data["BaseEndpoint"] = ""
	c.Data["AuthType"] = "none"
	c.Data["APIKey"] = ""
	c.Data["BearerToken"] = ""
	c.Data["BasicUsername"] = ""
	c.Data["BasicPassword"] = ""
	c.Data["URL"] = ""
	c.Data["Token"] = ""
	c.Data["Org"] = ""
	c.Data["Bucket"] = ""
	c.Data["Timeout"] = 10

	if platform.Type == "REST" {
		var metadata model.RESTMetadata
		if err := json.Unmarshal([]byte(platform.Metadata), &metadata); err != nil {
			logs.Error("Invalid metadata format:", err)
			c.Data["Error"] = "Invalid metadata format"
			c.TplName = "platforms_edit.html"
			c.Layout = "layout.html"
			return
		}
		c.Data["BaseEndpoint"] = metadata.BaseEndpoint
		c.Data["AuthType"] = metadata.Auth.Type
		c.Data["APIKey"] = metadata.Auth.APIKey
		c.Data["BearerToken"] = metadata.Auth.BearerToken
		if metadata.Auth.BasicAuth != nil {
			c.Data["BasicUsername"] = metadata.Auth.BasicAuth.Username
			c.Data["BasicPassword"] = metadata.Auth.BasicAuth.Password
		} else {
			c.Data["BasicUsername"] = ""
			c.Data["BasicPassword"] = ""
		}
	} else if platform.Type == "InfluxDB" {
		var metadata model.InfluxDBMetadata
		if err := json.Unmarshal([]byte(platform.Metadata), &metadata); err != nil {
			logs.Error("Invalid metadata format:", err)
			c.Data["Error"] = "Invalid metadata format"
			c.TplName = "platforms_edit.html"
			c.Layout = "layout.html"
			return
		}
		c.Data["URL"] = metadata.URL
		c.Data["Token"] = metadata.Token
		c.Data["Org"] = metadata.Org
		c.Data["Bucket"] = metadata.Bucket
	}

	q := dal.Q
	// Fetch user's active API key
	apiKey, err := q.ApiKey.Where(
		q.ApiKey.UserID.Eq(userID.(uint)),
		q.ApiKey.IsActive.Is(true),
	).First()
	if err != nil || apiKey.ID == 0 {
		logs.Warn("No active API key found for user:", userID)
		c.Data["ApiToken"] = ""
		c.Data["TokenError"] = "No active API key found. Please generate one."
	} else {
		logs.Info("API key found for user:", userID, "Token:", apiKey.Token)
		c.Data["ApiToken"] = apiKey.Token
	}

	c.Data["Platform"] = platform
	c.TplName = "platforms/edit.html"
}

// ViewPlatform renders the view platform page (Web)
func (c *PlatformController) ViewPlatform() {
	logs.Info("Rendering view platform page for ID:", c.Ctx.Input.Param(":id"))
	userID := c.GetSession("user_id")
	if userID == nil {
		logs.Warn("Unauthenticated access to /platforms/%s/view, redirecting to /login", c.Ctx.Input.Param(":id"))
		c.Redirect("/login", 302)
		return
	}

	id, err := strconv.Atoi(c.Ctx.Input.Param(":id"))
	if err != nil {
		logs.Error("Invalid platform ID:", err)
		c.Data["Error"] = "Invalid platform ID"
		c.TplName = "platforms/view.html"
		return
	}

	platform, err := dal.Platform.Where(dal.Platform.ID.Eq(uint(id))).First()
	if err != nil {
		logs.Error("Platform not found:", err)
		c.Data["Error"] = "Platform not found"
		c.TplName = "platforms/view.html"
		return
	}

	q := dal.Q
	apiKey, err := dal.ApiKey.Where(
		q.ApiKey.UserID.Eq(userID.(uint)),
		q.ApiKey.IsActive.Is(true),
	).First()
	if err != nil || apiKey.ID == 0 {
		logs.Warn("No active API key found for user:", userID)
		c.Data["ApiToken"] = ""
		c.Data["TokenError"] = "No active API key found. Please generate one."
	} else {
		logs.Info("API key found for user:", userID, "Token:", apiKey.Token)
		c.Data["ApiToken"] = apiKey.Token
	}

	c.Data["Platform"] = platform
	c.TplName = "platforms/view.html"
}
