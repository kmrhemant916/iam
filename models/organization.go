package models

import "github.com/google/uuid"

type Organization struct {
	OrganizationID uuid.UUID
	Name  string
}