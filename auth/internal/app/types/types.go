package types

import (
	"time"

	"github.com/vterry/ddd-study/auth-server/internal/domain/session"
)

type LoginUserPayload struct {
	UserId   string `json:"userId" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type CreateLoginPayload struct {
	UserId   string `json:"userId" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResult struct {
	SessionID             string            `json:"sessionId"`
	JWTToken              session.JWTToken  `json:"accessToken"`
	RefreshToken          session.JWTToken  `json:"refreshToken"`
	CSRFToken             session.CSRFToken `json:"-"`
	AccessTokenExpiresAt  time.Time         `json:"accessToken_expires_at"`
	RefreshTokenExpiresAt time.Time         `json:"refreshToken_expires_at"`
}

type RenewTokenPayload struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}
