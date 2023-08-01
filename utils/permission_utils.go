package utils

import (
	"github.com/kmrhemant916/iam/entities"
	"github.com/kmrhemant916/iam/models"
)

func PermissionToEntity(model *models.Permission) *entities.Permission {
	var entity = &entities.Permission{}
	entity.ID = model.ID
	entity.Name = model.Name
	return entity
}

func PermissionToModel(entity *entities.Permission) *models.Permission {
	var model = &models.Permission{}
	model.ID = entity.ID
	model.Name = entity.Name
	return model
}