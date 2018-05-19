package model

// AuthRequest Auth Request
type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Token Auth Response
type Token struct {
	// jwt
	AccessToken string `json:"access_token"`
	// jwt
	RefreshToken string            `json:"refresh_token"`
	ExpiresAt    int64             `json:"expires_at"`
	Scopes       []string          `json:"scopes,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}