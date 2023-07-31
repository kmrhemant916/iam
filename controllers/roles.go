package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/kmrhemant916/iam/models"
)

type Role struct {
	Name string `json:"name"`
}

func (app *App)GetRoles(w http.ResponseWriter, r *http.Request) {
	var roles []models.Role
	app.DB.Find(&roles)
	var response []Role
	for _,role := range roles {
		response = append(response, Role{role.Name})
	}
	res := map[string]interface{}{
		"roles": response,
	}
	jsonResponse, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
