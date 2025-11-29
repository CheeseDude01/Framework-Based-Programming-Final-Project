package routes

import (
	"thriftshop/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	// Serve register.html
	r.GET("/register", func(c *gin.Context) {
		c.File("./templates/register.html")
	})

	// Serve login.html
	r.GET("/login", func(c *gin.Context) {
		c.File("./templates/login.html")
	})

	// Removed root GET as it is now handled in main.go to avoid conflict

	r.GET("/items", controllers.GetItems)
	r.GET("/items/:id", controllers.GetItem)
	r.POST("/items", controllers.AddItem)
	r.PUT("/items/:id", controllers.UpdateItem)
	r.DELETE("/items/:id", controllers.DeleteItem)

	r.POST("/items/:id/buy", controllers.BuyItem)

	// Serve upload.html
	r.GET("/upload", func(c *gin.Context) {
		c.File("./templates/upload.html")
	})

	// Serve profile.html
	r.GET("/profile", func(c *gin.Context) {
		c.File("./templates/profile.html")
	})
}
