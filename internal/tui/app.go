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

	"github.com/Jericoz-JC/flowState-CLI/internal/config"
	embeddings "github.com/Jericoz-JC/flowState-CLI/internal/embeddings"
	"github.com/Jericoz-JC/flowState-CLI/internal/search"
	"github.com/Jericoz-JC/flowState-CLI/internal/storage/sqlite"
	"github.com/Jericoz-JC/flowState-CLI/internal/tui/keymap"
	"github.com/Jericoz-JC/flowState-CLI/internal/tui/screens"
	"github.com/Jericoz-JC/flowState-CLI/internal/tui/styles"
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
	ScreenMindMap
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
//
// Phase 3: Linking System
//   - linkScreen: Modal for creating/viewing links
//   - Ctrl+L opens link modal for selected item
//
// Phase 4: UX Overhaul
//   - quickCaptureScreen: Global quick note capture via Ctrl+X
//
// Phase 5: Focus Sessions
//   - focusScreen: Pomodoro-style focus timer with session tracking
type Model struct {
	width              int
	height             int
	currentScreen      Screen
	config             *config.Config
	store              *sqlite.Store
	embedder           *embeddings.Embedder
	semantic           *search.SemanticSearch
	notesScreen        *screens.NotesListModel
	todosScreen        *screens.TodosListModel
	focusScreen        *screens.FocusModel
	searchScreen       *screens.SearchModel
	mindMapScreen      *screens.MindMapModel
	linkScreen         *screens.LinkModel
	quickCaptureScreen *screens.QuickCaptureModel
	showHelpModal      bool
	status             string
	lastUpdate         time.Time
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

	embedder, err := embeddings.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create embedder: %w", err)
	}

	semantic := search.New(embedder, store)
	// Best-effort initial indexing (can be re-run later).
	_ = semantic.IndexAllNotes()

	notesScreen := screens.NewNotesListModel(store)
	todosScreen := screens.NewTodosListModel(store)
	focusScreen := screens.NewFocusModel(store)
	linkScreen := screens.NewLinkModel(store)
	quickCaptureScreen := screens.NewQuickCaptureModel(store)
	searchScreen := screens.NewSearchModel(store, semantic)
	mindMapScreen := screens.NewMindMapModel(store)

	return &Model{
		currentScreen:      ScreenHome,
		config:             cfg,
		store:              store,
		embedder:           embedder,
		semantic:           semantic,
		notesScreen:        &notesScreen,
		todosScreen:        &todosScreen,
		focusScreen:        &focusScreen,
		searchScreen:       &searchScreen,
		mindMapScreen:      &mindMapScreen,
		linkScreen:         &linkScreen,
		quickCaptureScreen: &quickCaptureScreen,
		showHelpModal:      false,
		status:             "Ready",
		lastUpdate:         time.Now(),
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
	if m.linkScreen != nil {
		m.linkScreen.SetSize(width, height)
	}
	if m.quickCaptureScreen != nil {
		m.quickCaptureScreen.SetSize(width, height)
	}
	if m.focusScreen != nil {
		m.focusScreen.SetSize(width, height)
	}
	if m.searchScreen != nil {
		m.searchScreen.SetSize(width, height)
	}
	if m.mindMapScreen != nil {
		m.mindMapScreen.SetSize(width, height)
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
	// Help modal has highest priority when open.
	if m.showHelpModal {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "?", "esc", "q":
				m.showHelpModal = false
				return m, nil
			}
		}
		return m, nil
	}

	// Handle quick capture modal if open
	if m.quickCaptureScreen != nil && m.quickCaptureScreen.IsOpen() {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			updatedQC, cmd := m.quickCaptureScreen.Update(msg)
			m.quickCaptureScreen = &updatedQC
			// If closed, reload notes in case we're on that screen
			if !m.quickCaptureScreen.IsOpen() {
				m.status = "Ready"
				if m.currentScreen == ScreenNotes {
					m.notesScreen.LoadNotes()
				}
			}
			return m, cmd
		}
	}

	// Handle link modal if open
	if m.linkScreen != nil && m.linkScreen.IsOpen() {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			updatedLink, cmd := m.linkScreen.Update(msg)
			m.linkScreen = &updatedLink
			return m, cmd
		}
	}

	switch msg := msg.(type) {
	case screens.OpenNoteMsg:
		// Open the note from search results by navigating to Notes and selecting it.
		m.currentScreen = ScreenNotes
		m.status = "Notes"
		if m.notesScreen != nil {
			_ = m.notesScreen.LoadNotes()
			m.notesScreen.SelectNoteByID(msg.NoteID)
		}
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "?":
			m.showHelpModal = true
			return m, nil
		}

		// Use cross-platform key bindings
		// IMPORTANT: Return early after handling global shortcuts to prevent
		// the key event from being passed to screen components (which might consume it)
		if keymap.IsModH(msg) {
			// Ctrl+H: Go Home - highest priority navigation
			m.currentScreen = ScreenHome
			m.status = "Home"
			return m, nil
		} else if keymap.IsModX(msg) {
			// Open quick capture modal from anywhere
			if m.quickCaptureScreen != nil {
				m.quickCaptureScreen.Open()
				m.status = "Quick Capture"
			}
			return m, nil
		} else if keymap.IsModN(msg) {
			m.currentScreen = ScreenNotes
			m.status = "Notes"
			m.notesScreen.LoadNotes()
			return m, nil
		} else if keymap.IsModT(msg) {
			m.currentScreen = ScreenTodos
			m.status = "Todos"
			m.todosScreen.LoadTodos()
			return m, nil
		} else if keymap.IsModF(msg) {
			m.currentScreen = ScreenFocus
			m.status = "Focus"
			if m.focusScreen != nil {
				m.focusScreen.LoadHistory()
			}
			return m, nil
		} else if keymap.IsModSlash(msg) {
			m.currentScreen = ScreenSearch
			m.status = "Search"
			return m, nil
		} else if keymap.IsModG(msg) {
			m.currentScreen = ScreenMindMap
			m.status = "Mind Map"
			if m.mindMapScreen != nil {
				_ = m.mindMapScreen.LoadGraph()
			}
			return m, nil
		} else if keymap.IsModL(msg) {
			// Open link modal for currently selected item
			if m.currentScreen == ScreenNotes && m.notesScreen != nil {
				if selected := m.notesScreen.GetSelectedNote(); selected != nil {
					m.linkScreen.Open("note", selected.ID, selected.Title)
					m.status = "Links"
				}
			} else if m.currentScreen == ScreenTodos && m.todosScreen != nil {
				if selected := m.todosScreen.GetSelectedTodo(); selected != nil {
					m.linkScreen.Open("todo", selected.ID, selected.Title)
					m.status = "Links"
				}
			}
			return m, nil
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
	case ScreenFocus:
		if m.focusScreen != nil {
			updatedFocus, cmd := m.focusScreen.Update(msg)
			m.focusScreen = &updatedFocus
			return m, cmd
		}
	case ScreenSearch:
		if m.searchScreen != nil {
			updatedSearch, cmd := m.searchScreen.Update(msg)
			m.searchScreen = &updatedSearch
			return m, cmd
		}
	case ScreenMindMap:
		if m.mindMapScreen != nil {
			updatedMM, cmd := m.mindMapScreen.Update(msg)
			m.mindMapScreen = &updatedMM
			return m, cmd
		}
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
		if m.focusScreen != nil {
			content = m.focusScreen.View()
		} else {
			content = m.focusView()
		}
	case ScreenSearch:
		if m.searchScreen != nil {
			content = m.searchScreen.View()
		} else {
			content = "Search unavailable"
		}
	case ScreenMindMap:
		if m.mindMapScreen != nil {
			content = m.mindMapScreen.View()
		} else {
			content = "Mind map unavailable"
		}
	default:
		content = m.homeView()
	}

	// Overlay link modal if open
	if m.linkScreen != nil && m.linkScreen.IsOpen() {
		content = m.linkScreen.View()
	}

	// Overlay quick capture modal if open
	if m.quickCaptureScreen != nil && m.quickCaptureScreen.IsOpen() {
		content = m.quickCaptureScreen.View()
	}

	// Overlay help modal last (highest priority)
	if m.showHelpModal {
		content = m.helpModalView()
	}

	// Build status bar with platform-appropriate shortcuts
	mod := keymap.ModKeyDisplay()
	statusBar := styles.StatusBarStyle.Render(
		fmt.Sprintf(" %s | [%s+X] Capture [%s+N] Notes [%s+T] Todos [%s+G] Map [%s+L] Link [%s+H] Home [q] Quit ",
			m.status, mod, mod, mod, mod, mod, mod),
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		content,
		"",
		statusBar,
	)
}

