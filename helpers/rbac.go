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
	// for _, roleRolePermission := range c.RolePermissions {
	// 	r.AssignPermissions(roleRolePermission.Role, roleRolePermission.Permissions)
	// }
	r.CreateGroups(c.Groups)
	// for _, groupRole := range c.GroupRoles {
	// 	r.AssignGroupRoles(groupRole.Group, groupRole.Roles)
	// }
}
