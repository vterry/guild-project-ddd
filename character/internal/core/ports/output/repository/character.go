package repository

import "github.com/vterry/ddd-study/character/internal/core/domain/character"

type CharacterRepository interface {
	Save(character.Character) error
	FindCharacterById(character.CharacterID) (*character.Character, error)
	Update(character.Character) error
}
