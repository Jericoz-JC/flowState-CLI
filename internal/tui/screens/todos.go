// Package screens provides TUI screen implementations for flowState-cli.
//
// Phase 2: Notes & Todos
//   - TodosListModel: Todo management UI
//   - Create, read, update, delete operations
//   - Status tracking and priority levels
//
// Phase 3: Notion-Inspired Overhaul (v0.1.6)
//   - Sort modes: Date↓, Priority, Date↑, Alphabetical
//   - Tag support: Extract #hashtags from description
//   - Priority filter: Cycle through priority levels
//   - Preview mode: View full todo details
//   - Visual cards: Colored status badges, dates
package screens

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Jericoz-JC/flowState-CLI/internal/models"
	"github.com/Jericoz-JC/flowState-CLI/internal/storage/sqlite"
	"github.com/Jericoz-JC/flowState-CLI/internal/tui/components"
	"github.com/Jericoz-JC/flowState-CLI/internal/tui/keymap"
	"github.com/Jericoz-JC/flowState-CLI/internal/tui/styles"
)

// TodoSortMode defines how todos are sorted.
type TodoSortMode int

const (
	// TodoSortByDate sorts by creation date descending (newest first)
	TodoSortByDate TodoSortMode = iota
	// TodoSortByPriority sorts by priority (high first)
	TodoSortByPriority
	// TodoSortByDateAsc sorts by creation date ascending (oldest first)
	TodoSortByDateAsc
	// TodoSortByTitle sorts alphabetically by title
	TodoSortByTitle
	// TodoSortByDueDate sorts by due date (earliest first)
	TodoSortByDueDate
)

// String returns a display name for the sort mode.
func (s TodoSortMode) String() string {
	switch s {
	case TodoSortByDate:
		return "Date↓"
	case TodoSortByPriority:
		return "Priority"
	case TodoSortByDateAsc:
		return "Date↑"
	case TodoSortByTitle:
		return "A-Z"
	case TodoSortByDueDate:
		return "Due Date"
	default:
		return "Date↓"
	}
}

// tagPattern matches #hashtags in text
var tagPattern = regexp.MustCompile(`#(\w+)`)

