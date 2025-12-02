package models

import "gorm.io/gorm"

type Item struct {
    gorm.Model
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
    Stock       int     `json:"stock"`
    ImageURL    string  `json:"imageUrl"`
    Status      string  `json:"status" gorm:"default:'available'"`
    OwnerID     uint    `json:"owner_id"`
}
