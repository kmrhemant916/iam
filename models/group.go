package models

import "github.com/google/uuid"

type Group struct {
	GroupID uuid.UUID
	Name    string
}