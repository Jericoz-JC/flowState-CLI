package screens

import (
	"path/filepath"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Jericoz-JC/flowState-CLI/internal/config"
	"github.com/Jericoz-JC/flowState-CLI/internal/storage/sqlite"
)

func newTestNotesModel(t *testing.T) NotesListModel {
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

	model := NewNotesListModel(store)
	model.SetSize(100, 40)
	return model
}

func TestNotesScreenRender(t *testing.T) {
	t.Parallel()

	m := newTestNotesModel(t)
	v := m.View()
	if v == "" {
		t.Fatalf("expected non-empty view")
	}
}

func TestNotesCreateMode(t *testing.T) {
	t.Parallel()

	m := newTestNotesModel(t)

	// Press 'c' to enter create mode
	mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
	m = *mm.(*NotesListModel)

	if !m.showCreate {
		t.Fatalf("expected showCreate to be true after pressing 'c'")
	}
	if !m.titleInput.Focused() {
		t.Fatalf("expected title input to be focused in create mode")
	}
}

func TestNotesEnterInBodyCreatesNewline(t *testing.T) {
	t.Parallel()

	m := newTestNotesModel(t)

	// Enter create mode
	mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
	m = *mm.(*NotesListModel)

	// Type a title
	m.titleInput.SetValue("Test Note")

	// Tab to body
	mm, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = *mm.(*NotesListModel)

	if m.titleInput.Focused() {
		t.Fatalf("expected title input to NOT be focused after Tab")
	}
	if !m.bodyInput.Focused() {
		t.Fatalf("expected body input to be focused after Tab")
	}

	// Type some text in body
	mm, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'H'}})
	m = *mm.(*NotesListModel)
	mm, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'i'}})
	m = *mm.(*NotesListModel)

	// Press Enter - should NOT save (body is focused), should create newline
	mm, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = *mm.(*NotesListModel)

	// Should still be in create mode
	if !m.showCreate {
		t.Fatalf("expected to still be in create mode after Enter in body (Enter should create newline, not save)")
	}

	// The body should have the newline (textarea handles this)
	// We verify we didn't save by checking we're still in create mode
}

func TestNotesEnterInTitleSaves(t *testing.T) {
	t.Parallel()

	m := newTestNotesModel(t)

	// Enter create mode
	mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
	m = *mm.(*NotesListModel)

	// Type a title (title is focused by default)
	m.titleInput.SetValue("Test Note Title")
	m.bodyInput.SetValue("Test body content")

	// Press Enter while title is focused - should save
	mm, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = *mm.(*NotesListModel)

	// Should exit create mode after save
	if m.showCreate {
		t.Fatalf("expected to exit create mode after Enter in title (should save)")
	}
}

func TestNotesCtrlSSaves(t *testing.T) {
	t.Parallel()

	m := newTestNotesModel(t)

	// Enter create mode
	mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
	m = *mm.(*NotesListModel)

	// Type a title
	m.titleInput.SetValue("Test Note")

	// Tab to body
	mm, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = *mm.(*NotesListModel)

	// Type in body
	m.bodyInput.SetValue("Some body text")

	// Press Ctrl+S while body is focused - should still save
	mm, _ = m.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
	m = *mm.(*NotesListModel)

	// Should exit create mode
	if m.showCreate {
		t.Fatalf("expected Ctrl+S to save and exit create mode even when body is focused")
	}
}

func TestNotesEscCancels(t *testing.T) {
	t.Parallel()

	m := newTestNotesModel(t)

	// Enter create mode
	mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
	m = *mm.(*NotesListModel)

	// Press Esc to cancel
	mm, _ = m.Update(tea.KeyMsg{Type: tea.KeyEscape})
	m = *mm.(*NotesListModel)

	if m.showCreate {
		t.Fatalf("expected Esc to cancel and exit create mode")
	}
}
