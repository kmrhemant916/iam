package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/kmrhemant916/iam/entities"
	"github.com/kmrhemant916/iam/global"
	"github.com/kmrhemant916/iam/helpers"
	"github.com/kmrhemant916/iam/models"
	"github.com/kmrhemant916/iam/repositories"
	"github.com/kmrhemant916/iam/service"
	"github.com/kmrhemant916/iam/utils"
	"gorm.io/gorm"
)

type Profile struct {
	Email string `json:"email"`
	ID uuid.UUID `json:"user_id"`
	Groups []string `json:"groups"`
	Roles []string `json:"roles"`
	Permission []string `json:"permission"`
}

const (
	userFindQueryByID = "SELECT * FROM `users` WHERE user_id = ?"
	userGroupFindQueryByID = "SELECT * FROM `user_groups` WHERE user_id = ?"
	groupFindQueryByID = "SELECT * FROM `groups` WHERE group_id = ?"
	groupRoleFindQueryByID = "SELECT * FROM `group_roles` WHERE group_id = ?"
	roleFindQueryByID = "SELECT * FROM `roles` WHERE id = ?"
	rolePermissionFindQueryByID = "SELECT * FROM `role_permissions` WHERE role_id = ?"
	permissionFindQueryByID = "SELECT * FROM `permissions` WHERE id = ?"
)

func (app *App)GetUserProfile(w http.ResponseWriter, r *http.Request) {
	var userProfile Profile
	var user models.User
	userRepository := repositories.NewGenericRepository[entities.User](app.DB)
	userService := service.NewGenericService[entities.User](userRepository)
	userEntity, err := userService.FindOne((utils.UserToEntity(&user)), userFindQueryByID, chi.URLParam(r, "id"))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response := map[string]interface{}{
				"message": "User doesn't exist",
			}
			helpers.SendResponse(w, response, http.StatusNotFound)
			return
		} else {
			helpers.SendResponse(w, global.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}
	}
	userProfile.Email = userEntity.Email
	userProfile.ID = userEntity.UserID
	var userGroups []entities.UserGroup
	userGroupRepository := repositories.NewGenericRepository[entities.UserGroup](app.DB)
	userGroupService := service.NewGenericService[entities.UserGroup](userGroupRepository)
	err = userGroupService.FindMany(&userGroups, userGroupFindQueryByID, userEntity.UserID)
	if err != nil {
		helpers.SendResponse(w, global.InternalServerErrorMessage, http.StatusInternalServerError)
		return
	}
	for _, userGroup := range userGroups {
		var dbGroup models.Group
		groupRepository := repositories.NewGenericRepository[entities.Group](app.DB)
		groupService := service.NewGenericService[entities.Group](groupRepository)
		groupEntity, err := groupService.FindOne((utils.GroupToEntity(&dbGroup)), groupFindQueryByID, userGroup.GroupID)
		if err != nil {
			helpers.SendResponse(w, global.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}
		userProfile.Groups = append(userProfile.Groups, groupEntity.Name)
	}
	var groupRoles []models.GroupRole
	for _, userGroup := range userGroups {
		var groupRole models.GroupRole
		groupRoleRepository := repositories.NewGenericRepository[entities.GroupRole](app.DB)
		groupRoleService := service.NewGenericService[entities.GroupRole](groupRoleRepository)
		groupRoleEntity, err := groupRoleService.FindOne((utils.GroupRoleToEntity(&groupRole)), groupRoleFindQueryByID, userGroup.GroupID)
		if err != nil {
			helpers.SendResponse(w, global.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}
		app.DB.Where("group_id = ?", userGroup.GroupID).Find(&groupRole)
		groupRoles = append(groupRoles, *utils.GroupRoleToModel(groupRoleEntity))
	}
	for _, groupRole := range groupRoles {
		var dbRole models.Role
		roleRepository := repositories.NewGenericRepository[entities.Role](app.DB)
		roleService := service.NewGenericService[entities.Role](roleRepository)
		roleEntity, err := roleService.FindOne((utils.RoleToEntity(&dbRole)), roleFindQueryByID, groupRole.RoleID)
		if err != nil {
			helpers.SendResponse(w, global.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}
		userProfile.Roles = append(userProfile.Roles, roleEntity.Name)
		var rolePermissions []entities.RolePermission
		rolePermissionRepository := repositories.NewGenericRepository[entities.RolePermission](app.DB)
		rolePermissionService := service.NewGenericService[entities.RolePermission](rolePermissionRepository)
		err = rolePermissionService.FindMany(&rolePermissions, rolePermissionFindQueryByID, roleEntity.ID)
		if err != nil {
			helpers.SendResponse(w, global.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}
		for _, rolePermission := range rolePermissions {
			var permission models.Permission
			permissionRepository := repositories.NewGenericRepository[entities.Permission](app.DB)
			permissionService := service.NewGenericService[entities.Permission](permissionRepository)
			permissionEntitiy, err := permissionService.FindOne((utils.PermissionToEntity(&permission)), permissionFindQueryByID, rolePermission.PermissionID)
			if err != nil {
				helpers.SendResponse(w, global.InternalServerErrorMessage, http.StatusInternalServerError)
				return
			}
			userProfile.Permission = append(userProfile.Permission, permissionEntitiy.Name)
		}
	}
	response := map[string]interface{}{
		"message": userProfile,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
