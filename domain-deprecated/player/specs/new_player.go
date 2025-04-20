package specs

import (
	"errors"
	"regexp"

	"github.com/vterry/guild-project-ddd/domain/common"
)

var (
	ErrInvalidNicknameSize = errors.New("nickname must by between 4 and 15 characters")
	ErrInvalidNickname     = errors.New("nickname must not contain special characters")
	ErrNullPassword        = errors.New("a password must be provided")
)

type NewPlayerParams struct {
	nickname string
	password string
}

func NewPlayer(nickname string, password string) error {
	params := NewPlayerParams{
		nickname: nickname,
		password: password,
	}
	spec := NewPlayerSpec()
	return spec(common.Base[NewPlayerParams]{Entity: &params})
}

func NewPlayerSpec() common.Specification[NewPlayerParams] {
	return common.And(
		NameAndPassNotEmptySpec(),
		NotSpecialCharacterSpec(),
	)
}

func NameAndPassNotEmptySpec() common.Specification[NewPlayerParams] {
	return func(b common.Base[NewPlayerParams]) error {
		nickname := b.Entity.nickname
		password := b.Entity.password

		if len(password) == 0 {
			return ErrNullPassword
		}

		if len(nickname) < 4 || len(nickname) > 15 {
			return ErrInvalidNicknameSize
		}
		return nil
	}
}

func NotSpecialCharacterSpec() common.Specification[NewPlayerParams] {
	return func(b common.Base[NewPlayerParams]) error {
		if hasSpecialCharacters(b.Entity.nickname) {
			return ErrInvalidNickname
		}

		return nil
	}
}

func hasSpecialCharacters(input string) bool {
	re := regexp.MustCompile(`[^a-zA-Z0-9]`)
	return re.MatchString(input)
}
