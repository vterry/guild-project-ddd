package repository

import (
	"context"

	"github.com/vterry/ddd-study/character/internal/core/domain/character"
)

type CharacterRepository interface {
	FindCharacterById(ctx context.Context, characterId character.CharacterID) (*character.Character, error)
	Save(ctx context.Context, character character.Character) error
	Update(ctx context.Context, character character.Character) error
}
