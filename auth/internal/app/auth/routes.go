package auth

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/vterry/ddd-study/auth-server/internal/app/types"
	"github.com/vterry/ddd-study/auth-server/internal/app/utils"
	middleware "github.com/vterry/ddd-study/auth-server/internal/infra/middlware"
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
	mux.Handle("POST /token/renew", middleware.Chain(http.HandlerFunc(h.handleRenew), middleware.Auhtentication()))
	mux.Handle("POST /logout", middleware.Chain(http.HandlerFunc(h.handleRevoke), middleware.Auhtentication()))
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

	sessionId := utils.RecoverSessionId(r)
	if sessionId == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("a session id must be informed"))
		return
	}

	result, err := h.service.Renew(sessionId, payload.RefreshToken)
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
	sessionId := utils.RecoverSessionId(r)
	if sessionId == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("a session id must be informed"))
		return
	}

	if err := h.service.Revoke(sessionId); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("failed to revoke session: %w", err))
		return
	}
}
