// Package config provides configuration management for flowState-cli.
//
// Phase 1: Core Infrastructure
// - Loads configuration from ~/.config/flowState/config.yaml
// - Creates data directory if it doesn't exist
// - Provides sensible defaults for all configuration options
// - Stores: data directory, database path, Qdrant URL, model path
//
// Configuration Fields:
//   - DataDir: Base directory for all application data (~/.config/flowState)
//   - DbPath: SQLite database file path
//   - QdrantUrl: Vector database URL for semantic search
//   - ModelPath: Path to store embedding models
//   - EmbeddingsEnabled: Toggle semantic search features
//
// Usage:
//
//	cfg, err := config.Load()
//	if err != nil { ... }
//	dataDir := cfg.DataDir
package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	DataDir           string `mapstructure:"data_dir"`
	DbPath            string `mapstructure:"db_path"`
	QdrantUrl         string `mapstructure:"qdrant_url"`
	ModelPath         string `mapstructure:"model_path"`
	EmbeddingsEnabled bool   `mapstructure:"embeddings_enabled"`
}

var cfg *Config

// Load initializes configuration with sensible defaults.
// Phase 1: Creates ~/.config/flowState directory structure.
// Returns cached config on subsequent calls.
func Load() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	// Phase 1: Create application data directory
	dataDir := filepath.Join(homeDir, ".config", "flowState")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, err
	}

	cfg = &Config{
		DataDir:           dataDir,
		DbPath:            filepath.Join(dataDir, "flowState.db"),
		QdrantUrl:         "localhost:6333",
		ModelPath:         filepath.Join(dataDir, "models"),
		EmbeddingsEnabled: true,
	}

	return cfg, nil
}

// Get returns the cached configuration instance.
func Get() *Config {
	return cfg
}
