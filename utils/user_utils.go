package utils

import (
	"github.com/kmrhemant916/iam/entities"
	"github.com/kmrhemant916/iam/models"
)

func UserToEntity(model *models.User) *entities.User {
	var entity = &entities.User{}
	entity.UserID = model.UserID
	entity.Email = model.Email
	entity.Password = model.Password
	entity.IsRoot = model.IsRoot
	entity.OrganizationID = model.OrganizationID
	return entity
}

func UserToModel(entity *entities.User) *models.User {
	var model = &models.User{}
	model.UserID = entity.UserID
	model.Email = entity.Email
	model.Password = entity.Password
	model.IsRoot = entity.IsRoot
	model.OrganizationID = entity.OrganizationID
	return model
}