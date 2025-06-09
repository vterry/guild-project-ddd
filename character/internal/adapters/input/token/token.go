package token

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/vterry/ddd-study/character/internal/adapters/input/keycloak"
)

type TokenValidationAdapter struct {
	keycloakClient *keycloak.KeycloakClient
}

func NewTokenValidator(keycloakClient *keycloak.KeycloakClient) *TokenValidationAdapter {
	return &TokenValidationAdapter{
		keycloakClient: keycloakClient,
	}
}

func (t *TokenValidationAdapter) TokenValidation(ctx context.Context, token string) (bool, error) {
	keySet := t.keycloakClient.Provider.VerifierContext(ctx, &oidc.Config{
		SkipClientIDCheck: true,
	})

	if keySet == nil {
		return false, fmt.Errorf("cannot verify provider access token")
	}

	jwt, err := keySet.Verify(ctx, token)
	if err != nil {
		return false, fmt.Errorf("cannot verify access token: %w", err)
	}

	var claims TokenClaims

	if err := jwt.Claims(&claims); err != nil {
		return false, fmt.Errorf("cannot verify token claims: %w", err)
	}
	return true, nil
}
