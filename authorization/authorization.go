package authorization

import (
	"errors"

	"github.com/google/uuid"
	"github.com/kmrhemant916/iam/models"
	"gorm.io/gorm"
)

var (
	ErrPermissionInUse     = errors.New("cannot delete assigned permission")
	ErrPermissionNotFound  = errors.New("permission not found")
	ErrRoleAlreadyAssigned = errors.New("this role is already assigned to the user")
	ErrRoleInUse           = errors.New("cannot delete assigned role")
	ErrRoleNotFound        = errors.New("role not found")
	ErrGroupNotFound       = errors.New("group not found")
	ErrUserNotFound        = errors.New("user not found")
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

func (r *Rbac) AssignGroupRoles(group string, roles []string) (error) {
	var dbGroup models.Group
	res := r.DB.Where("name = ?", group).First(&dbGroup)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return ErrGroupNotFound
		}
	}
	var dbRoles []models.Role
	for _, role := range roles {
		var dbRole models.Role
		res := r.DB.Where("name = ?", role).First(&dbRole)
		if res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				return ErrRoleNotFound
			}
		}
		dbRoles = append(dbRoles, dbRole)
	}
	for _, dbRole := range dbRoles {
		var groupRole models.GroupRole
		result := r.DB.Where("role_id = ?", dbRole.ID).Where("group_id = ?", dbGroup.ID).First(&groupRole)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				r.DB.Create(&models.GroupRole{GroupID: dbGroup.ID, RoleID: dbRole.ID})
			} else {
				return result.Error
			}
		}
	}
	return nil
}

func (r *Rbac) CheckGroupRole(groupID uint, roleName string) (bool, error) {
	var role models.Role
	res := r.DB.Where("name = ?", roleName).First(&role)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return false, ErrRoleNotFound
		}
	}
	var groupRole models.GroupRole
	res = r.DB.Where("group_id = ?", groupID).Where("role_id = ?", role.ID).First(&groupRole)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
	}
	return true, nil
}

func (r *Rbac) CheckGroupPermission(groupID uint, permName string) (bool, error) {
	var groupRoles []models.GroupRole
	res := r.DB.Where("user_id = ?", groupID).Find(&groupRoles)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
	}
	var roleIDs []uint
	for _, r := range groupRoles {
		roleIDs = append(roleIDs, r.RoleID)
	}
	var perm models.Permission
	res = r.DB.Where("name = ?", permName).First(&perm)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return false, ErrPermissionNotFound
		}
	}
	var rolePermission models.RolePermission
	res = r.DB.Where("role_id IN (?)", roleIDs).Where("permission_id = ?", perm.ID).First(&rolePermission)
	if res.Error != nil {
		return false, nil
	}
	return true, nil
}

func (r *Rbac) AssignGroups(userID uuid.UUID, groups []string) (error) {
	var user models.User
	res := r.DB.Where("id = ?", userID).First(&user)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
	}
	var groupIDs []uint
	for _, group := range groups {
		var g models.Group
		result := r.DB.Where("name = ?", group).First(&g)
		if result.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				return ErrGroupNotFound
			}
		}
		groupIDs = append(groupIDs, g.ID)
	}
	for _, groupID := range groupIDs {
		var userGroup models.UserGroup
		res := r.DB.Where("user_id = ? AND group_id = ?", user.UserID, groupID).First(&userGroup)
		if res.Error != nil {
			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
				r.DB.Create(&models.UserGroup{UserID: user.UserID, GroupID: groupID})
			} else {
				return res.Error
			}
		}
	}
	return nil
}