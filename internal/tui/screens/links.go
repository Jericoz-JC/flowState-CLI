// Package screens provides TUI screen implementations for flowState-cli.
//
// Phase 3: Linking System
//   - LinkModel: Link management UI
//   - Create bidirectional links between notes and todos
//   - View and delete existing links
//   - Support for related, contains, references link types
package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"flowState-cli/internal/models"
	"flowState-cli/internal/storage/sqlite"
	"flowState-cli/internal/tui/styles"
)

// LinkMode represents what the link modal is doing
type LinkMode int

const (
	LinkModeSelectSource LinkMode = iota
	LinkModeSelectType
	LinkModeSelectTarget
	LinkModeViewLinks
)

// LinkModel implements the link creation and management modal.
//
// Phase 3: Linking System
//   - Shows existing links for selected item
//   - Allows creating new links with type selection
//   - Supports bidirectional link queries
type LinkModel struct {
	store        *sqlite.Store
	mode         LinkMode
	sourceType   string // "note" or "todo"
	sourceID     int64
	sourceTitle  string
	targetType   string
	targetID     int64
	selectedType models.LinkType
	linkTypes    []models.LinkType
	typeIndex    int
	targetList   list.Model
	linkList     list.Model
	notes        []models.Note
	todos        []models.Todo
	links        []models.Link
	width        int
	height       int
	showModal    bool
}

// NewLinkModel creates a new link management model.
func NewLinkModel(store *sqlite.Store) LinkModel {
	targetDelegate := list.NewDefaultDelegate()
	targetList := list.New([]list.Item{}, targetDelegate, 0, 0)
	targetList.Title = "Select Target"
	targetList.SetShowHelp(false)

	linkDelegate := list.NewDefaultDelegate()
	linkList := list.New([]list.Item{}, linkDelegate, 0, 0)
	linkList.Title = "Linked Items"
	linkList.SetShowHelp(false)

	return LinkModel{
		store:      store,
		mode:       LinkModeViewLinks,
		linkTypes:  []models.LinkType{models.LinkTypeRelated, models.LinkTypeContains, models.LinkTypeReferences},
		typeIndex:  0,
		targetList: targetList,
		linkList:   linkList,
		showModal:  false,
	}
}

// SetSize updates the model dimensions.
func (m *LinkModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.targetList.SetSize(width-10, height-15)
	m.linkList.SetSize(width-10, height-15)
}

// Open opens the link modal for a specific source item.
func (m *LinkModel) Open(sourceType string, sourceID int64, sourceTitle string) {
	m.sourceType = sourceType
	m.sourceID = sourceID
	m.sourceTitle = sourceTitle
	m.showModal = true
	m.mode = LinkModeViewLinks
	m.loadLinks()
}

// Close closes the link modal.
func (m *LinkModel) Close() {
	m.showModal = false
	m.mode = LinkModeViewLinks
}

// IsOpen returns whether the modal is currently open.
func (m *LinkModel) IsOpen() bool {
	return m.showModal
}

// loadLinks loads existing links for the current source item.
func (m *LinkModel) loadLinks() {
	links, err := m.store.GetLinksForItem(m.sourceType, m.sourceID)
	if err != nil {
		return
	}
	m.links = links

	items := make([]list.Item, 0, len(links))
	for _, link := range links {
		items = append(items, LinkItem{link: link, store: m.store, sourceType: m.sourceType, sourceID: m.sourceID})
	}
	m.linkList.SetItems(items)
}

// loadTargets loads potential link targets (notes and todos).
func (m *LinkModel) loadTargets() {
	notes, _ := m.store.ListNotes()
	todos, _ := m.store.ListTodos()
	m.notes = notes
	m.todos = todos

	items := make([]list.Item, 0)

	// Add notes as potential targets (if source is not this note)
	for _, note := range notes {
		if m.sourceType == "note" && note.ID == m.sourceID {
			continue // Skip self
		}
		items = append(items, TargetItem{
			itemType: "note",
			id:       note.ID,
			title:    note.Title,
		})
	}

	// Add todos as potential targets (if source is not this todo)
	for _, todo := range todos {
		if m.sourceType == "todo" && todo.ID == m.sourceID {
			continue // Skip self
		}
		items = append(items, TargetItem{
			itemType: "todo",
			id:       todo.ID,
			title:    todo.Title,
		})
	}

	m.targetList.SetItems(items)
}

// Update handles messages for the link modal.
func (m *LinkModel) Update(msg tea.Msg) (LinkModel, tea.Cmd) {
	if !m.showModal {
		return *m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.mode {
		case LinkModeViewLinks:
			switch msg.String() {
			case "esc":
				m.Close()
				return *m, nil
			case "c", "n": // Create new link
				m.mode = LinkModeSelectType
				m.typeIndex = 0
				return *m, nil
			case "d": // Delete selected link
				if len(m.linkList.Items()) > 0 {
					if selected, ok := m.linkList.SelectedItem().(LinkItem); ok {
						m.store.DeleteLink(selected.link.ID)
						m.loadLinks()
					}
				}
				return *m, nil
			}
			var cmd tea.Cmd
			m.linkList, cmd = m.linkList.Update(msg)
			return *m, cmd

		case LinkModeSelectType:
			switch msg.String() {
			case "esc":
				m.mode = LinkModeViewLinks
				return *m, nil
			case "up", "k":
				if m.typeIndex > 0 {
					m.typeIndex--
				}
				return *m, nil
			case "down", "j":
				if m.typeIndex < len(m.linkTypes)-1 {
					m.typeIndex++
				}
				return *m, nil
			case "enter":
				m.selectedType = m.linkTypes[m.typeIndex]
				m.mode = LinkModeSelectTarget
				m.loadTargets()
				return *m, nil
			}

		case LinkModeSelectTarget:
			switch msg.String() {
			case "esc":
				m.mode = LinkModeSelectType
				return *m, nil
			case "enter":
				if len(m.targetList.Items()) > 0 {
					if selected, ok := m.targetList.SelectedItem().(TargetItem); ok {
						// Create the link
						link := &models.Link{
							SourceType: m.sourceType,
							SourceID:   m.sourceID,
							TargetType: selected.itemType,
							TargetID:   selected.id,
							LinkType:   m.selectedType,
						}
						m.store.CreateLink(link)
						m.mode = LinkModeViewLinks
						m.loadLinks()
					}
				}
				return *m, nil
			}
			var cmd tea.Cmd
			m.targetList, cmd = m.targetList.Update(msg)
			return *m, cmd
		}
	}

	return *m, nil
}

