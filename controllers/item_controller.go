package controllers

import (
	"net/http"
	"strconv"
	"thriftshop/config"
	"thriftshop/models"

	"github.com/gin-gonic/gin"
)

// Create item
func AddItem(c *gin.Context) {
	var item models.Item
	if err := c.ShouldBind(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// For multipart, ownerId is sent as string, need to resolve to user ID
	ownerIdStr := c.PostForm("ownerId")
	if ownerIdStr != "" {
		var user models.User
		if err := config.DB.Where("username = ?", ownerIdStr).First(&user).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid owner"})
			return
		}
		item.OwnerID = user.ID
	}

	if err := config.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

func GetItems(c *gin.Context) {
	var items []models.Item
	config.DB.Find(&items)
	c.JSON(http.StatusOK, items)
}

func GetItem(c *gin.Context) {
	id := c.Param("id")
	var item models.Item
	if err := config.DB.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

func UpdateItem(c *gin.Context) {
	id := c.Param("id")
	var item models.Item
	if err := config.DB.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
		return
	}

	var input models.Item
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item.Name = input.Name
	item.Description = input.Description
	item.Price = input.Price
	item.Stock = input.Stock

	config.DB.Save(&item)
	c.JSON(http.StatusOK, item)
}

func DeleteItem(c *gin.Context) {
	id := c.Param("id")
	var item models.Item
	if err := config.DB.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
		return
	}
	config.DB.Delete(&item)
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// Buy item (creates an order and reduces stock)
func BuyItem(c *gin.Context) {
	id := c.Param("id")
	qtyStr := c.DefaultQuery("qty", "1")
	qty, err := strconv.Atoi(qtyStr)
	if err != nil || qty <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid qty"})
		return
	}

	var item models.Item
	if err := config.DB.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
		return
	}

	if item.Stock < qty {
		c.JSON(http.StatusBadRequest, gin.H{"error": "not enough stock"})
		return
	}

	item.Stock -= qty
	if err := config.DB.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	order := models.Order{
		UserID:   0, // anonymous/buyer id; replace with authenticated user id later
		ItemID:   item.ID,
		Quantity: qty,
		Total:    float64(qty) * item.Price,
	}
	config.DB.Create(&order)

	c.JSON(http.StatusOK, gin.H{"message": "purchase success", "order": order})
}
