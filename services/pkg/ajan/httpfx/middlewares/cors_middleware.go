package middlewares

import (
	"strings"

	"github.com/eser/aya.is-services/pkg/ajan/httpfx"
)

// Constants for CORS headers.
const (
	AccessControlAllowOriginHeader      = "Access-Control-Allow-Origin"
	AccessControlAllowCredentialsHeader = "Access-Control-Allow-Credentials"
	AccessControlAllowHeadersHeader     = "Access-Control-Allow-Headers"
	AccessControlAllowMethodsHeader     = "Access-Control-Allow-Methods"
)

// corsConfig holds the configuration for CORS headers.
// It is unexported as it's an internal detail of the CorsMiddleware.
type corsConfig struct {
	allowOrigin      string
	allowHeaders     []string
	allowMethods     []string
	allowCredentials bool
}

// CorsOption is a function type that modifies the corsConfig.
type CorsOption func(*corsConfig)

// WithAllowOrigin sets the Access-Control-Allow-Origin header.
// If not set, defaults to "*".
func WithAllowOrigin(origin string) CorsOption {
	return func(cfg *corsConfig) {
		cfg.allowOrigin = origin
	}
}

// WithAllowCredentials sets the Access-Control-Allow-Credentials header.
func WithAllowCredentials(allow bool) CorsOption {
	return func(cfg *corsConfig) {
		cfg.allowCredentials = allow
	}
}

// WithAllowHeaders sets the Access-Control-Allow-Headers header.
func WithAllowHeaders(headers []string) CorsOption {
	return func(cfg *corsConfig) {
		cfg.allowHeaders = headers
	}
}

// WithAllowMethods sets the Access-Control-Allow-Methods header.
func WithAllowMethods(methods []string) CorsOption {
	return func(cfg *corsConfig) {
		cfg.allowMethods = methods
	}
}

// CorsMiddleware creates a CORS middleware using functional options.
func CorsMiddleware(options ...CorsOption) httpfx.Handler {
	// Start with default configuration
	cfg := &corsConfig{
		allowOrigin:      "*", // Default to allow all origins
		allowCredentials: false,
		allowHeaders:     []string{},
		allowMethods:     []string{},
	}

	// Apply all provided options
	for _, option := range options {
		option(cfg)
	}

	return func(ctx *httpfx.Context) httpfx.Result {
		result := ctx.Next()

		headers := ctx.ResponseWriter.Header()

		headers.Set(AccessControlAllowOriginHeader, cfg.allowOrigin)

		if cfg.allowCredentials {
			headers.Set(AccessControlAllowCredentialsHeader, "true")
		}

		if len(cfg.allowHeaders) > 0 {
			headers.Set(AccessControlAllowHeadersHeader, strings.Join(cfg.allowHeaders, ", "))
		}

		if len(cfg.allowMethods) > 0 {
			headers.Set(AccessControlAllowMethodsHeader, strings.Join(cfg.allowMethods, ", "))
		}

		return result
	}
}
