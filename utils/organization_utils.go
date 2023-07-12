package utils

import (
	"github.com/kmrhemant916/iam/entities"
	"github.com/kmrhemant916/iam/models"
)

func OrganizationToEntity(model *models.Organization) *entities.Organization {
	var entity = &entities.Organization{}
	entity.OrganizationID = model.OrganizationID
	entity.Name = model.Name
	return entity
}

func OrganizationToModel(entity *entities.Organization) *models.Organization {
	var model = &models.Organization{}
	model.OrganizationID = entity.OrganizationID
	model.Name = entity.Name
	return model
}