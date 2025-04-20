package auth

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vterry/ddd-study/auth-server/internal/app/password"
	"github.com/vterry/ddd-study/auth-server/internal/domain/common/valueobjects"
	"github.com/vterry/ddd-study/auth-server/internal/domain/login"
	"github.com/vterry/ddd-study/auth-server/internal/domain/session"
)

func TestNewAuthService(t *testing.T) {
	sessionService := new(mockSessionService)
	loginRepo := new(mockLoginRepository)

	authService := NewAuthService(sessionService, loginRepo)

	assert.NotNil(t, authService)
	assert.Equal(t, sessionService, authService.sessionService)
	assert.Equal(t, loginRepo, authService.loginRepo)

}

func TestAuthService(t *testing.T) {
	t.Run("test login in - sucess", func(t *testing.T) {
		sessionService := new(mockSessionService)
		loginRepo := new(mockLoginRepository)

		authService := NewAuthService(sessionService, loginRepo)

		userId := uuid.New()
		hashedPass, _ := password.HashePassword("password")
		userLogin, err := login.CreateLogin(valueobjects.NewUserID(userId), hashedPass)
		assert.NoError(t, err)
		loginRepo.On("FindLoginByUserID", valueobjects.NewUserID(userId)).Return(userLogin, nil)

		sessionService.On("CreateNewSession", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&session.Session{}, nil)

		result, err := authService.LoginIn(userId.String(), "password")
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotEmpty(t, result.SessionID)
		assert.NotEmpty(t, result.JWTToken)
		assert.NotEmpty(t, result.RefreshToken)
		assert.NotEmpty(t, result.CSRFToken)
		assert.False(t, result.AccessTokenExpiresAt.IsZero())
		assert.False(t, result.RefreshTokenExpiresAt.IsZero())
	})

	t.Run("test login in - invalid user id", func(t *testing.T) {
		sessionService := new(mockSessionService)
		loginRepo := new(mockLoginRepository)

		authService := NewAuthService(sessionService, loginRepo)

		result, err := authService.LoginIn("failed", "password")
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user id format")

	})

	t.Run("test login in - invalid user not found", func(t *testing.T) {
		sessionService := new(mockSessionService)
		loginRepo := new(mockLoginRepository)

		authService := NewAuthService(sessionService, loginRepo)

		userId := uuid.New()
		loginRepo.On("FindLoginByUserID", valueobjects.NewUserID(userId)).Return(nil, errors.New("user not found"))

		result, err := authService.LoginIn(userId.String(), "password")
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid userid or password")

	})

	t.Run("test login in - invalid password", func(t *testing.T) {
		sessionService := new(mockSessionService)
		loginRepo := new(mockLoginRepository)

		authService := NewAuthService(sessionService, loginRepo)

		userId := uuid.New()
		userLogin, err := login.CreateLogin(valueobjects.NewUserID(userId), "password")
		assert.NoError(t, err)
		loginRepo.On("FindLoginByUserID", valueobjects.NewUserID(userId)).Return(userLogin, nil)

		result, err := authService.LoginIn(userId.String(), "password")
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid userid or password")

	})

	t.Run("test login in - failure on create session", func(t *testing.T) {
		sessionService := new(mockSessionService)
		loginRepo := new(mockLoginRepository)

		authService := NewAuthService(sessionService, loginRepo)

		userId := uuid.New()
		hashedPass, _ := password.HashePassword("password")
		userLogin, err := login.CreateLogin(valueobjects.NewUserID(userId), hashedPass)
		assert.NoError(t, err)
		loginRepo.On("FindLoginByUserID", valueobjects.NewUserID(userId)).Return(userLogin, nil)

		sessionService.On("CreateNewSession", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("failure on session's creation"))

		result, err := authService.LoginIn(userId.String(), "password")
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to create session")
	})

	t.Run("test create login - success", func(t *testing.T) {
		sessionService := new(mockSessionService)
		loginRepo := new(mockLoginRepository)

		authService := NewAuthService(sessionService, loginRepo)

		userId := uuid.New()
		loginRepo.On("FindLoginByUserID", valueobjects.NewUserID(userId)).Return(nil, errors.New("userid not found"))
		loginRepo.On("Save", mock.AnythingOfType("login.Login")).Return(&login.Login{}, nil)

		err := authService.CreateUserLogin(userId.String(), "password")
		assert.NoError(t, err)
	})

	t.Run("test create login - invalid user id", func(t *testing.T) {
		sessionService := new(mockSessionService)
		loginRepo := new(mockLoginRepository)

		authService := NewAuthService(sessionService, loginRepo)

		err := authService.CreateUserLogin("invalid id", "password")
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "invalid user id format:")
	})

	t.Run("test create login - user id already existis", func(t *testing.T) {
		sessionService := new(mockSessionService)
		loginRepo := new(mockLoginRepository)

		authService := NewAuthService(sessionService, loginRepo)

		userId := uuid.New()
		loginRepo.On("FindLoginByUserID", valueobjects.NewUserID(userId)).Return(nil, nil)

		err := authService.CreateUserLogin(userId.String(), "password")
		assert.NotNil(t, err)
		assert.Equal(t, err.Error(), fmt.Errorf("user with id %s already exists", userId).Error())
	})

	t.Run("test create login - failure on save login", func(t *testing.T) {
		sessionService := new(mockSessionService)
		loginRepo := new(mockLoginRepository)

		authService := NewAuthService(sessionService, loginRepo)

		userId := uuid.New()
		loginRepo.On("FindLoginByUserID", valueobjects.NewUserID(userId)).Return(nil, errors.New("userid not found"))
		loginRepo.On("Save", mock.AnythingOfType("login.Login")).Return(nil, errors.New("failure to save"))

		err := authService.CreateUserLogin(userId.String(), "password")
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "error while save login")
	})

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

type mockSessionService struct {
	mock.Mock
}

func (m *mockSessionService) CreateNewSession(userId valueobjects.UserID, jwtToken session.JWTToken, refreshToken session.JWTToken, csrfToken session.CSRFToken, expiresAt time.Time) (*session.Session, error) {
	args := m.Called(userId, jwtToken, csrfToken, expiresAt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*session.Session), args.Error(1)
}

func (m *mockSessionService) FindSessionById(sessionId session.SessionID) (*session.Session, error) {
	args := m.Called(sessionId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*session.Session), args.Error(1)
}

func (m *mockSessionService) RevokeSession(sessionId session.SessionID) error {
	args := m.Called(sessionId)
	return args.Error(0)
}

func (m *mockSessionService) RenewSession(sessionId session.SessionID, jwtToken session.JWTToken, renewToken session.JWTToken, expiresAt time.Time) (*session.Session, error) {
	args := m.Called(sessionId, jwtToken, renewToken, expiresAt)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*session.Session), args.Error(1)
}

func (m *mockSessionService) IsSessionValid(sessionID session.SessionID) bool {
	args := m.Called(sessionID)
	return args.Bool(0)
}
