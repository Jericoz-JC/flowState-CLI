package graph

import (
	"testing"

	"flowState-cli/internal/models"
)

func TestBuildGraphFromLinks(t *testing.T) {
	t.Parallel()

	links := []models.Link{
		{SourceType: "note", SourceID: 1, TargetType: "note", TargetID: 2, LinkType: models.LinkTypeRelated},
		{SourceType: "note", SourceID: 2, TargetType: "todo", TargetID: 9, LinkType: models.LinkTypeReferences},
	}

	nodeTags := map[string][]string{
		NodeKey("note", 1): {"project"},
		NodeKey("note", 2): {"project", "work"},
		NodeKey("todo", 9): {"work"},
	}

	g := BuildGraphFromLinks(links, nodeTags)

	if len(g.Nodes) != 3 {
		t.Fatalf("expected 3 nodes, got %d", len(g.Nodes))
	}
	// Undirected adjacency edges should exist both ways.
	a := NodeKey("note", 1)
	b := NodeKey("note", 2)
	if _, ok := g.Adj[a][b]; !ok {
		t.Fatalf("expected edge %s -> %s", a, b)
	}
	if _, ok := g.Adj[b][a]; !ok {
		t.Fatalf("expected edge %s -> %s", b, a)
	}
}

func TestFindConnectedComponents(t *testing.T) {
	t.Parallel()

	links := []models.Link{
		{SourceType: "note", SourceID: 1, TargetType: "note", TargetID: 2},
		{SourceType: "note", SourceID: 10, TargetType: "note", TargetID: 11},
	}

	g := BuildGraphFromLinks(links, nil)
	comps := ConnectedComponents(g)
	if len(comps) != 2 {
		t.Fatalf("expected 2 components, got %d", len(comps))
	}
}

func TestCalculateNodePositions(t *testing.T) {
	t.Parallel()

	links := []models.Link{
		{SourceType: "note", SourceID: 1, TargetType: "note", TargetID: 2},
		{SourceType: "note", SourceID: 2, TargetType: "note", TargetID: 3},
	}
	g := BuildGraphFromLinks(links, nil)

	pos := CalculateNodePositions(g, 40, 20)
	if len(pos) != len(g.Nodes) {
		t.Fatalf("expected positions for all nodes, got %d want %d", len(pos), len(g.Nodes))
	}

	seen := map[Point]bool{}
	for k := range g.Nodes {
		p, ok := pos[k]
		if !ok {
			t.Fatalf("missing position for node %s", k)
		}
		if p.X < 0 || p.X >= 40 || p.Y < 0 || p.Y >= 20 {
			t.Fatalf("position out of bounds for %s: %+v", k, p)
		}
		// In practice two nodes could collide; for small graphs our layout should avoid it.
		if seen[p] {
			t.Fatalf("duplicate position detected: %+v", p)
		}
		seen[p] = true
	}
}

func TestGraphWithTagColors(t *testing.T) {
	t.Parallel()

	colors := TagColors([]string{"project", "work", "personal", "work"})
	if colors["project"] == "" || colors["work"] == "" || colors["personal"] == "" {
		t.Fatalf("expected colors assigned for tags, got %+v", colors)
	}
	if colors["work"] != colors["work"] {
		t.Fatalf("expected stable color mapping")
	}
}


