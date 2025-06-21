package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/vterry/ddd-study/character/internal/adapters/output/gateway/dto"
	"github.com/vterry/ddd-study/character/internal/core/domain/common/login"
	"github.com/vterry/ddd-study/character/internal/infra/keycloak"
)

type LoginGateway struct {
	keycloakClient *keycloak.KeycloakClient
	Client         *http.Client
}

func NewLoginGateway(keycloakClient *keycloak.KeycloakClient) *LoginGateway {
	return &LoginGateway{
		keycloakClient: keycloakClient,
	}
}

func (l *LoginGateway) IsLoginValid(ctx context.Context, loginId login.LoginID) (bool, error) {

	// Get admin token for accessing the admin API
	adminToken, err := l.getAdminToken(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to get admin token: %w", err)
	}

	// Make request to Keycloak Admin API
	url := fmt.Sprintf("%s/admin/realms/%s/users/%s", l.keycloakClient.BaseURL(), l.keycloakClient.Realm(), loginId.ID().String())
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+adminToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := l.Client.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to get user info: %w", err)
	}

	defer closeBodyWithError(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return false, fmt.Errorf("failed to get user info: status code %d, body: %s", resp.StatusCode, string(body))
	}

	var user dto.Login
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return false, fmt.Errorf("failed to decode user info: %w", err)
	}

	if !user.Enabled {
		return false, nil
	}

	return true, nil
}

// getAdminToken retrieves an admin token for accessing the Keycloak Admin API
func (l *LoginGateway) getAdminToken(ctx context.Context) (string, error) {
	// Create form data
	formData := url.Values{}
	formData.Set("grant_type", "client_credentials")
	formData.Set("client_id", l.keycloakClient.Oauth.ClientID)
	formData.Set("client_secret", l.keycloakClient.Oauth.ClientSecret)

	url := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", l.keycloakClient.BaseURL(), l.keycloakClient.Realm())
	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(formData.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := l.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get admin token: %w", err)
	}

	defer closeBodyWithError(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to get admin token: status code %d, body: %s", resp.StatusCode, string(body))
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}

	return tokenResp.AccessToken, nil
}

func closeBodyWithError(body io.ReadCloser) {
	if body == nil {
		return
	}
	if err := body.Close(); err != nil {
		fmt.Printf("error closing response body: %v\n", err)
	}
}
