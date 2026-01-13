package graph

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"

	"github.com/Jericoz-JC/flowState-CLI/internal/models"
)

func TestRenderSingleNode(t *testing.T) {
	t.Parallel()

	out := RenderSingleNode("Hello", 12, "")
	if !strings.Contains(out, "┌") || !strings.Contains(out, "┐") {
		t.Fatalf("expected box drawing chars, got:\n%s", out)
	}
	if !strings.Contains(out, "Hello") {
		t.Fatalf("expected label in output, got:\n%s", out)
	}
}

func TestRenderEdge(t *testing.T) {
	t.Parallel()

	grid := [][]rune{[]rune(strings.Repeat(" ", 20))}
	drawHorizontalArrow(grid, 2, 8, 0)
	out := string(grid[0])
	if !strings.Contains(out, "▶") {
		t.Fatalf("expected arrow head, got: %q", out)
	}
	if !strings.Contains(out, "──") {
		t.Fatalf("expected line dashes, got: %q", out)
	}
}

func TestRenderSmallGraph(t *testing.T) {
	t.Parallel()

	links := []models.Link{
		{SourceType: "note", SourceID: 1, TargetType: "note", TargetID: 2, LinkType: models.LinkTypeRelated},
	}
	g := BuildGraphFromLinks(links, nil)

	labels := map[string]string{
		NodeKey("note", 1): "Alpha",
		NodeKey("note", 2): "Beta",
	}
	positions := map[string]Point{
		NodeKey("note", 1): {X: 0, Y: 0},
		NodeKey("note", 2): {X: 25, Y: 0}, // same Y so edge is drawn
	}

	out := RenderGraphASCII(g, labels, positions, 80, 10, nil)
	if !strings.Contains(out, "Alpha") || !strings.Contains(out, "Beta") {
		t.Fatalf("expected labels in output, got:\n%s", out)
	}
	if !strings.Contains(out, "▶") {
		t.Fatalf("expected edge arrow, got:\n%s", out)
	}
}

func TestFitToTerminalWidth(t *testing.T) {
	t.Parallel()

	links := []models.Link{
		{SourceType: "note", SourceID: 1, TargetType: "note", TargetID: 2},
	}
	g := BuildGraphFromLinks(links, nil)
	labels := map[string]string{
		NodeKey("note", 1): "AlphaAlphaAlphaAlpha",
		NodeKey("note", 2): "BetaBetaBetaBeta",
	}
	positions := map[string]Point{
		NodeKey("note", 1): {X: 0, Y: 0},
		NodeKey("note", 2): {X: 10, Y: 0},
	}

	width := 30
	out := RenderGraphASCII(g, labels, positions, width, 6, nil)
	for _, line := range strings.Split(out, "\n") {
		if len(line) > width {
			t.Fatalf("line exceeds width %d: %q", width, line)
		}
	}
}

func TestTagColorCoding(t *testing.T) {
	t.Parallel()

	// Force ANSI color output in tests.
	prev := lipgloss.ColorProfile()
	lipgloss.SetColorProfile(termenv.TrueColor)
	t.Cleanup(func() { lipgloss.SetColorProfile(prev) })

	links := []models.Link{
		{SourceType: "note", SourceID: 1, TargetType: "note", TargetID: 2},
	}
	g := BuildGraphFromLinks(links, nil)
	labels := map[string]string{
		NodeKey("note", 1): "Alpha",
		NodeKey("note", 2): "Beta",
	}
	positions := map[string]Point{
		NodeKey("note", 1): {X: 0, Y: 0},
		NodeKey("note", 2): {X: 25, Y: 0},
	}
	colors := map[string]string{
		NodeKey("note", 1): "#22D3EE",
	}

	out := RenderGraphASCII(g, labels, positions, 80, 10, colors)
	if !strings.Contains(out, "\x1b[") {
		t.Fatalf("expected ANSI color codes in output, got:\n%s", out)
	}
}
