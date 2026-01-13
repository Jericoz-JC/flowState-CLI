// Package components provides reusable TUI components for flowState-cli.
//
// Phase 4: UX Overhaul
//   - HelpBar: Context-sensitive keyboard hints at bottom of screen
//   - Provides consistent navigation guidance across all screens
package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
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

	// Styles
	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F472B6")). // Pink accent
		Bold(true)

	primaryKeyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#22D3EE")). // Cyan for primary
		Bold(true)

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6C7086")) // Muted text

	separatorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#313244")) // Border color

	barStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#1E1E2E")). // Surface color
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

	separator := separatorStyle.Render(" • ")
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
		{Key: "d", Description: "Delete"},
		{Key: "/", Description: "Filter"},
		{Key: "Ctrl+L", Description: "Link"},
		{Key: "Ctrl+H", Description: "Home"},
	}

	// NotesEditHints are the hints when editing a note
	NotesEditHints = []HelpHint{
		{Key: "Tab", Description: "Switch Field"},
		{Key: "Ctrl+S", Description: "Save", Primary: true},
		{Key: "Esc", Description: "Cancel"},
	}

	// TodosListHints are the hints for the todos list view
	TodosListHints = []HelpHint{
		{Key: "c", Description: "Create", Primary: true},
		{Key: "e", Description: "Edit"},
		{Key: "d", Description: "Delete"},
		{Key: "Space", Description: "Toggle"},
		{Key: "Ctrl+L", Description: "Link"},
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
		{Key: "Enter", Description: "Save", Primary: true},
		{Key: "Esc", Description: "Cancel"},
	}

	// LinksHints are the hints for the links modal
	LinksHints = []HelpHint{
		{Key: "c", Description: "Create Link", Primary: true},
		{Key: "d", Description: "Delete"},
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
	FocusDurationHints = []HelpHint{
		{Key: "←/→", Description: "Select"},
		{Key: "Tab", Description: "Switch"},
		{Key: "Enter", Description: "Confirm", Primary: true},
		{Key: "Esc", Description: "Cancel"},
	}
)
