// Package components provides reusable TUI components for flowState-cli.
//
// Phase 4: UX Overhaul
//   - HelpBar: Context-sensitive keyboard hints at bottom of screen
//   - Provides consistent navigation guidance across all screens
package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/Jericoz-JC/flowState-CLI/internal/tui/styles"
)

// HelpHint represents a single keyboard shortcut hint.
type HelpHint struct {
	Key         string // e.g., "c", "Ctrl+N"
	Description string // e.g., "Create", "Notes"
	Primary     bool   // Primary actions are highlighted
}

// HelpBar provides context-sensitive keyboard hints.
type HelpBar struct {
	hints []HelpHint
	width int
}

// NewHelpBar creates a new help bar with the given hints.
func NewHelpBar(hints []HelpHint) HelpBar {
	return HelpBar{
		hints: hints,
		width: 80,
	}
}

// SetWidth updates the width constraint for rendering.
func (h *HelpBar) SetWidth(width int) {
	h.width = width
}

// SetHints replaces the current hints.
func (h *HelpBar) SetHints(hints []HelpHint) {
	h.hints = hints
}

// View renders the help bar.
func (h *HelpBar) View() string {
	if len(h.hints) == 0 {
		return ""
	}

	// Styles using centralized ARCHWAVE theme
	keyStyle := lipgloss.NewStyle().
		Foreground(styles.AccentColor). // Hot pink
		Bold(true)

	primaryKeyStyle := lipgloss.NewStyle().
		Foreground(styles.SecondaryColor). // Neon cyan
		Bold(true)

	descStyle := lipgloss.NewStyle().
		Foreground(styles.MutedColor) // Pale blue

	barStyle := lipgloss.NewStyle().
		Background(styles.SurfaceColor). // Dark purple
		Padding(0, 1).
		Width(h.width)

	// Build hint strings
	var parts []string
	for _, hint := range h.hints {
		ks := keyStyle
		if hint.Primary {
			ks = primaryKeyStyle
		}
		part := ks.Render("["+hint.Key+"]") + " " + descStyle.Render(hint.Description)
		parts = append(parts, part)
	}

	// Use vaporwave separator
	separator := styles.VaporwaveSeparator()
	content := strings.Join(parts, separator)

	return barStyle.Render(content)
}

