package auth

import (
	"time"

	"github.com/eser/aya.is/services/pkg/api/business/users"
)

// Config types.
type GitHubAuthProviderConfig struct {
	ClientID     string `conf:"client_id"`
	ClientSecret string `conf:"client_secret"`
	Scope        string `conf:"scope"         default:"read:user user:email"`
}

type Config struct {
	JwtSecret string        `conf:"jwt_secret" default:"secret"`
	TokenTTL  time.Duration `conf:"token_ttl"  default:"365d"`

	GitHub GitHubAuthProviderConfig `conf:"github"`
}

// Auth types

type OAuthState struct {
	State       string
	RedirectURI string
}

type AuthResult struct {
	ExpiresAt   time.Time
	User        *users.User
	SessionID   string
	JWT         string
	RedirectURI string
}

type JWTClaims struct {
	UserID    string
	SessionID string
	ExpiresAt int64
}
