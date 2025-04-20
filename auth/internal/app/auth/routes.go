package auth

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/vterry/ddd-study/auth-server/internal/app/types"
	"github.com/vterry/ddd-study/auth-server/internal/app/utils"
	"github.com/vterry/ddd-study/auth-server/internal/domain/session"
)

type Handler struct {
	service AuthService
}

func NewHandler(service AuthService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /login", h.handleLogin)
	mux.HandleFunc("POST /login/create", h.handleCreateLogin)
	mux.HandleFunc("POST /token/renew", h.handleRenew) // TODO - Protect
	mux.HandleFunc("POST /logout", h.handleRevoke)     // TODO - Protect
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var payload types.LoginUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	result, err := h.service.LoginIn(payload.UserId, payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	for _, cookie := range h.service.GetAuthCookies(result) {
		http.SetCookie(w, cookie)
	}

	utils.WriteJSON(w, http.StatusOK, result)

}

func (h *Handler) handleRenew(w http.ResponseWriter, r *http.Request) {
	var payload types.RenewTokenPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	reqSessionId, err := utils.RecoverSessionId(r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	sessionId, err := uuid.Parse(reqSessionId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	result, err := h.service.Renew(session.NewSessionID(sessionId), payload.RefreshToken)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	for _, cookie := range h.service.GetAuthCookies(result) {
		http.SetCookie(w, cookie)
	}

	utils.WriteJSON(w, http.StatusOK, result)

}

func (h *Handler) handleRevoke(w http.ResponseWriter, r *http.Request) {
	reqSessionId, err := utils.RecoverSessionId(r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	sessionId, err := uuid.Parse(reqSessionId)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("failed to parse session id: %w", err))
		return
	}

	if err := h.service.sessionService.RevokeSession(session.NewSessionID(sessionId)); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to revoke session: %w", err))
		return
	}
}

func (h *Handler) handleCreateLogin(w http.ResponseWriter, r *http.Request) {
	var payload types.CreateLoginPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		if errors, ok := err.(validator.ValidationErrors); ok {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
			return
		}
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload"))
		return
	}

	if err := h.service.CreateUserLogin(payload.UserId, payload.Password); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, nil)
}