func (m *Model) helpModalView() string {
	mod := keymap.ModKeyDisplay()
	box := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(styles.AccentColor). // Hot pink border
		Padding(1, 2).
		Width(52)

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(styles.SecondaryColor) // Neon cyan
	keyStyle := lipgloss.NewStyle().Foreground(styles.AccentColor)                 // Hot pink
	descStyle := lipgloss.NewStyle().Foreground(styles.TextColor)                  // Off-white
	mutedStyle := lipgloss.NewStyle().Foreground(styles.MutedColor)                // Pale blue

	title := titleStyle.Render(styles.DecoStar + " Keyboard Shortcuts " + styles.DecoStar)
	lines := []string{
		title,
		"",
		keyStyle.Render(mod+"+X") + descStyle.Render("  Quick Capture"),
		keyStyle.Render(mod+"+N") + descStyle.Render("  Notes"),
		keyStyle.Render(mod+"+T") + descStyle.Render("  Todos"),
		keyStyle.Render(mod+"+F") + descStyle.Render("  Focus"),
		keyStyle.Render(mod+"+/") + descStyle.Render("  Search"),
		keyStyle.Render(mod+"+G") + descStyle.Render("  Mind Map"),
		keyStyle.Render(mod+"+L") + descStyle.Render("  Links"),
		keyStyle.Render(mod+"+H") + descStyle.Render("  Home"),
		"",
		keyStyle.Render("q") + descStyle.Render("      Quit"),
		keyStyle.Render("?") + descStyle.Render("      Toggle this help"),
		"",
		mutedStyle.Render("Press Esc or ? to close"),
	}
	content := box.Render(lipgloss.JoinVertical(lipgloss.Left, lines...))

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
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
	return nil
}
