package session

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/vterry/ddd-study/auth-server/internal/domain/common/valueobjects"
)

var (
	ErrCreateSession        = errors.New("error while creating session")
	ErrSessionAlreadyRevoke = errors.New("session already revoked")
)

type JWTToken string
type CSRFToken string

type Session struct {
	SessionID
	userId     valueobjects.UserID
	jwtToken   JWTToken
	renewToken JWTToken
	csrfToken  CSRFToken
	expiresAt  time.Time
	revoked    bool
}

func Create(userId valueobjects.UserID, jwt JWTToken, renewToken JWTToken, csrf CSRFToken, expiresAt time.Time) (*Session, error) {
	if err := ValidateNewSession(userId, jwt, renewToken, csrf, expiresAt); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCreateSession, err)
	}

	session := &Session{
		SessionID:  NewSessionID(uuid.New()),
		userId:     userId,
		jwtToken:   jwt,
		renewToken: renewToken,
		csrfToken:  csrf,
		expiresAt:  expiresAt,
		revoked:    false,
	}

	return session, nil
}

func (s *Session) UserId() valueobjects.UserID {
	return s.userId
}

func (s *Session) JWTToken() JWTToken {
	return s.jwtToken
}

func (s *Session) RenewToken() JWTToken {
	return s.renewToken
}

func (s *Session) CSRFToken() CSRFToken {
	return s.csrfToken
}

func (s *Session) ExpiresAt() time.Time {
	return s.expiresAt
}

func (s *Session) IsRevoked() bool {
	return s.revoked
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.expiresAt)
}

func (s *Session) Revoke() error {
	if s.revoked {
		return ErrSessionAlreadyRevoke
	}
	s.revoked = true
	return nil
}
