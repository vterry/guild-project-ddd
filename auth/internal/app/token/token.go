package token

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/vterry/ddd-study/auth-server/internal/infra/config"
)

func GenerateJWTToken(secret []byte, userID string, duration time.Duration) (string, *LoginClaims, error) {
	claims, err := NewLoginClaims(userID, duration)
	if err != nil {
		return "", nil, fmt.Errorf("error while generation jwt claims: %w", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", nil, fmt.Errorf("error while signing token: %w", err)
	}
	return tokenString, claims, nil
}

func ValidateJWT(tokenString string) (*LoginClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &LoginClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("invalid token signing method")
		}

		return []byte(config.Envs.JWTSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	claims, ok := token.Claims.(*LoginClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

func GenerateCSRFToken() (string, error) {
	bytes := make([]byte, config.Envs.CSRFTokenLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
