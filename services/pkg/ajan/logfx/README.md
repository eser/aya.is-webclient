# ajan/logfx

## Overview

**logfx** package is a configurable logging solution that leverages the
`log/slog` of the standard library for structured logging. It includes
pretty-printing options and **centralized OTLP connection management** through
`connfx` for log export to modern observability platforms.
The package supports OpenTelemetry-compatible severity levels and provides
extensive test coverage to ensure reliability and correctness.

### Key Features

- 🎯 **Extended Log Levels** - OpenTelemetry-compatible levels while using standard `log/slog` under the hood
- 🔄 **Automatic Correlation IDs** - Request tracing across your entire application
- 🌐 **Centralized OTLP Integration** - Uses `connfx` registry for shared OTLP connections
- 📊 **Structured Logging** - JSON output for production, pretty printing for development
- 🎨 **Pretty Printing** - Colored output for development
- ⚡ **Performance Optimized** - Asynchronous exports, structured logging

## 🚀 **Extended Log Levels**

**The Problem**: Go's standard `log/slog` package provides only 4 log levels (Debug, Info, Warn, Error), which is insufficient for modern observability and OpenTelemetry compatibility.

**The Solution**: logfx extends the standard library to provide **7 OpenTelemetry-compatible log levels** while maintaining full compatibility with `log/slog`:

```go
// Standard Go slog levels (limited)
slog.LevelDebug  // -4
slog.LevelInfo   //  0
slog.LevelWarn   //  4
slog.LevelError  //  8

// logfx extended levels (OpenTelemetry compatible)
logfx.LevelTrace // -8  ← Additional
logfx.LevelDebug // -4
logfx.LevelInfo  //  0
logfx.LevelWarn  //  4
logfx.LevelError //  8
logfx.LevelFatal // 12  ← Additional
logfx.LevelPanic // 16  ← Additional
```

### Why This Matters

1. **OpenTelemetry Compatibility** - Maps perfectly to OpenTelemetry log severity levels
2. **Better Observability** - More granular log levels for better debugging and monitoring
3. **Standard Library Foundation** - Built on `log/slog`, not a replacement
4. **Zero Breaking Changes** - Existing slog code works unchanged
5. **Proper Severity Mapping** - Correct OTLP export with appropriate severity levels

### Extended Level Usage

```go
import "github.com/eser/aya.is-services/pkg/ajan/logfx"

logger := logfx.NewLogger(
    logfx.WithLevel(logfx.LevelTrace), // Now supports all 7 levels
)

// Use all OpenTelemetry-compatible levels
logger.Trace("Detailed debugging info")           // Most verbose
logger.Debug("Debug information")                 // Development debugging
logger.Info("General information")                // Standard info
logger.Warn("Warning message")                    // Potential issues
logger.Error("Error occurred")                    // Errors that don't stop execution
logger.Fatal("Fatal error")                       // Critical errors
logger.Panic("Panic condition")                   // Most severe
```

**Colored Output** (development mode):
```bash
23:45:12.123 TRACE Detailed debugging info
23:45:12.124 DEBUG Debug information
23:45:12.125 INFO General information
23:45:12.126 WARN Warning message
23:45:12.127 ERROR Error occurred
23:45:12.128 FATAL Fatal error
23:45:12.129 PANIC Panic condition
```

**Structured Output** (production mode):
```json
{"time":"2024-01-15T23:45:12.123Z","level":"TRACE","msg":"Detailed debugging info"}
{"time":"2024-01-15T23:45:12.124Z","level":"DEBUG","msg":"Debug information"}
{"time":"2024-01-15T23:45:12.125Z","level":"INFO","msg":"General information"}
{"time":"2024-01-15T23:45:12.126Z","level":"WARN","msg":"Warning message"}
{"time":"2024-01-15T23:45:12.127Z","level":"ERROR","msg":"Error occurred"}
{"time":"2024-01-15T23:45:12.128Z","level":"FATAL","msg":"Fatal error"}
{"time":"2024-01-15T23:45:12.129Z","level":"PANIC","msg":"Panic condition"}
```

