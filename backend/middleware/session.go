package middleware

import (
	"net/http"

	"github.com/cat-spmog/iothubgo/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"gorm.io/gorm"
)

func SessionMiddleware(store *sessions.CookieStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		session, _ := store.Get(c.Request, "session")
		c.Set("session", session)
		c.Next()
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, exists := c.Get("session")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: No session"})
			c.Abort()
			return
		}
		sess := session.(*sessions.Session)
		userID := sess.Values["user_id"]
		if userID == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: No user session"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// RequireRole ensures the user has a role with the required access level for a specific resource type
func RequireRole(db *gorm.DB, requiredLevel string, resourceType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// First ensure authentication
		session, exists := c.Get("session")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: No session"})
			c.Abort()
			return
		}

		sess := session.(*sessions.Session)
		userID := sess.Values["user_id"]
		if userID == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: No user session"})
			c.Abort()
			return
		}

		// Load user with roles from database
		var user models.User
		if err := db.Preload("Roles").First(&user, userID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Invalid user"})
			c.Abort()
			return
		}

		// Check if user has required permission
		if !user.HasPermission(requiredLevel, resourceType) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Insufficient permissions"})
			c.Abort()
			return
		}

		// Store user in context for later use
		c.Set("user", user)
		c.Next()
	}
}

// GetCurrentUser retrieves the authenticated user from the gin context
func GetCurrentUser(c *gin.Context) (models.User, bool) {
	user, exists := c.Get("user")
	if !exists {
		return models.User{}, false
	}
	typedUser, ok := user.(models.User)
	return typedUser, ok
}
