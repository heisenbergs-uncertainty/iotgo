package controllers

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	model "github.com/cat-spmog/iothubgo/models"
	"github.com/cat-spmog/iothubgo/query"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DeviceController struct {
	DB *gorm.DB
}

func NewDeviceController(db *gorm.DB) *DeviceController {
	return &DeviceController{DB: db}
}

func (dc *DeviceController) CreateDevice(c *gin.Context) {
	q := query.Device
	session := sessions.Default(c)
	if session.Get("user_id") == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: No user session"})
		return
	}

	var device model.Device
	if err := c.ShouldBindJSON(&device); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := q.WithContext(c).Create(&device); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusCreated, device)
}

func (dc *DeviceController) GetDevice(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("user_id") == nil {
		log.Printf("GetDevice: Unauthorized access")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: No user session"})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	q := query.Device

	device, err := q.WithContext(c).GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}
	c.JSON(http.StatusOK, device)
}

func (dc *DeviceController) GetAllDevices(c *gin.Context) {
	var devices []model.Device
	query := dc.DB.Model(&model.Device{})

	// Handle query parameters
	if fields := c.Query("fields"); fields != "" {
		query = query.Select(strings.Split(fields, ","))
	}
	if limit, err := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64); err == nil {
		query = query.Limit(int(limit))
	}
	if offset, err := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64); err == nil {
		query = query.Offset(int(offset))
	}
	if sortby := c.Query("sortby"); sortby != "" {
		orders := strings.Split(c.Query("order"), ",")
		sortFields := strings.Split(sortby, ",")
		for i, field := range sortFields {
			order := "asc"
			if i < len(orders) && orders[i] == "desc" {
				order = "desc"
			}
			query = query.Order(field + " " + order)
		}
	}
	if q := c.Query("query"); q != "" {
		for _, cond := range strings.Split(q, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query key/value pair"})
				return
			}
			query = query.Where(kv[0]+" = ?", kv[1])
		}
	}

	if err := query.Find(&devices).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, devices)
}

func (dc *DeviceController) UpdateDevice(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var device model.Device
	if err := dc.DB.First(&device, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}
	if err := c.ShouldBindJSON(&device); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := dc.DB.Save(&device).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, device)
}

func (dc *DeviceController) DeleteDevice(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := dc.DB.Delete(&model.Device{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Device deleted"})
}
