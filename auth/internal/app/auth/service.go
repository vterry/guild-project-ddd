package auth

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/vterry/ddd-study/auth-server/internal/app/password"
	"github.com/vterry/ddd-study/auth-server/internal/app/token"
	"github.com/vterry/ddd-study/auth-server/internal/app/types"
	"github.com/vterry/ddd-study/auth-server/internal/domain/common/valueobjects"
	"github.com/vterry/ddd-study/auth-server/internal/domain/login"
	"github.com/vterry/ddd-study/auth-server/internal/domain/session"
	"github.com/vterry/ddd-study/auth-server/internal/infra/config"
)

var secret = []byte(config.Envs.JWTSecret)

type AuthService struct {
	sessionService session.SessionService
	loginRepo      login.Repository
}

func NewAuthService(sessionService session.SessionService, loginRepo login.Repository) *AuthService {
	return &AuthService{
		sessionService: sessionService,
		loginRepo:      loginRepo,
	}
}

func (a *AuthService) LoginIn(userId string, pass string) (*types.LoginResult, error) {
	accessDuration := time.Second * time.Duration(config.Envs.AccessDuration)
	refreshDuration := time.Second * time.Duration(config.Envs.RefreshTokenDuration)

	parsedId, err := uuid.Parse(userId)
	if err != nil {
		return nil, fmt.Errorf("invalid user id format: %w", err)
	}

	login, err := a.loginRepo.FindLoginByUserId(valueobjects.NewUserID(parsedId))
	if err != nil {
		return nil, fmt.Errorf("invalid userid or password")
	}

	if !password.ComparePasswords(login.Password(), []byte(pass)) {
		return nil, fmt.Errorf("invalid userid or password")
	}

	jwtToken, accessClaim, err := token.GenerateJWTToken(secret, login.UserId().ID().String(), accessDuration)
	if err != nil {
		return nil, fmt.Errorf("failed to create JWT Token: %w", err)
	}

	refreshToken, refreshClaim, err := token.GenerateJWTToken(secret, login.UserId().ID().String(), refreshDuration)
	if err != nil {
		return nil, fmt.Errorf("failed to create Refresh Token: %w", err)
	}

	csrfToken, err := token.GenerateCSRFToken()
	if err != nil {
		return nil, fmt.Errorf("failed to create CSRF Token: %w", err)
	}

	expiresAt := refreshClaim.ExpiresAt.Time

	newSession, err := a.sessionService.CreateNewSession(valueobjects.NewUserID(parsedId), session.JWTToken(jwtToken), session.JWTToken(refreshToken), session.CSRFToken(csrfToken), expiresAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return &types.LoginResult{
		SessionID:             newSession.SessionID.ID().String(),
		JWTToken:              session.JWTToken(jwtToken),
		RefreshToken:          session.JWTToken(refreshToken),
		CSRFToken:             session.CSRFToken(csrfToken),
		AccessTokenExpiresAt:  accessClaim.ExpiresAt.Time,
		RefreshTokenExpiresAt: refreshClaim.ExpiresAt.Time,
	}, nil
}

func (a *AuthService) CreateUserLogin(userId string, pass string) error {
	parsedId, err := uuid.Parse(userId)
	if err != nil {
		return fmt.Errorf("invalid user id format: %w", err)
	}

	_, err = a.loginRepo.FindLoginByUserId(valueobjects.NewUserID(parsedId))
	if err == nil {
		return fmt.Errorf("user with id %s already exists", userId)
	}

	hashedPass, err := password.HashePassword(pass)
	if err != nil {
		return fmt.Errorf("error while hashing password")
	}

	userLogin, err := login.CreateLogin(valueobjects.NewUserID(parsedId), hashedPass)
	if err != nil {
		return fmt.Errorf("error while create login: %w", err)
	}

	if err := a.loginRepo.Save(*userLogin); err != nil {
		return fmt.Errorf("error while save login: %w", err)
	}

	return nil
}

func (a *AuthService) Renew(reqSessionId string, refreshToken string) (*types.LoginResult, error) {
	accessDuration := time.Second * time.Duration(config.Envs.AccessDuration)
	refreshDuration := time.Second * time.Duration(config.Envs.RefreshTokenDuration)

	sessionId, err := uuid.Parse(reqSessionId)
	if err != nil {
		return nil, fmt.Errorf("invalid session id format: %w", err)
	}

	userSession, err := a.sessionService.FindSessionById(session.NewSessionID(sessionId))
	if err != nil {
		return nil, fmt.Errorf("failed to renew session: %w", err)
	}

	refreshTokenClaims, err := token.ValidateJWT(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("error verifying token: %w", err)
	}

	if refreshTokenClaims.ExpiresAt.Before(time.Now()) {
		log.Println(refreshTokenClaims.ExpiresAt)
		return nil, fmt.Errorf("refresh token is expired")
	}

	userId, err := uuid.Parse(refreshTokenClaims.UserId)
	if err != nil {
		return nil, fmt.Errorf("failed to recovery user id from claims: %w", err)
	}

	if !userSession.UserId().Equals(valueobjects.NewUserID(userId)) {
		return nil, fmt.Errorf("invalid session/refresh token")
	}

	newjwtToken, accessClaim, err := token.GenerateJWTToken(secret, userSession.UserId().ID().String(), accessDuration)
	if err != nil {
		return nil, fmt.Errorf("failed to create new JWT Token: %w", err)
	}

	newRefreshToken, refreshClaim, err := token.GenerateJWTToken(secret, userSession.UserId().ID().String(), refreshDuration)
	if err != nil {
		return nil, fmt.Errorf("failed to create Refresh Token: %w", err)
	}

	newExpiresAt := refreshClaim.ExpiresAt.Time

	newSession, err := a.sessionService.RenewSession(userSession.SessionID, session.JWTToken(newjwtToken), session.JWTToken(newRefreshToken), newExpiresAt)
	if err != nil {
		return nil, fmt.Errorf("failed to renew session id: %v - %w", userSession.SessionID.ID(), err)
	}

	return &types.LoginResult{
		SessionID:             newSession.SessionID.ID().String(),
		JWTToken:              session.JWTToken(newjwtToken),
		RefreshToken:          session.JWTToken(newRefreshToken),
		CSRFToken:             newSession.CSRFToken(),
		AccessTokenExpiresAt:  accessClaim.ExpiresAt.Time,
		RefreshTokenExpiresAt: refreshClaim.ExpiresAt.Time,
	}, nil
}

func (a *AuthService) Revoke(reqSessionId string) error {
	sessionId, err := uuid.Parse(reqSessionId)
	if err != nil {
		return fmt.Errorf("invalid session id format: %w", err)
	}

	return a.sessionService.RevokeSession(session.NewSessionID(sessionId))
}

func (a *AuthService) GetAuthCookies(result *types.LoginResult) []*http.Cookie {

	return []*http.Cookie{
		{
			Name:     "session_id",
			Value:    string(result.SessionID),
			Expires:  result.AccessTokenExpiresAt,
			HttpOnly: true,
		},
		{
			Name:     "csrf_token",
			Value:    string(result.CSRFToken),
			Expires:  result.AccessTokenExpiresAt,
			HttpOnly: false,
			SameSite: http.SameSiteLaxMode,
		},
	}
}
