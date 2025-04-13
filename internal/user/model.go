package user

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Phone    string `gorm:"index"`
	Password string `json:"password" validate:"required"`
	Name     string `json:"name" validate:"required"`
	Token    string `json:"token"`
}
