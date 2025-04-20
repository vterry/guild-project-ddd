package specs

import (
	"errors"
	"regexp"

	"github.com/vterry/guild-project-ddd/domain/common"
	"github.com/vterry/guild-project-ddd/domain/player"
)

var (
	ErrInvalidGuildName    = errors.New("guild's name must by between 4 and 15 characters")
	ErrInvalidCharName     = errors.New("name cannot contain special characteres")
	ErrMustInformGuidOwner = errors.New("a guild master must be inform")
)

type CreateGuildParams struct {
	guildName  string
	guildOwner *player.Player
}

func ValidateGuildCreation(guildName string, guildOwner *player.Player) error {
	params := CreateGuildParams{
		guildName:  guildName,
		guildOwner: guildOwner,
	}
	spec := NewCreateGuildSpecification()
	return spec(common.Base[CreateGuildParams]{Entity: &params})
}

// Still reflective if for that example this implementation worth -- It feel It could be implemented in a very simple way

func NewCreateGuildSpecification() common.Specification[CreateGuildParams] {
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
	return common.And(
		OwnerNotEmptySpec(),
		PlayerNotInAnotherGuildSpec(
			func(p *CreateGuildParams) *player.Player {
				return p.guildOwner
			},
		),
	)
}

func hasSpecialCharacters(input string) bool {
	re := regexp.MustCompile(`[^a-zA-Z0-9]`)
	return re.MatchString(input)
}
