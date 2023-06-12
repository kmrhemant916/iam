package helpers

import (
	"encoding/json"
	"net/http"
)

func StatusInternalServerErrorResponse(w http.ResponseWriter, _ *http.Request) {
	response := map[string]interface{}{
		"message": "Internal server error",
	}
	jsonResponse, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(jsonResponse)
}