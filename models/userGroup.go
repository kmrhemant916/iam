package models

import "github.com/google/uuid"

type UserGroup struct {
	UserGroupID uuid.UUID  
	UserID 		uuid.UUID 
	GroupID  	uuid.UUID 
}