**OpenTelemetry Export** (automatic severity mapping):
```json
{
  "logRecords": [
    {"body": {"stringValue": "Detailed debugging info"}, "severityNumber": 1, "severityText": "TRACE"},
    {"body": {"stringValue": "Debug information"}, "severityNumber": 5, "severityText": "DEBUG"},
    {"body": {"stringValue": "General information"}, "severityNumber": 9, "severityText": "INFO"},
    {"body": {"stringValue": "Warning message"}, "severityNumber": 13, "severityText": "WARN"},
    {"body": {"stringValue": "Error occurred"}, "severityNumber": 17, "severityText": "ERROR"},
    {"body": {"stringValue": "Fatal error"}, "severityNumber": 21, "severityText": "FATAL"},
    {"body": {"stringValue": "Panic condition"}, "severityNumber": 24, "severityText": "PANIC"}
  ]
}
```

## Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "log/slog"
    "os"

    "github.com/eser/aya.is-services/pkg/ajan/connfx"
    "github.com/eser/aya.is-services/pkg/ajan/logfx"
)

func main() {
    ctx := context.Background()

    // Create logger first
    logger := logfx.NewLogger()

    // Create connection registry and configure OTLP connection
    registry := connfx.NewRegistryWithDefaults(logger)

    // Configure OTLP connection once, use everywhere
    otlpConfig := &connfx.ConfigTarget{
        Protocol: "otlp",
        DSN:      "otel-collector:4318",
        Properties: map[string]any{
            "service_name":    "my-service",
            "service_version": "1.0.0",
            "insecure":        true,
        },
    }

    // Add OTLP connection to registry
    _, err := registry.AddConnection(ctx, "otel", otlpConfig)
    if err != nil {
        panic(err)
    }

    // Create logger with connection registry (enables OTLP export)
    logger = logfx.NewLogger(
        logfx.WithConfig(&logfx.Config{
            Level:              "TRACE",
            OTLPConnectionName: "otel", // Reference the connection
        }),
        logfx.WithRegistry(registry), // Pass the registry
    )

    // Use structured logging with extended levels
    logger.Info("Application started",
        slog.String("service", "my-service"),
        slog.String("version", "1.0.0"),
    )

    // Extended levels for better observability
    logger.Trace("Connection pool initialized")     // Very detailed
    logger.Debug("Processing user request")         // Debug info
    logger.Warn("High memory usage detected")       // Warnings
    logger.Fatal("Database connection failed")      // Critical errors
}
```

### Complete Observability Stack Integration

```go
package main

import (
    "context"
    "log/slog"
    "net/http"
    "os"

    "github.com/eser/aya.is-services/pkg/ajan/connfx"
    "github.com/eser/aya.is-services/pkg/ajan/httpfx"
    "github.com/eser/aya.is-services/pkg/ajan/httpfx/middlewares"
    "github.com/eser/aya.is-services/pkg/ajan/logfx"
)

