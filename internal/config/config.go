package config

import (
	"encoding/json"

	"github.com/spf13/viper"
)

type Config struct {
	Env             string
	Addr            string
	DBDSN           string
	JWTSecret       string
	APIVersion      string
	MigrationsTable string
	FeatureFlags    map[string]bool
}

func Load() (*Config, error) {
	viper.SetDefault("ENV", "dev")
	viper.SetDefault("ADDR", ":8080")
	viper.SetDefault("API_VERSION", "v1")
	viper.SetDefault("MIGRATIONS_TABLE", "schema_migrations")
	viper.AutomaticEnv()

	featureFlags := map[string]bool{}
	if rawMap := viper.GetStringMap("FEATURE_FLAGS"); len(rawMap) > 0 {
		for k, v := range rawMap {
			if enabled, ok := v.(bool); ok {
				featureFlags[k] = enabled
			}
		}
	} else if raw := viper.GetString("FEATURE_FLAGS"); raw != "" {
		_ = json.Unmarshal([]byte(raw), &featureFlags)
	}

	cfg := &Config{
		Env:             viper.GetString("ENV"),
		Addr:            viper.GetString("ADDR"),
		DBDSN:           viper.GetString("DB_DSN"),
		JWTSecret:       viper.GetString("JWT_SECRET"),
		APIVersion:      viper.GetString("API_VERSION"),
		MigrationsTable: viper.GetString("MIGRATIONS_TABLE"),
		FeatureFlags:    featureFlags,
	}
	return cfg, nil
}
