package controllers

import (
	"net/http"
	"strconv"

	"github.com/cat-spmog/iothubgo/models"
	model "github.com/cat-spmog/iothubgo/models"
	"github.com/cat-spmog/iothubgo/query"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{DB: db}
}

func (uc *UserController) GetUsers(c *gin.Context) {
	q := query.User
	users, err := q.WithContext(c).Find()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (uc *UserController) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	q := query.User
	user, err := q.WithContext(c).GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (uc *UserController) CreateUser(c *gin.Context) {
	q := query.User
	password := c.GetString("password")
	if password == "" || &password == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password is required"})
		return
	}

	user := model.User{
		Name:   c.GetString("name"),
		Email:  c.GetString("email"),
		Active: true,
		Roles:  nil,
	}

	user.SetNewPassword(password)

	err := q.WithContext(c).Create(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"msg": "ok"})
}

func (uc *UserController) UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	q := query.User
	// First, retrieve the existing user
	user, err := q.WithContext(c).GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	// Bind only the fields that were sent in the request
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Handle password separately if provided
	password, passwordExists := updates["password"].(string)
	if passwordExists && password != "" {
		user.SetNewPassword(password)
		delete(updates, "password") // Remove password from updates map
	}

	// Only update fields that were sent in the request
	if len(updates) > 0 {
		if err := q.WithContext(c).
			Where(q.ID.Eq(int32(user.ID))).
			Where(q.Active.Is(true)).
			Updates(updates); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Get the updated user to return
	updatedUser, err := q.WithContext(c).GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated user"})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

func (uc *UserController) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	if err := uc.DB.Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}
