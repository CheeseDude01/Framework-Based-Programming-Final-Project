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

    r.Run(":8081")
}
