package migrations

import (
    "thriftshop/config"
    "thriftshop/models"
)

func RunMigrations() {
    db := config.DB
    db.AutoMigrate(&models.User{}, &models.Item{}, &models.Order{})
}
