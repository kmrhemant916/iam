package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID uuid.UUID `gorm:"type:char(36);primary_key;" json:"id"`
	Email  string `gorm:"unique" json:"email"`
	Password string `json:"password"`
}