// Package app implements the main TUI application for flowState-cli.
//
// Phase 1: Core Infrastructure
//   - Bubble Tea program initialization
//   - Screen management (Home, Notes, Todos, Focus, Search)
//   - Key binding navigation
//   - Status bar display
//
// Phase 2: Notes & Todos
//   - Notes screen integration
//   - Todos screen integration
//   - Screen-specific keyboard handlers
//
// Navigation Keys:
//   - Ctrl+N: Notes screen
//   - Ctrl+T: Todos screen
//   - Ctrl+F: Focus sessions
//   - Ctrl+S: Semantic search
//   - Ctrl+H: Home screen / Help
//   - q: Quit application
package app

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"flowState-cli/internal/config"
	"flowState-cli/internal/search"
	"flowState-cli/internal/storage/qdrant"
	"flowState-cli/internal/storage/sqlite"
	"flowState-cli/internal/tui/screens"
	"flowState-cli/internal/tui/styles"
)

// Screen represents the current visible screen.
//
// Phase 1: Core Infrastructure
//   - ScreenHome: Main menu
//   - ScreenNotes: Note management (Phase 2)
//   - ScreenTodos: Todo management (Phase 2)
//   - ScreenFocus: Focus timer (Phase 4)
//   - ScreenSearch: Semantic search (Phase 5)
type Screen int

const (
	ScreenHome Screen = iota
	ScreenNotes
	ScreenTodos
	ScreenFocus
	ScreenSearch
)

// Model is the main application model.
//
// Phase 1: Core Infrastructure
//   - currentScreen: Active screen
//   - store: SQLite database connection
//   - vectorStore: Vector storage for search
//   - notesScreen/todosScreen: Screen models
//
// Phase 2: Notes & Todos
//   - Notes and todos screens are initialized
//   - Screen-specific updates delegated to screen models
type Model struct {
	width         int
	height        int
	currentScreen Screen
	config        *config.Config
	store         *sqlite.Store
	vectorStore   *qdrant.VectorStore
	semantic      *search.SemanticSearch
	notesScreen   *screens.NotesListModel
	todosScreen   *screens.TodosListModel
	status        string
	lastUpdate    time.Time
}

// New creates and initializes the application.
//
// Phase 1: Core Infrastructure
//   - Opens SQLite database connection
//   - Initializes vector store
//   - Creates screen models
//   - Sets initial screen to Home
func New(cfg *config.Config) (*Model, error) {
	store, err := sqlite.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to open store: %w", err)
	}

	vectorStore, err := qdrant.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create vector store: %w", err)
	}

	notesScreen := screens.NewNotesListModel(store)
	todosScreen := screens.NewTodosListModel(store)

	return &Model{
		currentScreen: ScreenHome,
		config:        cfg,
		store:         store,
		vectorStore:   vectorStore,
		notesScreen:   &notesScreen,
		todosScreen:   &todosScreen,
		status:        "Ready",
		lastUpdate:    time.Now(),
	}, nil
}

// SetSize updates the model dimensions when window is resized.
//
// Phase 1: Core Infrastructure
//   - Called by Bubble Tea on window resize events
//   - Propagates size to child screens
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	if m.notesScreen != nil {
		m.notesScreen.SetSize(width, height)
	}
	if m.todosScreen != nil {
		m.todosScreen.SetSize(width, height)
	}
}

// Update handles incoming messages and updates the model.
//
// Phase 1: Core Infrastructure
//   - KeyMsg: Handle navigation keys
//   - WindowSizeMsg: Handle window resize
//
// Phase 2: Notes & Todos
//   - Delegates to notesScreen or todosScreen when active
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "ctrl+n":
			m.currentScreen = ScreenNotes
			m.status = "Notes"
			m.notesScreen.LoadNotes()
		case "ctrl+t":
			m.currentScreen = ScreenTodos
			m.status = "Todos"
			m.todosScreen.LoadTodos()
		case "ctrl+f":
			m.currentScreen = ScreenFocus
			m.status = "Focus"
		case "ctrl+/":
			m.currentScreen = ScreenSearch
			m.status = "Search"
		case "ctrl+h":
			m.currentScreen = ScreenHome
			m.status = "Home"
		}
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
	}

	switch m.currentScreen {
	case ScreenNotes:
		_, cmd := m.notesScreen.Update(msg)
		return m, cmd
	case ScreenTodos:
		_, cmd := m.todosScreen.Update(msg)
		return m, cmd
	}

	return m, nil
}

