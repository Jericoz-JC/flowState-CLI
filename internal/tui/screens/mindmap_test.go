package screens

import (
	"path/filepath"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Jericoz-JC/flowState-CLI/internal/config"
	"github.com/Jericoz-JC/flowState-CLI/internal/graph"
	"github.com/Jericoz-JC/flowState-CLI/internal/models"
	"github.com/Jericoz-JC/flowState-CLI/internal/storage/sqlite"
)

func TestMindMapScreenRender(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	cfg := &config.Config{DbPath: filepath.Join(tmpDir, "test.db")}
	store, err := sqlite.New(cfg)
	if err != nil {
		t.Fatalf("sqlite.New() err = %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })

	n1 := &models.Note{Title: "A"}
	n2 := &models.Note{Title: "B"}
	_ = store.CreateNote(n1)
	_ = store.CreateNote(n2)
	_ = store.CreateLink(&models.Link{SourceType: "note", SourceID: n1.ID, TargetType: "note", TargetID: n2.ID, LinkType: models.LinkTypeRelated})

	m := NewMindMapModel(store)
	m.SetSize(100, 40)
	if err := m.LoadGraph(); err != nil {
		t.Fatalf("LoadGraph() err = %v", err)
	}
	if v := m.View(); v == "" {
		t.Fatalf("expected non-empty view")
	}
}

func TestNodeSelection(t *testing.T) {
	t.Parallel()

	m := NewMindMapModel(nil)
	m.nodeOrder = []string{"note:1", "note:2"}
	m.selected = 0
	if m.nodeOrder[m.selected] != "note:1" {
		t.Fatalf("unexpected selection: %s", m.nodeOrder[m.selected])
	}
}

func TestNavigateBetweenNodes(t *testing.T) {
	t.Parallel()

	m := NewMindMapModel(nil)
	m.nodeOrder = []string{"note:1", "note:2"}
	m.positions = map[string]graph.Point{
		"note:1": {X: 0, Y: 0},
		"note:2": {X: 10, Y: 0},
	}
	m.selected = 0

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
	if updated.selected != 1 {
		t.Fatalf("expected selection to move right to idx=1, got %d", updated.selected)
	}
}

func TestOpenNoteFromGraph(t *testing.T) {
	t.Parallel()

	m := NewMindMapModel(nil)
	m.nodeOrder = []string{"note:42"}
	m.positions = map[string]graph.Point{"note:42": {X: 0, Y: 0}}
	m.selected = 0

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatalf("expected cmd")
	}
	msg := cmd()
	open, ok := msg.(OpenNoteMsg)
	if !ok {
		t.Fatalf("expected OpenNoteMsg, got %T", msg)
	}
	if open.NoteID != 42 {
		t.Fatalf("expected note id 42, got %d", open.NoteID)
	}
}

func TestZoomLevel(t *testing.T) {
	t.Parallel()

	m := NewMindMapModel(nil)
	if m.zoom != 1 {
		t.Fatalf("expected initial zoom 1, got %d", m.zoom)
	}

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'+'}})
	if updated.zoom != 2 {
		t.Fatalf("expected zoom 2 after +, got %d", updated.zoom)
	}

	updated2, _ := updated.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'-'}})
	if updated2.zoom != 1 {
		t.Fatalf("expected zoom 1 after -, got %d", updated2.zoom)
	}
}
