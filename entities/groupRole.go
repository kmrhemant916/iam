package entities

import "github.com/google/uuid"

type GroupRole struct {
    GroupRoleID   uuid.UUID `gorm:"type:char(36);primary_key"`
    GroupID       uuid.UUID `gorm:"not null"`
    RoleID        uint      `gorm:"not null"`
}