package config

import (
	"encoding/json"
	"fmt"
	"strings"

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
	ServiceName     string
	MetricsPath     string
	OTLPEndpoint    string
	OTLPHeaders     map[string]string
}

func Load() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	_ = viper.ReadInConfig()

	viper.SetDefault("ENV", "dev")
	viper.SetDefault("ADDR", ":8080")
	viper.SetDefault("API_VERSION", "v1")
	viper.SetDefault("MIGRATIONS_TABLE", "schema_migrations")
	viper.SetDefault("SERVICE_NAME", "workpulse-api")
	viper.SetDefault("METRICS_PATH", "/metrics")
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

	otlpHeaders := map[string]string{}
	if rawHeaders := viper.GetString("OTEL_EXPORTER_OTLP_HEADERS"); rawHeaders != "" {
		pairs := strings.Split(rawHeaders, ",")
		for _, pair := range pairs {
			items := strings.SplitN(pair, "=", 2)
			if len(items) == 2 {
				otlpHeaders[strings.TrimSpace(items[0])] = strings.TrimSpace(items[1])
			}
		}
	}

	cfg := &Config{
		Env:             viper.GetString("ENV"),
		Addr:            viper.GetString("ADDR"),
		DBDSN:           viper.GetString("DB_DSN"),
		JWTSecret:       viper.GetString("JWT_SECRET"),
		APIVersion:      viper.GetString("API_VERSION"),
		MigrationsTable: viper.GetString("MIGRATIONS_TABLE"),
		FeatureFlags:    featureFlags,
		ServiceName:     viper.GetString("SERVICE_NAME"),
		MetricsPath:     viper.GetString("METRICS_PATH"),
		OTLPEndpoint:    viper.GetString("OTEL_EXPORTER_OTLP_ENDPOINT"),
		OTLPHeaders:     otlpHeaders,
	}

	if cfg.DBDSN == "" {
		return nil, fmt.Errorf("DB_DSN is required (set it in .env or the environment)")
	}

	return cfg, nil
}
