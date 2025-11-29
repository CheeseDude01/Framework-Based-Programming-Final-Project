package migrations

import (
    "thriftshop/config"
    "thriftshop/models"
)

func RunMigrations() {
    // This file is optional; ConnectDatabase already AutoMigrates.
    // You can add seed data here if you want.
    db := config.DB
    db.AutoMigrate(&models.User{}, &models.Item{}, &models.Order{})
}