// extractTagsFromTodo extracts #hashtags from todo title and description.
func extractTagsFromTodo(todo *models.Todo) []string {
	text := todo.Title + " " + todo.Description
	matches := tagPattern.FindAllStringSubmatch(text, -1)
	seen := make(map[string]bool)
	var tags []string
	for _, match := range matches {
		if len(match) > 1 {
			tag := strings.ToLower(match[1])
			if !seen[tag] {
				seen[tag] = true
				tags = append(tags, tag)
			}
		}
	}
	return tags
}

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
// Phase 3: Notion-Inspired Overhaul (v0.1.6)
//   - Sort modes: 's' key cycles through Date↓ → Priority → Date↑ → A-Z → Due Date
//   - Tag support: Extract #hashtags, 't' key filters by tag
//   - Priority filter: 'p' key cycles through priority levels
//   - Preview mode: 'v' key shows full todo details
//
// Keyboard Shortcuts (when viewing list):
//   - c: Create new todo
//   - e: Edit selected todo
//   - d: Delete selected todo
//   - space: Toggle completion status
//   - s: Cycle sort mode
//   - t: Toggle tag filter
//   - p: Cycle priority filter
//   - v: Toggle preview mode
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

	// Phase 3: Notion-inspired features
	sortMode       TodoSortMode           // Current sort mode
	allTags        []string               // All unique tags across todos
	selectedTags   map[string]bool        // Selected tags for filtering
	priorityFilter models.TodoPriority    // Filter by priority: -1 = all, 0-2 = specific
	showPreview    bool                   // Whether preview mode is active
	previewTodo    *models.Todo           // Todo being previewed
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
		descInput:        components.NewTextArea("Description (optional, supports #tags)"),
		header:           components.NewHeader("✅", "Todos"),
		helpBar:          components.NewHelpBar(components.TodosListHints),
		// Phase 3: Notion-inspired features
		sortMode:       TodoSortByDate,
		allTags:        []string{},
		selectedTags:   make(map[string]bool),
		priorityFilter: -1, // -1 = all priorities
		showPreview:    false,
		previewTodo:    nil,
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

	// Collect all unique tags for the tag filter UI
	allTagsSet := make(map[string]bool)
	for _, todo := range todos {
		for _, tag := range extractTagsFromTodo(&todo) {
			allTagsSet[tag] = true
		}
	}
	m.allTags = make([]string, 0, len(allTagsSet))
	for tag := range allTagsSet {
		m.allTags = append(m.allTags, tag)
	}
	sort.Strings(m.allTags)

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

		// Filter by priority (Phase 3)
		if m.priorityFilter >= 0 && todo.Priority != m.priorityFilter {
			continue
		}

		// Filter by selected tags (Phase 3)
		if len(m.selectedTags) > 0 {
			todoTags := extractTagsFromTodo(&todo)
			hasMatchingTag := false
			for _, tag := range todoTags {
				if m.selectedTags[tag] {
					hasMatchingTag = true
					break
				}
			}
			if !hasMatchingTag {
				continue
			}
		}

		filtered = append(filtered, todo)
	}

	// Apply sorting (Phase 3)
	switch m.sortMode {
	case TodoSortByDate:
		sort.Slice(filtered, func(i, j int) bool {
			return filtered[i].CreatedAt.After(filtered[j].CreatedAt)
		})
	case TodoSortByDateAsc:
		sort.Slice(filtered, func(i, j int) bool {
			return filtered[i].CreatedAt.Before(filtered[j].CreatedAt)
		})
	case TodoSortByPriority:
		sort.Slice(filtered, func(i, j int) bool {
			// Higher priority first (3=high, 1=low)
			if filtered[i].Priority != filtered[j].Priority {
				return filtered[i].Priority > filtered[j].Priority
			}
			// Same priority: sort by date descending
			return filtered[i].CreatedAt.After(filtered[j].CreatedAt)
		})
	case TodoSortByTitle:
		sort.Slice(filtered, func(i, j int) bool {
			return strings.ToLower(filtered[i].Title) < strings.ToLower(filtered[j].Title)
		})
	case TodoSortByDueDate:
		sort.Slice(filtered, func(i, j int) bool {
			// Nil due dates go to the end
			if filtered[i].DueDate == nil && filtered[j].DueDate == nil {
				return filtered[i].CreatedAt.After(filtered[j].CreatedAt)
			}
			if filtered[i].DueDate == nil {
				return false
			}
			if filtered[j].DueDate == nil {
				return true
			}
			return filtered[i].DueDate.Before(*filtered[j].DueDate)
		})
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
		// Handle filter input with search-as-you-type
		if m.showFilter {
			switch msg.String() {
			case "enter":
				// Enter closes filter but keeps the filter value
				m.showFilter = false
				m.filterInput.Blur()
				return m, nil
			case "esc":
				// Esc clears filter and closes
				m.showFilter = false
				m.filter = ""
				m.filterInput.SetValue("")
				m.filterInput.Blur()
				m.LoadTodos()
				return m, nil
			default:
				var cmd tea.Cmd
				m.filterInput, cmd = m.filterInput.Update(msg)
				// Search-as-you-type: update filter and reload on every keystroke
				m.filter = m.filterInput.Value()
				m.LoadTodos()
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
				// Only save if title input is focused (allow newlines in description)
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
					return m, nil
				}
				// When description is focused, DON'T return - let Enter pass through
				// to the textarea for newline handling (falls through to input update below)
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

		// Handle preview mode keys first
		if m.showPreview {
			switch msg.String() {
			case "esc", "v", "q":
				m.showPreview = false
				m.previewTodo = nil
				return m, nil
			case "e":
				// Edit from preview
				if m.previewTodo != nil {
					m.showPreview = false
					m.showCreate = true
					m.editingID = m.previewTodo.ID
					m.titleInput.SetValue(m.previewTodo.Title)
					m.descInput.SetValue(m.previewTodo.Description)
					m.titleInput.Focus()
					m.previewTodo = nil
				}
				return m, nil
			case "d":
				// Delete from preview
				if m.previewTodo != nil {
					m.showPreview = false
					m.confirmingDelete = true
					m.deleteTargetID = m.previewTodo.ID
					m.previewTodo = nil
				}
				return m, nil
			}
			return m, nil
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
		case "s":
			// Phase 3: Cycle through sort modes
			m.sortMode = (m.sortMode + 1) % 5 // 5 sort modes total
			m.LoadTodos()
			return m, nil
		case "p":
			// Phase 3: Cycle through priority filters (-1 = all, 0 = low, 1 = med, 2 = high)
			switch m.priorityFilter {
			case -1:
				m.priorityFilter = models.TodoPriorityHigh
			case models.TodoPriorityHigh:
				m.priorityFilter = models.TodoPriorityMedium
			case models.TodoPriorityMedium:
				m.priorityFilter = models.TodoPriorityLow
			case models.TodoPriorityLow:
				m.priorityFilter = -1
			default:
				m.priorityFilter = -1
			}
			m.LoadTodos()
			return m, nil
		case "t":
			// Phase 3: Cycle through tag filters
			if len(m.allTags) == 0 {
				return m, nil // No tags available
			}
			// If no tags selected, select first tag
			if len(m.selectedTags) == 0 {
				m.selectedTags[m.allTags[0]] = true
			} else {
				// Find current tag index and move to next
				currentIdx := -1
				for i, tag := range m.allTags {
					if m.selectedTags[tag] {
						currentIdx = i
						break
					}
				}
				// Clear current selection
				m.selectedTags = make(map[string]bool)
				// Select next tag or clear all
				if currentIdx >= 0 && currentIdx < len(m.allTags)-1 {
					m.selectedTags[m.allTags[currentIdx+1]] = true
				}
				// If we've cycled through all, selectedTags stays empty (show all)
			}
			m.LoadTodos()
			return m, nil
		case "v":
			// Phase 3: Toggle preview mode
			if len(m.list.VisibleItems()) > 0 {
				if selected, ok := m.list.SelectedItem().(TodoItem); ok {
					m.showPreview = true
					m.previewTodo = &selected.todo
				}
			}
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
			m.priorityFilter = -1
			m.selectedTags = make(map[string]bool)
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
//
// Phase 3: Notion-inspired additions
//   - Preview mode for viewing full todo details
//   - Sort and filter indicators in help bar
func (m *TodosListModel) View() string {
	// Phase 3: Preview mode
	if m.showPreview && m.previewTodo != nil {
		return m.renderPreview()
	}

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
		descLabel := styles.SubtitleStyle.Render("Description (supports #tags)")
		if m.titleInput.Focused() {
			titleLabel = styles.SelectedItemStyle.Render("▶ Title")
		} else {
			descLabel = styles.SelectedItemStyle.Render("▶ Description (supports #tags)")
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

	// Update help hints (with platform-appropriate mod key)
	mod := keymap.ModKeyDisplay()

	// Get current status filter display
	var statusDesc string
	switch m.statusFilter {
	case "":
		statusDesc = "All"
	case models.TodoStatusPending:
		statusDesc = "Pend"
	case models.TodoStatusInProgress:
		statusDesc = "Prog"
	case models.TodoStatusCompleted:
		statusDesc = "Done"
	}

	// Get current priority filter display (Phase 3)
	var priorityDesc string
	switch m.priorityFilter {
	case -1:
		priorityDesc = "All"
	case models.TodoPriorityHigh:
		priorityDesc = "High"
	case models.TodoPriorityMedium:
		priorityDesc = "Med"
	case models.TodoPriorityLow:
		priorityDesc = "Low"
	}

	// Get current tag filter display (Phase 3)
	tagDesc := "All"
	for tag := range m.selectedTags {
		tagDesc = "#" + tag
		break
	}

	listHints := []components.HelpHint{
		{Key: "c", Description: "Create", Primary: true},
		{Key: "e", Description: "Edit"},
		{Key: "v", Description: "View"},
		{Key: "Space", Description: "Toggle"},
		{Key: "s", Description: m.sortMode.String()},
		{Key: "f", Description: statusDesc},
		{Key: "p", Description: priorityDesc},
		{Key: "t", Description: tagDesc},
		{Key: mod + "+H", Description: "Home"},
	}
	m.helpBar.SetHints(listHints)

	// Build active filters status line (Phase 3 enhanced)
	var filterParts []string
	if m.filter != "" {
		filterParts = append(filterParts, fmt.Sprintf("search:%q", m.filter))
	}
	if m.statusFilter != "" {
		filterParts = append(filterParts, "status:"+string(m.statusFilter))
	}
	if m.priorityFilter >= 0 {
		filterParts = append(filterParts, "priority:"+priorityDesc)
	}
	if len(m.selectedTags) > 0 {
		filterParts = append(filterParts, "tag:"+tagDesc)
	}

	var filterStatus string
	if len(filterParts) > 0 {
		filterStatusStyle := lipgloss.NewStyle().
			Foreground(styles.CreamYellow).
			Background(styles.SurfaceColor).
			Padding(0, 1)
		filterStatus = filterStatusStyle.Render("🔎 " + strings.Join(filterParts, " • ") + " [" + mod + "+R reset]")
	}

	// Sort indicator
	sortIndicator := lipgloss.NewStyle().
		Foreground(styles.SecondaryColor).
		Render("⬡ Sort: " + m.sortMode.String())

	// Empty state
	if len(m.list.Items()) == 0 {
		emptyMsg := "No todos yet. Add something to get done!"
		if m.filter != "" || m.statusFilter != "" || m.priorityFilter >= 0 || len(m.selectedTags) > 0 {
			emptyMsg = "No todos match your filters. Press [" + mod + "+R] to reset."
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
		sortIndicator,
	)
	if filterStatus != "" {
		content = lipgloss.JoinVertical(lipgloss.Left, content, filterStatus)
	}
	content = lipgloss.JoinVertical(
		lipgloss.Left,
		content,
		"",
		m.list.View(),
		"",
		m.helpBar.View(),
	)
	return content
}

// renderPreview renders the full todo details in preview mode (Phase 3).
func (m *TodosListModel) renderPreview() string {
	todo := m.previewTodo
	if todo == nil {
		return ""
	}

	// Preview hints
	previewHints := []components.HelpHint{
		{Key: "e", Description: "Edit", Primary: true},
		{Key: "d", Description: "Delete"},
		{Key: "Esc", Description: "Close"},
	}
	m.helpBar.SetHints(previewHints)

	// Status badge
	var statusBadge string
	statusStyle := lipgloss.NewStyle().Padding(0, 1).Bold(true)
	switch todo.Status {
	case models.TodoStatusPending:
		statusBadge = statusStyle.Background(styles.CreamYellow).Foreground(lipgloss.Color("#000")).Render("PENDING")
	case models.TodoStatusInProgress:
		statusBadge = statusStyle.Background(styles.SecondaryColor).Foreground(lipgloss.Color("#000")).Render("IN PROGRESS")
	case models.TodoStatusCompleted:
		statusBadge = statusStyle.Background(styles.SuccessColor).Foreground(lipgloss.Color("#000")).Render("COMPLETED")
	}

	// Priority badge
	var priorityBadge string
	priorityStyle := lipgloss.NewStyle().Padding(0, 1)
	switch todo.Priority {
	case models.TodoPriorityHigh:
		priorityBadge = priorityStyle.Background(styles.ErrorColor).Foreground(lipgloss.Color("#fff")).Render("HIGH")
	case models.TodoPriorityMedium:
		priorityBadge = priorityStyle.Background(styles.CreamYellow).Foreground(lipgloss.Color("#000")).Render("MEDIUM")
	case models.TodoPriorityLow:
		priorityBadge = priorityStyle.Background(styles.MutedColor).Foreground(lipgloss.Color("#000")).Render("LOW")
	}

	// Tags
	tags := extractTagsFromTodo(todo)
	var tagsLine string
	if len(tags) > 0 {
		tagStyle := lipgloss.NewStyle().
			Foreground(styles.AccentColor).
			Background(styles.SurfaceColor).
			Padding(0, 1)
		tagStrs := make([]string, len(tags))
		for i, tag := range tags {
			tagStrs[i] = tagStyle.Render("#" + tag)
		}
		tagsLine = strings.Join(tagStrs, " ")
	}

	// Dates
	createdStr := todo.CreatedAt.Format("Jan 2, 2006 3:04 PM")
	var dueStr string
	if todo.DueDate != nil {
		dueStr = todo.DueDate.Format("Jan 2, 2006")
		// Add relative time
		daysUntil := int(time.Until(*todo.DueDate).Hours() / 24)
		if daysUntil < 0 {
			dueStr += fmt.Sprintf(" (overdue by %d days)", -daysUntil)
		} else if daysUntil == 0 {
			dueStr += " (today!)"
		} else if daysUntil == 1 {
			dueStr += " (tomorrow)"
		} else if daysUntil <= 7 {
			dueStr += fmt.Sprintf(" (in %d days)", daysUntil)
		}
	}

	// Build preview content
	titleStyle := lipgloss.NewStyle().
		Foreground(styles.TextColor).
		Bold(true)

	labelStyle := lipgloss.NewStyle().
		Foreground(styles.MutedColor).
		Bold(true)

	descStyle := lipgloss.NewStyle().
		Foreground(styles.TextColor).
		Width(m.width - 10)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		styles.TitleStyle.Render("📋 Todo Preview"),
		"",
		titleStyle.Render(todo.Title),
		"",
		lipgloss.JoinHorizontal(lipgloss.Top, statusBadge, " ", priorityBadge),
		"",
	)

	if tagsLine != "" {
		content = lipgloss.JoinVertical(lipgloss.Left, content, tagsLine, "")
	}

	content = lipgloss.JoinVertical(
		lipgloss.Left,
		content,
		labelStyle.Render("Description"),
	)

	if todo.Description != "" {
		content = lipgloss.JoinVertical(lipgloss.Left, content, descStyle.Render(todo.Description))
	} else {
		content = lipgloss.JoinVertical(lipgloss.Left, content, lipgloss.NewStyle().Foreground(styles.MutedColor).Italic(true).Render("No description"))
	}

	content = lipgloss.JoinVertical(
		lipgloss.Left,
		content,
		"",
		labelStyle.Render("Created"),
		styles.SubtitleStyle.Render(createdStr),
	)

	if dueStr != "" {
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			content,
			"",
			labelStyle.Render("Due Date"),
			styles.SubtitleStyle.Render(dueStr),
		)
	}

	content = lipgloss.JoinVertical(
		lipgloss.Left,
		content,
		"",
		m.helpBar.View(),
	)

	return styles.PanelStyle.Render(content)
}

// TodoItem implements list.Item for displaying todos in the list.
//
// Phase 2: Todos
//   - Title: Shows status indicator and title with priority
//   - Description: Shows description preview
//   - FilterValue: Used for search/filter
//
// Phase 3: Notion-inspired enhancements
//   - Colored status badges
//   - Tag display
//   - Due date with relative time
type TodoItem struct {
	todo models.Todo
}

func (t TodoItem) Title() string {
	// Status indicator with color hint
	status := "○"  // Pending (hollow circle)
	if t.todo.Status == models.TodoStatusCompleted {
		status = "✓" // Completed (checkmark)
	} else if t.todo.Status == models.TodoStatusInProgress {
		status = "◐" // In progress (half circle)
	}

	// Priority indicator
	priority := ""
	if t.todo.Priority == models.TodoPriorityHigh {
		priority = " 🔴"
	} else if t.todo.Priority == models.TodoPriorityLow {
		priority = " 🟢"
	}

	// Due date indicator
	dueIndicator := ""
	if t.todo.DueDate != nil {
		daysUntil := int(time.Until(*t.todo.DueDate).Hours() / 24)
		if daysUntil < 0 {
			dueIndicator = " ⚠️" // Overdue
		} else if daysUntil == 0 {
			dueIndicator = " 📅" // Due today
		} else if daysUntil <= 3 {
			dueIndicator = " ⏰" // Due soon
		}
	}

	return fmt.Sprintf("%s %s%s%s", status, t.todo.Title, priority, dueIndicator)
}

func (t TodoItem) Description() string {
	parts := []string{}

	// Tags (Phase 3)
	tags := extractTagsFromTodo(&t.todo)
	if len(tags) > 0 {
		tagStrs := make([]string, len(tags))
		for i, tag := range tags {
			tagStrs[i] = "#" + tag
		}
		// Show first 3 tags max
		if len(tagStrs) > 3 {
			tagStrs = tagStrs[:3]
			tagStrs = append(tagStrs, "...")
		}
		parts = append(parts, strings.Join(tagStrs, " "))
	}

	// Due date (Phase 3)
	if t.todo.DueDate != nil {
		daysUntil := int(time.Until(*t.todo.DueDate).Hours() / 24)
		var dueStr string
		if daysUntil < 0 {
			dueStr = fmt.Sprintf("Overdue %d days", -daysUntil)
		} else if daysUntil == 0 {
			dueStr = "Due today"
		} else if daysUntil == 1 {
			dueStr = "Due tomorrow"
		} else if daysUntil <= 7 {
			dueStr = fmt.Sprintf("Due in %d days", daysUntil)
		} else {
			dueStr = "Due " + t.todo.DueDate.Format("Jan 2")
		}
		parts = append(parts, dueStr)
	}

	// Description preview
	if t.todo.Description != "" {
		// Remove hashtags from preview (already shown separately)
		preview := tagPattern.ReplaceAllString(t.todo.Description, "")
		preview = strings.TrimSpace(preview)
		if len(preview) > 40 {
			preview = preview[:40] + "..."
		}
		if preview != "" {
			parts = append(parts, preview)
		}
	}

	if len(parts) == 0 {
		return "No description"
	}

	return strings.Join(parts, " • ")
}

func (t TodoItem) FilterValue() string {
	return t.todo.Title + " " + t.todo.Description
}
