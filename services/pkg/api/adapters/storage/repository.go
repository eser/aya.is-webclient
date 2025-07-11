package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/eser/aya.is-services/pkg/ajan/connfx"
	"github.com/eser/aya.is-services/pkg/ajan/logfx"
	"github.com/eser/aya.is-services/pkg/lib/caching"
	"github.com/eser/aya.is-services/pkg/lib/cursors"
	"github.com/pressly/goose/v3"
)

const (
	DefaultCacheTTL = 1 * time.Hour
)

var (
	ErrDatasourceNotFound = errors.New("datasource not found")
	ErrMigrationFailed    = errors.New("migration failed")
	ErrSeedingFailed      = errors.New("seeding failed")
	ErrInvalidFilePath    = errors.New("invalid file path")
)

type Repository struct {
	db       *sql.DB
	queries  *Queries
	cache    *caching.Cache
	logger   *logfx.Logger
	cacheTTL time.Duration
}

func NewRepositoryFromDefault(
	logger *logfx.Logger,
	dataRegistry *connfx.Registry,
) (*Repository, error) {
	return NewRepositoryFromNamed(logger, dataRegistry, connfx.DefaultConnection)
}

func NewRepositoryFromNamed(
	logger *logfx.Logger,
	dataRegistry *connfx.Registry,
	name string,
) (*Repository, error) {
	sqlDB, err := connfx.GetTypedConnection[*sql.DB](dataRegistry, name)
	if err != nil {
		return nil, err
	}

	repository := &Repository{ //nolint:exhaustruct
		db:       sqlDB,
		queries:  &Queries{db: sqlDB},
		cacheTTL: DefaultCacheTTL,
		logger:   logger,
	}

	repository.cache = caching.NewCache(
		func(ctx context.Context, key string, target any) (bool, error) {
			cachedMessage, err := repository.CacheGetSince(
				ctx,
				key,
				time.Now().Add(-1*repository.cacheTTL),
			)
			if err != nil {
				return false, err
			}

			if cachedMessage == nil {
				return false, nil
			}

			unmarshallErr := json.Unmarshal(*cachedMessage, target)
			if unmarshallErr != nil {
				return false, unmarshallErr //nolint:wrapcheck
			}

			return true, nil
		},
		func(ctx context.Context, key string, value any) error {
			message, marshallErr := json.Marshal(value)
			if marshallErr != nil {
				return marshallErr //nolint:wrapcheck
			}

			err := repository.CacheSet(ctx, key, message)
			if err != nil {
				return err
			}

			return nil
		},
	)

	return repository, nil
}

func (r *Repository) RunMigrations(ctx context.Context, migrationsDir string) error {
	r.logger.InfoContext(
		ctx,
		"Running database migrations",
		slog.String("dir", migrationsDir),
	)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("%w: failed to set goose dialect: %w", ErrMigrationFailed, err)
	}

	if err := goose.Up(r.db, migrationsDir); err != nil {
		return fmt.Errorf("%w: failed to run migrations: %w", ErrMigrationFailed, err)
	}

	version, err := goose.GetDBVersion(r.db)
	if err != nil {
		return fmt.Errorf("%w: failed to get DB version: %w", ErrMigrationFailed, err)
	}

	r.logger.InfoContext(ctx, "Migrations completed", slog.Int64("version", version))

	return nil
}

func (r *Repository) SeedData(ctx context.Context, seedFilePath string) error {
	r.logger.InfoContext(ctx, "Checking if seed data is needed")

	cursor := &cursors.Cursor{ //nolint:exhaustruct
		Limit: 1, // Just need to check if any profiles exist
	}

	profilesResult, err := r.ListProfiles(ctx, "en", cursor)
	if err != nil {
		return fmt.Errorf("%w: failed to check existing profiles: %w", ErrSeedingFailed, err)
	}

	if len(profilesResult.Data) > 0 {
		r.logger.InfoContext(
			ctx,
			"Profiles already exist, skipping seed data",
			slog.Int("count", len(profilesResult.Data)),
		)

		return nil
	}

	r.logger.InfoContext(ctx, "No profiles found, seeding data", slog.String("file", seedFilePath))

	seedFilePathClean := filepath.Clean(seedFilePath)

	seedSQL, err := os.ReadFile(seedFilePathClean)
	if err != nil {
		return fmt.Errorf("%w: failed to read seed file: %w", ErrSeedingFailed, err)
	}

	if _, err := r.db.ExecContext(ctx, string(seedSQL)); err != nil {
		return fmt.Errorf("%w: failed to execute seed data: %w", ErrSeedingFailed, err)
	}

	r.logger.InfoContext(ctx, "Seed data executed successfully")

	return nil
}
