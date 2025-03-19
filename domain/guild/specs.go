package guild

import (
	"github.com/vterry/guild-project-ddd/domain/common"
	"github.com/vterry/guild-project-ddd/domain/player"
)

type CreateGuildParams struct {
	guildName  string
	guildOwner *player.Player
}

// Totaly overengineer here, but I wanted to implement an example of Specification Pattern in Go

func NewGuildSpecification() common.Specification[CreateGuildParams] {
	return common.And(
		NameNotEmptySpec(),
		NotSpecialCharacterSpec(),
		OwnerNotEmptySpec(),
		NotBeingAnotherGuildMember(),
	)
}

func NameNotEmptySpec() common.Specification[CreateGuildParams] {
	return func(b common.Base[CreateGuildParams]) error {
		guildName := b.Entity.guildName
		if len(guildName) < 4 || len(guildName) > 15 {
			return ErrInvalidGuildName
		}
		return nil
	}
}

func NotSpecialCharacterSpec() common.Specification[CreateGuildParams] {
	return func(b common.Base[CreateGuildParams]) error {
		if hasSpecialCharacters(b.Entity.guildName) {
			return ErrInvalidCharName
		}

		return nil
	}
}

func OwnerNotEmptySpec() common.Specification[CreateGuildParams] {
	return func(b common.Base[CreateGuildParams]) error {
		if b.Entity.guildOwner == nil {
			return ErrMustInformGuidOwner
		}
		return nil
	}
}

func NotBeingAnotherGuildMember() common.Specification[CreateGuildParams] {
	return func(b common.Base[CreateGuildParams]) error {

		spec := OwnerNotEmptySpec()
		if err := spec(b); err != nil {
			return err
		}

		if b.Entity.guildOwner.GetCurrentGuild() != "" {
			return ErrAnotherGuildMember
		}
		return nil
	}
}

func hasSpecialCharacters(input string) bool {
	for _, char := range input {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9')) {
			return true
		}
	}
	return false
}
