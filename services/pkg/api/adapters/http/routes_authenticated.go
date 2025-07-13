package http

import (
	"log/slog"
	"net/http"

	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
	"github.com/eser/aya.is/services/pkg/api/business/profiles"
	"github.com/eser/aya.is/services/pkg/api/business/stories"
	"github.com/eser/aya.is/services/pkg/api/business/users"
)

// RegisterAuthenticatedRoutes registers routes that require authentication.
func RegisterAuthenticatedRoutes( //nolint:funlen
	routes *httpfx.Router,
	logger *logfx.Logger,
	authService *auth.Service,
	userService *users.Service,
	profileService *profiles.Service,
	storyService *stories.Service,
) {
	// Register authenticated route directly with auth middleware
	routes.Route(
		"GET /{locale}/protected/user-info",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			// Get user ID from context (set by auth middleware)
			userID, ok := ctx.Request.Context().Value(ContextKeyUserID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText("User ID not found in context"),
				)
			}

			// Get user data
			user, err := userService.GetByID(ctx.Request.Context(), userID)
			if err != nil {
				return ctx.Results.Error(
					http.StatusNotFound,
					httpfx.WithPlainText("User not found"),
				)
			}

			// Prepare response with user data
			response := map[string]any{
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
				if profileErr == nil && profile != nil {
					response["individual_profile"] = map[string]any{
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

			// Wrap response in the expected format for the frontend fetcher
			wrappedResponse := map[string]any{
				"data":  response,
				"error": nil,
			}

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("Get Current User").
		HasDescription("Returns the current authenticated user with profile data.").
		HasResponse(http.StatusOK)
}
