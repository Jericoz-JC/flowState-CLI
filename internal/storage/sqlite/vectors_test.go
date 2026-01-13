package sqlite

import (
	"path/filepath"
	"testing"

	"flowState-cli/internal/config"
	"flowState-cli/internal/models"
)

func TestCreateVectorTable(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	cfg := &config.Config{DbPath: dbPath}

	store, err := New(cfg)
	if err != nil {
		t.Fatalf("New() err = %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })

	var count int
	if err := store.db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='note_vectors'").Scan(&count); err != nil {
		t.Fatalf("sqlite_master query err = %v", err)
	}
	if count != 1 {
		t.Fatalf("expected note_vectors table to exist, count=%d", count)
	}
}

func TestStoreEmbedding(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	cfg := &config.Config{DbPath: dbPath}

	store, err := New(cfg)
	if err != nil {
		t.Fatalf("New() err = %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })

	note := &models.Note{Title: "vec note", Body: "hello"}
	if err := store.CreateNote(note); err != nil {
		t.Fatalf("CreateNote() err = %v", err)
	}

	emb := make([]float32, 384)
	emb[0] = 1
	emb[1] = 2
	if err := store.UpsertNoteEmbedding(note.ID, emb); err != nil {
		t.Fatalf("UpsertNoteEmbedding() err = %v", err)
	}

	got, ok, err := store.GetNoteEmbedding(note.ID)
	if err != nil {
		t.Fatalf("GetNoteEmbedding() err = %v", err)
	}
	if !ok {
		t.Fatalf("expected embedding to exist")
	}
	if len(got) != 384 || got[0] != 1 || got[1] != 2 {
		t.Fatalf("unexpected embedding roundtrip: len=%d got[0]=%v got[1]=%v", len(got), got[0], got[1])
	}
}

func TestSearchByVector(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	cfg := &config.Config{DbPath: dbPath}

	store, err := New(cfg)
	if err != nil {
		t.Fatalf("New() err = %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })

	n1 := &models.Note{Title: "n1"}
	n2 := &models.Note{Title: "n2"}
	if err := store.CreateNote(n1); err != nil {
		t.Fatalf("CreateNote(n1) err = %v", err)
	}
	if err := store.CreateNote(n2); err != nil {
		t.Fatalf("CreateNote(n2) err = %v", err)
	}

	e1 := make([]float32, 384)
	e2 := make([]float32, 384)
	e1[0] = 1
	e2[1] = 1

	if err := store.UpsertNoteEmbedding(n1.ID, e1); err != nil {
		t.Fatalf("UpsertNoteEmbedding(n1) err = %v", err)
	}
	if err := store.UpsertNoteEmbedding(n2.ID, e2); err != nil {
		t.Fatalf("UpsertNoteEmbedding(n2) err = %v", err)
	}

	results, err := store.SearchNoteEmbeddings(e1, 10)
	if err != nil {
		t.Fatalf("SearchNoteEmbeddings() err = %v", err)
	}
	if len(results) < 2 {
		t.Fatalf("expected at least 2 results, got %d", len(results))
	}
	if results[0].NoteID != n1.ID {
		t.Fatalf("expected best match to be n1, got note_id=%d", results[0].NoteID)
	}
	if results[0].Score <= results[1].Score {
		t.Fatalf("expected first score > second score, got %v <= %v", results[0].Score, results[1].Score)
	}
}

func TestUpdateEmbeddingOnNoteEdit(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	cfg := &config.Config{DbPath: dbPath}

	store, err := New(cfg)
	if err != nil {
		t.Fatalf("New() err = %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })

	note := &models.Note{Title: "n"}
	if err := store.CreateNote(note); err != nil {
		t.Fatalf("CreateNote() err = %v", err)
	}

	e1 := make([]float32, 384)
	e2 := make([]float32, 384)
	e1[0] = 1
	e2[0] = 2

	if err := store.UpsertNoteEmbedding(note.ID, e1); err != nil {
		t.Fatalf("UpsertNoteEmbedding(e1) err = %v", err)
	}
	if err := store.UpsertNoteEmbedding(note.ID, e2); err != nil {
		t.Fatalf("UpsertNoteEmbedding(e2) err = %v", err)
	}

	got, ok, err := store.GetNoteEmbedding(note.ID)
	if err != nil {
		t.Fatalf("GetNoteEmbedding() err = %v", err)
	}
	if !ok || got[0] != 2 {
		t.Fatalf("expected updated embedding, ok=%v got[0]=%v", ok, got[0])
	}
}

func TestDeleteEmbeddingOnNoteDelete(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")
	cfg := &config.Config{DbPath: dbPath}

	store, err := New(cfg)
	if err != nil {
		t.Fatalf("New() err = %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })

	note := &models.Note{Title: "n"}
	if err := store.CreateNote(note); err != nil {
		t.Fatalf("CreateNote() err = %v", err)
	}

	e := make([]float32, 384)
	e[0] = 1
	if err := store.UpsertNoteEmbedding(note.ID, e); err != nil {
		t.Fatalf("UpsertNoteEmbedding() err = %v", err)
	}

	if err := store.DeleteNote(note.ID); err != nil {
		t.Fatalf("DeleteNote() err = %v", err)
	}

	_, ok, err := store.GetNoteEmbedding(note.ID)
	if err != nil {
		t.Fatalf("GetNoteEmbedding() err = %v", err)
	}
	if ok {
		t.Fatalf("expected embedding to be deleted when note is deleted")
	}
}


