package config

import (
	"path/filepath"
	"strings"
	"testing"
)

// TestResolveDerivesPathsFromBase verifies all paths are derived from the
// platform config base directory and nested under a single flowState dir.
func TestResolveDerivesPathsFromBase(t *testing.T) {
	base := filepath.Join("tmp", "cfgbase")
	c := Resolve(base)

	wantDataDir := filepath.Join(base, "flowState")
	if c.DataDir != wantDataDir {
		t.Errorf("DataDir = %q, want %q", c.DataDir, wantDataDir)
	}
	if c.DbPath != filepath.Join(wantDataDir, "flowState.db") {
		t.Errorf("DbPath = %q, want it under DataDir", c.DbPath)
	}
	if c.ModelPath != filepath.Join(wantDataDir, "models") {
		t.Errorf("ModelPath = %q, want it under DataDir", c.ModelPath)
	}
	if c.LogPath != filepath.Join(wantDataDir, "logs", "debug.log") {
		t.Errorf("LogPath = %q, want %q", c.LogPath, filepath.Join(wantDataDir, "logs", "debug.log"))
	}
}

// TestResolveLogPathIsDebugLog verifies the log file is named debug.log and
// lives in a dedicated logs directory (not the current working directory).
func TestResolveLogPathIsDebugLog(t *testing.T) {
	c := Resolve(filepath.Join("any", "base"))
	if filepath.Base(c.LogPath) != "debug.log" {
		t.Errorf("log basename = %q, want debug.log", filepath.Base(c.LogPath))
	}
	if filepath.Base(filepath.Dir(c.LogPath)) != "logs" {
		t.Errorf("log parent dir = %q, want logs", filepath.Base(filepath.Dir(c.LogPath)))
	}
	if strings.HasPrefix(c.LogPath, ".") || c.LogPath == "debug.log" {
		t.Errorf("LogPath %q should be an absolute-ish app path, not CWD-relative", c.LogPath)
	}
}

// TestResolveEmbeddingsDefaultEnabled documents the default toggle.
func TestResolveEmbeddingsDefaultEnabled(t *testing.T) {
	c := Resolve("base")
	if !c.EmbeddingsEnabled {
		t.Error("EmbeddingsEnabled should default to true")
	}
}

// TestBaseConfigDirNonEmpty verifies the platform base resolver returns a path.
func TestBaseConfigDirNonEmpty(t *testing.T) {
	dir, err := baseConfigDir()
	if err != nil {
		t.Fatalf("baseConfigDir() error: %v", err)
	}
	if dir == "" {
		t.Error("baseConfigDir() returned empty path")
	}
}
