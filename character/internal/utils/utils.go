package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
)

var (
	Validate = validator.New()
)

func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}
	return json.NewDecoder(r.Body).Decode(payload)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	if err := WriteJSON(w, status, map[string]string{"error": err.Error()}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func RecoverSessionId(r *http.Request) string {
	if reqSessionId := r.Header.Get("session_id"); reqSessionId != "" {
		return reqSessionId
	}

	cookie, err := r.Cookie("session_id")
	if err != nil {
		return cookie.Value
	}
	return ""
}
