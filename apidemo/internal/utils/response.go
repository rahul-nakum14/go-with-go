package utils

import (
	"encoding/json"
	"net/http"
)

func SendResponse(w http.ResponseWriter, statusCode int, message string, data interface{}, err string) {
	response := APIResponse{
		Message: message,
		Data:    data,
		Error:   err,
	}

	w.WriteHeader(statusCode)

	w.Header().Set("Content-Type", "application/json")

	if jsonErr := json.NewEncoder(w).Encode(response); jsonErr != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
