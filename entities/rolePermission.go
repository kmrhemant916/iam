package entities

import "github.com/google/uuid"

type RolePermission struct {
    ID           uuid.UUID `gorm:"type:char(36);primary_key"`
    RoleID       uuid.UUID `gorm:"not null"`
    PermissionID uuid.UUID `gorm:"not null"`
}