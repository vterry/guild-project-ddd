package dao

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/vterry/ddd-study/auth-server/internal/domain/common/valueobjects"
	"github.com/vterry/ddd-study/auth-server/internal/domain/login"
	"github.com/vterry/ddd-study/auth-server/internal/domain/session"
)

func LoginToDAO(l login.Login) Login {
	login := Login{
		LoginID:  l.LoginID.ID().String(),
		UserId:   l.UserId().ID().String(),
		Password: l.Password(),
	}
	return login
}

func SessionToDAO(s session.Session) Session {
	return Session{
		SessionId:  s.SessionID.ID().String(),
		UserId:     s.UserId().ID().String(),
		JwtToken:   string(s.JWTToken()),
		RenewToken: string(s.RenewToken()),
		CsrfToken:  string(s.CSRFToken()),
		ExpiresAt:  s.ExpiresAt(),
		Revoked:    s.IsRevoked(),
	}
}

func DAOtoLogin(l Login) (login.Login, error) {
	loginId, err := uuid.Parse(l.LoginID)
	if err != nil {
		return login.Login{}, fmt.Errorf("error while parsing login id: %w", err)
	}

	userId, err := uuid.Parse(l.UserId)
	if err != nil {
		return login.Login{}, fmt.Errorf("error while parsing user id: %w", err)
	}

	logindId := login.NewLoginID(loginId)

	return login.NewLoginBuilder().
		WithLoginId(logindId).
		WithUserId(valueobjects.NewUserID(userId)).
		WithPassword(l.Password).
		Build(), nil
}

func DAOtoSession(s Session) (session.Session, error) {
	sessionId, err := uuid.Parse(s.SessionId)
	if err != nil {
		return session.Session{}, fmt.Errorf("error while parsing session id: %w", err)
	}

	userId, err := uuid.Parse(s.UserId)
	if err != nil {
		return session.Session{}, fmt.Errorf("error while parsing user id: %w", err)
	}

	session := session.NewSessionBuilder().
		WithSessionId(session.NewSessionID(sessionId)).
		WithUserId(valueobjects.NewUserID(userId)).
		WithJWTToken(session.JWTToken(s.JwtToken)).
		WithRenewToken(session.JWTToken(s.RenewToken)).
		WithCSRFToken(session.CSRFToken(s.CsrfToken)).
		WithExpirationAt(s.ExpiresAt).
		WithRevokeStatus(s.Revoked).
		Build()

	return session, nil
}
