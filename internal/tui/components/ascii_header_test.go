package components

import (
	"strings"
	"testing"
)

func TestNewASCIIHeader(t *testing.T) {
	h := NewASCIIHeader("📝", "Notes")

	if h.icon != "📝" {
		t.Errorf("expected icon '📝', got '%s'", h.icon)
	}

	if h.title != "Notes" {
		t.Errorf("expected title 'Notes', got '%s'", h.title)
	}
}

func TestASCIIHeaderView(t *testing.T) {
	h := NewASCIIHeader("📝", "Notes")
	h.SetWidth(60)

	view := h.View()
	if view == "" {
		t.Fatalf("expected non-empty view")
	}

	// View should contain the title
	if !strings.Contains(view, "Notes") && !strings.Contains(view, "N O T E S") {
		t.Errorf("expected view to contain title")
	}
}

func TestASCIIHeaderStyles(t *testing.T) {
	tests := []struct {
		name  string
		style HeaderStyle
	}{
		{"minimal", HeaderStyleMinimal},
		{"boxed", HeaderStyleBoxed},
		{"banner", HeaderStyleBanner},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewASCIIHeader("🔍", "Search")
			h.SetStyle(tt.style)
			h.SetWidth(60)

			view := h.View()
			if view == "" {
				t.Errorf("expected non-empty view for style %s", tt.name)
			}
		})
	}
}

func TestASCIIHeaderBoxedHasBorders(t *testing.T) {
	h := NewASCIIHeader("✅", "Todos")
	h.SetStyle(HeaderStyleBoxed)
	h.SetWidth(60)

	view := h.View()

	// Boxed style should have border characters
	hasBorder := strings.Contains(view, "╔") ||
		strings.Contains(view, "║") ||
		strings.Contains(view, "╚") ||
		strings.Contains(view, "─") ||
		strings.Contains(view, "│")

	if !hasBorder {
		t.Errorf("boxed style should have border characters")
	}
}

func TestASCIIHeaderSetWidth(t *testing.T) {
	h := NewASCIIHeader("🧠", "Mind Map")

	h.SetWidth(40)
	if h.width != 40 {
		t.Errorf("expected width 40, got %d", h.width)
	}

	h.SetWidth(80)
	if h.width != 80 {
		t.Errorf("expected width 80, got %d", h.width)
	}
}

func TestASCIIHeaderSetItemCount(t *testing.T) {
	h := NewASCIIHeader("📝", "Notes")
	h.SetWidth(60)

	h.SetItemCount(5)
	view := h.View()

	// Should show item count
	if !strings.Contains(view, "5") {
		t.Errorf("expected view to contain item count '5'")
	}
}

func TestASCIIHeaderItemCountHidden(t *testing.T) {
	h := NewASCIIHeader("📝", "Notes")
	h.SetWidth(60)

	// -1 means hidden
	h.SetItemCount(-1)
	view := h.View()

	// Should not contain "items" text
	if strings.Contains(view, "items") {
		t.Errorf("item count should be hidden when set to -1")
	}
}

func TestASCIIHeaderSetSubtitle(t *testing.T) {
	h := NewASCIIHeader("🍅", "Focus")
	h.SetSubtitle("Work Session")
	h.SetWidth(60)

	view := h.View()

	// Should contain subtitle
	if !strings.Contains(view, "Work Session") {
		t.Errorf("expected view to contain subtitle")
	}
}

func TestScreenASCIIHeaders(t *testing.T) {
	screens := []string{"notes", "todos", "focus", "search", "mindmap"}

	for _, screen := range screens {
		t.Run(screen, func(t *testing.T) {
			ascii, exists := ScreenASCII[screen]
			if !exists {
				t.Errorf("expected ScreenASCII to contain '%s'", screen)
				return
			}
			if ascii == "" {
				t.Errorf("ScreenASCII['%s'] should not be empty", screen)
			}
		})
	}
}

func TestASCIIHeaderBannerStyle(t *testing.T) {
	h := NewASCIIHeader("🍅", "Focus")
	h.SetStyle(HeaderStyleBanner)
	h.SetWidth(60)

	view := h.View()

	// Banner should span full width (have decorative elements)
	if len(view) < 20 {
		t.Errorf("banner style should have substantial content")
	}
}
