package session

import (
	"errors"
	"time"

	"github.com/vterry/ddd-study/auth-server/internal/domain/common"
	"github.com/vterry/ddd-study/auth-server/internal/domain/common/valueobjects"
)

var (
	ErrEmptySessionID     = errors.New("session id cannot be empty")
	ErrEmptyUserId        = errors.New("user id cannot be empty")
	ErrEmptyJWT           = errors.New("jwt token cannot be empty")
	ErrEmptyRenewToken    = errors.New("renew token cannot be empty")
	ErrEmptyCRSF          = errors.New("csrf token cannot be empty")
	ErrEmptyExpirionTime  = errors.New("expires at cannot be empty")
	ErrExpiresAtInThePast = errors.New("session cannot be expired")
)

type SessionParams struct {
	SessionID  SessionID
	UserID     valueobjects.UserID
	JWTToken   JWTToken
	RenewToken JWTToken
	CSRFToken  CSRFToken
	ExpiresAt  time.Time
}

func NewSessionParams(session *Session) *SessionParams {
	return &SessionParams{
		SessionID:  session.SessionID,
		UserID:     session.userId,
		JWTToken:   session.jwtToken,
		RenewToken: session.renewToken,
		CSRFToken:  session.csrfToken,
		ExpiresAt:  session.expiresAt,
	}
}

func ValidateSession(session *Session) error {
	return SessionSpecification()(common.Base[SessionParams]{
		Entity: NewSessionParams(session),
	})
}

func ValidateNewSession(userId valueobjects.UserID, jwt JWTToken, renewToken JWTToken, csrf CSRFToken, expiresAt time.Time) error {
	params := SessionParams{
		UserID:     userId,
		JWTToken:   jwt,
		RenewToken: renewToken,
		CSRFToken:  csrf,
		ExpiresAt:  expiresAt,
	}

	spec := NewSessionSpecification()
	return spec(common.Base[SessionParams]{Entity: &params})
}

func NewSessionSpecification() common.Specification[SessionParams] {
	return common.And(
		UserIdNotEmptySpec(),
		CSRFTokenNotEmptySpec(),
		JWTTokenNotEmptySpec(),
		RenewTokenNotEmptySpec(),
		ExpiresAtNotEmptySpec(),
		ExpiresCannotBeInPast(),
	)
}

func SessionSpecification() common.Specification[SessionParams] {
	return common.And(
		SessionIDNotEmptySpec(),
		NewSessionSpecification(),
	)
}

func UserIdNotEmptySpec() common.Specification[SessionParams] {
	return func(b common.Base[SessionParams]) error {
		if b.Entity.UserID.Equals(valueobjects.UserID{}) {
			return ErrEmptyUserId
		}
		return nil
	}
}

func CSRFTokenNotEmptySpec() common.Specification[SessionParams] {
	return func(b common.Base[SessionParams]) error {
		if b.Entity.CSRFToken == "" {
			return ErrEmptyCRSF
		}
		return nil
	}
}

func JWTTokenNotEmptySpec() common.Specification[SessionParams] {
	return func(b common.Base[SessionParams]) error {
		if b.Entity.JWTToken == "" {
			return ErrEmptyJWT
		}
		return nil
	}
}

func RenewTokenNotEmptySpec() common.Specification[SessionParams] {
	return func(b common.Base[SessionParams]) error {
		if b.Entity.RenewToken == "" {
			return ErrEmptyRenewToken
		}
		return nil
	}
}

func ExpiresAtNotEmptySpec() common.Specification[SessionParams] {
	return func(b common.Base[SessionParams]) error {
		if b.Entity.ExpiresAt.IsZero() {
			return ErrEmptyExpirionTime
		}
		return nil
	}
}

func ExpiresCannotBeInPast() common.Specification[SessionParams] {
	return func(b common.Base[SessionParams]) error {
		if b.Entity.ExpiresAt.Before(time.Now()) {
			return ErrExpiresAtInThePast
		}
		return nil
	}
}

func SessionIDNotEmptySpec() common.Specification[SessionParams] {
	return func(b common.Base[SessionParams]) error {
		if b.Entity.SessionID.Equals(SessionID{}) {
			return ErrEmptySessionID
		}
		return nil
	}
}
