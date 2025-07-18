package http

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
	"github.com/eser/aya.is/services/pkg/api/business/profiles"
	"github.com/eser/aya.is/services/pkg/api/business/stories"
	"github.com/eser/aya.is/services/pkg/api/business/users"
	"github.com/eser/aya.is/services/pkg/lib/cursors"
)

func RegisterHTTPRoutesForProfiles( //nolint:funlen,cyclop,maintidx
	routes *httpfx.Router,
	logger *logfx.Logger,
	authService *auth.Service,
	userService *users.Service,
	profileService *profiles.Service,
	storyService *stories.Service,
) {
	routes.
		Route("GET /{locale}/profiles", func(ctx *httpfx.Context) httpfx.Result {
			// get variables from path
			localeParam := ctx.Request.PathValue("locale")
			cursor := cursors.NewCursorFromRequest(ctx.Request)

			filterKind, filterKindOk := cursor.Filters["kind"]
			if !filterKindOk {
				return ctx.Results.BadRequest(httpfx.WithPlainText("filter_kind is required"))
			}

			kinds := strings.SplitSeq(filterKind, ",")
			for kind := range kinds {
				if kind != "individual" && kind != "organization" && kind != "product" {
					return ctx.Results.BadRequest(httpfx.WithPlainText("filter_kind is invalid"))
				}
			}

			records, err := profileService.List(ctx.Request.Context(), localeParam, cursor)
			if err != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText(err.Error()),
				)
			}

			return ctx.Results.JSON(records)
		}).
		HasSummary("List profiles").
		HasDescription("List profiles.").
		HasResponse(http.StatusOK)

	routes.
		Route("GET /{locale}/profiles/{slug}", func(ctx *httpfx.Context) httpfx.Result {
			// get variables from path
			localeParam := ctx.Request.PathValue("locale")
			slugParam := ctx.Request.PathValue("slug")

			if slugParam == "" {
				return ctx.Results.BadRequest(httpfx.WithPlainText("slug parameter is required"))
			}

			record, err := profileService.GetBySlugEx(
				ctx.Request.Context(),
				localeParam,
				slugParam,
			)
			if err != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText(err.Error()),
				)
			}

			wrappedResponse := cursors.WrapResponseWithCursor(record, nil)

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("Get profile by slug").
		HasDescription("Get profile by slug.").
		HasResponse(http.StatusOK)

	routes.Route(
		"GET /{locale}/profiles/{slug}/_check",
		func(ctx *httpfx.Context) httpfx.Result {
			// get variables from path
			slugParam := ctx.Request.PathValue("slug")

			if slugParam == "" {
				return ctx.Results.BadRequest(httpfx.WithPlainText("slug parameter is required"))
			}

			exists, err := profileService.CheckSlugExists(ctx.Request.Context(), slugParam)
			if err != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText(err.Error()),
				)
			}

			result := map[string]bool{
				"exists": exists,
			}

			wrappedResponse := cursors.WrapResponseWithCursor(result, nil)

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("Check if profile slug exists").
		HasDescription("Check if a profile slug is already taken in the database.").
		HasResponse(http.StatusOK)

	routes.
		Route("GET /{locale}/profiles/{slug}/pages", func(ctx *httpfx.Context) httpfx.Result {
			// get variables from path
			localeParam := ctx.Request.PathValue("locale")
			slugParam := ctx.Request.PathValue("slug")

			records, err := profileService.ListPagesBySlug(
				ctx.Request.Context(),
				localeParam,
				slugParam,
			)
			if err != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText(err.Error()),
				)
			}

			wrappedResponse := cursors.WrapResponseWithCursor(records, nil)

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("List profile pages by profile slug").
		HasDescription("List profile pages by profile slug.").
		HasResponse(http.StatusOK)

	routes.
		Route(
			"GET /{locale}/profiles/{slug}/pages/{pageSlug}",
			func(ctx *httpfx.Context) httpfx.Result {
				// get variables from path
				localeParam := ctx.Request.PathValue("locale")
				slugParam := ctx.Request.PathValue("slug")
				pageSlugParam := ctx.Request.PathValue("pageSlug")

				records, err := profileService.GetPageBySlug(
					ctx.Request.Context(),
					localeParam,
					slugParam,
					pageSlugParam,
				)
				if err != nil {
					return ctx.Results.Error(
						http.StatusInternalServerError,
						httpfx.WithPlainText(err.Error()),
					)
				}

				wrappedResponse := cursors.WrapResponseWithCursor(records, nil)

				return ctx.Results.JSON(wrappedResponse)
			},
		).
		HasSummary("List profile page by profile slug and page slug").
		HasDescription("List profile page by profile slug and page slug.").
		HasResponse(http.StatusOK)

	routes.
		Route("GET /{locale}/profiles/{slug}/links", func(ctx *httpfx.Context) httpfx.Result {
			// get variables from path
			localeParam := ctx.Request.PathValue("locale")
			slugParam := ctx.Request.PathValue("slug")

			records, err := profileService.ListLinksBySlug(
				ctx.Request.Context(),
				localeParam,
				slugParam,
			)
			if err != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText(err.Error()),
				)
			}

			wrappedResponse := cursors.WrapResponseWithCursor(records, nil)

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("List profile links by profile slug").
		HasDescription("List profile links by profile slug.").
		HasResponse(http.StatusOK)

	routes.
		Route("GET /{locale}/profiles/{slug}/stories", func(ctx *httpfx.Context) httpfx.Result {
			// get variables from path
			localeParam := ctx.Request.PathValue("locale")
			slugParam := ctx.Request.PathValue("slug")
			cursor := cursors.NewCursorFromRequest(ctx.Request)

			records, err := storyService.ListByPublicationProfileSlug(
				ctx.Request.Context(),
				localeParam,
				slugParam,
				cursor,
			)
			if err != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText(err.Error()),
				)
			}

			return ctx.Results.JSON(records)
		}).
		HasSummary("List stories published to profile slug").
		HasDescription("List stories published to profile slug.").
		HasResponse(http.StatusOK)

	routes.
		Route(
			"GET /{locale}/profiles/{slug}/stories/{storySlug}",
			func(ctx *httpfx.Context) httpfx.Result {
				// get variables from path
				localeParam := ctx.Request.PathValue("locale")
				// slugParam := ctx.Request.PathValue("slug")
				storySlugParam := ctx.Request.PathValue("storySlug")

				// TODO(@eser) pass profile slug too for getting story by profile slug and story slug
				record, err := storyService.GetBySlug(
					ctx.Request.Context(),
					localeParam,
					storySlugParam,
				)
				if err != nil {
					return ctx.Results.Error(
						http.StatusInternalServerError,
						httpfx.WithPlainText(err.Error()),
					)
				}

				// if record == nil {
				// 	return ctx.Results.NotFound(httpfx.WithPlainText("story not found"))
				// }

				wrappedResponse := cursors.WrapResponseWithCursor(record, nil)

				return ctx.Results.JSON(wrappedResponse)
			},
		).
		HasSummary("Get story by profile slug and story slug").
		HasDescription("Get story by profile slug and story slug.").
		HasResponse(http.StatusOK)

	routes.
		Route(
			"GET /{locale}/profiles/{slug}/contributions",
			func(ctx *httpfx.Context) httpfx.Result {
				// get variables from path
				localeParam := ctx.Request.PathValue("locale")
				slugParam := ctx.Request.PathValue("slug")
				cursor := cursors.NewCursorFromRequest(ctx.Request)

				records, err := profileService.ListProfileContributionsBySlug(
					ctx.Request.Context(),
					localeParam,
					slugParam,
					cursor,
				)
				if err != nil {
					return ctx.Results.Error(
						http.StatusInternalServerError,
						httpfx.WithPlainText(err.Error()),
					)
				}

				return ctx.Results.JSON(records)
			},
		).
		HasSummary("List profile contributions by profile slug").
		HasDescription("List profile contributions by profile slug.").
		HasResponse(http.StatusOK)

	routes.
		Route(
			"GET /{locale}/profiles/{slug}/members",
			func(ctx *httpfx.Context) httpfx.Result {
				// get variables from path
				localeParam := ctx.Request.PathValue("locale")
				slugParam := ctx.Request.PathValue("slug")
				cursor := cursors.NewCursorFromRequest(ctx.Request)

				records, err := profileService.ListProfileMembersBySlug(
					ctx.Request.Context(),
					localeParam,
					slugParam,
					cursor,
				)
				if err != nil {
					return ctx.Results.Error(
						http.StatusInternalServerError,
						httpfx.WithPlainText(err.Error()),
					)
				}

				return ctx.Results.JSON(records)
			},
		).
		HasSummary("List profile members by profile slug").
		HasDescription("List profile members by profile slug.").
		HasResponse(http.StatusOK)

	// Register profile creation route (protected, requires authentication)
	routes.Route(
		"POST /{locale}/profiles/_create",
		AuthMiddleware(authService, userService),
		func(ctx *httpfx.Context) httpfx.Result {
			// Get session ID from context (set by auth middleware)
			sessionID, ok := ctx.Request.Context().Value(ContextKeySessionID).(string)
			if !ok {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText("Session ID not found in context"),
				)
			}

			// Get locale from path
			locale := ctx.Request.PathValue("locale")

			// Parse request body
			var requestBody struct {
				Kind        string `json:"kind"`
				Slug        string `json:"slug"`
				Title       string `json:"title"`
				Description string `json:"description"`
			}

			if err := ctx.ParseJSONBody(&requestBody); err != nil {
				return ctx.Results.BadRequest(httpfx.WithPlainText("Invalid request body"))
			}

			// Validate required fields
			if requestBody.Kind == "" || requestBody.Slug == "" || requestBody.Title == "" ||
				requestBody.Description == "" {
				return ctx.Results.BadRequest(
					httpfx.WithPlainText(
						"All fields (kind, slug, title, description) are required",
					),
				)
			}

			// Validate kind
			if requestBody.Kind != "individual" && requestBody.Kind != "organization" &&
				requestBody.Kind != "product" {
				return ctx.Results.BadRequest(httpfx.WithPlainText("Invalid profile kind"))
			}

			session, sessionErr := userService.GetSessionByID(ctx.Request.Context(), sessionID)
			if sessionErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText("Failed to get session information"),
				)
			}

			// Get current user to check individual profile restriction
			user, userErr := userService.GetByID(ctx.Request.Context(), *session.LoggedInUserID)
			if userErr != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText("Failed to get user information"),
				)
			}

			// Check individual profile restriction
			if requestBody.Kind == "individual" && user.IndividualProfileID != nil {
				return ctx.Results.BadRequest(
					httpfx.WithPlainText("User already has an individual profile"),
				)
			}

			// Check if slug already exists
			exists, err := profileService.CheckSlugExists(ctx.Request.Context(), requestBody.Slug)
			if err != nil {
				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText("Failed to check slug availability"),
				)
			}

			if exists {
				return ctx.Results.BadRequest(httpfx.WithPlainText("Slug is already taken"))
			}

			// Create the profile
			profile, err := profileService.Create(
				ctx.Request.Context(),
				locale,
				requestBody.Slug,
				requestBody.Kind,
				requestBody.Title,
				requestBody.Description,
				nil, // customDomain
				nil, // profilePictureURI
				nil, // pronouns
				nil, // properties
			)
			if err != nil {
				logger.ErrorContext(ctx.Request.Context(), "Profile creation failed",
					slog.String("error", err.Error()),
					slog.String("session_id", sessionID),
					slog.String("user_id", user.ID),
					slog.String("slug", requestBody.Slug))

				return ctx.Results.Error(
					http.StatusInternalServerError,
					httpfx.WithPlainText("Failed to create profile"),
				)
			}

			// If creating an individual profile, set it as the user's individual profile
			if requestBody.Kind == "individual" {
				err = userService.SetIndividualProfileID(ctx.Request.Context(), user.ID, profile.ID)
				if err != nil {
					logger.ErrorContext(
						ctx.Request.Context(),
						"Failed to set user individual profile ID",
						slog.String("error", err.Error()),
						slog.String("session_id", sessionID),
						slog.String("user_id", user.ID),
						slog.String("profile_id", profile.ID),
					)

					// Note: We don't return an error here because the profile was created successfully
					// The user can still access their profile, they just need to manually link it
					logger.WarnContext(
						ctx.Request.Context(),
						"Profile created but failed to link to user",
						slog.String("session_id", sessionID),
						slog.String("user_id", user.ID),
						slog.String("profile_id", profile.ID),
					)
				}
			}

			// Prepare response
			response := map[string]any{
				"id":          profile.ID,
				"slug":        profile.Slug,
				"kind":        profile.Kind,
				"title":       profile.Title,
				"description": profile.Description,
				"created_at":  profile.CreatedAt,
			}

			// Wrap response in the expected format for the frontend fetcher
			wrappedResponse := map[string]any{
				"data":  response,
				"error": nil,
			}

			return ctx.Results.JSON(wrappedResponse)
		}).
		HasSummary("Create Profile").
		HasDescription("Create a new profile for the authenticated user.").
		HasResponse(http.StatusOK)
}
