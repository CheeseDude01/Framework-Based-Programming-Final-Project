package config

import (
    "log"
    "thriftshop/models"

    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
    db, err := gorm.Open(sqlite.Open("thriftshop.db"), &gorm.Config{})
    if err != nil {
        log.Fatal("failed to connect database:", err)
    }

    // Run migrations
    err = db.AutoMigrate(&models.User{}, &models.Item{}, &models.Order{})
    if err != nil {
        log.Fatal("migration failed:", err)
    }

    DB = db
}
