package authorization

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/kmrhemant916/iam/entities"
	"github.com/kmrhemant916/iam/global"
	"github.com/kmrhemant916/iam/models"
	"github.com/kmrhemant916/iam/repositories"
	"github.com/kmrhemant916/iam/service"
	"github.com/kmrhemant916/iam/utils"
	"gorm.io/gorm"
)

type Rbac struct {
	DB *gorm.DB
}

// func (r *Rbac) CreateRole(roles []string) (error) {
// 	for _, role := range roles {
// 		var dbRole models.Role
// 		res := r.DB.Where("name = ?", role).First(&dbRole)
// 		if res.Error != nil {
// 			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
// 				r.DB.Create(&models.Role{Name: role})
// 			} else {
// 				return res.Error
// 			}
// 		}
// 	}
// 	return nil
// }

// func (r *Rbac) CreatePermission(permissions []string) (error) {
// 	for _, permission := range permissions {
// 		var dbPermission models.Permission
// 		res := r.DB.Where("name = ?", permission).First(&dbPermission)
// 		if res.Error != nil {
// 			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
// 				r.DB.Create(&models.Permission{Name: permission})
// 			} else {
// 				return res.Error
// 			}
// 		}
// 	}
// 	return nil
// }

// func (r *Rbac) AssignPermissions(roleName string, permNames []string) (error) {
// 	var role models.Role
// 	rRes := r.DB.Where("name = ?", roleName).First(&role)
// 	if rRes.Error != nil {
// 		if errors.Is(rRes.Error, gorm.ErrRecordNotFound) {
// 			return ErrRoleNotFound
// 		}
// 	}
// 	var perms []models.Permission
// 	for _, permName := range permNames {
// 		var perm models.Permission
// 		rRes := r.DB.Where("name = ?", permName).First(&perm)
// 		if rRes.Error != nil {
// 			if errors.Is(rRes.Error, gorm.ErrRecordNotFound) {
// 				return ErrPermissionNotFound
// 			}
// 		}
// 		perms = append(perms, perm)
// 	}
// 	for _, perm := range perms {
// 		var rolePerm models.RolePermission
// 		rRes := r.DB.Where("role_id = ?", role.ID).Where("permission_id = ?", perm.ID).First(&rolePerm)
// 		if rRes.Error != nil {
// 			res := r.DB.Create(&models.RolePermission{RoleID: role.ID, PermissionID: perm.ID})
// 			if res.Error != nil {
// 				return res.Error
// 			}
// 		}
// 	}
// 	return nil
// }

func (r *Rbac) CreateGroups(groups []string) (error) {
	for _, group := range groups {
		var dbGroup models.Group
		query := "SELECT * FROM `groups` WHERE name = ?"
		groupRepository := repositories.NewGenericRepository[entities.Group](r.DB)
		groupService := service.NewGenericService[entities.Group](groupRepository)
		_, err := groupService.FindOne((utils.GroupToEntity(&dbGroup)), query, group)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				newGroup := models.Group{
					GroupID: uuid.New(),
					Name:    group,
				}
				groupService.Create(utils.GroupToEntity(&newGroup))
			} else {
				return err
			}
		}
	}
	return nil
}

// func (r *Rbac) AssignGroupRoles(group string, roles []string) (error) {
// 	var dbGroup models.Group
// 	res := r.DB.Where("name = ?", group).First(&dbGroup)
// 	if res.Error != nil {
// 		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
// 			return ErrGroupNotFound
// 		}
// 	}
// 	var dbRoles []models.Role
// 	for _, role := range roles {
// 		var dbRole models.Role
// 		res := r.DB.Where("name = ?", role).First(&dbRole)
// 		if res.Error != nil {
// 			if errors.Is(res.Error, gorm.ErrRecordNotFound) {
// 				return ErrRoleNotFound
// 			}
// 		}
// 		dbRoles = append(dbRoles, dbRole)
// 	}
// 	for _, dbRole := range dbRoles {
// 		var groupRole models.GroupRole
// 		result := r.DB.Where("role_id = ?", dbRole.ID).Where("group_id = ?", dbGroup.ID).First(&groupRole)
// 		if result.Error != nil {
// 			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
// 				r.DB.Create(&models.GroupRole{GroupID: dbGroup.ID, RoleID: dbRole.ID})
// 			} else {
// 				return result.Error
// 			}
// 		}
// 	}
// 	return nil
// }

