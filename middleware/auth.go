package middleware

import (
	"app/dal"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/beego/beego/v2/server/web/context"
	"golang.org/x/crypto/bcrypt"
)

func WebAuthMiddleware(ctx *context.Context) {
	if strings.HasPrefix(ctx.Request.URL.Path, "/api") {
		return
	}

	if ctx.Input.Session("user_id") == nil && ctx.Request.URL.Path != "/login" {
		ctx.Redirect(302, "/login")
	}
}

// ApiAuthFilter validates Bearer tokens for API routes
func ApiAuthFilter(ctx *context.Context) {
	authHeader := ctx.Input.Header("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		log.Printf("Missing or invalid Authorization header")
		ctx.Output.SetStatus(401)
		ctx.Output.JSON(map[string]interface{}{"error": "Unauthorized", "code": 401}, false, false)
		return
	}

	q := dal.Q
	token := strings.TrimPrefix(authHeader, "Bearer ")
	apiKey, err := q.ApiKey.Where(q.ApiKey.Token.Eq(token)).First()
	if err != nil || apiKey.ID == 0 {
		log.Printf("Invalid token: %v", err)
		ctx.Output.SetStatus(401)
		ctx.Output.JSON(map[string]interface{}{"error": "Invalid token", "code": 401}, false, false)
		return
	}

	if !apiKey.IsActive {
		log.Printf("Inactive token: %s", token)
		ctx.Output.SetStatus(401)
		ctx.Output.JSON(map[string]interface{}{"error": "Token is inactive", "code": 401}, false, false)
		return
	}

	if apiKey.ExpiresAt != nil && apiKey.ExpiresAt.Before(time.Now()) {
		log.Printf("Expired token: %s", token)
		ctx.Output.SetStatus(401)
		ctx.Output.JSON(map[string]interface{}{"error": "Token has expired", "code": 401}, false, false)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(apiKey.KeyHash), []byte(token))
	if err != nil {
		log.Printf("Token hash mismatch: %v", err)
		ctx.Output.SetStatus(401)
		ctx.Output.JSON(map[string]interface{}{"error": "Invalid token", "code": 401}, false, false)
		return
	}

	// Check scopes
	var metadata map[string]interface{}
	if err := json.Unmarshal([]byte(apiKey.Metadata), &metadata); err != nil {
		log.Printf("Invalid metadata: %v", err)
		ctx.Output.SetStatus(403)
		ctx.Output.JSON(map[string]interface{}{"error": "Invalid metadata", "code": 403}, false, false)
		return
	}

	scopes, ok := metadata["scopes"].([]interface{})
	if !ok {
		log.Printf("No scopes defined in metadata")
		ctx.Output.SetStatus(403)
		ctx.Output.JSON(map[string]interface{}{"error": "No scopes defined", "code": 403}, false, false)
		return
	}

	// Example: Require "write" scope for POST/PUT/DELETE
	if ctx.Input.Method() != "GET" {
		hasWrite := false
		for _, scope := range scopes {
			if scope == "write" {
				hasWrite = true
				break
			}
		}
		if !hasWrite {
			log.Printf("Write scope required for method: %s", ctx.Input.Method())
			ctx.Output.SetStatus(403)
			ctx.Output.JSON(map[string]interface{}{"error": "Write scope required", "code": 403}, false, false)
			return
		}
	}

	ctx.Input.SetData("user_id", apiKey.UserID)
}
