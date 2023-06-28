package controllers

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

type Profile struct {
	Email string `json:"email"`
	ID uuid.UUID `json:"user_id"`
	Groups []string `json:"groups"`
	// Roles []string `json:"roles"`
	// Permission []string `json:"permission"`
}

func (app *App)GetUserProfile(w http.ResponseWriter, r *http.Request) {
	// var userProfile Profile
	userID := chi.URLParam(r, "id")
	fmt.Println("User ID:", userID)
	fmt.Fprintf(w, "User ID: %s", userID)
	// var user models.User
	// res := app.DB.Where("id = ?", userID).Find(&user)
	// if res.Error != nil {
	// 	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
	// 		response := map[string]interface{}{
	// 			"message": "User not found",
	// 		}
	// 		helpers.SendResponse(w, response, http.StatusNotFound)
	// 		return
	// 	}
	// 	helpers.SendResponse(w, nil, http.StatusInternalServerError)
	// 	return
	// }
	// userProfile.Email = user.Email
	// userProfile.ID = user.ID
	// var userGroups []models.UserGroup
	// app.DB.Where("user_id = ?", user.ID).Find(&userGroups)
	// for _, userGroup := range userGroups {
	// 	var group models.Group
	// 	app.DB.Where("group_id = ?", userGroup.GroupID).Find(&group)
	// 	userProfile.Groups = append(userProfile.Groups, group.Name)
	// }
	// response := map[string]interface{}{
	// 	"message": userProfile,
	// }
	// jsonResponse, _ := json.Marshal(response)
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// w.Write(jsonResponse)
}