// View renders the current screen.
//
// Phase 1: Core Infrastructure
//   - Renders status bar at bottom
//   - Renders current screen content
//   - Uses Lip Gloss for styling
//
// Phase 2: Notes & Todos
//   - Notes screen with note list/create/edit
//   - Todos screen with todo list/create/toggle
func (m *Model) View() string {
	if m.width == 0 {
		return "Initializing..."
	}

	var content string
	switch m.currentScreen {
	case ScreenHome:
		content = m.homeView()
	case ScreenNotes:
		content = m.notesScreen.View()
	case ScreenTodos:
		content = m.todosScreen.View()
	case ScreenFocus:
		content = m.focusView()
	case ScreenSearch:
		content = m.searchView()
	default:
		content = m.homeView()
	}

	statusBar := styles.StatusBarStyle.Render(
		fmt.Sprintf(" %s | [Ctrl+N] Notes [Ctrl+T] Todos [Ctrl+F] Focus [Ctrl+/] Search [Ctrl+H] Home [q] Quit ", m.status),
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		content,
		"",
		statusBar,
	)
}

// homeView renders the home screen.
//
// Phase 1: Core Infrastructure
//   - Application title with ASCII art
//   - List of available screens with styled shortcuts
//   - Keyboard shortcuts reference
func (m *Model) homeView() string {
	// ASCII art logo
	logo := styles.LogoStyle.Render(styles.LogoASCII)

	// Subtitle
	subtitle := styles.SubtitleStyle.Render("Your unified terminal productivity system")

	// Menu items with styled shortcuts
	menuItems := lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		styles.MenuItemStyle.Render(styles.KeyHint("Ctrl+N", "Notes")+"   - Capture and organize your thoughts"),
		styles.MenuItemStyle.Render(styles.KeyHint("Ctrl+T", "Todos")+"   - Track your tasks and priorities"),
		styles.MenuItemStyle.Render(styles.KeyHint("Ctrl+F", "Focus")+"   - Pomodoro timer for deep work"),
		styles.MenuItemStyle.Render(styles.KeyHint("Ctrl+/", "Search")+"  - Find anything with semantic search"),
		"",
	)

	// Quick tips
	tips := styles.HelpStyle.Render("Press " + styles.KeyStyle.Render("q") + " to quit • " + styles.KeyStyle.Render("Ctrl+H") + " for help")

	return lipgloss.JoinVertical(
		lipgloss.Center,
		logo,
		subtitle,
		menuItems,
		tips,
	)
}

// focusView placeholder for focus session screen.
//
// Phase 4: Focus Sessions (upcoming)
//   - Pomodoro timer
//   - Session tracking
//   - Statistics display
func (m *Model) focusView() string {
	title := styles.TitleStyle.Render("🍅 Focus Session")

	timer := styles.TimerStyle.Render("25:00")

	progress := styles.ProgressBarStyle.Render("████████████░░░░░░░░")

	help := lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		styles.HelpStyle.Render("This feature is coming soon!"),
		"",
		styles.MenuItemStyle.Render(styles.KeyHint("Ctrl+H", "Return to Home")),
	)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		timer,
		progress,
		help,
	)

	return styles.PanelStyle.Render(content)
}

// searchView placeholder for semantic search screen.
//
// Phase 5: Semantic Search (upcoming)
//   - Natural language query input
//   - Results with similarity scores
//   - Filter by tags/dates
func (m *Model) searchView() string {
	title := styles.TitleStyle.Render("🔍 Semantic Search")

	inputPlaceholder := styles.InputStyle.Render("Type your search query...")

	help := lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		styles.HelpStyle.Render("This feature is coming soon!"),
		"",
		styles.MenuItemStyle.Render(styles.KeyHint("Ctrl+H", "Return to Home")),
	)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		title,
		"",
		inputPlaceholder,
		help,
	)

	return styles.PanelStyle.Render(content)
}

// Init is called once at program start.
//
// Phase 1: Core Infrastructure
//   - Returns nil (no initial command)
func (m *Model) Init() tea.Cmd {
	return nil
}

// Close cleans up resources on exit.
//
// Phase 1: Core Infrastructure
//   - Closes SQLite database
//   - Closes vector store
func (m *Model) Close() error {
	if m.store != nil {
		m.store.Close()
	}
	if m.vectorStore != nil {
		m.vectorStore.Close()
	}
	return nil
}
