package response

import (
	"encoding/json"
	"net/http"
)

type SuccessResponse struct {
	Success bool   `json:"success"`
	URL     string `json:"url"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func Success(w http.ResponseWriter, url string) {
	writeJSON(w, http.StatusOK, SuccessResponse{
		Success: true,
		URL:     url,
	})
}

func Error(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, ErrorResponse{
		Success: false,
		Error:   code,
		Message: message,
	})
}
