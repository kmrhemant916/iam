package utils

import (
	"github.com/kmrhemant916/iam/entities"
	"github.com/kmrhemant916/iam/models"
)

func RolePermissionToEntity(model *models.RolePermission) *entities.RolePermission {
	var entity = &entities.RolePermission{}
	entity.ID = model.ID
	entity.RoleID = model.RoleID
    entity.PermissionID = model.PermissionID
	return entity
}

func RolePermissionToModel(entity *entities.RolePermission) *models.RolePermission {
	var model = &models.RolePermission{}
	model.ID = entity.ID
	model.RoleID = entity.RoleID
    model.PermissionID = entity.PermissionID
	return model
}