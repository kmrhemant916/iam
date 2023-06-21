package helpers

import (
	"github.com/kmrhemant916/iam/rbac"
	"gorm.io/gorm"
)

func InitialiseRbac(db *gorm.DB, role []string, permission []string) {
	r := &rbac.Rbac{
		DB: db,
	}
	r.CreateRole(role)
	r.CreatePermission(permission)
}
