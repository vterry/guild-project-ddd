package login

import (
	"errors"
	"regexp"
)

var (
	ErrUserIdIsEmpty  = errors.New("user id cannot be empty")
	ErrEmailIsEmpty   = errors.New("email cannot be empty")
	ErrEmailIsInvalid = errors.New("email is not in a valid format")
	emailRegex        = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
)

type Login struct {
	userID string
	email  string
}

func NewLogin(userID string, email string) (*Login, error) {

	if userID == "" {
		return nil, ErrUserIdIsEmpty
	}

	if email == "" {
		return nil, ErrEmailIsEmpty
	}
	if !emailRegex.MatchString(email) {
		return nil, ErrEmailIsInvalid
	}

	return &Login{
		userID: userID,
		email:  email,
	}, nil
}

func (l *Login) Equals(other *Login) bool {
	if other == nil {
		return false
	}
	return l.userID == other.userID && l.email == other.email
}
