package authorization

import (
	"errors"

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

func (r *Rbac) CreateRole(roles []string) (error) {
	for _, role := range roles {
		var dbRole models.Role
		roleRepository := repositories.NewGenericRepository[entities.Role](r.DB)
		roleService := service.NewGenericService[entities.Role](roleRepository)
		_, err := roleService.FindOne((utils.RoleToEntity(&dbRole)), global.RoleFindQueryByName, role)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				newRole := models.Role{
					ID: uuid.New(),
					Name: role,
				}
				roleService.Create((utils.RoleToEntity(&newRole)))
			} else {
				return err
			}
		}
	}
	return nil
}

func (r *Rbac) CreatePermission(permissions []string) (error) {
	for _, permission := range permissions {
		var dbPermission models.Permission
		permissionRepository := repositories.NewGenericRepository[entities.Permission](r.DB)
		permissionService := service.NewGenericService[entities.Permission](permissionRepository)
		_, err := permissionService.FindOne((utils.PermissionToEntity(&dbPermission)), global.PermissionFindQueryByName, permission)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				newPermission := models.Permission{
					ID: uuid.New(),
					Name: permission,
				}
				permissionService.Create((utils.PermissionToEntity(&newPermission)))
			} else {
				return err
			}
		}
	}
	return nil
}

func (r *Rbac) AssignRolePermissions(roleName string, permNames []string) (error) {
	var role models.Role
	roleRepository := repositories.NewGenericRepository[entities.Role](r.DB)
	roleService := service.NewGenericService[entities.Role](roleRepository)
	roleEntity, err := roleService.FindOne((utils.RoleToEntity(&role)), global.RoleFindQueryByName, roleName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return global.ErrRoleNotFound
		} else {
			return err
		}
	}
	var perms []models.Permission
	for _, permName := range permNames {
		var perm models.Permission
		permissionRepository := repositories.NewGenericRepository[entities.Permission](r.DB)
		permissionService := service.NewGenericService[entities.Permission](permissionRepository)
		entitiy, err := permissionService.FindOne((utils.PermissionToEntity(&perm)), global.PermissionFindQueryByName, permName)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return global.ErrPermissionNotFound
			} else {
				return err
			}
		}
		perms = append(perms, *utils.PermissionToModel(entitiy))
	}
	for _, perm := range perms {
		var rolePerm models.RolePermission
		rolePermissionRepository := repositories.NewGenericRepository[entities.RolePermission](r.DB)
		rolePermissionService := service.NewGenericService[entities.RolePermission](rolePermissionRepository)
		_, err := rolePermissionService.FindOne((utils.RolePermissionToEntity(&rolePerm)), global.RolePermissionFindQueryByID, roleEntity.ID, perm.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				newRolePermission := models.RolePermission{
					ID: uuid.New(),
					RoleID: roleEntity.ID,
					PermissionID: perm.ID,
				}
				rolePermissionService.Create(utils.RolePermissionToEntity(&newRolePermission))
			} else {
				return err
			}
		}
	}
	return nil
}

func (r *Rbac) CreateGroups(groups []string) (error) {
	for _, group := range groups {
		var dbGroup models.Group
		groupRepository := repositories.NewGenericRepository[entities.Group](r.DB)
		groupService := service.NewGenericService[entities.Group](groupRepository)
		_, err := groupService.FindOne((utils.GroupToEntity(&dbGroup)), global.GroupFindQueryByName, group)
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

func (r *Rbac) AssignGroupRoles(group string, roles []string) (error) {
	var dbGroup models.Group
	groupRepository := repositories.NewGenericRepository[entities.Group](r.DB)
	groupService := service.NewGenericService[entities.Group](groupRepository)
	groupEntity, err := groupService.FindOne((utils.GroupToEntity(&dbGroup)), global.GroupFindQueryByName, group)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return global.ErrGroupNotFound
		} else {
			return err
		}
	}
	var dbRoles []models.Role
	for _, role := range roles {
		var dbRole models.Role
		roleRepository := repositories.NewGenericRepository[entities.Role](r.DB)
		roleService := service.NewGenericService[entities.Role](roleRepository)
		entitiy, err := roleService.FindOne((utils.RoleToEntity(&dbRole)), global.RoleFindQueryByName, role)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return global.ErrRoleNotFound
			} else {
				return err
			}
		}
		dbRoles = append(dbRoles, *utils.RoleToModel(entitiy))
	}
	for _, dbRole := range dbRoles {
		var groupRole models.GroupRole
		groupRoleRepository := repositories.NewGenericRepository[entities.GroupRole](r.DB)
		groupRoleService := service.NewGenericService[entities.GroupRole](groupRoleRepository)
		_, err := groupRoleService.FindOne((utils.GroupRoleToEntity(&groupRole)), global.GroupRoleFindQueryByGroupRoleID, dbRole.ID, groupEntity.GroupID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				newGroupRole := models.GroupRole{
					GroupRoleID: uuid.New(),
					GroupID: groupEntity.GroupID,
					RoleID: dbRole.ID,
				}
				groupRoleService.Create(utils.GroupRoleToEntity(&newGroupRole))
			} else {
				return err
			}
		}
	}
	return nil
}

