package character

import (
	"errors"
	"regexp"

	"github.com/vterry/ddd-study/character/internal/core/domain/common/class"
	"github.com/vterry/ddd-study/character/internal/core/domain/common/guild"
	"github.com/vterry/ddd-study/character/internal/core/domain/common/login"
	"github.com/vterry/ddd-study/character/internal/core/domain/common/specifications"
	"github.com/vterry/ddd-study/character/internal/core/domain/common/vault"
	"github.com/vterry/ddd-study/character/internal/core/domain/inventory"
)

var (
	ErrInvalidNicknameSize  = errors.New("invalid nickname size -  must by between 4 and 15 characters")
	ErrInvalidNicknameChars = errors.New("invalid nickname charecters -  must not contain special characters")
	ErrInvalidClass         = errors.New("invalid class was provided")
	ErrInvalidLoginId       = errors.New("loginid not provided")
	ErrEmptyCharacterID     = errors.New("player id cannot be empty")
)

type CharacterParams struct {
	characterID CharacterID
	loginID     login.LoginID
	nickname    string
	class       class.Class
	inventory   inventory.Inventory
	guild       guild.GuildID
	vault       vault.VaultID
}

func NewCharacterParams(character *Character) *CharacterParams {
	return &CharacterParams{
		characterID: character.CharacterID,
		loginID:     character.loginID,
		nickname:    character.nickname,
		class:       character.class,
		inventory:   character.inventory,
		guild:       character.guild,
		vault:       character.vault,
	}
}

func ValidateNewCharacter(nickname string, loginID login.LoginID, class class.Class) error {
	params := CharacterParams{
		nickname: nickname,
		loginID:  loginID,
		class:    class,
	}

	spec := NewCharacterSpecification()
	return spec(specifications.Base[CharacterParams]{Entity: &params})
}

func NewCharacterSpecification() specifications.Specification[CharacterParams] {
	return specifications.And(
		NicknameSizeSpec(),
		NotSpecialCharacterSpec(),
		LoginNotEmptySpec(),
	)
}

func NicknameSizeSpec() specifications.Specification[CharacterParams] {
	return func(b specifications.Base[CharacterParams]) error {
		nickname := b.Entity.nickname
		if len(nickname) < 4 || len(nickname) > 15 {
			return ErrInvalidNicknameSize
		}
		return nil
	}
}

func NotSpecialCharacterSpec() specifications.Specification[CharacterParams] {
	return func(b specifications.Base[CharacterParams]) error {
		if hasSpecialCharacters(b.Entity.nickname) {
			return ErrInvalidNicknameChars
		}
		return nil
	}
}

func CharacterIDNotEmptySpec() specifications.Specification[CharacterParams] {
	return func(b specifications.Base[CharacterParams]) error {
		if b.Entity.characterID.Equals(CharacterID{}) {
			return ErrEmptyCharacterID
		}
		return nil
	}
}

func LoginNotEmptySpec() specifications.Specification[CharacterParams] {
	return func(b specifications.Base[CharacterParams]) error {
		if b.Entity.loginID == (login.LoginID{}) {
			return ErrInvalidLoginId
		}
		return nil
	}
}

func hasSpecialCharacters(input string) bool {
	re := regexp.MustCompile(`[^a-zA-Z0-9]`)
	return re.MatchString(input)
}
