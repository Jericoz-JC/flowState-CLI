// Package screens provides TUI screen implementations for flowState-cli.
//
// Phase 2: Notes & Todos
//   - NotesListModel: Note management UI
//   - TodosListModel: Todo management UI
//   - Create, read, update, delete operations
//   - Auto-tagging from #hashtag syntax
package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Jericoz-JC/flowState-CLI/internal/models"
	"github.com/Jericoz-JC/flowState-CLI/internal/storage/sqlite"
	"github.com/Jericoz-JC/flowState-CLI/internal/tui/components"
	"github.com/Jericoz-JC/flowState-CLI/internal/tui/keymap"
	"github.com/Jericoz-JC/flowState-CLI/internal/tui/styles"
)

// NotesListModel implements the notes management screen.
//
// Phase 2: Notes
//   - Displays of all notes sorted by update date
//   - Shows date list, title, and tags for each note
//   - Create new notes with title and body
//   - Edit existing notes
//   - Delete notes
//   - Auto-extract tags from #hashtag syntax
//
// Keyboard Shortcuts (when viewing list):
//   - c: Create new note
//   - e: Edit selected note
//   - d: Delete selected note
//   - j/down: Move selection down
//   - k/up: Move selection up
//   - esc: Cancel/create mode
//   - enter: Save note (in create/edit mode)
//
// Keyboard Shortcuts (when creating/editing):
//   - enter: Save and return to list
//   - esc: Cancel and return to list
type NotesListModel struct {
	list             list.Model
	store            *sqlite.Store
	filter           string
	filterInput      components.TextInputModel
	showFilter       bool
	selectedTags     []string // Tags to filter by
	showCreate       bool
	showPreview      bool         // Preview mode (read-only markdown)
	previewNote      *models.Note // Note being previewed
	editingID        int64        // 0 = creating new, >0 = editing existing
	confirmingDelete bool
	deleteTargetID   int64
	titleInput       components.TextInputModel
	bodyInput        components.TextAreaModel
	header           components.Header
	helpBar          components.HelpBar
	width            int
	height           int
}

// NewNotesListModel creates a new notes list screen.
func NewNotesListModel(store *sqlite.Store) NotesListModel {
	items := []list.Item{}
	delegate := list.NewDefaultDelegate()

	l := list.New(items, delegate, 0, 0)
	l.Title = ""
	l.SetShowHelp(false) // We'll use our own help bar
	l.SetShowTitle(false)
	l.SetFilteringEnabled(false) // We handle filtering ourselves

	filterInput := components.NewTextInput("Type to filter...")
	filterInput.Blur()

	return NotesListModel{
		list:             l,
		store:            store,
		filter:           "",
		filterInput:      filterInput,
		showFilter:       false,
		selectedTags:     []string{},
		showCreate:       false,
		showPreview:      false,
		previewNote:      nil,
		editingID:        0,
		confirmingDelete: false,
		deleteTargetID:   0,
		titleInput:       components.NewTextInput("Note title"),
		bodyInput:        components.NewTextArea("Note body"),
		header:           components.NewHeader("📝", "Notes"),
		helpBar:          components.NewHelpBar(components.NotesListHints),
	}
}

// Init implements tea.Model.
func (m *NotesListModel) Init() tea.Cmd {
	return func() tea.Msg { return nil }
}

// SetSize updates the list dimensions.
func (m *NotesListModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.list.SetSize(width-4, height-14) // Account for header and help bar
	m.header.SetWidth(width - 4)
	m.helpBar.SetWidth(width - 4)
}

// GetSelectedNote returns the currently selected note, or nil if none selected.
func (m *NotesListModel) GetSelectedNote() *models.Note {
	if len(m.list.Items()) == 0 {
		return nil
	}
	if selected, ok := m.list.SelectedItem().(NoteItem); ok {
		return &selected.note
	}
	return nil
}