// Common hint sets for reuse
var (
	// HomeHints are the hints shown on the home screen
	HomeHints = []HelpHint{
		{Key: "Ctrl+N", Description: "Notes", Primary: true},
		{Key: "Ctrl+T", Description: "Todos", Primary: true},
		{Key: "Ctrl+X", Description: "Quick Capture", Primary: true},
		{Key: "Ctrl+F", Description: "Focus"},
		{Key: "q", Description: "Quit"},
	}

	// NotesListHints are the hints for the notes list view
	NotesListHints = []HelpHint{
		{Key: "c", Description: "Create", Primary: true},
		{Key: "e", Description: "Edit"},
		{Key: "p", Description: "Preview"},
		{Key: "d", Description: "Delete"},
		{Key: "/", Description: "Filter"},
		{Key: "Ctrl+R", Description: "Reset"},
		{Key: "Ctrl+H", Description: "Home"},
	}

	// NotesEditHints are the hints when editing a note
	NotesEditHints = []HelpHint{
		{Key: "Tab", Description: "Switch Field"},
		{Key: "Ctrl+S", Description: "Save", Primary: true},
		{Key: "Esc", Description: "Cancel"},
	}

	// NotesPreviewHints are the hints when previewing a note
	NotesPreviewHints = []HelpHint{
		{Key: "e", Description: "Edit", Primary: true},
		{Key: "Esc", Description: "Close"},
		{Key: "p", Description: "Close"},
	}

	// TodosListHints are the hints for the todos list view
	TodosListHints = []HelpHint{
		{Key: "c", Description: "Create", Primary: true},
		{Key: "e", Description: "Edit"},
		{Key: "d", Description: "Delete"},
		{Key: "Space", Description: "Toggle"},
		{Key: "?", Description: "Help"},
		{Key: "Ctrl+H", Description: "Home"},
	}

	// TodosEditHints are the hints when editing a todo
	TodosEditHints = []HelpHint{
		{Key: "Tab", Description: "Switch Field"},
		{Key: "Ctrl+S", Description: "Save", Primary: true},
		{Key: "Esc", Description: "Cancel"},
	}

	// QuickCaptureHints are the hints for quick capture modal
	QuickCaptureHints = []HelpHint{
		{Key: "Ctrl+S", Description: "Save", Primary: true},
		{Key: "?", Description: "Help"},
		{Key: "Esc", Description: "Cancel"},
	}

	// LinksHints are the hints for the links modal
	LinksHints = []HelpHint{
		{Key: "c", Description: "Create Link", Primary: true},
		{Key: "d", Description: "Delete"},
		{Key: "?", Description: "Help"},
		{Key: "Esc", Description: "Close"},
	}

	// ConfirmHints are the hints for confirmation dialogs
	ConfirmHints = []HelpHint{
		{Key: "y", Description: "Yes", Primary: true},
		{Key: "n", Description: "No"},
	}

	// FocusIdleHints are the hints when focus timer is idle
	FocusIdleHints = []HelpHint{
		{Key: "s", Description: "Start", Primary: true},
		{Key: "d", Description: "Duration"},
		{Key: "h", Description: "History"},
		{Key: "Ctrl+H", Description: "Home"},
	}

	// FocusRunningHints are the hints when focus timer is running
	FocusRunningHints = []HelpHint{
		{Key: "p", Description: "Pause", Primary: true},
		{Key: "c", Description: "Cancel"},
		{Key: "b", Description: "Skip to Break"},
	}

	// FocusPausedHints are the hints when focus timer is paused
	FocusPausedHints = []HelpHint{
		{Key: "s", Description: "Resume", Primary: true},
		{Key: "c", Description: "Cancel"},
	}

	// FocusBreakHints are the hints during break time
	FocusBreakHints = []HelpHint{
		{Key: "b", Description: "Skip Break", Primary: true},
		{Key: "c", Description: "Cancel"},
		{Key: "Esc", Description: "End Break"},
	}

	// FocusHistoryHints are the hints for session history view
	FocusHistoryHints = []HelpHint{
		{Key: "d", Description: "Delete"},
		{Key: "Esc", Description: "Back", Primary: true},
		{Key: "h", Description: "Back"},
	}

	// FocusDurationHints are the hints for duration picker
	// UX: Arrow keys update live with visual feedback, Tab switches work/break, Enter exits
	FocusDurationHints = []HelpHint{
		{Key: "←/→", Description: "Adjust (auto-saves)", Primary: true},
		{Key: "Tab", Description: "Work/Break"},
		{Key: "Enter", Description: "Done"},
		{Key: "Esc", Description: "Cancel"},
	}

	// SearchInputHints are the hints for the semantic search screen (query entry).
	SearchInputHints = []HelpHint{
		{Key: "Enter", Description: "Search", Primary: true},
		{Key: "?", Description: "Help"},
		{Key: "Ctrl+H", Description: "Home"},
		{Key: "Esc", Description: "Back"},
	}

	// SearchResultsHints are the hints for the semantic search screen (results navigation).
	SearchResultsHints = []HelpHint{
		{Key: "j/k", Description: "Navigate"},
		{Key: "Enter", Description: "Open Note", Primary: true},
		{Key: "?", Description: "Help"},
		{Key: "Esc", Description: "Edit Query"},
		{Key: "Ctrl+H", Description: "Home"},
	}

	// MindMapHints are the hints for the mind map screen.
	MindMapHints = []HelpHint{
		{Key: "h/j/k/l", Description: "Move"},
		{Key: "+/-", Description: "Zoom"},
		{Key: "Enter", Description: "Open Note", Primary: true},
		{Key: "?", Description: "Help"},
		{Key: "Ctrl+H", Description: "Home"},
	}
)
