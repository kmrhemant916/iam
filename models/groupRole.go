package models

import "github.com/google/uuid"

type GroupRole struct {
    GroupRoleID   uuid.UUID
    GroupID       uuid.UUID
    RoleID        uuid.UUID
}