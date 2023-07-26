package utils

import (
	"github.com/kmrhemant916/iam/entities"
	"github.com/kmrhemant916/iam/models"
)

func GroupRoleToEntity(model *models.GroupRole) *entities.GroupRole {
	var entity = &entities.GroupRole{}
	entity.GroupRoleID = model.GroupRoleID
	entity.GroupID = model.GroupID
	entity.RoleID = model.RoleID
	return entity
}

func GroupRoleToModel(entity *entities.GroupRole) *models.GroupRole {
	var model = &models.GroupRole{}
	model.GroupRoleID = entity.GroupRoleID
	model.GroupID = entity.GroupID
	model.RoleID = entity.RoleID
	return model
}