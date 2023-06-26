package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/kmrhemant916/iam/helpers"
	"github.com/kmrhemant916/iam/models"
	"golang.org/x/crypto/bcrypt"
)

type SigninPayload struct {
	Email  string `gorm:"unique" json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var jwtKey = []byte("my_secret_key")

func (app *App)Signin(w http.ResponseWriter, r *http.Request) {
	var signinPayload SigninPayload
	err := json.NewDecoder(r.Body).Decode(&signinPayload)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}
	var user models.User
	if err := app.DB.Where("email = ?", signinPayload.Email).First(&user).Error; err != nil {
		response := map[string]interface{}{
			"message": "User doesn't exist",
		}
		helpers.SendResponse(w, response, http.StatusNotFound)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(signinPayload.Password))
	if err != nil {
		response := map[string]interface{}{
			"message": "Wrong password",
		}
		helpers.SendResponse(w, response, http.StatusUnauthorized)
		return
	}
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: signinPayload.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		response := map[string]interface{}{
			"message": "Internal server error",
		}
		helpers.SendResponse(w, response, http.StatusInternalServerError)
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

