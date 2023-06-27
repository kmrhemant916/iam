package controllers

import (
	"encoding/json"
	"net/http"
)

func (app *App)Roles(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"token": "hello world",
	}
	jsonResponse, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

