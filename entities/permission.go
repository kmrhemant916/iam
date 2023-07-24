package entities

import "github.com/google/uuid"

type Permission struct {
	ID 		uuid.UUID 	`gorm:"type:char(36);primary_key"`
	Name  	string 		`gorm:"unique"`
}