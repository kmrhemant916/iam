package authorization

import (
	"errors"

	"github.com/kmrhemant916/iam/models"
	"gorm.io/gorm"
)

var (
	ErrPermissionInUse     = errors.New("Cannot delete assigned permission")
	ErrPermissionNotFound  = errors.New("Permission not found")
	ErrRoleAlreadyAssigned = errors.New("This role is already assigned to the user")
	ErrRoleInUse           = errors.New("Cannot delete assigned role")
	ErrRoleNotFound        = errors.New("Role not found")
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

func (r *Rbac) AssignPermissions(roleName string, permNames []string) (error) {
	var role models.Role
	rRes := r.DB.Where("name = ?", roleName).First(&role)
	if rRes.Error != nil {
		if errors.Is(rRes.Error, gorm.ErrRecordNotFound) {
			return ErrRoleNotFound
		}
	}
	var perms []models.Permission
	for _, permName := range permNames {
		var perm models.Permission
		rRes := r.DB.Where("name = ?", permName).First(&perm)
		if rRes.Error != nil {
			if errors.Is(rRes.Error, gorm.ErrRecordNotFound) {
				return ErrPermissionNotFound
			}
		}
		perms = append(perms, perm)
	}
	for _, perm := range perms {
		var rolePerm models.RolePermission
		rRes := r.DB.Where("role_id = ?", role.ID).Where("permission_id = ?", perm.ID).First(&rolePerm)
		if rRes.Error != nil {
			res := r.DB.Create(&models.RolePermission{RoleID: role.ID, PermissionID: perm.ID})
			if res.Error != nil {
				return res.Error
			}
		}
	}
	return nil
}

func (r *Rbac) CreateGroups(groups []string) (error) {
	for _, group := range groups {
		var dbGroup models.Group
		res := r.DB.Where("name = ?", group).First(&dbGroup)
		if res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				r.DB.Create(&models.Group{Name: group})
			} else {
				return res.Error
			}
		}
	}
	return nil
}