package session

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vterry/ddd-study/auth-server/internal/domain/common/valueobjects"
	"github.com/vterry/ddd-study/auth-server/internal/domain/login"
)

func TestNewSessionService(t *testing.T) {
	sessionRepo := new(mockSessionRepository)
	loginRepo := new(mockLoginRepository)

	service := NewSessionService(sessionRepo, loginRepo) // lessons learned: ensure that repositories implement its intefaces correctly
	assert.NotNil(t, service)

	serviceImpl, ok := service.(*SessionServiceImpl)
	assert.True(t, ok)
	assert.Equal(t, sessionRepo, serviceImpl.sessionRepository)
	assert.Equal(t, loginRepo, serviceImpl.loginRepository)

}

func TestSessionService(t *testing.T) {
	t.Run("create session - success", func(t *testing.T) {
		sessionRepo := new(mockSessionRepository)
		loginRepo := new(mockLoginRepository)
		service := NewSessionService(sessionRepo, loginRepo)

		userId := valueobjects.NewUserID(uuid.New())
		userLogin, err := login.CreateLogin(userId, "password")
		assert.NoError(t, err)
		loginRepo.On("FindLoginByUserID", userId).Return(userLogin, nil)

		sessionRepo.On("Save", mock.AnythingOfType("session.Session")).Return(&Session{}, nil)

		result, err := service.CreateNewSession(userId, "jwt-token", "renew-token", "csrf-token", time.Now().Add(24*time.Hour))
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.False(t, result.IsRevoked())
		assert.False(t, result.IsExpired())
	})

	t.Run("create session - error while creating session", func(t *testing.T) {
		sessionRepo := new(mockSessionRepository)
		loginRepo := new(mockLoginRepository)
		service := NewSessionService(sessionRepo, loginRepo)

		userId := valueobjects.NewUserID(uuid.New())
		userLogin, err := login.CreateLogin(userId, "password")
		assert.NoError(t, err)
		loginRepo.On("FindLoginByUserID", userId).Return(userLogin, nil)

		sessionRepo.On("Save", mock.AnythingOfType("session.Session")).Return(nil, ErrWhileCreatingSession)

		result, err := service.CreateNewSession(userId, "jwt-token", "renew-token", "csrf-token", time.Now().Add(24*time.Hour))
		assert.Nil(t, result)
		assert.ErrorIs(t, err, ErrWhileCreatingSession)
	})

	t.Run("create session - login user not found", func(t *testing.T) {
		sessionRepo := new(mockSessionRepository)
		loginRepo := new(mockLoginRepository)
		service := NewSessionService(sessionRepo, loginRepo)

		userId := valueobjects.NewUserID(uuid.New())
		loginRepo.On("FindLoginByUserID", mock.Anything).Return(nil, nil)

		result, err := service.CreateNewSession(userId, "jwt-token", "renew-token", "csrf-token", time.Now().Add(24*time.Hour))
		assert.Nil(t, result)
		assert.ErrorIs(t, err, ErrCannotCreateSessionWithoutLogin)

	})

	t.Run("create session - error while creating session", func(t *testing.T) {
		sessionRepo := new(mockSessionRepository)
		loginRepo := new(mockLoginRepository)
		service := NewSessionService(sessionRepo, loginRepo)

		userId := valueobjects.NewUserID(uuid.New())
		userLogin, err := login.CreateLogin(userId, "password")
		assert.NoError(t, err)
		loginRepo.On("FindLoginByUserID", userId).Return(userLogin, nil)

		result, err := service.CreateNewSession(userId, "", "", "", time.Time{})
		assert.Nil(t, result)
		assert.ErrorIs(t, err, ErrWhileCreatingSession)
	})

	t.Run("find session by id - success", func(t *testing.T) {
		sessionRepo := new(mockSessionRepository)
		loginRepo := new(mockLoginRepository)
		service := NewSessionService(sessionRepo, loginRepo)

		userId := valueobjects.NewUserID(uuid.New())
		session, err := Create(userId, "jwt-token", "renew-token", "csrf-token", time.Now().Add(2*time.Hour))
		assert.NoError(t, err)

		sessionRepo.On("FindSessionByID", session.SessionID).Return(session, nil)
		_, err = service.FindSessionById(session.SessionID)
		assert.NoError(t, err)
	})

	t.Run("find session by id - session not found", func(t *testing.T) {
		sessionRepo := new(mockSessionRepository)
		loginRepo := new(mockLoginRepository)
		service := NewSessionService(sessionRepo, loginRepo)

		sessionId := NewSessionID(uuid.New())

		sessionRepo.On("FindSessionByID", sessionId).Return(nil, errors.New("session id not found"))
		_, err := service.FindSessionById(sessionId)

		assert.ErrorIs(t, err, ErrCannotRecoverySession)
		assert.ErrorContains(t, err, "session id not found")
	})

	t.Run("find session by id - failure on recovery", func(t *testing.T) {
		sessionRepo := new(mockSessionRepository)
		loginRepo := new(mockLoginRepository)
		service := NewSessionService(sessionRepo, loginRepo)

		sessionId := NewSessionID(uuid.New())

		sessionRepo.On("FindSessionByID", sessionId).Return(&Session{}, nil)
		_, err := service.FindSessionById(sessionId)

		assert.ErrorIs(t, err, ErrCannotRecoverySession)
		assert.ErrorContains(t, err, ErrEmptySessionID.Error())
		assert.ErrorContains(t, err, ErrEmptyUserId.Error())
		assert.ErrorContains(t, err, ErrEmptyJWT.Error())
		assert.ErrorContains(t, err, ErrEmptyRenewToken.Error())
		assert.ErrorContains(t, err, ErrEmptyCRSF.Error())
		assert.ErrorContains(t, err, ErrEmptyExpirionTime.Error())
	})

	t.Run("renew session - success", func(t *testing.T) {
		newExpiresAt := time.Now().Add(4 * time.Hour)

		sessionRepo := new(mockSessionRepository)
		loginRepo := new(mockLoginRepository)
		service := NewSessionService(sessionRepo, loginRepo)

		userId := valueobjects.NewUserID(uuid.New())
		session, err := Create(userId, "jwt-token", "renew-token", "csrf-token", time.Now().Add(2*time.Hour))
		assert.NoError(t, err)

		sessionRepo.On("FindSessionByID", session.SessionID).Return(session, nil)

		stubSession := *session
		stubSession.jwtToken = "new-jwt-token"
		stubSession.renewToken = "new-refresh-token"
		stubSession.expiresAt = newExpiresAt

		sessionRepo.On("Update", stubSession).Return(&stubSession, nil)

		_, err = service.RenewSession(session.SessionID, "new-jwt-token", "new-refresh-token", newExpiresAt)
		assert.NoError(t, err)
	})

	t.Run("renew session - session not found", func(t *testing.T) {
		newExpiresAt := time.Now().Add(4 * time.Hour)

		sessionRepo := new(mockSessionRepository)
		loginRepo := new(mockLoginRepository)
		service := NewSessionService(sessionRepo, loginRepo)

		sessionId := NewSessionID(uuid.New())

		sessionRepo.On("FindSessionByID", sessionId).Return(nil, errors.New("session id not found"))
		_, err := service.RenewSession(sessionId, "new-token", "refresh-token", newExpiresAt)
		assert.ErrorIs(t, err, ErrCannotRenewSession)
		assert.ErrorContains(t, err, "session id not found")
	})

	t.Run("renew session - invalid jwt token", func(t *testing.T) {
		newExpiresAt := time.Now().Add(4 * time.Hour)

		sessionRepo := new(mockSessionRepository)
		loginRepo := new(mockLoginRepository)
		service := NewSessionService(sessionRepo, loginRepo)

		userId := valueobjects.NewUserID(uuid.New())
		session, err := Create(userId, "jwt-token", "renew-token", "csrf-token", time.Now().Add(2*time.Hour))
		assert.NoError(t, err)

		sessionRepo.On("FindSessionByID", session.SessionID).Return(session, nil)

		_, err = service.RenewSession(session.SessionID, "jwt-token", "new-refresh-token", newExpiresAt)
		assert.ErrorIs(t, err, ErrCannotRenewSession)
		assert.ErrorContains(t, err, ErrJwtMustBeDifferent.Error())
	})

	t.Run("renew session - invalid refresh token", func(t *testing.T) {
		newExpiresAt := time.Now().Add(4 * time.Hour)

		sessionRepo := new(mockSessionRepository)
		loginRepo := new(mockLoginRepository)
		service := NewSessionService(sessionRepo, loginRepo)

		userId := valueobjects.NewUserID(uuid.New())
		session, err := Create(userId, "jwt-token", "renew-token", "csrf-token", time.Now().Add(2*time.Hour))
		assert.NoError(t, err)

		sessionRepo.On("FindSessionByID", session.SessionID).Return(session, nil)

		_, err = service.RenewSession(session.SessionID, "new-jwt-token", "renew-token", newExpiresAt)
		assert.ErrorIs(t, err, ErrCannotRenewSession)
		assert.ErrorContains(t, err, ErrJwtMustBeDifferent.Error())
	})

	t.Run("renew session - invalid duration token", func(t *testing.T) {
		invalidExpireTime := time.Now().Add(1 * time.Hour)

		sessionRepo := new(mockSessionRepository)
		loginRepo := new(mockLoginRepository)
		service := NewSessionService(sessionRepo, loginRepo)

		userId := valueobjects.NewUserID(uuid.New())
		session, err := Create(userId, "jwt-token", "renew-token", "csrf-token", time.Now().Add(24*time.Hour))
		assert.NoError(t, err)

		sessionRepo.On("FindSessionByID", session.SessionID).Return(session, nil)

		_, err = service.RenewSession(session.SessionID, "new-jwt-token", "new-refresh-token", invalidExpireTime)
		assert.ErrorIs(t, err, ErrCannotRenewSession)
		assert.ErrorContains(t, err, ErrExpirationTimeMustBeGreater.Error())
	})

	t.Run("renew session - failed to update session", func(t *testing.T) {
		newExpiresAt := time.Now().Add(4 * time.Hour)

		sessionRepo := new(mockSessionRepository)
		loginRepo := new(mockLoginRepository)
		service := NewSessionService(sessionRepo, loginRepo)

		userId := valueobjects.NewUserID(uuid.New())
		session, err := Create(userId, "jwt-token", "renew-token", "csrf-token", time.Now().Add(2*time.Hour))
		assert.NoError(t, err)

		sessionRepo.On("FindSessionByID", session.SessionID).Return(session, nil)

		stubSession := *session
		stubSession.jwtToken = "new-jwt-token"
		stubSession.renewToken = "new-refresh-token"
		stubSession.expiresAt = newExpiresAt

		sessionRepo.On("Update", stubSession).Return(&stubSession, errors.New("failed to run update command"))

		_, err = service.RenewSession(session.SessionID, "new-jwt-token", "new-refresh-token", newExpiresAt)

		assert.ErrorIs(t, err, ErrCannotRenewSession)
		assert.ErrorContains(t, err, "failed to run update command")
	})

	t.Run("revoke session - success", func(t *testing.T) {
		sessionRepo := new(mockSessionRepository)
		loginRepo := new(mockLoginRepository)
		service := NewSessionService(sessionRepo, loginRepo)

		userId := valueobjects.NewUserID(uuid.New())
		session, err := Create(userId, "jwt-token", "renew-token", "csrf-token", time.Now().Add(24*time.Hour))
		assert.NoError(t, err)

		sessionRepo.On("FindSessionByID", session.SessionID).Return(session, nil)

		stubSession := *session
		stubSession.Revoke()
		sessionRepo.On("Update", stubSession).Return(&stubSession, nil)

		err = service.RevokeSession(session.SessionID)
		assert.NoError(t, err)
	})

	t.Run("revoke session - session not found", func(t *testing.T) {
		sessionRepo := new(mockSessionRepository)
		loginRepo := new(mockLoginRepository)
		service := NewSessionService(sessionRepo, loginRepo)

		sessionId := NewSessionID(uuid.New())

		sessionRepo.On("FindSessionByID", sessionId).Return(nil, errors.New("session id not found"))

		err := service.RevokeSession(sessionId)
		assert.ErrorIs(t, err, ErrCannotRevokeSession)
		assert.ErrorContains(t, err, "session id not found")
	})

	t.Run("revoke session - session already revoked", func(t *testing.T) {
		sessionRepo := new(mockSessionRepository)
		loginRepo := new(mockLoginRepository)
		service := NewSessionService(sessionRepo, loginRepo)

		userId := valueobjects.NewUserID(uuid.New())
		session, err := Create(userId, "jwt-token", "renew-token", "csrf-token", time.Now().Add(24*time.Hour))
		assert.NoError(t, err)

		sessionRepo.On("FindSessionByID", session.SessionID).Return(session, nil)
		session.Revoke()

		err = service.RevokeSession(session.SessionID)
		assert.ErrorIs(t, err, ErrCannotRevokeSession)
		assert.ErrorContains(t, err, ErrSessionAlreadyRevoke.Error())
	})

	t.Run("revoke session - failed to update session", func(t *testing.T) {
		sessionRepo := new(mockSessionRepository)
		loginRepo := new(mockLoginRepository)
		service := NewSessionService(sessionRepo, loginRepo)

		userId := valueobjects.NewUserID(uuid.New())
		session, err := Create(userId, "jwt-token", "renew-token", "csrf-token", time.Now().Add(24*time.Hour))
		assert.NoError(t, err)

		sessionRepo.On("FindSessionByID", session.SessionID).Return(session, nil)

		stubSession := *session
		stubSession.Revoke()
		sessionRepo.On("Update", stubSession).Return(nil, errors.New("failed to run update command"))

		err = service.RevokeSession(session.SessionID)
		assert.ErrorIs(t, err, ErrCannotRevokeSession)
		assert.ErrorContains(t, err, "failed to run update command")
	})

	t.Run("validate session - true", func(t *testing.T) {
		sessionRepo := new(mockSessionRepository)
		loginRepo := new(mockLoginRepository)
		service := NewSessionService(sessionRepo, loginRepo)

		userId := valueobjects.NewUserID(uuid.New())
		session, err := Create(userId, "jwt-token", "renew-token", "csrf-token", time.Now().Add(24*time.Hour))
		assert.NoError(t, err)
		sessionRepo.On("FindSessionByID", session.SessionID).Return(session, nil)

		ok := service.IsSessionValid(session.SessionID)
		assert.True(t, ok)

	})

	t.Run("validate session - invalid ", func(t *testing.T) {
		sessionRepo := new(mockSessionRepository)
		loginRepo := new(mockLoginRepository)
		service := NewSessionService(sessionRepo, loginRepo)

		userId := valueobjects.NewUserID(uuid.New())
		session, err := Create(userId, "jwt-token", "renew-token", "csrf-token", time.Now().Add(24*time.Hour))
		assert.NoError(t, err)
		session.Revoke()
		sessionRepo.On("FindSessionByID", session.SessionID).Return(session, nil)

		ok := service.IsSessionValid(session.SessionID)
		assert.False(t, ok)

	})

	t.Run("validate session - invalid (not found)", func(t *testing.T) {
		sessionRepo := new(mockSessionRepository)
		loginRepo := new(mockLoginRepository)
		service := NewSessionService(sessionRepo, loginRepo)

		sessionRepo.On("FindSessionByID", mock.Anything).Return(nil, errors.New("session not found"))

		ok := service.IsSessionValid(NewSessionID(uuid.New()))
		assert.False(t, ok)

	})

}

type mockSessionRepository struct {
	mock.Mock
}

func (m *mockSessionRepository) Save(sess Session) error {
	args := m.Called(sess)
	if args.Get(0) == nil {
		return args.Error(1)
	}
	return args.Error(1)
}

func (m *mockSessionRepository) Update(sess Session) (*Session, error) {
	args := m.Called(sess)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Session), args.Error(1)
}

func (m *mockSessionRepository) FindSessionByID(sesId SessionID) (*Session, error) {
	args := m.Called(sesId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Session), args.Error(1)
}

type mockLoginRepository struct {
	mock.Mock
}

func (m *mockLoginRepository) Save(inputLogin login.Login) error {
	args := m.Called(inputLogin)
	if args.Get(0) == nil {
		return args.Error(1)
	}
	return args.Error(1)
}

func (m *mockLoginRepository) FindLoginByUserID(userID valueobjects.UserID) (*login.Login, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*login.Login), args.Error(1)
}

func (m *mockLoginRepository) UpdatePassword(userId valueobjects.UserID, password string) error {
	args := m.Called(userId, password)
	if args.Get(0) == nil {
		return args.Error(1)
	}
	return args.Error(1)
}
