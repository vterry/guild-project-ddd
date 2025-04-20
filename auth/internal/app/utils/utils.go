package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
)

var (
	ErrSessionIDNotInformed = errors.New("session id was not informed")
	Validate                = validator.New()
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
	WriteJSON(w, status, map[string]string{"error": err.Error()})
}

func RecoverSessionId(r *http.Request) (string, error) {
	if sessionID := r.Header.Get("session_id"); sessionID != "" {
		return sessionID, nil
	}

	cookie, err := r.Cookie("session_id")
	if err != nil {
		return "", ErrSessionIDNotInformed
	}
	return cookie.Value, nil
}
