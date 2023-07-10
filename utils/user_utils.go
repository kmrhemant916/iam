package utils

import (
	"github.com/kmrhemant916/iam/entities"
	"github.com/kmrhemant916/iam/models"
)

func UserToEntity(model *models.User) *entities.User {
	var entity = &entities.User{}
	entity.ID = model.ID
	entity.Email = model.Email
	entity.Password = model.Password
	entity.IsRoot = model.IsRoot
	return entity
}

func UserToModel(entity *entities.User) *models.User {
	var model = &models.User{}
	model.ID = entity.ID
	model.Email = entity.Email
	model.Password = entity.Password
	model.IsRoot = entity.IsRoot
	return model
}