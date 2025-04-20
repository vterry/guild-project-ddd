package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type LoginClaims struct {
	UserId string `json:"userId"`
	jwt.RegisteredClaims
}

func NewLoginClaims(userId string, duration time.Duration) (*LoginClaims, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("error generation token id: %w", err)
	}

	return &LoginClaims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID.String(),
			Subject:   userId,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}, nil

}
