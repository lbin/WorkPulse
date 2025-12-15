package config

import "github.com/spf13/viper"

type Config struct {
	Env       string
	Addr      string
	DBDSN     string
	JWTSecret string
}

func Load() (*Config, error) {
	viper.SetDefault("ENV", "dev")
	viper.SetDefault("ADDR", ":8080")
	viper.AutomaticEnv()

	cfg := &Config{
		Env:       viper.GetString("ENV"),
		Addr:      viper.GetString("ADDR"),
		DBDSN:     viper.GetString("DB_DSN"),
		JWTSecret: viper.GetString("JWT_SECRET"),
	}
	return cfg, nil
}
