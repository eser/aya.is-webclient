package http

import (
	"github.com/eser/aya.is/services/pkg/ajan/httpfx"
	"github.com/eser/aya.is/services/pkg/ajan/logfx"
	"github.com/eser/aya.is/services/pkg/api/business/auth"
	"github.com/eser/aya.is/services/pkg/api/business/profiles"
	"github.com/eser/aya.is/services/pkg/api/business/stories"
	"github.com/eser/aya.is/services/pkg/api/business/users"
)

// RegisterAuthenticatedRoutes registers routes that require authentication.
func RegisterAuthenticatedRoutes(
	routes *httpfx.Router,
	logger *logfx.Logger,
	authService *auth.Service,
	userService *users.Service,
	profileService *profiles.Service,
	storyService *stories.Service,
) {
	// Create a sub-router with auth middleware
	authRoutes := routes.Group("/")
	authRoutes.Use(AuthMiddleware(authService, userService))
	// Example: Protected user profile endpoints
	// authRoutes.Route("PUT /{locale}/users/{id}", updateUserHandler)
	// authRoutes.Route("DELETE /{locale}/users/{id}", deleteUserHandler)
	// Example: Protected story creation/editing endpoints
	// authRoutes.Route("POST /{locale}/stories", createStoryHandler)
	// authRoutes.Route("PUT /{locale}/stories/{id}", updateStoryHandler)
	// authRoutes.Route("DELETE /{locale}/stories/{id}", deleteStoryHandler)
	// Add more protected routes here as needed
}
