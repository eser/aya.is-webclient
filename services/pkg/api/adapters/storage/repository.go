package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/eser/aya.is-services/pkg/ajan/connfx"
	"github.com/eser/aya.is-services/pkg/ajan/logfx"
	"github.com/eser/aya.is-services/pkg/lib/caching"
)

const (
	DefaultCacheTTL = 1 * time.Hour
)

var ErrDatasourceNotFound = errors.New("datasource not found")

type Repository struct {
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
