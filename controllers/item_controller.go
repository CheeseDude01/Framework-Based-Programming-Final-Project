package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"thriftshop/config"
	"thriftshop/models"
	"time"

	"github.com/gin-gonic/gin"
)

func AddItem(c *gin.Context) {
	ownerIdStr := c.PostForm("ownerId")
	name := c.PostForm("name")
	description := c.PostForm("description")
	priceStr := c.PostForm("price")

	if ownerIdStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "owner ID is required"})
		return
	}

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid price"})
		return
	}

	var user models.User
	if err := config.DB.Where("username = ?", ownerIdStr).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found: " + ownerIdStr})
		return
	}

	imageURL := ""
	file, err := c.FormFile("image")
	if err == nil && file != nil {
		ext := filepath.Ext(file.Filename)
		filename := fmt.Sprintf("%d_%s%s", time.Now().Unix(), ownerIdStr, ext)
		savePath := filepath.Join("static", "uploads", filename)
		
		if err := c.SaveUploadedFile(file, savePath); err == nil {
			imageURL = "/static/uploads/" + filename
		}
	}

	item := models.Item{
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       1,
		ImageURL:    imageURL,
		Status:      "available",
		OwnerID:     user.ID,
	}

	if err := config.DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":          item.ID,
		"name":        item.Name,
		"description": item.Description,
		"price":       item.Price,
		"stock":       item.Stock,
		"imageUrl":    item.ImageURL,
		"ownerId":     item.OwnerID,
	})
}

func GetItems(c *gin.Context) {
	var items []models.Item
	config.DB.Find(&items)

	type ItemResponse struct {
		ID          uint    `json:"id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Stock       int     `json:"stock"`
		OwnerID     uint    `json:"owner_id"`
		OwnerName   string  `json:"ownerId"`
		ImageURL    string  `json:"imageUrl"`
		Status      string  `json:"status"`
	}

	var response []ItemResponse
	for _, item := range items {
		var user models.User
		ownerName := "unknown"
		if err := config.DB.First(&user, item.OwnerID).Error; err == nil {
			ownerName = user.Username
		}
		response = append(response, ItemResponse{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			Price:       item.Price,
			Stock:       item.Stock,
			OwnerID:     item.OwnerID,
			OwnerName:   ownerName,
			ImageURL:    item.ImageURL,
			Status:      item.Status,
		})
	}

	c.JSON(http.StatusOK, response)
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

	contentType := c.GetHeader("Content-Type")
	if strings.Contains(contentType, "multipart/form-data") {
		name := c.PostForm("name")
		description := c.PostForm("description")
		priceStr := c.PostForm("price")

		if name != "" {
			item.Name = name
		}
		if description != "" {
			item.Description = description
		}
		if priceStr != "" {
			price, err := strconv.ParseFloat(priceStr, 64)
			if err == nil {
				item.Price = price
			}
		}

		file, err := c.FormFile("image")
		if err == nil && file != nil {
			ext := filepath.Ext(file.Filename)
			filename := fmt.Sprintf("%d_update%s", time.Now().Unix(), ext)
			savePath := filepath.Join("static", "uploads", filename)
			if err := c.SaveUploadedFile(file, savePath); err == nil {
				item.ImageURL = "/static/uploads/" + filename
			}
		}
	} else {
		var input models.Item
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		item.Name = input.Name
		item.Description = input.Description
		item.Price = input.Price
		item.Stock = input.Stock
	}

	config.DB.Save(&item)
	c.JSON(http.StatusOK, gin.H{
		"id":          item.ID,
		"name":        item.Name,
		"description": item.Description,
		"price":       item.Price,
		"stock":       item.Stock,
		"imageUrl":    item.ImageURL,
		"status":      item.Status,
		"ownerId":     item.OwnerID,
	})
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

	if item.Status == "sold" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "item already sold"})
		return
	}

	if item.Stock < qty {
		c.JSON(http.StatusBadRequest, gin.H{"error": "not enough stock"})
		return
	}

	item.Stock -= qty
	item.Status = "sold"
	if err := config.DB.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	order := models.Order{
		UserID:   0,
		ItemID:   item.ID,
		Quantity: qty,
		Total:    float64(qty) * item.Price,
	}
	config.DB.Create(&order)

	c.JSON(http.StatusOK, gin.H{"message": "purchase success", "order": order})
}
