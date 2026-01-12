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

	"flowState-cli/internal/models"
	"flowState-cli/internal/storage/sqlite"
	"flowState-cli/internal/tui/components"
	"flowState-cli/internal/tui/styles"
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
	list       list.Model
	store      *sqlite.Store
	filter     string
	showCreate bool
	titleInput components.TextInputModel
	bodyInput  components.TextAreaModel
}

// NewNotesListModel creates a new notes list screen.
func NewNotesListModel(store *sqlite.Store) NotesListModel {
	items := []list.Item{}
	delegate := list.NewDefaultDelegate()

	l := list.New(items, delegate, 0, 0)
	l.Title = "Notes"
	l.SetShowHelp(true)

	return NotesListModel{
		list:       l,
		store:      store,
		filter:     "",
		showCreate: false,
		titleInput: components.NewTextInput("Note title"),
		bodyInput:  components.NewTextArea("Note body"),
	}
}

// Init implements tea.Model.
func (m *NotesListModel) Init() tea.Cmd {
	return func() tea.Msg { return nil }
}

// SetSize updates the list dimensions.
func (m *NotesListModel) SetSize(width, height int) {
	m.list.SetSize(width-4, height-10)
}

// LoadNotes refreshes the note list from the database.
func (m *NotesListModel) LoadNotes() error {
	notes, err := m.store.ListNotes()
	if err != nil {
		return err
	}

	items := make([]list.Item, 0, len(notes))
	for _, note := range notes {
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
						note := &models.Note{
							Title: title,
							Body:  body,
							Tags:  extractTags(body),
						}
						if err := m.store.CreateNote(note); err != nil {
							return m, nil
						}
						m.showCreate = false
						m.titleInput.SetValue("")
						m.bodyInput.SetValue("")
						m.LoadNotes()
					}
				}
				return m, nil
			case "ctrl+s":
				// Alternative save shortcut
				title := strings.TrimSpace(m.titleInput.Value())
				body := strings.TrimSpace(m.bodyInput.Value())
				if title != "" {
					note := &models.Note{
						Title: title,
						Body:  body,
						Tags:  extractTags(body),
					}
					if err := m.store.CreateNote(note); err != nil {
						return m, nil
					}
					m.showCreate = false
					m.titleInput.SetValue("")
					m.bodyInput.SetValue("")
					m.LoadNotes()
				}
				return m, nil
			case "esc":
				m.showCreate = false
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

		// Handle keys when viewing list
		switch msg.String() {
		case "c":
			m.showCreate = true
			m.titleInput.Focus()
			m.bodyInput.Blur()
		case "e":
			if len(m.list.VisibleItems()) > 0 {
				if selected, ok := m.list.SelectedItem().(NoteItem); ok {
					m.showCreate = true
					m.titleInput.SetValue(selected.note.Title)
					m.bodyInput.SetValue(selected.note.Body)
					m.titleInput.Focus()
				}
			}
		case "d":
			if len(m.list.VisibleItems()) > 0 {
				if selected, ok := m.list.SelectedItem().(NoteItem); ok {
					// TODO: Add confirmation dialog
					m.store.DeleteNote(selected.note.ID)
					m.LoadNotes()
				}
			}
		}

		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View renders the notes screen.
//
// Phase 2: Notes
//   - Shows create/edit form when active
//   - Shows note list otherwise
func (m *NotesListModel) View() string {
	if m.showCreate {
		// Show which field is focused
		titleLabel := styles.SubtitleStyle.Render("Title")
		bodyLabel := styles.SubtitleStyle.Render("Body")
		if m.titleInput.Focused() {
			titleLabel = styles.SelectedItemStyle.Render("▶ Title")
		} else {
			bodyLabel = styles.SelectedItemStyle.Render("▶ Body")
		}

		form := lipgloss.JoinVertical(
			lipgloss.Left,
			styles.TitleStyle.Render("📝 Create Note"),
			"",
			titleLabel,
			m.titleInput.View(),
			"",
			bodyLabel,
			m.bodyInput.View(),
			"",
			styles.HelpStyle.Render(
				styles.KeyHint("Tab", "Switch field")+" • "+
					styles.KeyHint("Ctrl+S", "Save")+" • "+
					styles.KeyHint("Esc", "Cancel"),
			),
		)
		return styles.PanelStyle.Render(form)
	}

	return m.list.View()
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