// SelectNoteByID selects a note in the list by its ID (best-effort).
func (m *NotesListModel) SelectNoteByID(id int64) {
	items := m.list.Items()
	for i, it := range items {
		if ni, ok := it.(NoteItem); ok && ni.note.ID == id {
			m.list.Select(i)
			return
		}
	}
}

// LoadNotes refreshes the note list from the database.
func (m *NotesListModel) LoadNotes() error {
	notes, err := m.store.ListNotes()
	if err != nil {
		return err
	}

	// Apply filters
	filtered := make([]models.Note, 0)
	for _, note := range notes {
		// Filter by search text
		if m.filter != "" {
			searchText := strings.ToLower(m.filter)
			titleMatch := strings.Contains(strings.ToLower(note.Title), searchText)
			bodyMatch := strings.Contains(strings.ToLower(note.Body), searchText)
			if !titleMatch && !bodyMatch {
				continue
			}
		}

		// Filter by selected tags
		if len(m.selectedTags) > 0 {
			hasAllTags := true
			for _, selectedTag := range m.selectedTags {
				found := false
				for _, noteTag := range note.Tags {
					if noteTag == selectedTag {
						found = true
						break
					}
				}
				if !found {
					hasAllTags = false
					break
				}
			}
			if !hasAllTags {
				continue
			}
		}

		filtered = append(filtered, note)
	}

	items := make([]list.Item, 0, len(filtered))
	for _, note := range filtered {
		items = append(items, NoteItem{note: note})
	}

	m.list.SetItems(items)
	return nil
}

