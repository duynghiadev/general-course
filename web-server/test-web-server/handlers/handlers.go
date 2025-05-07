package handlers

import (
	"encoding/json"
	"net/http"
)

// APIHandler handles requests to the /api endpoint
func APIHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"message": "Hello from the server!"}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// NotFoundHandler handles 404 errors
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Page not found", http.StatusNotFound)
}