func main() {
    ctx := context.Background()

    // Step 1: Create connection registry with OTLP connection
    logger := logfx.NewLogger()
    registry := connfx.NewRegistryWithDefaults(logger)

    // Configure shared OTLP connection for all observability signals
    _, err := registry.AddConnection(ctx, "otel", &connfx.ConfigTarget{
        Protocol: "otlp",
        DSN:      "otel-collector:4318",
        Properties: map[string]any{
            "service_name":     "my-api",
            "service_version":  "1.0.0",
            "insecure":         true,
            "export_interval":  "15s",
            "batch_timeout":    "5s",
        },
    })
    if err != nil {
        panic(err)
    }

    // Step 2: Create observability stack using shared connection

    // Logging with extended levels
    logger = logfx.NewLogger(
        logfx.WithConfig(&logfx.Config{
            Level:              "TRACE",
            OTLPConnectionName: "otel",
        }),
        logfx.WithRegistry(registry),
    )

    // Metrics
    metricsProvider := metricsfx.NewMetricsProvider(&metricsfx.Config{
        ServiceName:        "my-api",
        ServiceVersion:     "1.0.0",
        OTLPConnectionName: "otel",
        ExportInterval:     15 * time.Second,
    }, registry)
    _ = metricsProvider.Init()

    // Tracing
    tracesProvider := tracesfx.NewTracesProvider(&tracesfx.Config{
        ServiceName:        "my-api",
        ServiceVersion:     "1.0.0",
        OTLPConnectionName: "otel",
        SampleRatio:        1.0,
    }, registry)
    _ = tracesProvider.Init()

    // Step 3: Setup HTTP service with observability middleware
    router := httpfx.NewRouter("/api")

    // Add trace middleware for automatic request tracking
    router.Use(middlewares.TraceMiddleware())
    router.Use(middlewares.LoggingMiddleware(logger))

    // Add metrics middleware
    httpMetrics, _ := metricsfx.NewHTTPMetrics(metricsProvider, "my-api", "1.0.0")
    router.Use(middlewares.MetricsMiddleware(httpMetrics))

    router.Route("GET /users/{id}", func(ctx *httpfx.Context) httpfx.Result {
        // All logs automatically include correlation_id and trace information
        logger.TraceContext(ctx.Request.Context(), "Starting user lookup")
        logger.InfoContext(ctx.Request.Context(), "Processing user request",
            slog.String("user_id", "123"),
        )

        return ctx.Results.JSON(map[string]string{"status": "success"})
    })

    http.ListenAndServe(":8080", router.GetMux())
}
```

**Log Output with Complete Correlation:**
```json
{"time":"2024-01-15T10:30:00Z","level":"INFO","msg":"HTTP request started","method":"GET","path":"/api/users/123","correlation_id":"abc-123-def"}
{"time":"2024-01-15T10:30:00Z","level":"TRACE","msg":"Starting user lookup","correlation_id":"abc-123-def","trace_id":"4bf92f3577b34da6","span_id":"00f067aa0bb902b7"}
{"time":"2024-01-15T10:30:00Z","level":"INFO","msg":"Processing user request","user_id":"123","correlation_id":"abc-123-def","trace_id":"4bf92f3577b34da6","span_id":"00f067aa0bb902b7"}
{"time":"2024-01-15T10:30:00Z","level":"INFO","msg":"HTTP request completed","method":"GET","status_code":200,"correlation_id":"abc-123-def","trace_id":"4bf92f3577b34da6"}
```

## Configuration

```go
type Config struct {
	Level  string `conf:"level"  default:"INFO"`    // Supports: TRACE, DEBUG, INFO, WARN, ERROR, FATAL, PANIC

	// Connection-based OTLP configuration (replaces direct endpoint config)
	OTLPConnectionName string `conf:"otlp_connection_name" default:""`

	DefaultLogger bool `conf:"default"    default:"false"`
	PrettyMode    bool `conf:"pretty"     default:"true"`
	AddSource     bool `conf:"add_source" default:"false"`
}
```

## Centralized Connection Management

### Why Use connfx for OTLP Connections?

The `logfx` architecture centralizes OTLP connection management through `connfx`, providing significant advantages:

```go
// Single OTLP connection shared across all packages
registry.AddConnection(ctx, "otel", &connfx.ConfigTarget{Protocol: "otlp", DSN: "otel-collector:4318"})

// All packages reference the same connection
logger := logfx.NewLogger(logfx.WithOTLP("otel"), logfx.WithRegistry(registry))
metrics := metricsfx.NewMetricsProvider(&metricsfx.Config{OTLPConnectionName: "otel"}, registry)
traces := tracesfx.NewTracesProvider(&tracesfx.Config{OTLPConnectionName: "otel"}, registry)
```

**Benefits:**
- 🔧 **Single Configuration Point** - Configure OTLP once, use everywhere
- 🔄 **Shared Connections** - Efficient resource usage and connection pooling
- 🎛️ **Centralized Management** - Health checks, monitoring, and lifecycle management
- 🔗 **Consistent Attribution** - Service name and version automatically applied to all signals
- 💰 **Cost Optimization** - Single connection reduces overhead
- 🛡️ **Error Handling** - Graceful fallbacks when connections are unavailable

### OTLP Connection Configuration

```go
// Configure OTLP connection with full options
otlpConfig := &connfx.ConfigTarget{
    Protocol: "otlp",
    DSN:      "otel-collector:4318",
    Properties: map[string]any{
        // Service identification (applied to all signals automatically)
        "service_name":    "my-service",
        "service_version": "1.0.0",

        // Connection settings
        "insecure":        true,                    // Use HTTP instead of HTTPS

        // Export configuration
        "export_interval": 30 * time.Second,       // Metrics export interval
        "batch_timeout":   5 * time.Second,        // Trace batch timeout
        "batch_size":      512,                    // Trace batch size
        "sample_ratio":    1.0,                    // Trace sampling ratio
    },
}

