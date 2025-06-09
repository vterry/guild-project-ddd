package token

type TokenClaims struct {
	// Common claims
	Subject  string `json:"sub"`
	Email    string `json:"email"`
	Username string `json:"preferred_username"`
	Name     string `json:"name"`

	// Access token specific claims
	Scope       string `json:"scope"`
	RealmAccess struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
	ResourceAccess map[string]struct {
		Roles []string `json:"roles"`
	} `json:"resource_access"`
}
