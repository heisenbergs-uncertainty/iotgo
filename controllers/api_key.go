package controllers

import (
	"app/dal"
	"app/model"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type ApiKeyController struct {
	BaseController
}

type GenerateApiKeyRequest struct {
	Name   string   `json:"name"`
	Scopes []string `json:"scopes"` // e.g., ["read", "write"]
}

func (c *ApiKeyController) Generate() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.JSONResponse(nil, errors.New("unauthenticated"))
		return
	}

	var req GenerateApiKeyRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSONResponse(nil, err)
		return
	}

	validScopes := map[string]bool{"read": true, "write": true}
	for _, scope := range req.Scopes {
		if !validScopes[scope] {
			c.JSONResponse(nil, fmt.Errorf("invalid Scope: %s", scope))
			return
		}
	}

	token := uuid.New().String()
	hash, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}

	metadata, _ := json.Marshal(map[string]interface{}{"scopes": req.Scopes})
	apiKey := model.ApiKey{
		Name:      req.Name,
		KeyID:     uuid.New().String(),
		Token:     token,
		KeyHash:   string(hash),
		IsActive:  true,
		UserID:    userID.(uint),
		Metadata:  string(metadata),
		ExpiresAt: nil,
	}

	q := dal.Q
	if err := q.ApiKey.Create(&apiKey); err != nil {
		c.JSONResponse(nil, err)
		return
	}

	c.JSONResponse(map[string]string{"token": token}, nil)
}

// Revoke deactivates an API key (API)
func (c *ApiKeyController) Revoke() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.JSONResponse(nil, errors.New("unauthenticated"))
		return
	}

	keyID := c.Ctx.Input.Param(":key_id")

	q := dal.Q
	info, err := q.ApiKey.Where(
		q.ApiKey.KeyID.Eq(keyID),
		q.ApiKey.UserID.Eq(userID.(uint)),
	).Update(q.ApiKey.IsActive, false)
	if err != nil {
		c.JSONResponse(nil, err)
		return
	}
	if info.RowsAffected == 0 {
		c.JSONResponse(nil, errors.New("no rows affected"))
		return
	}

	c.JSONResponse(map[string]string{"message": "API key revoked"}, info.Error)
}

// ManageKeys renders the API key management page (Web)
func (c *ApiKeyController) ManageKeys() {
	userID := c.GetSession("user_id")
	if userID == nil {
		c.Redirect("/login", 302)
		return
	}

	q := dal.Q
	keys, err := q.ApiKey.Where(q.ApiKey.UserID.Eq(userID.(uint))).Find()
	if err != nil {
		c.Data["Error"] = "Failed to load API keys"
	}

	// Fetch user's active API key for API calls
	apiKey, err := q.ApiKey.Where(q.ApiKey.UserID.Eq(userID.(uint)), q.ApiKey.IsActive.Is(true)).First()
	if err != nil || apiKey.ID == 0 {
		c.Data["ApiToken"] = ""
		c.Data["TokenError"] = "No Active API key found. Please generate one."
	} else {
		c.Data["ApiToken"] = apiKey.Token
	}

	c.Data["ApiKeys"] = keys
	c.TplName = "api_keys.html"
}
