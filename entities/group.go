package entities

import "github.com/google/uuid"

type Group struct {
	GroupID uuid.UUID `gorm:"type:char(36);primary_key;not null"`
	Name  	string 	  `gorm:"unique;not null"`
}