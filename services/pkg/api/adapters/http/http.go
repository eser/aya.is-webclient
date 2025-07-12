package http

import (
	"context"

	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/ajan/httpfx/middlewares"
	"github.com/eser/aya.is/services/pkg/ajan/httpfx/modules/healthcheck"
	"github.com/eser/aya.is/services/pkg/ajan/httpfx/modules/openapi"
	"github.com/eser/aya.is/services/pkg/ajan/httpfx/modules/profiling"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/profiles"
	"github.com/eser/aya.is/services/pkg/api/business/stories"
	"github.com/eser/aya.is/services/pkg/api/business/users"
)

func Run(
	ctx context.Context,
	config *httpfx.Config,
	logger *logfx.Logger,
	profilesService *profiles.Service,
	storiesService *stories.Service,
	usersService *users.Service,
) (func(), error) {
	routes := httpfx.NewRouter("/")
	httpService := httpfx.NewHTTPService(config, routes, logger)

	// http middlewares
	routes.Use(middlewares.ErrorHandlerMiddleware())
	routes.Use(middlewares.ResolveAddressMiddleware())
	routes.Use(middlewares.ResponseTimeMiddleware())
	routes.Use(middlewares.TracingMiddleware(logger)) //nolint:contextcheck
	routes.Use(middlewares.CorsMiddleware())
	routes.Use(middlewares.MetricsMiddleware(httpService.InnerMetrics)) //nolint:contextcheck

	// http modules
	healthcheck.RegisterHTTPRoutes(routes, config)
	openapi.RegisterHTTPRoutes(routes, config)
	profiling.RegisterHTTPRoutes(routes, config)

	// http routes
	RegisterHTTPRoutesForUsers( //nolint:contextcheck
		routes,
		logger,
		usersService,
	)
	RegisterHTTPRoutesForSite( //nolint:contextcheck
		routes,
		logger,
		profilesService,
	)
	RegisterHTTPRoutesForProfiles( //nolint:contextcheck
		routes,
		logger,
		profilesService,
		storiesService,
	)
	RegisterHTTPRoutesForStories( //nolint:contextcheck
		routes,
		logger,
		storiesService,
	)

	// Authenticated routes (wrapped with auth middleware)
	RegisterAuthenticatedRoutes( //nolint:contextcheck
		routes,
		logger,
		profilesService,
		storiesService,
		usersService,
	)

	// run
	return httpService.Start(ctx) //nolint:wrapcheck
}
