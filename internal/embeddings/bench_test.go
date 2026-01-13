package embedder

import (
	"path/filepath"
	"testing"

	"flowState-cli/internal/config"
)

func BenchmarkEmbedding(b *testing.B) {
	tmpDir := b.TempDir()
	cfg := &config.Config{ModelPath: filepath.Join(tmpDir, "models")}

	e, err := New(cfg)
	if err != nil {
		b.Fatalf("New() err = %v", err)
	}

	text := "quick brown fox jumps over the lazy dog"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := e.EmbedSingle(text); err != nil {
			b.Fatalf("EmbedSingle() err = %v", err)
		}
	}
}


