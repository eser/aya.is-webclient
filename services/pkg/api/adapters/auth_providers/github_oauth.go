package auth_providers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/eser/aya.is-services/pkg/ajan/logfx"
	"github.com/eser/aya.is-services/pkg/api/business/users"
	"github.com/golang-jwt/jwt/v5"
)

const (
	ExpirePeriod = 24 * time.Hour
)

var ErrFailedToGetAccessToken = errors.New("failed to get access token")

type Repository interface {
	CreateUser(ctx context.Context, user *users.User) error
	CreateSession(ctx context.Context, session *users.Session) error
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type GitHubAuthProvider struct {
	logger     *logfx.Logger
	httpClient HTTPClient
	repo       Repository

	ClientID     string
	ClientSecret string
	RedirectBase string
}

func NewGitHubAuthProvider(
	logger *logfx.Logger,
	httpClient HTTPClient,
	repo Repository,
) *GitHubAuthProvider {
	return &GitHubAuthProvider{
		logger:     logger,
		httpClient: httpClient,
		repo:       repo,

		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectBase: os.Getenv("GITHUB_REDIRECT_BASE"),
	}
}

func (g *GitHubAuthProvider) InitiateOAuth(
	ctx context.Context,
	redirectURI string,
) (string, users.OAuthState, error) {
	state := strconv.FormatInt(time.Now().UnixNano(), 10) // TODO: use secure random

	queryString := url.Values{}
	queryString.Set("client_id", g.ClientID)
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
	// 1. Exchange code for access token
	values := url.Values{
		"client_id":     {g.ClientID},
		"client_secret": {g.ClientSecret},
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
	if err != nil {
		return users.AuthResult{}, tokenRespErr //nolint:wrapcheck
	}

	defer func() {
		err = tokenResp.Body.Close()
	}()

	body, _ := io.ReadAll(tokenResp.Body)
	vals, _ := url.ParseQuery(string(body))

	accessToken := vals.Get("access_token")
	if accessToken == "" {
		return users.AuthResult{}, ErrFailedToGetAccessToken
	}

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

	createUserErr := g.repo.CreateUser(ctx, &user) // ignore error if already exists
	if createUserErr != nil {
		return users.AuthResult{}, createUserErr //nolint:wrapcheck
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

	createSessionErr := g.repo.CreateSession(ctx, &session)
	if createSessionErr != nil {
		return users.AuthResult{}, createSessionErr //nolint:wrapcheck
	}

	// 5. Issue JWT
	claims := users.JWTClaims{
		UserID:    user.ID,
		SessionID: session.ID,
		ExpiresAt: expiresAt.Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    claims.UserID,
		"session_id": claims.SessionID,
		"exp":        claims.ExpiresAt,
	})
	secret := os.Getenv("JWT_SECRET")

	tokenString, tokenStringErr := jwtToken.SignedString([]byte(secret))
	if tokenStringErr != nil {
		return users.AuthResult{}, tokenStringErr //nolint:wrapcheck
	}

	return users.AuthResult{
		User:      &user,
		SessionID: session.ID,
		JWT:       tokenString,
	}, nil
}
