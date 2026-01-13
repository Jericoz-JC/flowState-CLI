// Package screens provides TUI screen implementations for flowState-cli.
//
// Phase 2: Notes & Todos
//   - TodosListModel: Todo management UI
//   - Create, read, update, delete operations
//   - Status tracking and priority levels
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

// TodosListModel implements the todos management screen.
//
// Phase 2: Todos
//   - Displays list of all todos sorted by creation date
//   - Shows status indicator, title, and priority
//   - Create new todos with title and optional description
//   - Edit existing todos
//   - Delete todos
//   - Toggle completion with space bar
//   - Visual priority indicators (🔴 high, 🟢 low)
//
// Keyboard Shortcuts (when viewing list):
//   - c: Create new todo
//   - e: Edit selected todo
//   - d: Delete selected todo
//   - space: Toggle completion status
//   - j/down: Move selection down
//   - k/up: Move selection up
//   - esc: Cancel/create mode
//   - enter: Save todo (in create/edit mode)
//
// Status Indicators:
//   - [ ] Pending (default)
//   - [~] In progress
//   - [x] Completed
type TodosListModel struct {
	list             list.Model
	store            *sqlite.Store
	filter           string
	filterInput      components.TextInputModel
	showFilter       bool
	statusFilter     models.TodoStatus // Filter by status: "", "pending", "completed", "in_progress"
	showCreate       bool
	editingID        int64 // 0 = creating new, >0 = editing existing
	confirmingDelete bool
	deleteTargetID   int64
	titleInput       components.TextInputModel
	descInput        components.TextAreaModel
	header           components.Header
	helpBar          components.HelpBar
	width            int
	height           int
}

// NewTodosListModel creates a new todos list screen.
func NewTodosListModel(store *sqlite.Store) TodosListModel {
	items := []list.Item{}
	delegate := list.NewDefaultDelegate()

	l := list.New(items, delegate, 0, 0)
	l.Title = ""
	l.SetShowHelp(false) // We'll use our own help bar
	l.SetShowTitle(false)
	l.SetFilteringEnabled(false) // We handle filtering ourselves

	filterInput := components.NewTextInput("Type to filter...")
	filterInput.Blur()

	return TodosListModel{
		list:             l,
		store:            store,
		filter:           "",
		filterInput:      filterInput,
		showFilter:       false,
		statusFilter:     "",
		showCreate:       false,
		editingID:        0,
		confirmingDelete: false,
		deleteTargetID:   0,
		titleInput:       components.NewTextInput("Todo title"),
		descInput:        components.NewTextArea("Description (optional)"),
		header:           components.NewHeader("✅", "Todos"),
		helpBar:          components.NewHelpBar(components.TodosListHints),
	}
}

// Init implements tea.Model.
func (m *TodosListModel) Init() tea.Cmd {
	return func() tea.Msg { return nil }
}

// SetSize updates the list dimensions.
func (m *TodosListModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.list.SetSize(width-4, height-14) // Account for header and help bar
	m.header.SetWidth(width - 4)
	m.helpBar.SetWidth(width - 4)
}

// GetSelectedTodo returns the currently selected todo, or nil if none selected.
func (m *TodosListModel) GetSelectedTodo() *models.Todo {
	if len(m.list.Items()) == 0 {
		return nil
	}
	if selected, ok := m.list.SelectedItem().(TodoItem); ok {
		return &selected.todo
	}
	return nil
}

// LoadTodos refreshes the todo list from the database.
func (m *TodosListModel) LoadTodos() error {
	todos, err := m.store.ListTodos()
	if err != nil {
		return err
	}

	// Apply filters
	filtered := make([]models.Todo, 0)
	for _, todo := range todos {
		// Filter by search text
		if m.filter != "" {
			searchText := strings.ToLower(m.filter)
			titleMatch := strings.Contains(strings.ToLower(todo.Title), searchText)
			descMatch := strings.Contains(strings.ToLower(todo.Description), searchText)
			if !titleMatch && !descMatch {
				continue
			}
		}

		// Filter by status
		if m.statusFilter != "" && todo.Status != m.statusFilter {
			continue
		}

		filtered = append(filtered, todo)
	}

	items := make([]list.Item, 0, len(filtered))
	for _, todo := range filtered {
		items = append(items, TodoItem{todo: todo})
	}

	m.list.SetItems(items)
	return nil
}

