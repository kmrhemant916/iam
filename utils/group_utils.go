package utils

import (
	"github.com/kmrhemant916/iam/entities"
	"github.com/kmrhemant916/iam/models"
)

func GroupToEntity(model *models.Group) *entities.Group {
	var entity = &entities.Group{}
	entity.GroupID = model.GroupID
	entity.Name = model.Name
	return entity
}

func GroupToModel(entity *entities.Group) *models.Group {
	var model = &models.Group{}
	model.GroupID = entity.GroupID
	model.Name = entity.Name
	return model
}