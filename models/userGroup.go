package models

import "github.com/google/uuid"

type UserGroup struct {
	// gorm.Model
	ID 	  uint `gorm:"primary_key;autoIncrement;" json:"id"`
	UserID uuid.UUID `gorm:"type:char(36);" json:"user_id"`
	GroupID  uint `gorm:"unique" json:"group_id"`
}