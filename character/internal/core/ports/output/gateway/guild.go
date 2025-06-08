package gateway

import "github.com/vterry/ddd-study/character/internal/core/domain/character"

type Guild interface {
	LeaveGuild(characterId character.CharacterID) error
}
