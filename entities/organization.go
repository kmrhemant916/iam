package entities

import (
	"github.com/google/uuid"
)

type Organization struct {
	OrganizationID uuid.UUID `gorm:"type:char(36);primary_key"`
	Name           string    `gorm:"unique;not null"`
	Users          []User    `gorm:"constraint:OnDelete:CASCADE;"`
}