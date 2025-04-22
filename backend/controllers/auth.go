package controllers

import (
	"log"
	"net/http"

	model "github.com/cat-spmog/iothubgo/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (ac *AuthController) PostLogin(c *gin.Context) {
	log.Printf("PostLogin: Received login request")
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("PostLogin: Failed to parse request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user, err := model.AuthenticateUser(ac.DB, req.Email, req.Password)
	if err != nil || user.ID == 0 {
		log.Printf("PostLogin: Login failed for %s: %v", req.Email, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	session := sessions.Default(c)
	session.Set("email", req.Email)
	session.Set("roles", user.Roles)
	session.Set("user_id", user.ID)

	if err := session.Save(); err != nil {
		log.Printf("PostLogin: Failed to save session: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	log.Printf("User Session set for %s: user_id=%d, roles=%v", req.Email, user.ID, user.Roles)
	c.JSON(http.StatusOK, gin.H{
		"email":   req.Email,
		"roles":   user.Roles,
		"user_id": user.ID,
	})
}

func (ac *AuthController) GetMe(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		log.Printf("GetMe: No session found")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: No session"})
		return
	}
	username := session.Get("username")
	roles := session.Get("roles")
	if username == nil || roles == nil {
		log.Printf("GetMe: Incomplete session data: user_id=%v", userID)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: Incomplete session"})
		return
	}

	log.Printf("GetMe: Session found: user_id=%v, username=%v, roles=%v", userID, username, roles)
	c.JSON(http.StatusOK, gin.H{
		"username": username,
		"roles":    roles, // Changed from "role" to "roles" to match the actual data model
		"user_id":  userID,
	})
}

func (ac *AuthController) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	log.Printf("Logout: Session destroyed")
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
