package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/kmrhemant916/iam/helpers"
	"github.com/kmrhemant916/iam/models"
	"gorm.io/gorm"
)

type Profile struct {
	Email string `json:"email"`
	ID uuid.UUID `json:"user_id"`
	Groups []string `json:"groups"`
	Roles []string `json:"roles"`
	// Permission []string `json:"permission"`
}

func (app *App)GetUserProfile(w http.ResponseWriter, r *http.Request) {
	var userProfile Profile
	var user models.User
	res := app.DB.Where("id = ?", chi.URLParam(r, "id")).Find(&user)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			response := map[string]interface{}{
				"message": "User not found",
			}
			helpers.SendResponse(w, response, http.StatusNotFound)
			return
		}
		helpers.SendResponse(w, nil, http.StatusInternalServerError)
		return
	}
	userProfile.Email = user.Email
	userProfile.ID = user.ID
	var userGroups []models.UserGroup
	app.DB.Where("user_id = ?", user.ID).Find(&userGroups)
	for _, userGroup := range userGroups {
		var group models.Group
		app.DB.Where("id = ?", userGroup.GroupID).Find(&group)
		userProfile.Groups = append(userProfile.Groups, group.Name)
	}
	// var role models.GroupRole
	// app.DB.Where("group_id = ?", userGroup.GroupID).Find(&role)




	response := map[string]interface{}{
		"message": userProfile,
	}
	jsonResponse, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