// Update handles messages for the todos screen.
//
// Phase 2: Todos
//   - Key bindings for navigation
//   - Create/edit/delete operations
//   - Status toggle with space bar
//   - Form input handling
//   - Tab to switch between fields
func (m *TodosListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				m.LoadTodos()
				return m, nil
			case "esc":
				m.showFilter = false
				m.filter = ""
				m.filterInput.SetValue("")
				m.filterInput.Blur()
				m.LoadTodos()
				return m, nil
			default:
				var cmd tea.Cmd
				m.filterInput, cmd = m.filterInput.Update(msg)
				cmds = append(cmds, cmd)
				return m, tea.Batch(cmds...)
			}
		}

		// Handle delete confirmation dialog
		if m.confirmingDelete {
			switch msg.String() {
			case "y", "Y":
				m.store.DeleteTodo(m.deleteTargetID)
				m.confirmingDelete = false
				m.deleteTargetID = 0
				m.LoadTodos()
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
				// Toggle focus between title and description
				if m.titleInput.Focused() {
					m.titleInput.Blur()
					m.descInput.Focus()
				} else {
					m.descInput.Blur()
					m.titleInput.Focus()
				}
				return m, nil
			case "enter":
				// Only save if title input is focused
				if m.titleInput.Focused() {
					title := strings.TrimSpace(m.titleInput.Value())
					desc := strings.TrimSpace(m.descInput.Value())
					if title != "" {
						if m.editingID > 0 {
							// Update existing todo - fetch to preserve other fields
							existing, err := m.store.GetTodo(m.editingID)
							if err != nil || existing == nil {
								return m, nil
							}
							existing.Title = title
							existing.Description = desc
							if err := m.store.UpdateTodo(existing); err != nil {
								return m, nil
							}
						} else {
							// Create new todo
							todo := &models.Todo{
								Title:       title,
								Description: desc,
								Status:      models.TodoStatusPending,
								Priority:    models.TodoPriorityMedium,
							}
							if err := m.store.CreateTodo(todo); err != nil {
								return m, nil
							}
						}
						m.showCreate = false
						m.editingID = 0
						m.titleInput.SetValue("")
						m.descInput.SetValue("")
						m.LoadTodos()
					}
				}
				return m, nil
			}

			// Check for cross-platform save shortcut
			if keymap.IsModS(msg) {
				// Alternative save shortcut
				title := strings.TrimSpace(m.titleInput.Value())
				desc := strings.TrimSpace(m.descInput.Value())
				if title != "" {
					if m.editingID > 0 {
						// Update existing todo - fetch to preserve other fields
						existing, err := m.store.GetTodo(m.editingID)
						if err != nil || existing == nil {
							return m, nil
						}
						existing.Title = title
						existing.Description = desc
						if err := m.store.UpdateTodo(existing); err != nil {
							return m, nil
						}
					} else {
						// Create new todo
						todo := &models.Todo{
							Title:       title,
							Description: desc,
							Status:      models.TodoStatusPending,
							Priority:    models.TodoPriorityMedium,
						}
						if err := m.store.CreateTodo(todo); err != nil {
							return m, nil
						}
					}
					m.showCreate = false
					m.editingID = 0
					m.titleInput.SetValue("")
					m.descInput.SetValue("")
					m.LoadTodos()
				}
				return m, nil
			}

			if msg.String() == "esc" {
				m.showCreate = false
				m.editingID = 0
				m.titleInput.SetValue("")
				m.descInput.SetValue("")
				return m, nil
			}

			// Update the focused input
			var cmd tea.Cmd
			if m.titleInput.Focused() {
				m.titleInput, cmd = m.titleInput.Update(msg)
			} else {
				m.descInput, cmd = m.descInput.Update(msg)
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
		case "f":
			// Cycle through status filters: all -> pending -> in_progress -> completed -> all
			switch m.statusFilter {
			case "":
				m.statusFilter = models.TodoStatusPending
			case models.TodoStatusPending:
				m.statusFilter = models.TodoStatusInProgress
			case models.TodoStatusInProgress:
				m.statusFilter = models.TodoStatusCompleted
			case models.TodoStatusCompleted:
				m.statusFilter = ""
			default:
				m.statusFilter = ""
			}
			m.LoadTodos()
			return m, nil
		case "c":
			m.showCreate = true
			m.editingID = 0
			m.titleInput.SetValue("")
			m.descInput.SetValue("")
			m.titleInput.Focus()
			m.descInput.Blur()
			return m, nil // Return early to prevent list from processing
		case "e":
			if len(m.list.VisibleItems()) > 0 {
				if selected, ok := m.list.SelectedItem().(TodoItem); ok {
					m.showCreate = true
					m.editingID = selected.todo.ID
					m.titleInput.SetValue(selected.todo.Title)
					m.descInput.SetValue(selected.todo.Description)
					m.titleInput.Focus()
				}
			}
			return m, nil
		case "d":
			if len(m.list.VisibleItems()) > 0 {
				if selected, ok := m.list.SelectedItem().(TodoItem); ok {
					m.confirmingDelete = true
					m.deleteTargetID = selected.todo.ID
				}
			}
			return m, nil
		case " ":
			if len(m.list.VisibleItems()) > 0 {
				if selected, ok := m.list.SelectedItem().(TodoItem); ok {
					if selected.todo.Status == models.TodoStatusCompleted {
						selected.todo.Status = models.TodoStatusPending
					} else {
						selected.todo.Status = models.TodoStatusCompleted
					}
					m.store.UpdateTodo(&selected.todo)
					m.LoadTodos()
				}
			}
			return m, nil
		}

		// Check for cross-platform reset shortcut
		if keymap.IsModR(msg) {
			// Reset all filters
			m.filter = ""
			m.statusFilter = ""
			m.LoadTodos()
			return m, nil
		}

		// Pass other keys to list for navigation (j/k, up/down, etc.)
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View renders the todos screen.
//
// Phase 4: UX Overhaul
//   - Header with title and item count
//   - Context-sensitive help bar
//   - Shows create/edit form when active
//   - Filter input for searching
func (m *TodosListModel) View() string {
	// Filter input mode
	if m.showFilter {
		filterHints := []components.HelpHint{
			{Key: "Enter", Description: "Apply", Primary: true},
			{Key: "Esc", Description: "Cancel"},
		}
		m.helpBar.SetHints(filterHints)

		filterLabel := styles.TitleStyle.Render("🔍 Filter Todos")
		filterHelp := styles.SubtitleStyle.Render("Type to search by title or description")

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
			styles.TitleStyle.Render("⚠️ Delete Todo?"),
			"",
			styles.SubtitleStyle.Render("This action cannot be undone."),
			"",
			m.helpBar.View(),
		)
		return styles.PanelStyle.Render(confirmDialog)
	}

	if m.showCreate {
		m.helpBar.SetHints(components.TodosEditHints)

		// Show which field is focused
		titleLabel := styles.SubtitleStyle.Render("Title")
		descLabel := styles.SubtitleStyle.Render("Description")
		if m.titleInput.Focused() {
			titleLabel = styles.SelectedItemStyle.Render("▶ Title")
		} else {
			descLabel = styles.SelectedItemStyle.Render("▶ Description")
		}

		// Dynamic title for create vs edit
		formTitle := "✅ Create Todo"
		if m.editingID > 0 {
			formTitle = "✅ Edit Todo"
		}

		form := lipgloss.JoinVertical(
			lipgloss.Left,
			styles.TitleStyle.Render(formTitle),
			"",
			titleLabel,
			m.titleInput.View(),
			"",
			descLabel,
			m.descInput.View(),
			"",
			m.helpBar.View(),
		)
		return styles.PanelStyle.Render(form)
	}

	// Update header with item count
	m.header.SetItemCount(len(m.list.Items()))

	// Update help hints to include filter (with platform-appropriate mod key)
	mod := keymap.ModKeyDisplay()
	listHints := []components.HelpHint{
		{Key: "c", Description: "Create", Primary: true},
		{Key: "e", Description: "Edit"},
		{Key: "d", Description: "Delete"},
		{Key: "Space", Description: "Toggle"},
		{Key: "/", Description: "Filter"},
		{Key: "f", Description: "Status Filter"},
		{Key: mod + "+L", Description: "Link"},
		{Key: mod + "+H", Description: "Home"},
	}
	m.helpBar.SetHints(listHints)

	// Show active filters
	var filterStatus string
	if m.filter != "" || m.statusFilter != "" {
		filterParts := []string{}
		if m.filter != "" {
			filterParts = append(filterParts, fmt.Sprintf("search:%q", m.filter))
		}
		if m.statusFilter != "" {
			filterParts = append(filterParts, "status:"+string(m.statusFilter))
		}
		filterStatusStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F9E2AF")).
			Background(lipgloss.Color("#2E2E3E")).
			Padding(0, 1)
		filterStatus = filterStatusStyle.Render("🔎 Filtering: " + strings.Join(filterParts, ", ") + " [Ctrl+R to reset]")
	}

	// Empty state
	if len(m.list.Items()) == 0 {
		emptyMsg := "No todos yet. Add something to get done!"
		if m.filter != "" || m.statusFilter != "" {
			emptyMsg = "No todos match your filters. Press [Ctrl+R] to reset."
		}
		emptyState := lipgloss.JoinVertical(
			lipgloss.Left,
			m.header.View(),
			"",
			styles.SubtitleStyle.Render(emptyMsg),
			"",
			styles.HelpStyle.Render("Press [c] to create your first todo"),
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

// TodoItem implements list.Item for displaying todos in the list.
//
// Phase 2: Todos
//   - Title: Shows status indicator and title with priority
//   - Description: Shows description preview
//   - FilterValue: Used for search/filter
type TodoItem struct {
	todo models.Todo
}

func (t TodoItem) Title() string {
	status := "[ ]"
	if t.todo.Status == models.TodoStatusCompleted {
		status = "[x]"
	} else if t.todo.Status == models.TodoStatusInProgress {
		status = "[~]"
	}

	priority := ""
	if t.todo.Priority == models.TodoPriorityHigh {
		priority = " 🔴"
	} else if t.todo.Priority == models.TodoPriorityLow {
		priority = " 🟢"
	}

	return fmt.Sprintf("%s %s%s", status, t.todo.Title, priority)
}

func (t TodoItem) Description() string {
	if t.todo.Description == "" {
		return "No description"
	}
	preview := t.todo.Description
	if len(preview) > 50 {
		preview = preview[:50] + "..."
	}
	return preview
}

func (t TodoItem) FilterValue() string {
	return t.todo.Title + " " + t.todo.Description
}
