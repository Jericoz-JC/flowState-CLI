package search

import (
	"path/filepath"
	"testing"

	"flowState-cli/internal/config"
	embeddings "flowState-cli/internal/embeddings"
	"flowState-cli/internal/models"
	"flowState-cli/internal/storage/sqlite"
)

func newTestStoreAndSearcher(t *testing.T) (*sqlite.Store, *SemanticSearch) {
	t.Helper()

	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	cfg := &config.Config{
		DbPath:    dbPath,
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

	searcher := New(emb, store)
	return store, searcher
}

func TestIndexAllNotes(t *testing.T) {
	t.Parallel()

	store, searcher := newTestStoreAndSearcher(t)

	n1 := &models.Note{Title: "A", Body: "hello world", Tags: []string{"t1"}}
	n2 := &models.Note{Title: "B", Body: "goodbye world", Tags: []string{"t2"}}
	if err := store.CreateNote(n1); err != nil {
		t.Fatalf("CreateNote(n1) err = %v", err)
	}
	if err := store.CreateNote(n2); err != nil {
		t.Fatalf("CreateNote(n2) err = %v", err)
	}

	if err := searcher.IndexAllNotes(); err != nil {
		t.Fatalf("IndexAllNotes() err = %v", err)
	}

	if _, ok, err := store.GetNoteEmbedding(n1.ID); err != nil || !ok {
		t.Fatalf("expected embedding for n1, ok=%v err=%v", ok, err)
	}
	if _, ok, err := store.GetNoteEmbedding(n2.ID); err != nil || !ok {
		t.Fatalf("expected embedding for n2, ok=%v err=%v", ok, err)
	}
}

func TestSearchReturnsRankedResults(t *testing.T) {
	t.Parallel()

	store, searcher := newTestStoreAndSearcher(t)

	n1 := &models.Note{Title: "A", Body: "alpha beta gamma", Tags: []string{"x"}}
	n2 := &models.Note{Title: "B", Body: "delta epsilon zeta", Tags: []string{"y"}}
	if err := store.CreateNote(n1); err != nil {
		t.Fatalf("CreateNote(n1) err = %v", err)
	}
	if err := store.CreateNote(n2); err != nil {
		t.Fatalf("CreateNote(n2) err = %v", err)
	}

	if err := searcher.IndexNote(n1.ID, n1.Title+"\n"+n1.Body); err != nil {
		t.Fatalf("IndexNote(n1) err = %v", err)
	}
	if err := searcher.IndexNote(n2.ID, n2.Title+"\n"+n2.Body); err != nil {
		t.Fatalf("IndexNote(n2) err = %v", err)
	}

	// Query identical to n1's indexed text should rank n1 first (cosine ~ 1.0).
	results, err := searcher.Search(n1.Title+"\n"+n1.Body, 10)
	if err != nil {
		t.Fatalf("Search() err = %v", err)
	}
	if len(results) < 2 {
		t.Fatalf("expected >= 2 results, got %d", len(results))
	}
	if results[0].NoteID != n1.ID {
		t.Fatalf("expected n1 to be ranked first, got note_id=%d", results[0].NoteID)
	}
	if results[0].Score < results[1].Score {
		t.Fatalf("expected first score >= second score, got %v < %v", results[0].Score, results[1].Score)
	}
}

func TestSearchWithTagFilter(t *testing.T) {
	t.Parallel()

	store, searcher := newTestStoreAndSearcher(t)

	n1 := &models.Note{Title: "A", Body: "project planning", Tags: []string{"project", "work"}}
	n2 := &models.Note{Title: "B", Body: "project planning", Tags: []string{"personal"}}
	if err := store.CreateNote(n1); err != nil {
		t.Fatalf("CreateNote(n1) err = %v", err)
	}
	if err := store.CreateNote(n2); err != nil {
		t.Fatalf("CreateNote(n2) err = %v", err)
	}

	if err := searcher.IndexAllNotes(); err != nil {
		t.Fatalf("IndexAllNotes() err = %v", err)
	}

	results, err := searcher.SearchWithTagFilter("project planning", 10, []string{"project"})
	if err != nil {
		t.Fatalf("SearchWithTagFilter() err = %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 filtered result, got %d", len(results))
	}
	if results[0].NoteID != n1.ID {
		t.Fatalf("expected n1 to match tag filter, got note_id=%d", results[0].NoteID)
	}
}

func TestSearchEmptyQuery(t *testing.T) {
	t.Parallel()

	_, searcher := newTestStoreAndSearcher(t)

	results, err := searcher.Search("", 10)
	if err != nil {
		t.Fatalf("Search() err = %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results for empty query, got %d", len(results))
	}
}


