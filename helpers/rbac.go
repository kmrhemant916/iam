package helpers

import (
	"fmt"

	"github.com/kmrhemant916/iam/authorization"
	"gorm.io/gorm"
)

func InitialiseAuthorization(db *gorm.DB, role []string, permission []string, roleRolePermissions []RolePermission) {
	r := &authorization.Rbac{
		DB: db,
	}
	r.CreateRole(role)
	r.CreatePermission(permission)
	fmt.Println(roleRolePermissions)
	for _, roleRolePermission := range roleRolePermissions {
		r.AssignPermissions(roleRolePermission.Role, roleRolePermission.Permission)
	}
}
