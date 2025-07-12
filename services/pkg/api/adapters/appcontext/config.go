package appcontext

import (
	"github.com/eser/aya.is-services/pkg/ajan"
	"github.com/eser/aya.is-services/pkg/api/adapters/arcade"
	"github.com/eser/aya.is-services/pkg/api/business/users"
)

type DataConfig struct {
	MigrationsPath string `conf:"migration_path" default:"etc/data/default/migrations"`
	SeedFilePath   string `conf:"seed_file_path" default:"etc/data/default/seed/seed.sql"`
}

type FeatureFlags struct {
	Dummy bool `conf:"dummy" default:"false"` // dummy feature flag
}

type ExternalsConfig struct {
	Arcade arcade.Config `conf:"arcade"`
}

type AppConfig struct {
	Data      DataConfig       `conf:"data"`
	Auth      users.AuthConfig `conf:"auth"`
	Externals ExternalsConfig  `conf:"externals"`
	ajan.BaseConfig

	Features FeatureFlags `conf:"features"`
}
