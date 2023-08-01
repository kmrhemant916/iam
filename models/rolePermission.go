package models

import "github.com/google/uuid"

type RolePermission struct {
    ID           uuid.UUID
    RoleID       uuid.UUID
    PermissionID uuid.UUID
}