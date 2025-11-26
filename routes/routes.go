package routes

import (
    "thriftshop/controllers"

    "github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
    r.POST("/register", controllers.Register)
    r.POST("/login", controllers.Login)

    r.GET("/items", controllers.GetItems)
    r.GET("/items/:id", controllers.GetItem)
    r.POST("/items", controllers.AddItem)
    r.PUT("/items/:id", controllers.UpdateItem)
    r.DELETE("/items/:id", controllers.DeleteItem)

    r.POST("/items/:id/buy", controllers.BuyItem)
}
