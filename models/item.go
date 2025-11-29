package models

import "gorm.io/gorm"

type Item struct {
    gorm.Model
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
    Stock       int     `json:"stock"`
    OwnerID     uint    `json:"owner_id"`
}
