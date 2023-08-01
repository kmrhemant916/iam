package entities

import "github.com/google/uuid"

type UserGroup struct {
	UserGroupID uuid.UUID 	`gorm:"type:char(36);primary_key;not null"`
	UserID 	 	uuid.UUID 	`gorm:"type:char(36);not null"`
	GroupID  	uuid.UUID 	`gorm:"type:char(36);not null"`
}