package search

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/Jericoz-JC/flowState-CLI/internal/config"
	embeddings "github.com/Jericoz-JC/flowState-CLI/internal/embeddings"
	"github.com/Jericoz-JC/flowState-CLI/internal/models"
	"github.com/Jericoz-JC/flowState-CLI/internal/storage/sqlite"
)

func BenchmarkSearch1000Notes(b *testing.B) {
	tmpDir := b.TempDir()
	cfg := &config.Config{
		DbPath:    filepath.Join(tmpDir, "bench.db"),
		ModelPath: filepath.Join(tmpDir, "models"),
	}

	store, err := sqlite.New(cfg)
	if err != nil {
		b.Fatalf("sqlite.New() err = %v", err)
	}
	b.Cleanup(func() { _ = store.Close() })

	emb, err := embeddings.New(cfg)
	if err != nil {
		b.Fatalf("embeddings.New() err = %v", err)
	}
	searcher := New(emb, store)

	// Seed 1000 notes and index them once.
	for i := 0; i < 1000; i++ {
		n := &models.Note{
			Title: fmt.Sprintf("Note %d", i),
			Body:  fmt.Sprintf("This is note number %d about project planning and ideas.", i),
			Tags:  []string{"bench"},
		}
		if err := store.CreateNote(n); err != nil {
			b.Fatalf("CreateNote(%d) err = %v", i, err)
		}
	}
	if err := searcher.IndexAllNotes(); err != nil {
		b.Fatalf("IndexAllNotes() err = %v", err)
	}

	query := "Note 500 about project planning"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := searcher.Search(query, 10); err != nil {
			b.Fatalf("Search() err = %v", err)
		}
	}
}