// func (r *Rbac) CheckGroupRole(groupID uint, roleName string) (bool, error) {
// 	var role models.Role
// 	res := r.DB.Where("name = ?", roleName).First(&role)
// 	if res.Error != nil {
// 		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
// 			return false, ErrRoleNotFound
// 		}
// 	}
// 	var groupRole models.GroupRole
// 	res = r.DB.Where("group_id = ?", groupID).Where("role_id = ?", role.ID).First(&groupRole)
// 	if res.Error != nil {
// 		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
// 			return false, nil
// 		}
// 	}
// 	return true, nil
// }

// func (r *Rbac) CheckGroupPermission(groupID uint, permName string) (bool, error) {
// 	var groupRoles []models.GroupRole
// 	res := r.DB.Where("user_id = ?", groupID).Find(&groupRoles)
// 	if res.Error != nil {
// 		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
// 			return false, nil
// 		}
// 	}
// 	var roleIDs []uint
// 	for _, r := range groupRoles {
// 		roleIDs = append(roleIDs, r.RoleID)
// 	}
// 	var perm models.Permission
// 	res = r.DB.Where("name = ?", permName).First(&perm)
// 	if res.Error != nil {
// 		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
// 			return false, global.ErrPermissionNotFound
// 		}
// 	}
// 	var rolePermission models.RolePermission
// 	res = r.DB.Where("role_id IN (?)", roleIDs).Where("permission_id = ?", perm.ID).First(&rolePermission)
// 	if res.Error != nil {
// 		return false, nil
// 	}
// 	return true, nil
// }

func (r *Rbac) AssignGroups(userID uuid.UUID, groups []string) (error) {
	var user models.User
	userQuery := "SELECT * FROM `users` WHERE user_id = ?"
	userGroupRepository := repositories.NewGenericRepository[entities.User](r.DB)
	userGroupService := service.NewGenericService[entities.User](userGroupRepository)
	_, err := userGroupService.FindOne((utils.UserToEntity(&user)), userQuery, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return global.ErrUserNotFound
		}
		return err
	}
	var groupIDs []uuid.UUID
	for _, group := range groups {
		var g models.Group
		groupQuery := "SELECT * FROM `groups` WHERE name = ?"
		groupRepository := repositories.NewGenericRepository[entities.Group](r.DB)
		groupService := service.NewGenericService[entities.Group](groupRepository)
		entity, err := groupService.FindOne((utils.GroupToEntity(&g)), groupQuery, group)
		if err != nil {
			fmt.Println(err)
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return global.ErrGroupNotFound
			}
			return err
		}
		groupIDs = append(groupIDs, entity.GroupID)
	}
	for _, groupID := range groupIDs {
		var userGroup models.UserGroup
		userGroupQuery := "SELECT * FROM `user_groups` WHERE user_id = ? and group_id = ?"
		userGroupRepository := repositories.NewGenericRepository[entities.UserGroup](r.DB)
		userGroupService := service.NewGenericService[entities.UserGroup](userGroupRepository)
		_, err := userGroupService.FindOne((utils.UserGroupToEntity(&userGroup)), userGroupQuery, userID, groupID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				newuserGroup := models.UserGroup{
					UserGroupID: uuid.New(),
					UserID: userID,
					GroupID: groupID,
				}
				userGroupService.Create((utils.UserGroupToEntity(&newuserGroup)))
			}
			return err
		}
	}
	return nil
}