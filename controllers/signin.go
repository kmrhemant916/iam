package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/kmrhemant916/iam/authorization"
	"github.com/kmrhemant916/iam/entities"
	"github.com/kmrhemant916/iam/global"
	"github.com/kmrhemant916/iam/helpers"
	"github.com/kmrhemant916/iam/models"
	"github.com/kmrhemant916/iam/repositories"
	"github.com/kmrhemant916/iam/service"
	"github.com/kmrhemant916/iam/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type SigninPayload struct {
	Email  string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Claims struct {
	Username string
	OrganizationID uuid.UUID
	IsRoot bool
	Groups []string
	jwt.RegisteredClaims
}

func (app *App)Signin(w http.ResponseWriter, r *http.Request) {
	var signinPayload SigninPayload
	err := json.NewDecoder(r.Body).Decode(&signinPayload)
	if err != nil {
		helpers.SendResponse(w, global.InternalServerErrorMessage, http.StatusInternalServerError)
		return
	}
	errorsList, err := utils.ValidateJSON(signinPayload)
	if err != nil {
		for _, e := range errorsList {
			switch {
			case e.FailedField == "SigninPayload.Email" && (e.Tag == "required"):
					response := map[string]interface{}{
						"message": "email field is required",
					}
					helpers.SendResponse(w,response, http.StatusBadRequest)
					return
				case e.FailedField == "SigninPayload.Password" && (e.Tag == "required"):
					response := map[string]interface{}{
						"message": "password field is required",
					}
					helpers.SendResponse(w,response, http.StatusBadRequest)
					return
				default:
					helpers.SendResponse(w, global.InternalServerErrorMessage, http.StatusInternalServerError)
					return
			}
		}
		helpers.SendResponse(w, global.InvalidRequestPayloadMessage, http.StatusBadRequest)
		return
	}
	var user models.User
	userQuery := "SELECT * FROM `users` WHERE email = ?"
	userRepository := repositories.NewGenericRepository[entities.User](app.DB)
	userService := service.NewGenericService[entities.User](userRepository)
	userEntity, err := userService.FindOne((utils.UserToEntity(&user)), userQuery, signinPayload.Email)
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
	err = bcrypt.CompareHashAndPassword([]byte(userEntity.Password), []byte(signinPayload.Password))
	if err != nil {
		response := map[string]interface{}{
			"message": "Wrong password",
		}
		helpers.SendResponse(w, response, http.StatusUnauthorized)
		return
	}
	rbac := &authorization.Rbac{
		DB: app.DB,
	}
	var userGroups []entities.UserGroup
	userGroupRepository := repositories.NewGenericRepository[entities.UserGroup](app.DB)
	userGroupService := service.NewGenericService[entities.UserGroup](userGroupRepository)
	err = userGroupService.FindMany(&userGroups, global.UserGroupFindQueryByUserID, userEntity.UserID)
	if err != nil {
		helpers.SendResponse(w, global.InternalServerErrorMessage, http.StatusInternalServerError)
		return
	}
	var groupIDs []uuid.UUID
	for _, userGroup := range userGroups {
		groupIDs = append(groupIDs, userGroup.GroupID)
	}
	var groups []string
	for _, groupID := range groupIDs {
		group, err := rbac.GetGroupNameUsingID(groupID)
		if err != nil {
			helpers.SendResponse(w, global.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}
		groups = append(groups, group)
	}
	expirationTime := time.Now().Add(120 * time.Minute)
	claims := &Claims{
		Username: signinPayload.Email,
		OrganizationID: userEntity.OrganizationID,
		IsRoot: userEntity.IsRoot,
		Groups: groups,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(app.JWTKey)
	if err != nil {
		helpers.SendResponse(w, global.InternalServerErrorMessage, http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"token": tokenString,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
