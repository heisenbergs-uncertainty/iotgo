package middleware

import (
	"app/dal"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/beego/beego/v2/server/web/context"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// ApiAuthFilter validates Bearer tokens for API routes or JWT tokens
func ApiAuthFilter(ctx *context.Context) {
	authHeader := ctx.Input.Header("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		log.Printf("Missing or invalid Authorization header")
		ctx.Output.SetStatus(401)
		ctx.Output.JSON(map[string]interface{}{"error": "Unauthorized", "code": 401}, false, false)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	// Try API key first
	if authenticated := validateApiKey(ctx, token); authenticated {
		return // Authentication succeeded via API key
	}

	// If API key validation failed, try JWT
	if authenticated := validateJWT(ctx, token); authenticated {
		return // Authentication succeeded via JWT
	}

	// If both authentication methods failed
	ctx.Output.SetStatus(401)
	ctx.Output.JSON(map[string]interface{}{"error": "Invalid authentication token", "code": 401}, false, false)
}

// validateApiKey checks if the token is a valid API key
func validateApiKey(ctx *context.Context, token string) bool {
	q := dal.Q
	apiKey, err := q.ApiKey.Where(q.ApiKey.Token.Eq(token)).First()
	if err != nil || apiKey.ID == 0 {
		return false
	}

	if !apiKey.IsActive {
		return false
	}

	if apiKey.ExpiresAt != nil && apiKey.ExpiresAt.Before(time.Now()) {
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(apiKey.KeyHash), []byte(token))
	if err != nil {
		return false
	}

	// Check scopes
	var metadata map[string]interface{}
	if err := json.Unmarshal([]byte(apiKey.Metadata), &metadata); err != nil {
		ctx.Output.SetStatus(403)
		ctx.Output.JSON(map[string]interface{}{"error": "Invalid metadata", "code": 403}, false, false)
		return false
	}

	scopes, ok := metadata["scopes"].([]interface{})
	if !ok {
		ctx.Output.SetStatus(403)
		ctx.Output.JSON(map[string]interface{}{"error": "No scopes defined", "code": 403}, false, false)
		return false
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
			ctx.Output.SetStatus(403)
			ctx.Output.JSON(map[string]interface{}{"error": "Write scope required", "code": 403}, false, false)
			return false
		}
	}

	ctx.Input.SetData("user_id", apiKey.UserID)
	return true
}

// validateJWT checks if the token is a valid JWT
func validateJWT(ctx *context.Context, tokenString string) bool {
	secretKey := []byte("your-secret-key") // Replace with env variable

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		return false
	}

	// Store user ID in context
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		ctx.Input.SetData("user_id", int(claims["sub"].(float64)))
	}

	return true
}
