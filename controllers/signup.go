package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/kmrhemant916/iam/models"
	"gorm.io/gorm"
)

type UserInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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
	user := models.User{ID: id, Email: input.Email, Password: input.Password}
	result := app.DB.Create(&user)
	if result.Error != nil {
		response := map[string]interface{}{
			"message": "User already exist",
		}
		jsonResponse, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
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
