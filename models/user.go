package models

import (
	"github.com/google/uuid"
)

type User struct {
	UserID uuid.UUID
	Email  string
	Password string
	IsRoot bool
	OrganizationID uuid.UUID
}