func (r *Rbac) CheckGroupRole(groupID uuid.UUID, roleName string) (bool, error) {
	_, err := r.CheckGroupExistenceUsingID(groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, global.ErrGroupNotFound
		} else {
			return false, err
		}
	}
	var role models.Role
	roleRepository := repositories.NewGenericRepository[entities.Role](r.DB)
	roleService := service.NewGenericService[entities.Role](roleRepository)
	roleEntity, err := roleService.FindOne((utils.RoleToEntity(&role)), global.RoleFindQueryByName, roleName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, global.ErrRoleNotFound
		} else {
			return false, err
		}
	}
	var groupRole models.GroupRole
	groupRoleRepository := repositories.NewGenericRepository[entities.GroupRole](r.DB)
	groupRoleService := service.NewGenericService[entities.GroupRole](groupRoleRepository)
	_, err = groupRoleService.FindOne((utils.GroupRoleToEntity(&groupRole)), global.GroupRoleFindQueryByGroupRoleID, roleEntity.ID, groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
	}
	return true, nil
}

func (r *Rbac) CheckGroupPermission(groupID uuid.UUID, permName string) (bool, error) {
	_, err := r.CheckGroupExistenceUsingID(groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, global.ErrGroupNotFound
		} else {
			return false, err
		}
	}
	var perm models.Permission
	permissionRepository := repositories.NewGenericRepository[entities.Permission](r.DB)
	permissionService := service.NewGenericService[entities.Permission](permissionRepository)
	_, err = permissionService.FindOne((utils.PermissionToEntity(&perm)), global.PermissionFindQueryByName, permName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, global.ErrPermissionNotFound
		} else {
			return false, err
		}
	}
	var groupRole []entities.GroupRole
	groupRoleRepository := repositories.NewGenericRepository[entities.GroupRole](r.DB)
	groupRoleService := service.NewGenericService[entities.GroupRole](groupRoleRepository)
	err = groupRoleService.FindMany(&groupRole, global.GroupRoleFindQueryByRoleID, groupID)
	if err != nil {
		return false, err
	}
	var roleIDs []uuid.UUID
	for _, r := range groupRole {
		roleIDs = append(roleIDs, r.RoleID)
	}
	var rolePerm models.RolePermission
	rolePermissionRepository := repositories.NewGenericRepository[entities.RolePermission](r.DB)
	rolePermissionService := service.NewGenericService[entities.RolePermission](rolePermissionRepository)
	_, err = rolePermissionService.FindOne((utils.RolePermissionToEntity(&rolePerm)), global.RolePermissionFindQueryByWildcardRoleID, roleIDs, perm.ID)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (r *Rbac) AssignGroups(userID uuid.UUID, groups []string) (error) {
	var user models.User
	userGroupRepository := repositories.NewGenericRepository[entities.User](r.DB)
	userGroupService := service.NewGenericService[entities.User](userGroupRepository)
	_, err := userGroupService.FindOne((utils.UserToEntity(&user)), global.UserFindQueryByID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return global.ErrUserNotFound
		} else {
			return err
		}
	}
	var groupIDs []uuid.UUID
	for _, group := range groups {
		var g models.Group
		groupRepository := repositories.NewGenericRepository[entities.Group](r.DB)
		groupService := service.NewGenericService[entities.Group](groupRepository)
		entity, err := groupService.FindOne((utils.GroupToEntity(&g)), global.GroupFindQueryByName, group)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return global.ErrGroupNotFound
			} else {
				return err
			}
		}
		groupIDs = append(groupIDs, entity.GroupID)
	}
	for _, groupID := range groupIDs {
		var userGroup models.UserGroup
		userGroupRepository := repositories.NewGenericRepository[entities.UserGroup](r.DB)
		userGroupService := service.NewGenericService[entities.UserGroup](userGroupRepository)
		_, err := userGroupService.FindOne((utils.UserGroupToEntity(&userGroup)), global.UserGroupFindQueryByID, userID, groupID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				newUserGroup := models.UserGroup{
					UserGroupID: uuid.New(),
					UserID: userID,
					GroupID: groupID,
				}
				userGroupService.Create((utils.UserGroupToEntity(&newUserGroup)))
			} else {
				return err
			}
		}
	}
	return nil
}

func (r *Rbac) CheckGroupExistenceUsingID(groupID uuid.UUID) (bool, error) {
	var group models.Group
	groupRepository := repositories.NewGenericRepository[entities.Group](r.DB)
	groupService := service.NewGenericService[entities.Group](groupRepository)
	_, err := groupService.FindOne((utils.GroupToEntity(&group)), global.GroupFindQueryByID, groupID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, global.ErrGroupNotFound
		} else {
			return false, err
		}
	}
	return true, nil
}

func (r *Rbac) CheckGroupExistenceUsingName(groupName string) (bool, error) {
	var group models.Group
	groupRepository := repositories.NewGenericRepository[entities.Group](r.DB)
	groupService := service.NewGenericService[entities.Group](groupRepository)
	_, err := groupService.FindOne((utils.GroupToEntity(&group)), global.GroupFindQueryByName, groupName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, err
		}
	}
	return true, nil
}