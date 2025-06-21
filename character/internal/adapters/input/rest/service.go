package rest

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/vterry/ddd-study/character/internal/core/domain/common/class"
	"github.com/vterry/ddd-study/character/internal/core/domain/common/login"
	"github.com/vterry/ddd-study/character/internal/core/ports/input/service"
	"github.com/vterry/ddd-study/character/internal/core/ports/output/gateway"
)

var (
	ErrCannotCreateCharacter = errors.New("cannot create a character")
	ErrInvalidLoginInfo      = errors.New("invalid login info")
)

type CharacterService struct {
	charaterService service.CharacterService
	loginGateway    gateway.Login
}

func NewCharacterService(characterHandler service.CharacterService, loginGateway gateway.Login) *CharacterService {
	return &CharacterService{
		charaterService: characterHandler,
		loginGateway:    loginGateway,
	}
}

func (h *CharacterService) NewCharacter(ctx context.Context, loginId string, nickname string, characterClass string) error {

	parsedId, err := uuid.Parse(loginId)

	if err != nil {
		return fmt.Errorf("%v: %w", ErrCannotCreateCharacter, err)
	}

	loginID := login.NewLoginID(parsedId)
	ok, err := h.loginGateway.IsLoginValid(ctx, loginID)

	if err != nil {
		return fmt.Errorf("%v: %w", ErrCannotCreateCharacter, err)
	}

	if !ok {
		return fmt.Errorf("%v: %w", ErrCannotCreateCharacter, ErrInvalidLoginInfo)
	}

	classValue, err := class.ParseClass(characterClass)
	if err != nil {
		return fmt.Errorf("%v: %w", ErrCannotCreateCharacter, err)
	}

	if err := h.charaterService.CreateCharacter(ctx, loginID, nickname, classValue); err != nil {
		return fmt.Errorf("%v: %w", ErrCannotCreateCharacter, err)
	}

	return nil
}
