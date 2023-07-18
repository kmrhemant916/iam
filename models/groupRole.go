package models

import "github.com/google/uuid"

type GroupRole struct {
    GroupRoleID           uuid.UUID `json:"id"`
    GroupID       uuid.UUID `json:"group_id"`
    RoleID uint `json:"role_id"`
}