package rest

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/vterry/ddd-study/character/internal/utils"
)

type Handler struct {
	svc CharacterService
}

func NewHandler(svc CharacterService) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /character", h.handleCreateLogin)
}

func (h *Handler) handleCreateLogin(w http.ResponseWriter, r *http.Request) {
	var payload CreateCharacterRequest
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	if err := h.svc.NewCharacter(r.Context(), payload.UserID, payload.Nickname, payload.Class); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("invalid payload %v", err))
		return
	}

	if err := utils.WriteJSON(w, http.StatusOK, "character created"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
