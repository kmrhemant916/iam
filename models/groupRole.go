package models

type GroupRole struct {
    ID           uint `json:"id"`
    GroupID       uint `json:"group_id"`
    RoleID uint `json:"role_id"`
}