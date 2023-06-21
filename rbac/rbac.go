package rbac

import (
	"errors"

	"github.com/kmrhemant916/iam/models"
	"gorm.io/gorm"
)

type Rbac struct {
	DB *gorm.DB
}

func (r *Rbac) CreateRole(roles []string) (error) {
	for _, role := range roles {
		var dbRole models.Role
		res := r.DB.Where("name = ?", role).First(&dbRole)
		if res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				r.DB.Create(&models.Role{Name: role})
			} else {
				return res.Error
			}
		}
	}
	return nil
}

func (r *Rbac) CreatePermission(permissions []string) (error) {
	for _, permission := range permissions {
		var dbPermission models.Permission
		res := r.DB.Where("name = ?", permission).First(&dbPermission)
		if res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				r.DB.Create(&models.Permission{Name: permission})
			} else {
				return res.Error
			}
		}
	}
	return nil
}