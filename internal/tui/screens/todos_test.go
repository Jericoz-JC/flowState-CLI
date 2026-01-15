package screens

import (
	"path/filepath"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Jericoz-JC/flowState-CLI/internal/config"
	"github.com/Jericoz-JC/flowState-CLI/internal/storage/sqlite"
)

func newTestTodosModel(t *testing.T) *TodosListModel {
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

	model := NewTodosListModel(store)
	model.SetSize(100, 40)
	return &model
}

func TestTodosScreenRender(t *testing.T) {
	t.Parallel()

	m := newTestTodosModel(t)
	v := m.View()
	if v == "" {
		t.Fatalf("expected non-empty view")
	}
}

func TestTodosCreateModeEntry(t *testing.T) {
	t.Parallel()

	m := newTestTodosModel(t)

	// Press 'c' to enter create mode
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})

	if !m.showCreate {
		t.Fatalf("expected showCreate to be true after pressing 'c'")
	}

	// Title input should be focused
	if !m.titleInput.Focused() {
		t.Fatalf("expected titleInput to be focused in create mode")
	}
}

func TestTodosTitleInputCapture(t *testing.T) {
	t.Parallel()

	m := newTestTodosModel(t)

	// Enter create mode
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})

	// Type characters into title
	testTitle := "Test Todo Title"
	for _, char := range testTitle {
		m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{char}})
	}

	// Verify title was captured
	actualTitle := m.titleInput.Value()
	if actualTitle != testTitle {
		t.Fatalf("expected title %q, got %q", testTitle, actualTitle)
	}
}

func TestTodosCreateAndSave(t *testing.T) {
	t.Parallel()

	m := newTestTodosModel(t)

	// Get initial count
	initialTodos, _ := m.store.ListTodos()
	initialCount := len(initialTodos)

	// Enter create mode
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})

	// Type title
	for _, char := range "My New Todo" {
		m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{char}})
	}

	// Press Enter to save (title should be focused)
	if !m.titleInput.Focused() {
		t.Fatalf("title input should be focused before saving")
	}
	m.Update(tea.KeyMsg{Type: tea.KeyEnter})

	// Should exit create mode
	if m.showCreate {
		t.Fatalf("expected showCreate to be false after Enter")
	}

	// Todo should be saved
	todos, _ := m.store.ListTodos()
	if len(todos) != initialCount+1 {
		t.Fatalf("expected %d todos after create, got %d", initialCount+1, len(todos))
	}

	// Verify the title was saved correctly
	found := false
	for _, todo := range todos {
		if todo.Title == "My New Todo" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("created todo not found with expected title 'My New Todo'")
	}
}

func TestTodosEscCancelsCreate(t *testing.T) {
	t.Parallel()

	m := newTestTodosModel(t)

	// Enter create mode
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})

	// Type something
	for _, char := range "Test" {
		m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{char}})
	}

	// Press Esc to cancel
	m.Update(tea.KeyMsg{Type: tea.KeyEscape})

	// Should exit create mode
	if m.showCreate {
		t.Fatalf("expected showCreate to be false after Esc")
	}

	// Input should be cleared
	if m.titleInput.Value() != "" {
		t.Fatalf("expected title input to be cleared after cancel")
	}
}

func TestTodosTabSwitchesFocus(t *testing.T) {
	t.Parallel()

	m := newTestTodosModel(t)

	// Enter create mode
	m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})

	// Title should be focused initially
	if !m.titleInput.Focused() {
		t.Fatalf("expected title to be focused initially")
	}

	// Press Tab
	m.Update(tea.KeyMsg{Type: tea.KeyTab})

	// Description should now be focused
	if m.titleInput.Focused() {
		t.Fatalf("expected title to NOT be focused after Tab")
	}
	if !m.descInput.Focused() {
		t.Fatalf("expected description to be focused after Tab")
	}

	// Press Tab again
	m.Update(tea.KeyMsg{Type: tea.KeyTab})

	// Title should be focused again
	if !m.titleInput.Focused() {
		t.Fatalf("expected title to be focused after second Tab")
	}
}
