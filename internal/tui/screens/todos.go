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

	"flowState-cli/internal/models"
	"flowState-cli/internal/storage/sqlite"
	"flowState-cli/internal/tui/components"
	"flowState-cli/internal/tui/styles"
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
	showCreate       bool
	editingID        int64 // 0 = creating new, >0 = editing existing
	confirmingDelete bool
	deleteTargetID   int64
	titleInput       components.TextInputModel
	descInput        components.TextAreaModel
	filter           string
}

// NewTodosListModel creates a new todos list screen.
func NewTodosListModel(store *sqlite.Store) TodosListModel {
	items := []list.Item{}
	delegate := list.NewDefaultDelegate()

	l := list.New(items, delegate, 0, 0)
	l.Title = "Todos"
	l.SetShowHelp(true)

	return TodosListModel{
		list:             l,
		store:            store,
		showCreate:       false,
		editingID:        0,
		confirmingDelete: false,
		deleteTargetID:   0,
		titleInput:       components.NewTextInput("Todo title"),
		descInput:        components.NewTextArea("Description (optional)"),
		filter:           "",
	}
}

// Init implements tea.Model.
func (m *TodosListModel) Init() tea.Cmd {
	return func() tea.Msg { return nil }
}

// SetSize updates the list dimensions.
func (m *TodosListModel) SetSize(width, height int) {
	m.list.SetSize(width-4, height-10)
}

// LoadTodos refreshes the todo list from the database.
func (m *TodosListModel) LoadTodos() error {
	todos, err := m.store.ListTodos()
	if err != nil {
		return err
	}

	items := make([]list.Item, 0, len(todos))
	for _, todo := range todos {
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
			case "ctrl+s":
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
			case "esc":
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

		// Handle keys when viewing list
		switch msg.String() {
		case "c":
			m.showCreate = true
			m.titleInput.Focus()
			m.descInput.Blur()
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
		case "d":
			if len(m.list.VisibleItems()) > 0 {
				if selected, ok := m.list.SelectedItem().(TodoItem); ok {
					m.confirmingDelete = true
					m.deleteTargetID = selected.todo.ID
				}
			}
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
		}

		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View renders the todos screen.
//
// Phase 2: Todos
//   - Shows create/edit form when active
//   - Shows todo list otherwise
func (m *TodosListModel) View() string {
	// Delete confirmation dialog
	if m.confirmingDelete {
		confirmDialog := lipgloss.JoinVertical(
			lipgloss.Center,
			styles.TitleStyle.Render("⚠️ Delete Todo?"),
			"",
			styles.SubtitleStyle.Render("This action cannot be undone."),
			"",
			styles.HelpStyle.Render(
				styles.KeyHint("y", "Yes, delete")+" • "+
					styles.KeyHint("n", "No, cancel"),
			),
		)
		return styles.PanelStyle.Render(confirmDialog)
	}

	if m.showCreate {
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
			styles.HelpStyle.Render(
				styles.KeyHint("Tab", "Switch field")+" • "+
					styles.KeyHint("Ctrl+S", "Save")+" • "+
					styles.KeyHint("Esc", "Cancel"),
			),
		)
		return styles.PanelStyle.Render(form)
	}

	// Empty state
	if len(m.list.Items()) == 0 {
		emptyState := lipgloss.JoinVertical(
			lipgloss.Center,
			styles.TitleStyle.Render("✅ Todos"),
			"",
			styles.SubtitleStyle.Render("No todos yet. Add something to get done!"),
			"",
			styles.HelpStyle.Render(
				styles.KeyHint("c", "Create your first todo"),
			),
		)
		return styles.PanelStyle.Render(emptyState)
	}

	return m.list.View()
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
