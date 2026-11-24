package httpjson

import (
	"encoding/json"
	"net/http"
)

// WriteError sends a JSON error response with the correct content type.
func WriteError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}
