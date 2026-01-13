package context

import (
	"path/filepath"
	"testing"

	"github.com/Jericoz-JC/flowState-CLI/internal/config"
	embeddings "github.com/Jericoz-JC/flowState-CLI/internal/embeddings"
	"github.com/Jericoz-JC/flowState-CLI/internal/models"
	"github.com/Jericoz-JC/flowState-CLI/internal/search"
	"github.com/Jericoz-JC/flowState-CLI/internal/storage/sqlite"
)

func TestSuggestRelatedNotes(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	cfg := &config.Config{
		DbPath:    filepath.Join(tmpDir, "test.db"),
		ModelPath: filepath.Join(tmpDir, "models"),
	}

	store, err := sqlite.New(cfg)
	if err != nil {
		t.Fatalf("sqlite.New() err = %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })

	emb, err := embeddings.New(cfg)
	if err != nil {
		t.Fatalf("embeddings.New() err = %v", err)
	}
	searcher := search.New(emb, store)

	n1 := &models.Note{Title: "Work", Body: "planning project roadmap"}
	n2 := &models.Note{Title: "Groceries", Body: "buy milk eggs bread"}
	_ = store.CreateNote(n1)
	_ = store.CreateNote(n2)
	_ = searcher.IndexAllNotes()

	// Query should return n2 when excludeNoteID is n1 and query matches n2 exactly.
	got, err := SuggestRelatedNotes(searcher, store, n1.ID, n2.Title+"\n"+n2.Body, 1)
	if err != nil {
		t.Fatalf("SuggestRelatedNotes() err = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(got))
	}
	if got[0].ID != n2.ID {
		t.Fatalf("expected note %d, got %d", n2.ID, got[0].ID)
	}
}

func TestSuggestTagsFromContent(t *testing.T) {
	t.Parallel()

	tags := SuggestTagsFromContent("Hello #Work #work, also #Project!")
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %v", tags)
	}
	if tags[0] != "project" || tags[1] != "work" {
		t.Fatalf("unexpected tags: %v", tags)
	}
}

func TestSuggestLinksFromWikilinks(t *testing.T) {
	t.Parallel()

	links := SuggestLinksFromWikilinks("See [[Note A]] and [[Note B]] plus [[]] and [[  ]].")
	if len(links) != 2 {
		t.Fatalf("expected 2 links, got %v", links)
	}
	if links[0] != "Note A" || links[1] != "Note B" {
		t.Fatalf("unexpected links: %v", links)
	}
}
