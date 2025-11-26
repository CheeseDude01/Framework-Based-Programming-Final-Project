package models

import "gorm.io/gorm"

type Order struct {
    gorm.Model
    UserID uint
    ItemID uint
    Quantity int
    Total float64
}
