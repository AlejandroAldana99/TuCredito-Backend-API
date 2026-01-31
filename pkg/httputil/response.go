package httputil

import (
	"encoding/json"
	"net/http"
)

// The standard error response shape.
type ErrorBody struct {
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
	Details string `json:"details,omitempty"`
}

// Writes v as JSON with status code.
func JSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v != nil {
		_ = json.NewEncoder(w).Encode(v)
	}
}

// Writes an error response as JSON.
func Error(w http.ResponseWriter, status int, errMsg, code, details string) {
	JSON(w, status, ErrorBody{Error: errMsg, Code: code, Details: details})
}