// Update handles messages for the notes screen.
//
// Phase 2: Notes
//   - Key bindings for navigation
//   - Create/edit/delete operations
//   - Form input handling
//   - Tab to switch between fields
func (m *NotesListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle filter input
		if m.showFilter {
			switch msg.String() {
			case "enter":
				m.filter = m.filterInput.Value()
				m.showFilter = false
				m.filterInput.Blur()
				m.LoadNotes()
				return m, nil
			case "esc":
				m.showFilter = false
				m.filter = ""
				m.filterInput.SetValue("")
				m.filterInput.Blur()
				m.LoadNotes()
				return m, nil
			default:
				var cmd tea.Cmd
				m.filterInput, cmd = m.filterInput.Update(msg)
				cmds = append(cmds, cmd)
				return m, tea.Batch(cmds...)
			}
		}

		// Handle preview mode
		if m.showPreview {
			switch msg.String() {
			case "esc", "p", "q":
				m.showPreview = false
				m.previewNote = nil
				return m, nil
			}
			return m, nil
		}

		// Handle delete confirmation dialog
		if m.confirmingDelete {
			switch msg.String() {
			case "y", "Y":
				m.store.DeleteNote(m.deleteTargetID)
				m.confirmingDelete = false
				m.deleteTargetID = 0
				m.LoadNotes()
				return m, nil
			case "n", "N", "esc":
				m.confirmingDelete = false
				m.deleteTargetID = 0
				return m, nil
			}
			return m, nil
		}

		// Handle keys when in create/edit mode
		if m.showCreate {
			switch msg.String() {
			case "tab", "shift+tab":
				// Toggle focus between title and body
				if m.titleInput.Focused() {
					m.titleInput.Blur()
					m.bodyInput.Focus()
				} else {
					m.bodyInput.Blur()
					m.titleInput.Focus()
				}
				return m, nil
			case "enter":
				// Only save if title input is focused (allow newlines in body)
				if m.titleInput.Focused() {
					title := strings.TrimSpace(m.titleInput.Value())
					body := strings.TrimSpace(m.bodyInput.Value())
					if title != "" {
						tags := extractTags(body)
						wikilinks := parseWikilinks(body)

						if m.editingID > 0 {
							// Update existing note
							note := &models.Note{
								ID:    m.editingID,
								Title: title,
								Body:  body,
								Tags:  tags,
							}
							if err := m.store.UpdateNote(note); err != nil {
								return m, nil
							}
							// Create wikilinks
							m.createWikilinks(note.ID, wikilinks)
						} else {
							// Create new note
							note := &models.Note{
								Title: title,
								Body:  body,
								Tags:  tags,
							}
							if err := m.store.CreateNote(note); err != nil {
								return m, nil
							}
							// Create wikilinks
							m.createWikilinks(note.ID, wikilinks)
						}
						m.showCreate = false
						m.editingID = 0
						m.titleInput.SetValue("")
						m.bodyInput.SetValue("")
						m.LoadNotes()
					}
				}
				return m, nil
			}

			// Check for cross-platform save shortcut
			if keymap.IsModS(msg) {
				// Alternative save shortcut
				title := strings.TrimSpace(m.titleInput.Value())
				body := strings.TrimSpace(m.bodyInput.Value())
				if title != "" {
					tags := extractTags(body)
					wikilinks := parseWikilinks(body)

					if m.editingID > 0 {
						// Update existing note
						note := &models.Note{
							ID:    m.editingID,
							Title: title,
							Body:  body,
							Tags:  tags,
						}
						if err := m.store.UpdateNote(note); err != nil {
							return m, nil
						}
						// Create wikilinks
						m.createWikilinks(note.ID, wikilinks)
					} else {
						// Create new note
						note := &models.Note{
							Title: title,
							Body:  body,
							Tags:  tags,
						}
						if err := m.store.CreateNote(note); err != nil {
							return m, nil
						}
						// Create wikilinks
						m.createWikilinks(note.ID, wikilinks)
					}
					m.showCreate = false
					m.editingID = 0
					m.titleInput.SetValue("")
					m.bodyInput.SetValue("")
					m.LoadNotes()
				}
				return m, nil
			}

			if msg.String() == "esc" {
				m.showCreate = false
				m.editingID = 0
				m.titleInput.SetValue("")
				m.bodyInput.SetValue("")
				return m, nil
			}

			// Update the focused input
			var cmd tea.Cmd
			if m.titleInput.Focused() {
				m.titleInput, cmd = m.titleInput.Update(msg)
			} else {
				m.bodyInput, cmd = m.bodyInput.Update(msg)
			}
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		}

		// Handle keys when viewing list - process BEFORE passing to list
		switch msg.String() {
		case "/":
			// Open filter input
			m.showFilter = true
			m.filterInput.SetValue(m.filter)
			m.filterInput.Focus()
			return m, nil
		case "p":
			// Preview selected note
			if len(m.list.VisibleItems()) > 0 {
				if selected, ok := m.list.SelectedItem().(NoteItem); ok {
					fullNote, err := m.store.GetNote(selected.note.ID)
					if err != nil || fullNote == nil {
						return m, nil
					}
					m.showPreview = true
					m.previewNote = fullNote
				}
			}
			return m, nil
		case "t":
			// Toggle tag filter for selected note
			if len(m.list.VisibleItems()) > 0 {
				if selected, ok := m.list.SelectedItem().(NoteItem); ok {
					if len(selected.note.Tags) > 0 {
						// Add first tag to filter (in future, show tag selector)
						tag := selected.note.Tags[0]
						m.toggleTagFilter(tag)
						m.LoadNotes()
					}
				}
			}
			return m, nil
		case "c":
			m.showCreate = true
			m.editingID = 0
			m.titleInput.SetValue("")
			m.bodyInput.SetValue("")
			m.titleInput.Focus()
			m.bodyInput.Blur()
			return m, nil // Return early to prevent list from processing
		case "e":
			if len(m.list.VisibleItems()) > 0 {
				if selected, ok := m.list.SelectedItem().(NoteItem); ok {
					// Phase 4: Performance - Fetch full note content
					fullNote, err := m.store.GetNote(selected.note.ID)
					if err != nil || fullNote == nil {
						// TODO: Show error message
						return m, nil
					}

					m.showCreate = true
					m.editingID = fullNote.ID
					m.titleInput.SetValue(fullNote.Title)
					m.bodyInput.SetValue(fullNote.Body)
					m.titleInput.Focus()
				}
			}
			return m, nil
		case "d":
			if len(m.list.VisibleItems()) > 0 {
				if selected, ok := m.list.SelectedItem().(NoteItem); ok {
					m.confirmingDelete = true
					m.deleteTargetID = selected.note.ID
				}
			}
			return m, nil
		}

		// Check for cross-platform reset shortcut
		if keymap.IsModR(msg) {
			// Reset all filters
			m.filter = ""
			m.selectedTags = []string{}
			m.LoadNotes()
			return m, nil
		}

		// Pass other keys to list for navigation (j/k, up/down, etc.)
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View renders the notes screen.
//
// Phase 4: UX Overhaul
//   - Header with title and item count
//   - Context-sensitive help bar
//   - Shows create/edit form when active
//   - Preview mode for reading notes
//   - Filter input for searching
func (m *NotesListModel) View() string {
	// Preview mode
	if m.showPreview {
		return m.renderPreview()
	}

	// Filter input mode
	if m.showFilter {
		filterHints := []components.HelpHint{
			{Key: "Enter", Description: "Apply", Primary: true},
			{Key: "Esc", Description: "Cancel"},
		}
		m.helpBar.SetHints(filterHints)

		filterLabel := styles.TitleStyle.Render("🔍 Filter Notes")
		filterHelp := styles.SubtitleStyle.Render("Type to search by title or content")

		content := lipgloss.JoinVertical(
			lipgloss.Left,
			filterLabel,
			"",
			filterHelp,
			m.filterInput.View(),
			"",
			m.helpBar.View(),
		)
		return styles.PanelStyle.Render(content)
	}

	// Delete confirmation dialog
	if m.confirmingDelete {
		m.helpBar.SetHints(components.ConfirmHints)
		confirmDialog := lipgloss.JoinVertical(
			lipgloss.Center,
			styles.TitleStyle.Render("⚠️ Delete Note?"),
			"",
			styles.SubtitleStyle.Render("This action cannot be undone."),
			"",
			m.helpBar.View(),
		)
		return styles.PanelStyle.Render(confirmDialog)
	}

	if m.showCreate {
		m.helpBar.SetHints(components.NotesEditHints)

		// Show which field is focused
		titleLabel := styles.SubtitleStyle.Render("Title")
		bodyLabel := styles.SubtitleStyle.Render("Body (use #tags and [[links]])")
		if m.titleInput.Focused() {
			titleLabel = styles.SelectedItemStyle.Render("▶ Title")
		} else {
			bodyLabel = styles.SelectedItemStyle.Render("▶ Body (use #tags and [[links]])")
		}

		// Dynamic title for create vs edit
		formTitle := "📝 Create Note"
		if m.editingID > 0 {
			formTitle = "📝 Edit Note"
		}

		form := lipgloss.JoinVertical(
			lipgloss.Left,
			styles.TitleStyle.Render(formTitle),
			"",
			titleLabel,
			m.titleInput.View(),
			"",
			bodyLabel,
			m.bodyInput.View(),
			"",
			m.helpBar.View(),
		)
		return styles.PanelStyle.Render(form)
	}

	// Update header with item count and active filters
	m.header.SetItemCount(len(m.list.Items()))

	// Update help hints to include preview and filter (with platform-appropriate mod key)
	mod := keymap.ModKeyDisplay()
	listHints := []components.HelpHint{
		{Key: "c", Description: "Create", Primary: true},
		{Key: "e", Description: "Edit"},
		{Key: "p", Description: "Preview"},
		{Key: "d", Description: "Delete"},
		{Key: "/", Description: "Filter"},
		{Key: "t", Description: "Tag Filter"},
		{Key: mod + "+L", Description: "Link"},
		{Key: mod + "+H", Description: "Home"},
	}
	m.helpBar.SetHints(listHints)

	// Show active filters
	var filterStatus string
	if m.filter != "" || len(m.selectedTags) > 0 {
		filterParts := []string{}
		if m.filter != "" {
			filterParts = append(filterParts, fmt.Sprintf("search:%q", m.filter))
		}
		if len(m.selectedTags) > 0 {
			for _, tag := range m.selectedTags {
				filterParts = append(filterParts, "#"+tag)
			}
		}
		filterStatusStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F9E2AF")).
			Background(lipgloss.Color("#2E2E3E")).
			Padding(0, 1)
		filterStatus = filterStatusStyle.Render("🔎 Filtering: " + strings.Join(filterParts, ", ") + " [Ctrl+R to reset]")
	}

	// Empty state
	if len(m.list.Items()) == 0 {
		emptyMsg := "No notes yet. Start capturing your thoughts!"
		if m.filter != "" || len(m.selectedTags) > 0 {
			emptyMsg = "No notes match your filters. Press [Ctrl+R] to reset."
		}
		emptyState := lipgloss.JoinVertical(
			lipgloss.Left,
			m.header.View(),
			"",
			styles.SubtitleStyle.Render(emptyMsg),
			"",
			styles.HelpStyle.Render("Press [c] to create your first note"),
			"",
			m.helpBar.View(),
		)
		return styles.PanelStyle.Render(emptyState)
	}

	// Regular list view with header and help bar
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		m.header.View(),
		"",
	)
	if filterStatus != "" {
		content = lipgloss.JoinVertical(lipgloss.Left, content, filterStatus, "")
	}
	content = lipgloss.JoinVertical(
		lipgloss.Left,
		content,
		m.list.View(),
		"",
		m.helpBar.View(),
	)
	return content
}

