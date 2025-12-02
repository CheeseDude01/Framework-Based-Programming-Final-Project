package routes

import (
	"thriftshop/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.Static("/static", "./static")

	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/login")
	})

	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	r.GET("/register", func(c *gin.Context) {
		c.File("./templates/register.html")
	})

	r.GET("/login", func(c *gin.Context) {
		c.File("./templates/login.html")
	})

	r.GET("/explore", func(c *gin.Context) {
		c.File("./templates/index.html")
	})

	r.GET("/items", controllers.GetItems)
	r.GET("/items/:id", controllers.GetItem)
	r.POST("/items", controllers.AddItem)
	r.PUT("/items/:id", controllers.UpdateItem)
	r.DELETE("/items/:id", controllers.DeleteItem)

	r.POST("/items/:id/buy", controllers.BuyItem)

	r.GET("/upload", func(c *gin.Context) {
		c.File("./templates/upload.html")
	})

	r.GET("/profile", func(c *gin.Context) {
		c.File("./templates/profile.html")
	})
}
