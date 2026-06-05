// Package config provides configuration management for flowState-cli.
//
// Paths are resolved to platform-native locations via os.UserConfigDir():
//   - Windows: %AppData%\flowState
//   - macOS:   ~/Library/Application Support/flowState
//   - Linux:   $XDG_CONFIG_HOME/flowState (or ~/.config/flowState)
//
// Configuration Fields:
//   - DataDir: Base directory for all application data
//   - DbPath: SQLite database file path
//   - QdrantUrl: Vector database URL for semantic search (dev stub)
//   - ModelPath: Path to store embedding models
//   - LogPath: Absolute path to the debug log file
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
	LogPath           string `mapstructure:"log_path"`
	EmbeddingsEnabled bool   `mapstructure:"embeddings_enabled"`
}

var cfg *Config

// Load initializes configuration with platform-native defaults.
// It creates the data and log directories and returns a cached config on
// subsequent calls. Path resolution happens here so logging can be pointed
// at LogPath before any TUI initialization.
func Load() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	base, err := baseConfigDir()
	if err != nil {
		return nil, err
	}

	c := Resolve(base)

	// Best-effort migration of legacy ~/.config/flowState data to the
	// platform-native location (mainly affects macOS, where older builds
	// wrote to ~/.config instead of ~/Library/Application Support).
	migrateLegacyData(c.DataDir)

	if err := os.MkdirAll(c.DataDir, 0755); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(c.LogPath), 0755); err != nil {
		return nil, err
	}

	cfg = c
	return cfg, nil
}

// Resolve builds a Config from a platform config base directory.
// Kept pure (no filesystem side effects) so it is deterministic and testable.
func Resolve(baseConfigDir string) *Config {
	dataDir := filepath.Join(baseConfigDir, "flowState")
	return &Config{
		DataDir:           dataDir,
		DbPath:            filepath.Join(dataDir, "flowState.db"),
		QdrantUrl:         "localhost:6333",
		ModelPath:         filepath.Join(dataDir, "models"),
		LogPath:           filepath.Join(dataDir, "logs", "debug.log"),
		EmbeddingsEnabled: true,
	}
}

// baseConfigDir returns the platform-native config base directory, falling
// back to ~/.config if os.UserConfigDir is unavailable.
func baseConfigDir() (string, error) {
	if dir, err := os.UserConfigDir(); err == nil && dir != "" {
		return dir, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config"), nil
}

// migrateLegacyData moves data from the old ~/.config/flowState location to
// newDataDir when the new location does not yet exist. Best-effort: any error
// is ignored so a failed migration never blocks startup.
func migrateLegacyData(newDataDir string) {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}
	legacy := filepath.Join(home, ".config", "flowState")
	if legacy == newDataDir {
		return
	}
	if _, err := os.Stat(newDataDir); err == nil {
		return // new location already populated; don't clobber
	}
	if _, err := os.Stat(legacy); err != nil {
		return // nothing to migrate
	}
	_ = os.MkdirAll(filepath.Dir(newDataDir), 0755)
	_ = os.Rename(legacy, newDataDir)
}

// Get returns the cached configuration instance.
func Get() *Config {
	return cfg
}
