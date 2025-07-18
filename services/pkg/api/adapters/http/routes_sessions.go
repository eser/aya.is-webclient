package http

import (
	"log/slog"
	"net/http"

	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
	"github.com/eser/aya.is/services/pkg/api/business/profiles"
	"github.com/eser/aya.is/services/pkg/api/business/users"
)

func RegisterHTTPRoutesForSessions( //nolint:funlen
	routes *httpfx.Router,
	logger *logfx.Logger,
	authService *auth.Service,
	userService *users.Service,
	profileService *profiles.Service,
) {
	// Register authenticated route with auth middleware
	routes.Route(
		"GET /{locale}/sessions/current",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			// Get user ID from context (set by auth middleware)
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText("Session ID not found in context"),
				)
			}

			// Get session data
			session, err := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if err != nil {
				return ctx.Results.Error(
					http.StatusNotFound,
					httpfx.WithPlainText("Session not found"),
				)
			}

			// Prepare response with session data
			response := map[string]any{
				"id":               session.ID,
				"user":             nil,
				"selected_profile": nil,
			}

			if session.LoggedInUserID != nil { //nolint:nestif
				user, userErr := userService.GetByID(ctx.Request.Context(), *session.LoggedInUserID)
				if userErr != nil {
					return ctx.Results.Error(
						http.StatusNotFound,
						httpfx.WithPlainText("User not found"),
					)
				}

				if user != nil {
					response["user"] = map[string]any{
						"id":                    user.ID,
						"kind":                  user.Kind,
						"name":                  user.Name,
						"email":                 user.Email,
						"phone":                 user.Phone,
						"github_handle":         user.GithubHandle,
						"github_remote_id":      user.GithubRemoteID,
						"bsky_handle":           user.BskyHandle,
						"x_handle":              user.XHandle,
						"individual_profile_id": user.IndividualProfileID,
						"created_at":            user.CreatedAt,
						"updated_at":            user.UpdatedAt,
					}

					// If user has an individual profile, fetch it
					if user.IndividualProfileID != nil {
						// Get locale from path
						locale := ctx.Request.PathValue("locale")
						logger.InfoContext(ctx.Request.Context(), "Fetching profile",
							slog.String("locale", locale),
							slog.String("profile_id", *user.IndividualProfileID))

						profile, profileErr := profileService.GetByID(
							ctx.Request.Context(),
							locale,
							*user.IndividualProfileID,
						)
						if profileErr != nil {
							logger.ErrorContext(ctx.Request.Context(), "Profile fetch error",
								slog.String("error", profileErr.Error()),
								slog.String("profile_id", *user.IndividualProfileID))
						}

						if profile != nil {
							response["selected_profile"] = map[string]any{
								"id":                  profile.ID,
								"slug":                profile.Slug,
								"kind":                profile.Kind,
								"title":               profile.Title,
								"description":         profile.Description,
								"profile_picture_uri": profile.ProfilePictureURI,
								"custom_domain":       profile.CustomDomain,
								"pronouns":            profile.Pronouns,
								"properties":          profile.Properties,
								"created_at":          profile.CreatedAt,
								"updated_at":          profile.UpdatedAt,
							}
						}
					}
				}
			}

			// Wrap response in the expected format for the frontend fetcher
			wrappedResponse := map[string]any{
				"data":  response,
				"error": nil,
			}

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("Get Current Session").
		HasDescription("Returns the current authenticated session with user and profile data.").
		HasResponse(http.StatusOK)
}