// View renders the link modal.
func (m *LinkModel) View() string {
	if !m.showModal {
		return ""
	}

	var content string

	switch m.mode {
	case LinkModeViewLinks:
		content = m.viewLinksView()
	case LinkModeSelectType:
		content = m.selectTypeView()
	case LinkModeSelectTarget:
		content = m.selectTargetView()
	}

	return styles.PanelActiveStyle.Render(content)
}

func (m *LinkModel) viewLinksView() string {
	title := styles.TitleStyle.Render(fmt.Sprintf("🔗 Links for: %s", m.sourceTitle))

	var linkContent string
	if len(m.links) == 0 {
		linkContent = styles.SubtitleStyle.Render("No links yet. Press 'c' to create one.")
	} else {
		linkContent = m.linkList.View()
	}

	help := styles.HelpStyle.Render(
		styles.KeyHint("c", "Create link") + " • " +
			styles.KeyHint("d", "Delete") + " • " +
			styles.KeyHint("Esc", "Close"),
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		linkContent,
		"",
		help,
	)
}

func (m *LinkModel) selectTypeView() string {
	title := styles.TitleStyle.Render("🔗 Select Link Type")

	typeDescriptions := map[models.LinkType]string{
		models.LinkTypeRelated:    "General connection between items",
		models.LinkTypeContains:   "Parent/child relationship",
		models.LinkTypeReferences: "One-way citation or reference",
	}

	var typeOptions strings.Builder
	for i, lt := range m.linkTypes {
		prefix := "  "
		style := styles.MenuItemStyle
		if i == m.typeIndex {
			prefix = "▶ "
			style = styles.SelectedItemStyle
		}
		line := fmt.Sprintf("%s%s - %s", prefix, string(lt), typeDescriptions[lt])
		typeOptions.WriteString(style.Render(line))
		typeOptions.WriteString("\n")
	}

	help := styles.HelpStyle.Render(
		styles.KeyHint("↑/↓", "Navigate") + " • " +
			styles.KeyHint("Enter", "Select") + " • " +
			styles.KeyHint("Esc", "Back"),
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		styles.SubtitleStyle.Render(fmt.Sprintf("Creating link from: %s", m.sourceTitle)),
		"",
		typeOptions.String(),
		"",
		help,
	)
}

func (m *LinkModel) selectTargetView() string {
	title := styles.TitleStyle.Render("🔗 Select Target Item")

	subtitle := styles.SubtitleStyle.Render(
		fmt.Sprintf("Link type: %s | From: %s", m.selectedType, m.sourceTitle),
	)

	var targetContent string
	if len(m.targetList.Items()) == 0 {
		targetContent = styles.SubtitleStyle.Render("No items available to link.")
	} else {
		targetContent = m.targetList.View()
	}

	help := styles.HelpStyle.Render(
		styles.KeyHint("↑/↓", "Navigate") + " • " +
			styles.KeyHint("Enter", "Create link") + " • " +
			styles.KeyHint("Esc", "Back"),
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		subtitle,
		"",
		targetContent,
		"",
		help,
	)
}

// LinkItem represents an existing link in the list.
type LinkItem struct {
	link       models.Link
	store      *sqlite.Store
	sourceType string
	sourceID   int64
}

func (l LinkItem) Title() string {
	// Show the "other" item in the link
	var targetType string
	var targetID int64

	if l.link.SourceType == l.sourceType && l.link.SourceID == l.sourceID {
		targetType = l.link.TargetType
		targetID = l.link.TargetID
	} else {
		targetType = l.link.SourceType
		targetID = l.link.SourceID
	}

	// Get the target item's title
	var title string
	if targetType == "note" {
		if note, _ := l.store.GetNote(targetID); note != nil {
			title = note.Title
		} else {
			title = fmt.Sprintf("Note #%d", targetID)
		}
	} else {
		if todo, _ := l.store.GetTodo(targetID); todo != nil {
			title = todo.Title
		} else {
			title = fmt.Sprintf("Todo #%d", targetID)
		}
	}

	icon := "📝"
	if targetType == "todo" {
		icon = "✅"
	}

	return fmt.Sprintf("%s %s [%s]", icon, title, l.link.LinkType)
}

func (l LinkItem) Description() string {
	return fmt.Sprintf("Link type: %s", l.link.LinkType)
}

func (l LinkItem) FilterValue() string {
	return l.Title()
}

// TargetItem represents a potential link target.
type TargetItem struct {
	itemType string
	id       int64
	title    string
}

func (t TargetItem) Title() string {
	icon := "📝"
	if t.itemType == "todo" {
		icon = "✅"
	}
	return fmt.Sprintf("%s %s", icon, t.title)
}

func (t TargetItem) Description() string {
	return fmt.Sprintf("%s #%d", strings.Title(t.itemType), t.id)
}

func (t TargetItem) FilterValue() string {
	return t.title
}
