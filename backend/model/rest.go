package model

// RESTMetadata defines the structure for REST platform metadata
type RESTMetadata struct {
	BaseEndpoint string   `json:"base_endpoint"`
	Auth         RESTAuth `json:"auth"`
	Timeout      int      `json:"timeout,omitempty"`
}

// RESTAuth defines the authentication details for REST platforms
type RESTAuth struct {
	Type        string         `json:"type"` // none, api_key, bearer, basic
	APIKey      *string        `json:"api_key,omitempty"`
	BearerToken *string        `json:"bearer_token,omitempty"`
	BasicAuth   *RESTBasicAuth `json:"basic_auth,omitempty"`
}

// RESTBasicAuth defines the username and password for Basic Auth
type RESTBasicAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RESTResourceDetails defines the structure for REST endpoint details
type RESTResourceDetails struct {
	Method      string            `json:"method"` // GET, POST, PUT, DELETE
	Path        string            `json:"path"`   // e.g., "/assets"
	Headers     map[string]string `json:"headers,omitempty"`
	QueryParams map[string]string `json:"query_params,omitempty"`
	Body        string            `json:"body,omitempty"` // JSON string for request body template
}
