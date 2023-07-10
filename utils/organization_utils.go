package utils

import (
	"github.com/kmrhemant916/iam/entities"
	"github.com/kmrhemant916/iam/models"
)

func OrganizationToEntity(model *models.Organization) *entities.Organization {
	var entity = &entities.Organization{}
	entity.ID = model.ID
	entity.Name = model.Name
	return entity
}

func OrganizationToModel(entity *entities.Organization) *models.Organization {
	var model = &models.Organization{}
	model.ID = entity.ID
	model.Name = entity.Name
	return model
}