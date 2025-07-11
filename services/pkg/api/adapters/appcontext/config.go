package appcontext

import (
	"github.com/eser/aya.is-services/pkg/ajan"
	"github.com/eser/aya.is-services/pkg/api/adapters/arcade"
)

type Data struct {
	MigrationsPath string `conf:"MIGRATIONS_PATH" default:"etc/data/default/migrations"`
	SeedFilePath   string `conf:"SEED_FILE_PATH"  default:"etc/data/default/seed/seed.sql"`
}

type FeatureFlags struct {
	Dummy bool `conf:"DUMMY" default:"false"` // dummy feature flag
}

type Externals struct {
	Arcade arcade.Config `conf:"ARCADE"`
}

type AppConfig struct {
	Data      Data      `conf:"DATA"`
	Externals Externals `conf:"EXTERNALS"`
	ajan.BaseConfig

	Features FeatureFlags `conf:"FEATURES"`
}
