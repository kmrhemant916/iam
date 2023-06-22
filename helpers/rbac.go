package helpers

import (
	"github.com/kmrhemant916/iam/authorization"
	"gorm.io/gorm"
)

func InitialiseAuthorization(db *gorm.DB, c *Config) {
	r := &authorization.Rbac{
		DB: db,
	}
	r.CreateRole(c.Roles)
	r.CreatePermission(c.Permissions)
	for _, roleRolePermission := range c.RolePermissions {
		r.AssignPermissions(roleRolePermission.Role, roleRolePermission.Permission)
	}
	r.CreateGroups(c.Groups)
}
