package controllers

import (
	"net/http"
	"thriftshop/config"
	"thriftshop/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("alfa12345")

type registerInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Name     string `json:"name"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var input registerInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	user := models.User{Username: input.Username, Email: input.Email, Name: input.Name, Password: string(hash)}
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "registered", "user": user})
}

type loginInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var input loginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if input.Username != "" {
		if err := config.DB.First(&user, "username = ?", input.Username).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
	} else if input.Email != "" {
		if err := config.DB.First(&user, "email = ?", input.Email).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or email required"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "login success", "user": user, "token": tokenString})
}

func GetMe(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no token"})
		return
	}
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
		return
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user_id"})
		return
	}
	userID := uint(userIDFloat)

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
