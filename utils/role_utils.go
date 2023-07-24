package utils

import (
	"github.com/kmrhemant916/iam/entities"
	"github.com/kmrhemant916/iam/models"
)

func RoleToEntity(model *models.Role) *entities.Role {
	var entity = &entities.Role{}
	entity.ID = model.ID
	entity.Name = model.Name
	return entity
}

func RoleToModel(entity *entities.Role) *models.Role {
	var model = &models.Role{}
	model.ID = entity.ID
	model.Name = entity.Name
	return model
}