_, err := registry.AddConnection(ctx, "otel", otlpConfig)
```

### Environment-Based Configuration

```bash
# Connection configuration via environment
CONN_TARGETS_OTEL_PROTOCOL=otlp
CONN_TARGETS_OTEL_DSN=otel-collector:4318
CONN_TARGETS_OTEL_PROPERTIES_SERVICE_NAME=my-service
CONN_TARGETS_OTEL_PROPERTIES_SERVICE_VERSION=1.0.0
CONN_TARGETS_OTEL_PROPERTIES_INSECURE=true

# Package configuration references the connection
LOG_OTLP_CONNECTION_NAME=otel
METRICS_OTLP_CONNECTION_NAME=otel
TRACES_OTLP_CONNECTION_NAME=otel
```

### Multiple OTLP Endpoints

```go
// Different endpoints for different environments
_, err := registry.AddConnection(ctx, "otel-dev", &connfx.ConfigTarget{
    Protocol: "otlp",
    URL:      "http://dev-collector:4318",
    Properties: map[string]any{"service_name": "my-service-dev"},
})

_, err = registry.AddConnection(ctx, "otel-prod", &connfx.ConfigTarget{
    Protocol: "otlp",
    URL:      "https://prod-collector:4317",
    TLS:      true,
    Properties: map[string]any{
        "service_name": "my-service",
        "insecure":     false,
    },
})

// Use different connections in different packages
devLogger := logfx.NewLogger(logfx.WithOTLP("otel-dev"), logfx.WithRegistry(registry))
prodMetrics := metricsfx.NewMetricsProvider(&metricsfx.Config{OTLPConnectionName: "otel-prod"}, registry)
```

## Correlation IDs

### Automatic HTTP Correlation

When using with `httpfx`, correlation IDs are automatically:

- ✅ **Extracted** from `X-Correlation-ID` headers
- ✅ **Generated** if missing
- ✅ **Propagated** through Go context
- ✅ **Added** to all log entries
- ✅ **Included** in response headers

### Manual Correlation Access

```go
import "github.com/eser/aya.is-services/pkg/ajan/correlation"

func MyHandler(ctx *httpfx.Context) httpfx.Result {
    correlationID := correlation.GetCorrelationIDFromContext(ctx.Request.Context())

    // Use in external service calls
    externalReq.Header.Set("X-Correlation-ID", correlationID)

    return ctx.Results.JSON(map[string]string{
        "correlation_id": correlationID,
    })
}
```

## Advanced Usage

### Migration from Direct OTLP Configuration

**Old Code:**
```go
// Before: Direct OTLP configuration
logger := logfx.NewLogger(
    logfx.WithOTLP("otel-collector:4318", true),
)
```

**New Code:**
```go
// After: Connection-based configuration
registry := connfx.NewRegistryWithDefaults(logger)
_, err := registry.AddConnection(ctx, "otel", &connfx.ConfigTarget{
    Protocol: "otlp",
    DSN:      "otel-collector:4318",
    Properties: map[string]any{"insecure": true},
})

logger := logfx.NewLogger(
    logfx.WithOTLP("otel"),
    logfx.WithRegistry(registry),
)
```

### Level Configuration Examples

```go
// Development - verbose logging with all levels
devConfig := &logfx.Config{
    Level:              "TRACE",    // Most verbose - see everything
    PrettyMode:         true,
    AddSource:          true,
    OTLPConnectionName: "otel-dev",
}

// Production - structured output with appropriate level
prodConfig := &logfx.Config{
    Level:              "INFO",     // Production appropriate
    PrettyMode:         false,
    OTLPConnectionName: "otel-prod",
}

// Debug production issues - temporary verbose logging
debugConfig := &logfx.Config{
    Level:              "DEBUG",    // More detail for troubleshooting
    PrettyMode:         false,
    OTLPConnectionName: "otel-debug",
}
```

### Standard Library Compatibility

```go
// logfx extends slog.Level, so standard slog works unchanged
import "log/slog"

// This works exactly as before
slog.Info("Standard slog message")
slog.Debug("Debug with standard slog")

