package entity

import (
	"gorm.io/gorm"
)

type (
	AdminUser struct {
		gorm.Model
		Email        string `gorm:"unique;not null"`
		Password     string `gorm:"not null"`
		RefreshToken string
	}
)