// NoteItem implements list.Item for displaying notes in the list.
//
// Phase 2: Notes
//   - Title: Shows date and note title with tags
//   - Description: Shows body preview
//   - FilterValue: Used for search/filter
type NoteItem struct {
	note models.Note
}

func (n NoteItem) Title() string {
	date := n.note.UpdatedAt.Format("2006-01-02")
	tags := ""
	if len(n.note.Tags) > 0 {
		tags = " [" + strings.Join(n.note.Tags, ", ") + "]"
	}
	return fmt.Sprintf("%s %s%s", date, n.note.Title, tags)
}

func (n NoteItem) Description() string {
	preview := n.note.Body
	if len(preview) > 60 {
		preview = preview[:60] + "..."
	}
	return preview
}

func (n NoteItem) FilterValue() string {
	return n.note.Title + " " + n.note.Body
}

// toggleTagFilter adds or removes a tag from the filter list.
func (m *NotesListModel) toggleTagFilter(tag string) {
	for i, t := range m.selectedTags {
		if t == tag {
			// Remove tag
			m.selectedTags = append(m.selectedTags[:i], m.selectedTags[i+1:]...)
			return
		}
	}
	// Add tag
	m.selectedTags = append(m.selectedTags, tag)
}

// renderPreview renders a note in preview mode with markdown-like formatting.
func (m *NotesListModel) renderPreview() string {
	if m.previewNote == nil {
		return ""
	}

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#A78BFA")).
		Bold(true).
		Padding(0, 1)

	dateStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6C7086")).
		Italic(true)

	tagStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#F472B6")).
		Background(lipgloss.Color("#2E2E3E")).
		Padding(0, 1).
		MarginRight(1)

	bodyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#CDD6F4")).
		Padding(1, 2)

	wikilinkStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#22D3EE")).
		Underline(true)

	// Title
	title := titleStyle.Render(m.previewNote.Title)

	// Date
	date := dateStyle.Render(m.previewNote.UpdatedAt.Format("2006-01-02 15:04"))

	// Tags
	var tags string
	if len(m.previewNote.Tags) > 0 {
		tagParts := []string{}
		for _, tag := range m.previewNote.Tags {
			tagParts = append(tagParts, tagStyle.Render("#"+tag))
		}
		tags = strings.Join(tagParts, "")
	}

	// Body with wikilink highlighting
	body := m.previewNote.Body
	body = highlightWikilinks(body, wikilinkStyle)
	body = bodyStyle.Render(body)

	// Help hints
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#6C7086")).
		Padding(1, 2)
	help := helpStyle.Render("[p/esc] Close preview  [e] Edit")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		date,
		tags,
		"",
		body,
		"",
		help,
	)

	return styles.PanelStyle.Render(content)
}

