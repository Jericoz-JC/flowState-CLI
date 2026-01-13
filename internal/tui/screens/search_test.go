package screens

import (
	"path/filepath"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Jericoz-JC/flowState-CLI/internal/config"
	embeddings "github.com/Jericoz-JC/flowState-CLI/internal/embeddings"
	"github.com/Jericoz-JC/flowState-CLI/internal/models"
	"github.com/Jericoz-JC/flowState-CLI/internal/search"
	"github.com/Jericoz-JC/flowState-CLI/internal/storage/sqlite"
)

func newTestSearchModel(t *testing.T) SearchModel {
	t.Helper()

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

	semantic := search.New(emb, store)
	model := NewSearchModel(store, semantic)
	model.SetSize(100, 40)
	return model
}

func TestSearchScreenRender(t *testing.T) {
	t.Parallel()

	m := newTestSearchModel(t)
	v := m.View()
	if v == "" {
		t.Fatalf("expected non-empty view")
	}
}

func TestSearchInputHandling(t *testing.T) {
	t.Parallel()

	m := newTestSearchModel(t)

	mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
	m = mm
	mm, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
	m = mm

	if m.query.Value() != "hi" {
		t.Fatalf("expected query to be %q, got %q", "hi", m.query.Value())
	}
}

func TestResultsListNavigation(t *testing.T) {
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

	emb, _ := embeddings.New(cfg)
	semantic := search.New(emb, store)

	n1 := &models.Note{Title: "A", Body: "alpha"}
	n2 := &models.Note{Title: "B", Body: "beta"}
	_ = store.CreateNote(n1)
	_ = store.CreateNote(n2)
	_ = semantic.IndexAllNotes()

	m := NewSearchModel(store, semantic)
	m.SetSize(100, 40)
	m.query.SetValue("a")

	mm, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = mm
	if cmd == nil {
		t.Fatalf("expected search cmd")
	}
	// Execute async command and feed result back in.
	mm, _ = m.Update(cmd())
	m = mm
	if m.mode != searchModeResults {
		t.Fatalf("expected mode results after search")
	}

	old := m.selected
	mm, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m = mm
	if len(m.results) > 1 && m.selected == old {
		t.Fatalf("expected selection to change after j")
	}
}

func TestOpenNoteFromSearch(t *testing.T) {
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

	emb, _ := embeddings.New(cfg)
	semantic := search.New(emb, store)

	n := &models.Note{Title: "A", Body: "alpha"}
	_ = store.CreateNote(n)
	_ = semantic.IndexAllNotes()

	m := NewSearchModel(store, semantic)
	m.SetSize(100, 40)
	m.query.SetValue("alpha")

	mm, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = mm
	if cmd == nil {
		t.Fatalf("expected search cmd")
	}
	mm, _ = m.Update(cmd())
	m = mm
	if len(m.results) == 0 {
		t.Fatalf("expected results")
	}

	mm, cmd = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = mm
	if cmd == nil {
		t.Fatalf("expected cmd to open note")
	}
	msg := cmd()
	open, ok := msg.(OpenNoteMsg)
	if !ok {
		t.Fatalf("expected OpenNoteMsg, got %T", msg)
	}
	if open.NoteID != m.results[m.selected].NoteID {
		t.Fatalf("expected note_id %d, got %d", m.results[m.selected].NoteID, open.NoteID)
	}
}
