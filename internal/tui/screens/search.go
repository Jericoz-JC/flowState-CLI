package screens

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Jericoz-JC/flowState-CLI/internal/search"
	"github.com/Jericoz-JC/flowState-CLI/internal/storage/sqlite"
	"github.com/Jericoz-JC/flowState-CLI/internal/tui/components"
	"github.com/Jericoz-JC/flowState-CLI/internal/tui/styles"
)

// OpenNoteMsg is emitted by the Search screen when the user selects a result.
type OpenNoteMsg struct {
	NoteID int64
}

type searchMode int

const (
	searchModeInput searchMode = iota
	searchModeResults
)

type SearchModel struct {
	store    *sqlite.Store
	semantic *search.SemanticSearch

	mode     searchMode
	query    components.TextInputModel
	results  []search.SearchResult
	selected int
	loading  bool
	errText  string
	showHelp bool // Help modal state

	header  components.Header
	helpBar components.HelpBar
	width   int
	height  int
}

type searchCompletedMsg struct {
	results []search.SearchResult
	err     error
}

func NewSearchModel(store *sqlite.Store, semantic *search.SemanticSearch) SearchModel {
	return SearchModel{
		store:    store,
		semantic: semantic,
		mode:     searchModeInput,
		query:    components.NewTextInput("Search notes (semantic)..."),
		results:  nil,
		selected: 0,
		loading:  false,
		errText:  "",
		header:   components.NewHeader("🔍", "Search"),
		helpBar:  components.NewHelpBar(components.SearchInputHints),
	}
}

func (m *SearchModel) Init() tea.Cmd { return nil }

func (m *SearchModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.header.SetWidth(width - 4)
	m.helpBar.SetWidth(width - 4)
}

func (m *SearchModel) Update(msg tea.Msg) (SearchModel, tea.Cmd) {
	switch msg := msg.(type) {
	case searchCompletedMsg:
		m.loading = false
		if msg.err != nil {
			m.errText = msg.err.Error()
			m.results = nil
			m.mode = searchModeInput
			m.query.Focus()
			m.helpBar.SetHints(components.SearchInputHints)
			return *m, nil
		}
		m.errText = ""
		m.results = msg.results
		m.selected = 0
		m.mode = searchModeResults
		m.query.Blur()
		m.helpBar.SetHints(components.SearchResultsHints)
		return *m, nil
	case tea.KeyMsg:
		// Handle help modal
		if m.showHelp {
			m.showHelp = false
			return *m, nil
		}

		// ? opens help from any mode
		if msg.String() == "?" {
			m.showHelp = true
			return *m, nil
		}

		switch m.mode {
		case searchModeInput:
			switch msg.String() {
			case "enter":
				if m.loading {
					return *m, nil
				}
				q := strings.TrimSpace(m.query.Value())
				m.errText = ""
				if q == "" {
					m.results = nil
					return *m, nil
				}
				m.loading = true
				return *m, func() tea.Msg {
					results, err := m.semantic.Search(q, 20)
					return searchCompletedMsg{results: results, err: err}
				}
			default:
				var cmd tea.Cmd
				m.query, cmd = m.query.Update(msg)
				return *m, cmd
			}
		case searchModeResults:
			switch msg.String() {
			case "esc":
				m.mode = searchModeInput
				m.query.Focus()
				m.helpBar.SetHints(components.SearchInputHints)
				return *m, nil
			case "j", "down":
				if m.selected < len(m.results)-1 {
					m.selected++
				}
				return *m, nil
			case "k", "up":
				if m.selected > 0 {
					m.selected--
				}
				return *m, nil
			case "enter":
				if len(m.results) == 0 {
					return *m, nil
				}
				noteID := m.results[m.selected].NoteID
				return *m, func() tea.Msg { return OpenNoteMsg{NoteID: noteID} }
			}
		}
	}

	return *m, nil
}

func (m *SearchModel) View() string {
	panel := lipgloss.NewStyle().Padding(1, 2).Width(m.width).Height(m.height)

	// Show help modal if active
	if m.showHelp {
		return panel.Render(m.helpView())
	}

	bodyWidth := m.width - 4

	// Update header with result count
	if len(m.results) > 0 {
		m.header.SetItemCount(len(m.results))
	} else {
		m.header.SetItemCount(-1)
	}

	title := m.header.View()
	queryLine := styles.InputStyle.Render(m.query.View())

	var contentParts []string
	contentParts = append(contentParts, title)
	contentParts = append(contentParts, "")
	contentParts = append(contentParts, styles.SubtitleStyle.Render("Semantic search across your notes"))
	contentParts = append(contentParts, "")
	contentParts = append(contentParts, queryLine)

	if m.loading {
		loadingStyle := lipgloss.NewStyle().Foreground(styles.SecondaryColor)
		contentParts = append(contentParts, "")
		contentParts = append(contentParts, loadingStyle.Render("✦ Searching..."))
	}

	if m.errText != "" {
		errStyle := lipgloss.NewStyle().Foreground(styles.ErrorColor)
		contentParts = append(contentParts, "")
		contentParts = append(contentParts, errStyle.Render("Error: "+m.errText))
	}

	contentParts = append(contentParts, "")
	contentParts = append(contentParts, m.renderResults(bodyWidth))
	contentParts = append(contentParts, "")
	contentParts = append(contentParts, m.helpBar.View())

	return panel.Render(lipgloss.JoinVertical(lipgloss.Left, contentParts...))
}

func (m *SearchModel) renderResults(width int) string {
	if m.mode == searchModeInput && strings.TrimSpace(m.query.Value()) == "" {
		return styles.HelpStyle.Render("Type a query and press Enter to search.")
	}
	if len(m.results) == 0 {
		return styles.HelpStyle.Render("No results.")
	}

	rowStyle := lipgloss.NewStyle().Width(width).Foreground(styles.TextColor)
	selectedStyle := lipgloss.NewStyle().
		Width(width).
		Background(styles.SurfaceColor).
		Foreground(styles.SecondaryColor).
		Bold(true)

	lines := make([]string, 0, len(m.results))
	for i, r := range m.results {
		line := fmt.Sprintf("[%.2f] %s", r.Score, firstLine(r.NoteText))
		if i == m.selected && m.mode == searchModeResults {
			lines = append(lines, selectedStyle.Render(line))
		} else {
			lines = append(lines, rowStyle.Render(line))
		}
	}
	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func firstLine(s string) string {
	if idx := strings.IndexByte(s, '\n'); idx >= 0 {
		return s[:idx]
	}
	return s
}

// helpView renders the help modal for the search screen.
func (m *SearchModel) helpView() string {
	title := styles.TitleStyle.Render("🔍 SEARCH - Help")

	helpText := `Semantic search finds notes based on meaning, not just keywords.

` + styles.SelectedItemStyle.Render("How it Works:") + `
• Type a natural language query (e.g., "meeting notes from last week")
• Press Enter to search
• Results are ranked by semantic similarity

` + styles.SelectedItemStyle.Render("Navigation:") + `
• ` + styles.NeonStyle.Render("Enter") + `: Execute search / Open selected note
• ` + styles.NeonStyle.Render("j/k") + ` or Arrow Keys: Navigate results
• ` + styles.NeonStyle.Render("Esc") + `: Edit query / Go back

` + styles.SelectedItemStyle.Render("Tips:") + `
• Use descriptive queries for better results
• Search works across note titles and bodies
• Score indicates match quality (higher = better match)`

	help := styles.HelpStyle.Render("Press any key to close")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		helpText,
		"",
		help,
	)
}
