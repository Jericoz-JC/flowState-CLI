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
//   - n: Notes screen
//   - t: Todos screen
//   - f: Focus sessions
//   - s: Semantic search
//   - h: Home screen
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
		case "n":
			m.currentScreen = ScreenNotes
			m.status = "Notes"
			m.notesScreen.LoadNotes()
		case "t":
			m.currentScreen = ScreenTodos
			m.status = "Todos"
			m.todosScreen.LoadTodos()
		case "f":
			m.currentScreen = ScreenFocus
			m.status = "Focus"
		case "s":
			m.currentScreen = ScreenSearch
			m.status = "Search"
		case "h":
			m.currentScreen = ScreenHome
			m.status = "Home"
		}
	case tea.WindowSizeMsg:
		m.SetSize(msg.Width, msg.Height)
	}

	switch m.currentScreen {
	case ScreenNotes:
		return m.notesScreen.Update(msg)
	case ScreenTodos:
		return m.todosScreen.Update(msg)
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
		fmt.Sprintf(" %s | [n] Notes [t] Todos [f] Focus [s] Search [h] Home [q] Quit ", m.status),
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
//   - Application title and subtitle
//   - List of available screens
//   - Keyboard shortcuts reference
func (m *Model) homeView() string {
	return lipgloss.JoinVertical(
		lipgloss.Center,
		styles.TitleStyle.Render("flowState"),
		styles.SubtitleStyle.Render("Your unified terminal productivity system"),
		"",
		styles.MenuItemStyle.Render("[n] Notes"),
		styles.MenuItemStyle.Render("[t] Todos"),
		styles.MenuItemStyle.Render("[f] Focus Sessions"),
		styles.MenuItemStyle.Render("[s] Semantic Search"),
	)
}

// focusView placeholder for focus session screen.
//
// Phase 4: Focus Sessions (upcoming)
//   - Pomodoro timer
//   - Session tracking
//   - Statistics display
func (m *Model) focusView() string {
	return lipgloss.JoinVertical(
		lipgloss.Center,
		styles.TitleStyle.Render("Focus Session"),
		styles.SubtitleStyle.Render("Coming soon..."),
		"",
		styles.MenuItemStyle.Render("[h] Home"),
	)
}

// searchView placeholder for semantic search screen.
//
// Phase 5: Semantic Search (upcoming)
//   - Natural language query input
//   - Results with similarity scores
//   - Filter by tags/dates
func (m *Model) searchView() string {
	return lipgloss.JoinVertical(
		lipgloss.Center,
		styles.TitleStyle.Render("Semantic Search"),
		styles.SubtitleStyle.Render("Coming soon..."),
		"",
		styles.MenuItemStyle.Render("[h] Home"),
	)
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
