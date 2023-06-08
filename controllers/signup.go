package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/kmrhemant916/iam/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Organization string `json:"organization"`
}

type App struct {
	DB *gorm.DB
}

func (app *App)Signup(w http.ResponseWriter, r *http.Request) {
	var input UserInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}
	id := uuid.New()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	user := models.User{ID: id, Email: input.Email, Password: string(hashedPassword)}
	organization := models.Organization{Name: input.Organization}
	userResult := app.DB.Create(&user)
	organizationResult := app.DB.Create(&organization)
	if userResult.Error != nil ||  organizationResult.Error != nil {
		response := map[string]interface{}{
			"message": "Internal server error",
		}
		jsonResponse, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonResponse)
		return
	}
	response := map[string]interface{}{
		"message": "User stored successfully",
	}
	jsonResponse, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
