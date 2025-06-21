package keycloak

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/vterry/ddd-study/character/internal/infra/config"
	"golang.org/x/oauth2"
)

type KeycloakClient struct {
	Provider *oidc.Provider
	OIDC     *oidc.IDTokenVerifier
	Oauth    oauth2.Config
	config   *config.KeycloakConfig
}

type Option func(*config.KeycloakConfig)

func NewKeycloakClient(ctx context.Context, config *config.KeycloakConfig, options ...Option) (*KeycloakClient, error) {
	for _, opt := range options {
		opt(config)
	}

	providerURL := fmt.Sprintf("%s/realms/%s", config.BaseURL, config.Realm)

	provider, err := oidc.NewProvider(ctx, providerURL)
	if err != nil {
		return nil, fmt.Errorf("cannot get a oidc provider: %w", err)
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: config.ClientID})

	oauth2 := oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "roles"},
	}

	return &KeycloakClient{
		Oauth:    oauth2,
		OIDC:     verifier,
		Provider: provider,
		config:   config,
	}, nil
}

func (k *KeycloakClient) BaseURL() string {
	return k.config.BaseURL
}

func (k *KeycloakClient) Realm() string {
	return k.config.Realm
}
