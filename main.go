package main

import (
    "thriftshop/config"
    "thriftshop/routes"

    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    config.ConnectDatabase()
    routes.SetupRoutes(r)

    // Serve index.html at root URL
    r.GET("/", func(c *gin.Context) {
        c.File("./templates/index.html")
    })

    // Use port 8081 instead of default 8080 to avoid conflict
    r.Run(":8081")
}
