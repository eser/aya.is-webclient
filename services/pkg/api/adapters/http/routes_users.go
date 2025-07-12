package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/users"
	"github.com/eser/aya.is/services/pkg/lib/cursors"
)

func RegisterHTTPRoutesForUsers( //nolint:funlen,cyclop
	routes *httpfx.Router,
	logger *logfx.Logger,
	usersService *users.Service,
) {
	routes.
		Route("GET /{locale}/users", func(ctx *httpfx.Context) httpfx.Result {
			// get variables from path
			cursor := cursors.NewCursorFromRequest(ctx.Request)

			records, err := usersService.List(ctx.Request.Context(), cursor)
			if err != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText(err.Error()),
				)
			}

			return ctx.Results.JSON(records)
		}).
		HasSummary("List users").
		HasDescription("List users.").
		HasResponse(http.StatusOK)

	routes.
		Route("GET /{locale}/users/{id}", func(ctx *httpfx.Context) httpfx.Result {
			// get variables from path
			idParam := ctx.Request.PathValue("id")

			record, err := usersService.GetByID(ctx.Request.Context(), idParam)
			if err != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText(err.Error()),
				)
			}

			wrappedResponse := cursors.WrapResponseWithCursor(record, nil)

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("Get user by ID").
		HasDescription("Get user by ID.").
		HasResponse(http.StatusOK)

	// --- Auth endpoints ---
	routes.
		Route("GET /{locale}/auth/{authProvider}/login", func(ctx *httpfx.Context) httpfx.Result {
			// get auth provider from path
			authProviderName := ctx.Request.PathValue("authProvider")
			authProvider := usersService.GetAuthProvider(authProviderName)

			if authProvider == nil {
				return ctx.Results.NotFound(httpfx.WithPlainText("OAuth service not found"))
			}

			// Initiate OAuth flow
			redirectURI := ctx.Request.URL.Query().Get("redirect_uri")
			authURL, _, err := authProvider.InitiateOAuth(ctx.Request.Context(), redirectURI)
			if err != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText("OAuth initiation failed"),
				)
			}

			// Optionally set state in cookie/session
			return ctx.Results.Redirect(authURL)
		}).
		HasSummary("Auth Login").
		HasDescription("Redirects to auth provider login.").
		HasResponse(http.StatusFound)

	routes.
		Route(
			"GET /{locale}/auth/{authProvider}/callback",
			func(ctx *httpfx.Context) httpfx.Result {
				// get auth provider from path
				authProviderName := ctx.Request.PathValue("authProvider")
				authProvider := usersService.GetAuthProvider(authProviderName)

				if authProvider == nil {
					return ctx.Results.NotFound(httpfx.WithPlainText("OAuth service not found"))
				}

				url := ctx.Request.URL
				queryString := url.Query()
				code := queryString.Get("code")
				state := queryString.Get("state")

				result, err := authProvider.HandleOAuthCallback(ctx.Request.Context(), code, state)
				if err != nil {
					return ctx.Results.Unauthorized(httpfx.WithPlainText("OAuth callback failed"))
				}

				// Set JWT as cookie or return in response
				return ctx.Results.JSON(map[string]any{
					"token": result.JWT,
					"user":  result.User,
				})
			},
		).
		HasSummary("Auth Callback").
		HasDescription("Handles auth provider callback and returns JWT.").
		HasResponse(http.StatusOK)

	routes.
		Route("POST /{locale}/auth/logout", func(ctx *httpfx.Context) httpfx.Result {
			// Invalidate session logic (optional, e.g., remove session from DB)
			return ctx.Results.JSON(map[string]string{"status": "logged out"})
		}).
		HasSummary("Logout").
		HasDescription("Logs out the user.").
		HasResponse(http.StatusOK)

	routes.
		Route("POST /{locale}/auth/refresh", func(ctx *httpfx.Context) httpfx.Result {
			// Get current token from Authorization header
			auth := ctx.Request.Header.Get("Authorization")
			if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
				return ctx.Results.Unauthorized(httpfx.WithPlainText("No token provided"))
			}

			tokenStr := strings.TrimPrefix(auth, "Bearer ")

			// Use the service to refresh the token
			result, err := usersService.RefreshToken(ctx.Request.Context(), tokenStr)
			if err != nil {
				if errors.Is(err, users.ErrInvalidToken) {
					return ctx.Results.Unauthorized(httpfx.WithPlainText("Invalid token"))
				}
				if errors.Is(err, users.ErrSessionExpired) {
					return ctx.Results.Unauthorized(httpfx.WithPlainText("Session expired"))
				}

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText("Failed to refresh token"),
				)
			}

			return ctx.Results.JSON(map[string]any{
				"token":     result.JWT,
				"expiresAt": result.ExpiresAt.Unix(),
			})
		}).
		HasSummary("Refresh Token").
		HasDescription("Refreshes JWT token before expiration.").
		HasResponse(http.StatusOK)
}
