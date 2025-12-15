package service

import (
	"fmt"

	"workpulse/internal/config"
)

type ClientConfig struct {
	APIVersion   string          `json:"api_version"`
	FeatureFlags map[string]bool `json:"feature_flags"`
	Docs         DocsConfig      `json:"docs"`
}

type DocsConfig struct {
	OpenAPI    string `json:"openapi"`
	GraphQLSDL string `json:"graphql_sdl"`
}

type ConfigService struct {
	cfg *config.Config
}

func NewConfigService(cfg *config.Config) *ConfigService {
	return &ConfigService{cfg: cfg}
}

func (s *ConfigService) ClientConfig(basePath string) ClientConfig {
	docsBase := fmt.Sprintf("%s/docs", basePath)
	featureFlags := map[string]bool{
		"analytics":   true,
		"apiDocs":     true,
		"graphqlDocs": true,
	}
	for k, v := range s.cfg.FeatureFlags {
		featureFlags[k] = v
	}
	return ClientConfig{
		APIVersion:   s.cfg.APIVersion,
		FeatureFlags: featureFlags,
		Docs: DocsConfig{
			OpenAPI:    fmt.Sprintf("%s/openapi.json", docsBase),
			GraphQLSDL: fmt.Sprintf("%s/graphql.sdl", docsBase),
		},
	}
}
