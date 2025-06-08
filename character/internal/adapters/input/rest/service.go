package rest

import (
	"github.com/vterry/ddd-study/character/internal/core/domain/common/class"
	"github.com/vterry/ddd-study/character/internal/core/ports/input/service"
)

/*
	Para trabalhar aqui:
		- Controle transacional da criação de personagens
		- Transformações (DTO)
		- Publicação de eventos
*/

type CharacterHandler struct {
	charaterService service.CharacterService
}

func NewCharacterHandler(characterHandler service.CharacterService) *CharacterHandler {
	return &CharacterHandler{
		charaterService: characterHandler,
	}
}

func (h *CharacterHandler) NewCharacter(userId string, email string, nickname string, class class.Class) error {

	// TODO - validar informações do login
	// TODO - criar o personagem se o login for válido

	return nil
}
