package token

import "context"

type AuthService interface {
	TokenValidation(ctx context.Context, token string) (bool, error)
}