// highlightWikilinks finds [[text]] patterns and highlights them.
func highlightWikilinks(text string, style lipgloss.Style) string {
	// Simple regex-free approach
	result := ""
	inLink := false
	linkStart := 0

	for i := 0; i < len(text); i++ {
		if i < len(text)-1 && text[i] == '[' && text[i+1] == '[' {
			if !inLink {
				inLink = true
				linkStart = i
				i++ // Skip second [
				continue
			}
		}
		if i < len(text)-1 && text[i] == ']' && text[i+1] == ']' && inLink {
			// Found end of wikilink
			linkText := text[linkStart+2 : i]
			result += style.Render("[[" + linkText + "]]")
			inLink = false
			i++ // Skip second ]
			continue
		}
		if !inLink {
			result += string(text[i])
		}
	}
	return result
}

// parseWikilinks extracts all [[Note Name]] patterns from text.
func parseWikilinks(text string) []string {
	links := []string{}
	inLink := false
	linkStart := 0

	for i := 0; i < len(text); i++ {
		if i < len(text)-1 && text[i] == '[' && text[i+1] == '[' {
			if !inLink {
				inLink = true
				linkStart = i + 2
				i++
			}
		} else if i < len(text)-1 && text[i] == ']' && text[i+1] == ']' && inLink {
			linkText := strings.TrimSpace(text[linkStart:i])
			if linkText != "" {
				links = append(links, linkText)
			}
			inLink = false
			i++
		}
	}
	return links
}

