package utils

import (
	"github.com/kmrhemant916/iam/entities"
	"github.com/kmrhemant916/iam/models"
)

func UserGroupToEntity(model *models.UserGroup) *entities.UserGroup {
	var entity = &entities.UserGroup{}
	entity.UserGroupID = model.UserGroupID
	entity.UserID = model.UserID
	entity.GroupID = model.GroupID
	return entity
}

func UserGroupToModel(entity *entities.UserGroup) *models.UserGroup {
	var model = &models.UserGroup{}
	model.UserGroupID = entity.UserGroupID
	model.UserID = entity.UserID
	model.GroupID = entity.GroupID
	return model
}