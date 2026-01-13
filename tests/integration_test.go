package tests

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"flowState-cli/internal/config"
	embeddings "flowState-cli/internal/embeddings"
	"flowState-cli/internal/graph"
	"flowState-cli/internal/models"
	"flowState-cli/internal/search"
	"flowState-cli/internal/storage/sqlite"
)

func TestFullWorkflow(t *testing.T) {
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

	note := &models.Note{Title: "Integration Note", Body: "alpha beta gamma", Tags: []string{"integration"}}
	if err := store.CreateNote(note); err != nil {
		t.Fatalf("CreateNote() err = %v", err)
	}
	if err := searcher.IndexAllNotes(); err != nil {
		t.Fatalf("IndexAllNotes() err = %v", err)
	}

	results, err := searcher.Search(note.Title+"\n"+note.Body, 5)
	if err != nil {
		t.Fatalf("Search() err = %v", err)
	}
	if len(results) == 0 || results[0].NoteID != note.ID {
		t.Fatalf("expected to find created note, results=%v", results)
	}
}

func TestLinkingAndGraph(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{DbPath: filepath.Join(tmpDir, "test.db")}

	store, err := sqlite.New(cfg)
	if err != nil {
		t.Fatalf("sqlite.New() err = %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })

	n1 := &models.Note{Title: "A", Tags: []string{"project"}}
	n2 := &models.Note{Title: "B", Tags: []string{"project"}}
	_ = store.CreateNote(n1)
	_ = store.CreateNote(n2)

	_ = store.CreateLink(&models.Link{
		SourceType: "note",
		SourceID:   n1.ID,
		TargetType: "note",
		TargetID:   n2.ID,
		LinkType:   models.LinkTypeRelated,
	})

	links, err := store.ListLinks()
	if err != nil {
		t.Fatalf("ListLinks() err = %v", err)
	}
	notes, err := store.ListNotes()
	if err != nil {
		t.Fatalf("ListNotes() err = %v", err)
	}

	nodeTags := map[string][]string{}
	labels := map[string]string{}
	for _, n := range notes {
		k := graph.NodeKey("note", n.ID)
		nodeTags[k] = n.Tags
		labels[k] = n.Title
	}

	g := graph.BuildGraphFromLinks(links, nodeTags)
	pos := map[string]graph.Point{
		graph.NodeKey("note", n1.ID): {X: 0, Y: 0},
		graph.NodeKey("note", n2.ID): {X: 25, Y: 0},
	}

	out := graph.RenderGraphASCII(g, labels, pos, 80, 10, nil)
	if !containsAll(out, []string{"A", "B"}) {
		t.Fatalf("expected graph output to contain both titles, got:\n%s", out)
	}
}

func TestFocusSessionFlow(t *testing.T) {
	tmpDir := t.TempDir()
	cfg := &config.Config{DbPath: filepath.Join(tmpDir, "test.db")}

	store, err := sqlite.New(cfg)
	if err != nil {
		t.Fatalf("sqlite.New() err = %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })

	start := time.Now()
	end := start.Add(25 * time.Minute)
	s := &models.FocusSession{
		StartTime: start,
		EndTime:   &end,
		Duration:  25 * 60,
		Status:    models.SessionStatusCompleted,
	}
	if err := store.CreateSession(s); err != nil {
		t.Fatalf("CreateSession() err = %v", err)
	}

	stats, err := store.GetSessionStats()
	if err != nil {
		t.Fatalf("GetSessionStats() err = %v", err)
	}
	if stats.TotalSessions != 1 {
		t.Fatalf("expected 1 total session, got %d", stats.TotalSessions)
	}
}

func containsAll(s string, subs []string) bool {
	for _, sub := range subs {
		if !strings.Contains(s, sub) {
			return false
		}
	}
	return true
}