// createWikilinks creates links from the current note to notes mentioned in [[...]] syntax.
func (m *NotesListModel) createWikilinks(sourceNoteID int64, wikilinks []string) {
	if len(wikilinks) == 0 {
		return
	}

	// Get all notes to match titles
	allNotes, err := m.store.ListNotes()
	if err != nil {
		return
	}

	// For each wikilink, find or create the target note
	for _, linkTitle := range wikilinks {
		var targetID int64
		found := false

		// Search for existing note with this title
		for _, note := range allNotes {
			if strings.EqualFold(strings.TrimSpace(note.Title), strings.TrimSpace(linkTitle)) {
				targetID = note.ID
				found = true
				break
			}
		}

		// If not found, create a placeholder note
		if !found {
			placeholderNote := &models.Note{
				Title: linkTitle,
				Body:  "(Created from wikilink)",
				Tags:  []string{"placeholder"},
			}
			if err := m.store.CreateNote(placeholderNote); err != nil {
				continue
			}
			targetID = placeholderNote.ID
		}

		// Create the link
		link := &models.Link{
			SourceType: "note",
			SourceID:   sourceNoteID,
			TargetType: "note",
			TargetID:   targetID,
			LinkType:   "wikilink",
		}
		m.store.CreateLink(link)
	}
}

// extractTags finds all #hashtags in content and returns them as a slice.
//
// Phase 2: Notes
//   - Parses content for #word patterns
//   - Converts tags to lowercase
//   - Removes duplicates
func extractTags(content string) []string {
	tags := make(map[string]struct{})
	words := strings.Fields(content)
	for _, word := range words {
		if strings.HasPrefix(word, "#") {
			tag := strings.TrimPrefix(word, "#")
			tag = strings.TrimSpace(tag)
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
