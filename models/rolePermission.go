package models

type RolePermission struct {
    RoleID       uint `json:"role_id"`
    PermissionID uint `json:"permission_id"`
    Role         Role  `gorm:"foreignkey:RoleID"`
    Permission   Permission `gorm:"foreignkey:PermissionID"`
}