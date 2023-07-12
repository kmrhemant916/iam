package entities

import (
	"github.com/google/uuid"
)

type User struct {
	UserID         uuid.UUID    `gorm:"type:char(36);primary_key"`
	Email          string       `gorm:"unique;not null"`
	Password       string       `gorm:"not null"`
	IsRoot         bool         `gorm:"not null;default:false"`
	OrganizationID uuid.UUID	`gorm:"not null"`
}