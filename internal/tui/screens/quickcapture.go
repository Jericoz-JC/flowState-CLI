// Package screens provides TUI screen implementations for flowState-cli.
//
// Phase 4: UX Overhaul
//   - QuickCaptureModel: Instant note capture modal from anywhere
//   - Accessible via Ctrl+X global shortcut
//   - Auto-extracts title from first line
//   - Auto-tags with #quick for easy filtering
package screens

import (
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Jericoz-JC/flowState-CLI/internal/models"
	"github.com/Jericoz-JC/flowState-CLI/internal/storage/sqlite"
	"github.com/Jericoz-JC/flowState-CLI/internal/tui/components"
	"github.com/Jericoz-JC/flowState-CLI/internal/tui/styles"
)

// QuickCaptureModel implements a quick note capture modal.
type QuickCaptureModel struct {
	store   *sqlite.Store
	input   textarea.Model
	active  bool
	width   int
	height  int
	helpBar components.HelpBar
}

// NewQuickCaptureModel creates a new quick capture modal.
func NewQuickCaptureModel(store *sqlite.Store) QuickCaptureModel {
	ta := textarea.New()
	ta.Placeholder = "Type your thought...\n(First line becomes title, use #tags inline)"
	ta.Focus()
	ta.Prompt = ""
	ta.ShowLineNumbers = false
	ta.SetHeight(5)
	ta.SetWidth(50)
	ta.CharLimit = 5000

	return QuickCaptureModel{
		store:   store,
		input:   ta,
		active:  false,
		helpBar: components.NewHelpBar(components.QuickCaptureHints),
	}
}

// SetSize updates the modal dimensions.
func (m *QuickCaptureModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.input.SetWidth(width - 10)
	m.helpBar.SetWidth(width - 6)
}

// Open activates the quick capture modal.
func (m *QuickCaptureModel) Open() {
	m.active = true
	m.input.SetValue("")
	m.input.Focus()
}

// Close deactivates the quick capture modal.
func (m *QuickCaptureModel) Close() {
	m.active = false
	m.input.SetValue("")
	m.input.Blur()
}

// IsOpen returns whether the modal is currently active.
func (m *QuickCaptureModel) IsOpen() bool {
	return m.active
}

// Update handles messages for the quick capture modal.
func (m *QuickCaptureModel) Update(msg tea.Msg) (QuickCaptureModel, tea.Cmd) {
	if !m.active {
		return *m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.Close()
			return *m, nil
		case "ctrl+enter", "ctrl+s":
			// Save the note
			m.saveNote()
			m.Close()
			return *m, nil
		}
	}

	// Update textarea
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return *m, cmd
}

// saveNote creates a new note from the captured content.
func (m *QuickCaptureModel) saveNote() {
	content := strings.TrimSpace(m.input.Value())
	if content == "" {
		return
	}

	// Extract title from first line
	lines := strings.SplitN(content, "\n", 2)
	title := strings.TrimSpace(lines[0])
	body := ""
	if len(lines) > 1 {
		body = strings.TrimSpace(lines[1])
	}

	// If title is too long, truncate and put rest in body
	if len(title) > 50 {
		title = title[:50]
		body = content[50:] + "\n" + body
	}

	// Extract tags from content
	tags := extractQuickTags(content)

	// Always add #quick tag for filtering
	hasQuick := false
	for _, t := range tags {
		if t == "quick" {
			hasQuick = true
			break
		}
	}
	if !hasQuick {
		tags = append(tags, "quick")
	}

	note := &models.Note{
		Title: title,
		Body:  body,
		Tags:  tags,
	}

	m.store.CreateNote(note)
}

// extractQuickTags finds all #hashtags in content.
func extractQuickTags(content string) []string {
	tags := make(map[string]struct{})
	words := strings.Fields(content)
	for _, word := range words {
		if strings.HasPrefix(word, "#") {
			tag := strings.TrimPrefix(word, "#")
			// Clean up punctuation at end
			tag = strings.TrimRight(tag, ".,!?;:")
			tag = strings.ToLower(tag)
			if tag != "" {
				tags[tag] = struct{}{}
			}
		}
	}

	result := make([]string, 0, len(tags))
	for tag := range tags {
		result = append(result, tag)
	}
	return result
}

// View renders the quick capture modal.
func (m *QuickCaptureModel) View() string {
	if !m.active {
		return ""
	}

	// Styles using ARCHWAVE theme
	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(styles.AccentColor). // Hot pink
		Padding(1, 2).
		Width(m.width - 4)

	titleStyle := lipgloss.NewStyle().
		Foreground(styles.SecondaryColor). // Neon cyan
		Bold(true)

	tipStyle := lipgloss.NewStyle().
		Foreground(styles.MutedColor).
		Italic(true)

	// Build content
	title := titleStyle.Render(styles.DecoStar + " Quick Capture " + styles.DecoStar)

	tips := tipStyle.Render("Tip: First line → title • Use #tags inline • Ctrl+S to save")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		m.input.View(),
		"",
		tips,
		"",
		m.helpBar.View(),
	)

	return modalStyle.Render(content)
}
