package models

type RolePermission struct {
    ID           uint `json:"id"`
    RoleID       uint `json:"role_id"`
    PermissionID uint `json:"permission_id"`
}