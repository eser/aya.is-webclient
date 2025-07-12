package auth_providers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/eser/aya.is-services/pkg/ajan/logfx"
	"github.com/eser/aya.is-services/pkg/api/business/users"
)

const (
	ExpirePeriod = 24 * time.Hour
)

var (
	ErrFailedToGetAccessToken = errors.New("failed to get access token")
	ErrFailedToCreateUser     = errors.New("failed to create user")
	ErrFailedToCreateSession  = errors.New("failed to create session")
	ErrFailedToGetUserInfo    = errors.New("failed to get user info")
)

type Repository interface {
	CreateUser(ctx context.Context, user *users.User) error
	CreateSession(ctx context.Context, session *users.Session) error
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type GitHubAuthProvider struct {
	config       *users.GitHubAuthProviderConfig
	logger       *logfx.Logger
	httpClient   HTTPClient
	repo         Repository
	tokenService users.TokenService
}

func NewGitHubAuthProvider(
	config *users.GitHubAuthProviderConfig,
	logger *logfx.Logger,
	httpClient HTTPClient,
	repo Repository,
	tokenService users.TokenService,
) *GitHubAuthProvider {
	return &GitHubAuthProvider{
		config:       config,
		logger:       logger,
		httpClient:   httpClient,
		repo:         repo,
		tokenService: tokenService,
	}
}

func (g *GitHubAuthProvider) InitiateOAuth(
	ctx context.Context,
	redirectURI string,
) (string, users.OAuthState, error) {
	state := strconv.FormatInt(time.Now().UnixNano(), 10) // TODO: use secure random

	queryString := url.Values{}
	queryString.Set("client_id", g.config.ClientID)
	queryString.Set("redirect_uri", redirectURI)
	queryString.Set("state", state)
	queryString.Set("scope", "read:user user:email")

	oauthAuthorizeURL := url.URL{ //nolint:exhaustruct
		Scheme:   "https",
		Host:     "github.com",
		Path:     "/login/oauth/authorize",
		RawQuery: queryString.Encode(),
	}

	return oauthAuthorizeURL.String(), users.OAuthState{State: state, RedirectURI: redirectURI}, nil
}

func (g *GitHubAuthProvider) HandleOAuthCallback( //nolint:funlen
	ctx context.Context,
	code string,
	state string,
) (_ users.AuthResult, err error) {
	g.logger.DebugContext(ctx, "Starting GitHub OAuth callback",
		slog.String("state", state))

	// 1. Exchange code for access token
	values := url.Values{
		"client_id":     {g.config.ClientID},
		"client_secret": {g.config.ClientSecret},
		"code":          {code},
	}
	tokenReq, _ := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://github.com/login/oauth/access_token",
		strings.NewReader(values.Encode()),
	)
	tokenReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	tokenResp, tokenRespErr := g.httpClient.Do(tokenReq)
	if tokenRespErr != nil {
		g.logger.ErrorContext(ctx, "Failed to exchange code for access token",
			slog.String("error", tokenRespErr.Error()))

		return users.AuthResult{}, fmt.Errorf("%w: %w", ErrFailedToGetAccessToken, tokenRespErr)
	}

	defer func() {
		if closeErr := tokenResp.Body.Close(); closeErr != nil {
			g.logger.WarnContext(ctx, "Failed to close response body",
				slog.String("error", closeErr.Error()))
		}
	}()

	body, _ := io.ReadAll(tokenResp.Body)
	vals, _ := url.ParseQuery(string(body))

	accessToken := vals.Get("access_token")
	if accessToken == "" {
		g.logger.ErrorContext(ctx, "No access token received from GitHub")

		return users.AuthResult{}, ErrFailedToGetAccessToken
	}

	g.logger.DebugContext(ctx, "Successfully obtained GitHub access token")

	// 2. Fetch user info from GitHub
	userReq, _ := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://api.github.com/user",
		nil,
	)
	userReq.Header.Set("Authorization", "Bearer "+accessToken)

	userResp, userRespErr := g.httpClient.Do(userReq)
	if err != nil {
		return users.AuthResult{}, userRespErr //nolint:wrapcheck
	}

	defer func() {
		err = userResp.Body.Close()
	}()

	var ghUser struct {
		Login  string `json:"login"`
		Name   string `json:"name"`
		Email  string `json:"email"`
		Avatar string `json:"avatar_url"`
		ID     int64  `json:"id"`
	}

	ghUserErr := json.NewDecoder(userResp.Body).Decode(&ghUser)
	if ghUserErr != nil {
		return users.AuthResult{}, ghUserErr //nolint:wrapcheck
	}

	// 3. Upsert user in DB
	userID := fmt.Sprintf("github-%d", ghUser.ID)
	ghRemoteID := strconv.FormatInt(ghUser.ID, 10)
	now := time.Now()
	expiresAt := now.Add(ExpirePeriod)

	user := users.User{
		ID:                  userID,
		Kind:                "regular",
		Name:                ghUser.Name,
		Email:               &ghUser.Email,
		Phone:               nil,
		GithubHandle:        &ghUser.Login,
		GithubRemoteID:      &ghRemoteID,
		BskyHandle:          nil,
		XHandle:             nil,
		IndividualProfileID: nil,
		CreatedAt:           now,
		UpdatedAt:           nil,
		DeletedAt:           nil,
	}

	g.logger.DebugContext(ctx, "Creating/updating user in database",
		slog.String("user_id", user.ID),
		slog.String("github_handle", ghUser.Login))

	createUserErr := g.repo.CreateUser(ctx, &user) // ignore error if already exists
	if createUserErr != nil {
		g.logger.ErrorContext(ctx, "Failed to create user",
			slog.String("user_id", user.ID),
			slog.String("error", createUserErr.Error()))

		return users.AuthResult{}, fmt.Errorf("%w: %w", ErrFailedToCreateUser, createUserErr)
	}

	// 4. Create session in DB
	session := users.Session{
		ID:                       fmt.Sprintf("sess-%s-%d", user.ID, now.UnixNano()),
		Status:                   "active",
		OauthRequestState:        state,
		OauthRequestCodeVerifier: "", // PKCE not used here
		OauthRedirectURI:         nil,
		LoggedInUserID:           &user.ID,
		LoggedInAt:               &now,
		ExpiresAt:                &expiresAt,
		CreatedAt:                now,
		UpdatedAt:                nil,
	}

	g.logger.DebugContext(ctx, "Creating session",
		slog.String("session_id", session.ID),
		slog.String("user_id", user.ID))

	createSessionErr := g.repo.CreateSession(ctx, &session)
	if createSessionErr != nil {
		g.logger.ErrorContext(ctx, "Failed to create session",
			slog.String("session_id", session.ID),
			slog.String("error", createSessionErr.Error()))

		return users.AuthResult{}, fmt.Errorf("%w: %w", ErrFailedToCreateSession, createSessionErr)
	}

	// 5. Issue JWT
	claims := &users.JWTClaims{
		UserID:    user.ID,
		SessionID: session.ID,
		ExpiresAt: expiresAt.Unix(),
	}

	g.logger.DebugContext(ctx, "Generating JWT token",
		slog.String("user_id", user.ID),
		slog.String("session_id", session.ID))

	tokenString, tokenStringErr := g.tokenService.GenerateToken(claims)
	if tokenStringErr != nil {
		g.logger.ErrorContext(ctx, "Failed to generate JWT token",
			slog.String("user_id", user.ID),
			slog.String("error", tokenStringErr.Error()))

		return users.AuthResult{}, fmt.Errorf(
			"%w: %w",
			users.ErrFailedToGenerateToken,
			tokenStringErr,
		)
	}

	g.logger.DebugContext(ctx, "GitHub OAuth callback completed successfully",
		slog.String("user_id", user.ID),
		slog.String("session_id", session.ID),
		slog.String("github_handle", ghUser.Login))

	return users.AuthResult{
		User:      &user,
		SessionID: session.ID,
		JWT:       tokenString,
		ExpiresAt: expiresAt,
	}, nil
}