// But you can also use extended levels through logfx
logger.Trace("Extended trace level")    // Not available in standard slog
logger.Fatal("Extended fatal level")    // Not available in standard slog
logger.Panic("Extended panic level")    // Not available in standard slog
```

## Error Handling

The logger handles export failures gracefully:

```go
// Logger continues working even if OTLP connection fails
registry := connfx.NewRegistryWithDefaults(logger)

// If connection fails, logger falls back to local output only
logger := logfx.NewLogger(
    logfx.WithWriter(os.Stdout),
    logfx.WithOTLP("nonexistent-connection"),
    logfx.WithRegistry(registry),
)

// Logs always go to the primary writer (stdout/file)
// Connection failures are handled gracefully without affecting your app
logger.Info("This will always work, with or without OTLP")
```

## API Reference

### Logger Creation

#### NewLogger (Options Pattern)

```go
func NewLogger(options ...NewLoggerOption) *Logger
```

Create a logger using the flexible options pattern:

```go
// Basic logger with default configuration
logger := logfx.NewLogger()

// Logger with connection registry for OTLP export
logger := logfx.NewLogger(
    logfx.WithWriter(os.Stdout),
    logfx.WithConfig(&logfx.Config{
        Level:              "INFO",
        PrettyMode:         false,
        OTLPConnectionName: "otel",
    }),
    logfx.WithRegistry(registry),
)

// Logger with individual options
logger := logfx.NewLogger(
    logfx.WithLevel(slog.LevelDebug),
    logfx.WithPrettyMode(true),
    logfx.WithAddSource(true),
    logfx.WithOTLP("otel"),
    logfx.WithRegistry(registry),
    logfx.WithDefaultLogger(), // Set as default logger
)
```

#### Available Options

```go
// Configuration options
WithConfig(config *Config)                    // Full configuration
WithLevel(level slog.Level)                   // Set log level
WithPrettyMode(pretty bool)                   // Enable/disable pretty printing
WithAddSource(addSource bool)                 // Include source code location
WithDefaultLogger()                           // Set as default logger

// Output options
WithWriter(writer io.Writer)                  // Set output writer
WithFromSlog(slog *slog.Logger)              // Wrap existing slog.Logger

// Connection-based OTLP export (NEW)
WithOTLP(connectionName string)               // Reference OTLP connection by name
WithRegistry(registry ConnectionRegistry)     // Provide connection registry
```

#### Migration Guide

**Before:**
```go
// Old direct endpoint configuration
logger := logfx.NewLogger(
    logfx.WithOTLP("http://collector:4318", true),
)
```

**After:**
```go
// New connection-based configuration
registry := connfx.NewRegistryWithDefaults(logger)
registry.AddConnection(ctx, "otel", &connfx.ConfigTarget{
    Protocol: "otlp",
    URL:      "http://collector:4318",
    Properties: map[string]any{"insecure": true},
})

logger := logfx.NewLogger(
    logfx.WithOTLP("otel"),
    logfx.WithRegistry(registry),
)
```

## Best Practices

1. **Use Centralized Connections**: Configure OTLP connections once in `connfx`, use everywhere
2. **Connection Health Monitoring**: Use `registry.HealthCheck(ctx)` to monitor OTLP connection health
3. **Graceful Degradation**: Logger works with or without OTLP connections
4. **Correlation IDs**: Use with `httpfx` middleware for automatic request correlation
5. **Environment-Based Config**: Use environment variables for connection configuration
6. **Resource Attribution**: Set service name/version in connection properties for proper attribution
7. **Connection Lifecycle**: Use `registry.Close(ctx)` during shutdown to properly cleanup connections
8. **Multiple Environments**: Use different connection names for dev/staging/prod environments

## Architecture Benefits

- **Unified Configuration** - Single place to configure OTLP connections for all observability signals
- **Shared Resources** - Efficient connection pooling and resource utilization
- **Consistent Attribution** - Service information automatically applied to all logs
- **Health Monitoring** - Built-in connection health checks and monitoring
- **Graceful Fallbacks** - Continue working even when OTLP connections fail
- **Environment Flexibility** - Easy switching between different collectors/environments
- **Import Cycle Prevention** - Bridge pattern avoids circular dependencies
- **Thread Safety** - All connection operations are thread-safe
