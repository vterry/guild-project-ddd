package login

import (
	"errors"

	"github.com/google/uuid"
	"github.com/vterry/ddd-study/auth-server/internal/domain/common/valueobjects"
)

var ErrCreateLogin = errors.New("failed to create login - check input params")

type Login struct {
	LoginID
	userId   valueobjects.UserID
	password string
}

func CreateLogin(id valueobjects.UserID, password string) (*Login, error) {
	if id.Equals(valueobjects.UserID{}) || password == "" {
		return nil, ErrCreateLogin
	}

	login := Login{
		LoginID:  NewLoginID(uuid.New()),
		userId:   id,
		password: password,
	}

	return &login, nil
}

func (l *Login) UserId() valueobjects.UserID {
	return l.userId
}

func (l *Login) Password() string {
	return l.password
}
