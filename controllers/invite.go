package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/kmrhemant916/iam/authorization"
	"github.com/kmrhemant916/iam/entities"
	"github.com/kmrhemant916/iam/global"
	"github.com/kmrhemant916/iam/helpers"
	"github.com/kmrhemant916/iam/models"
	"github.com/kmrhemant916/iam/repositories"
	"github.com/kmrhemant916/iam/service"
	"github.com/kmrhemant916/iam/utils"
	"github.com/sethvargo/go-password/password"
	"gorm.io/gorm"
)

const (
	PasswordLength = 10
	PasswordNumDigits = 3
	PasswordNumSymbols = 3
	PasswordContainUpper = false
	PasswordContainRepeat = false
)

type InviteUserPayload struct {
    Email string `json:"email"`
    Group  []string `json:"group"`
}

type InviteResponse struct {
    Email string `json:"email"`
	Password string `json:"password"`
	OrganizationID uuid.UUID `json:"organizationID"`
	Group []string `json:"group"`
}

func (app *App)InviteUser(w http.ResponseWriter, r *http.Request) {
	var inviteUserPayload InviteUserPayload
	err := json.NewDecoder(r.Body).Decode(&inviteUserPayload)
	if err != nil {
		helpers.SendResponse(w, global.InternalServerErrorMessage, http.StatusInternalServerError)
		return
	}
	errorsList, err := utils.ValidateJSON(inviteUserPayload)
	if err != nil {
		for _, e := range errorsList {
			switch {
			case e.FailedField == "InviteUserPayload.Email" && (e.Tag == "required"):
					response := map[string]interface{}{
						"message": "email field is required",
					}
					helpers.SendResponse(w,response, http.StatusBadRequest)
					return
				case e.FailedField == "InviteUserPayload.Group" && (e.Tag == "required"):
					response := map[string]interface{}{
						"message": "group field is required",
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
	claims, ok := r.Context().Value("claims").(*Claims)
	if !ok {
		helpers.SendResponse(w, global.InternalServerErrorMessage, http.StatusInternalServerError)
		return
	}
	var user models.User
	userPassword, err := password.Generate(PasswordLength, PasswordNumDigits, PasswordNumSymbols, PasswordContainUpper, PasswordContainRepeat)
	if err != nil {
		helpers.SendResponse(w, global.InternalServerErrorMessage, http.StatusInternalServerError)
		return
	}
	hashedPassword, err := GeneratehashedPassword([]byte(userPassword))
	if err != nil {
		helpers.SendResponse(w, global.InternalServerErrorMessage, http.StatusInternalServerError)
		return
	}
	userRepository := repositories.NewGenericRepository[entities.User](app.DB)
	userService := service.NewGenericService[entities.User](userRepository)
	_, err = userService.FindOne((utils.UserToEntity(&user)), global.UserFindQueryByEmail, inviteUserPayload.Email)
	userID :=  uuid.New()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newUser := models.User{
				UserID: userID,
				Email: inviteUserPayload.Email,
				Password: hashedPassword,
				IsRoot: false,
				OrganizationID: claims.OrganizationID,
			}
			userService.Create((utils.UserToEntity(&newUser)))
		} else {
			helpers.SendResponse(w, global.InternalServerErrorMessage, http.StatusInternalServerError)
			return
		}
	}
	rbac := &authorization.Rbac{
		DB: app.DB,
	}
	var groups []string
	groups = append(groups, inviteUserPayload.Group...)
	rbac.AssignGroups(userID, groups)
	var inviteResponse InviteResponse
	inviteResponse.Email = inviteUserPayload.Email
	inviteResponse.Group = groups
	inviteResponse.OrganizationID = claims.OrganizationID
	inviteResponse.Password = userPassword
	response := map[string]interface{} {
		"message": inviteResponse,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
