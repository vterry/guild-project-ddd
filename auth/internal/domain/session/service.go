package session

import (
	"errors"
	"fmt"
	"time"

	"github.com/vterry/ddd-study/auth-server/internal/domain/common/valueobjects"
	"github.com/vterry/ddd-study/auth-server/internal/domain/login"
)

var (
	ErrCannotRecoverySession           = errors.New("cannot recover session")
	ErrCannotCreateSessionWithoutLogin = errors.New("cannot create session without login")
	ErrCannotRevokeSession             = errors.New("cannot revoke session")
	ErrCannotRenewSession              = errors.New("error while revew session")
	ErrWhileCreatingSession            = errors.New("error while creating session")
	ErrJwtMustBeDifferent              = errors.New("jwt token must be different from before")
	ErrExpirationTimeMustBeGreater     = errors.New("invalid new time duration - must be greater than before")
	ErrRevokedSession                  = errors.New("session is revoked")
)

type SessionService interface {
	CreateNewSession(userId valueobjects.UserID, jwtToken JWTToken, renewToken JWTToken, csrfToken CSRFToken, expiresAt time.Time) (*Session, error)
	FindSessionById(sessionId SessionID) (*Session, error)
	RevokeSession(sessionId SessionID) error
	RenewSession(sessionId SessionID, jwtToken JWTToken, renewToken JWTToken, expiresAt time.Time) (*Session, error)
	IsSessionValid(sessionID SessionID) bool
}

type SessionServiceImpl struct {
	sessionRepository Repository
	loginRepository   login.Repository
}

func NewSessionService(sessionRepo Repository, loginRepo login.Repository) SessionService {
	return &SessionServiceImpl{
		sessionRepository: sessionRepo,
		loginRepository:   loginRepo,
	}
}

func (s *SessionServiceImpl) CreateNewSession(userId valueobjects.UserID, jwtToken JWTToken, renewToken JWTToken, csrfToken CSRFToken, expiresAt time.Time) (*Session, error) {
	login, _ := s.loginRepository.FindLoginByUserID(userId)
	if login == nil {
		return nil, ErrCannotCreateSessionWithoutLogin
	}

	session, err := Create(userId, jwtToken, renewToken, csrfToken, expiresAt)
	if err != nil {
		return nil, fmt.Errorf("%w, %v", ErrWhileCreatingSession, err)
	}

	if err := s.sessionRepository.Save(*session); err != nil {
		return nil, fmt.Errorf("%w, %v", ErrWhileCreatingSession, err)
	}

	return session, nil
}

func (s *SessionServiceImpl) FindSessionById(sessionId SessionID) (*Session, error) {
	session, err := s.sessionRepository.FindSessionByID(sessionId)
	if err != nil {
		return nil, fmt.Errorf("%w, %v", ErrCannotRecoverySession, err)
	}

	if err := ValidateSession(session); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCannotRecoverySession, err)
	}

	return session, nil
}

func (s *SessionServiceImpl) RenewSession(sessionId SessionID, jwtToken JWTToken, renewToken JWTToken, expiresAt time.Time) (*Session, error) {
	userSession, err := s.sessionRepository.FindSessionByID(sessionId)
	if err != nil {
		return nil, fmt.Errorf("%w, %v", ErrCannotRenewSession, err)
	}

	if userSession.JWTToken() == jwtToken || jwtToken == "" {
		return nil, fmt.Errorf("%w: %v", ErrCannotRenewSession, ErrJwtMustBeDifferent)
	}

	if userSession.RenewToken() == renewToken || renewToken == "" {
		return nil, fmt.Errorf("%w: %v", ErrCannotRenewSession, ErrJwtMustBeDifferent)
	}

	if userSession.ExpiresAt() == expiresAt || userSession.ExpiresAt().After(expiresAt) {
		return nil, fmt.Errorf("%w: %v", ErrCannotRenewSession, ErrExpirationTimeMustBeGreater)
	}

	if userSession.IsRevoked() {
		return nil, fmt.Errorf("%w: %v", ErrCannotRenewSession, ErrRevokedSession)
	}

	userSession.jwtToken = jwtToken
	userSession.renewToken = renewToken
	userSession.expiresAt = expiresAt

	_, err = s.sessionRepository.Update(*userSession)
	if err != nil {
		return nil, fmt.Errorf("%w, %v", ErrCannotRenewSession, err)
	}

	return userSession, nil
}

func (s *SessionServiceImpl) RevokeSession(sessionId SessionID) error {
	userSession, err := s.sessionRepository.FindSessionByID(sessionId)
	if err != nil {
		return fmt.Errorf("%w, %v", ErrCannotRevokeSession, err)
	}

	err = userSession.Revoke()
	if err != nil {
		return fmt.Errorf("%w, %v", ErrCannotRevokeSession, err)
	}

	_, err = s.sessionRepository.Update(*userSession)
	if err != nil {
		return fmt.Errorf("%w, %v", ErrCannotRevokeSession, err)
	}

	return nil
}

func (s *SessionServiceImpl) IsSessionValid(sessionID SessionID) bool {

	userSession, err := s.sessionRepository.FindSessionByID(sessionID)
	if err != nil {
		return false
	}

	if userSession.IsRevoked() || userSession.IsExpired() {
		return false
	}

	return true
}
