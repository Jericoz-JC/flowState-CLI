// Package screens provides TUI screen implementations for flowState-cli.
//
// Phase 5: Focus Sessions
//   - FocusModel: Pomodoro-style focus timer UI
//   - Work sessions (default 25 min) and break sessions (5 min)
//   - Session tracking with statistics
//   - Streak calculation for consecutive focus days
package screens

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Jericoz-JC/flowState-CLI/internal/models"
	"github.com/Jericoz-JC/flowState-CLI/internal/storage/sqlite"
	"github.com/Jericoz-JC/flowState-CLI/internal/tui/components"
	"github.com/Jericoz-JC/flowState-CLI/internal/tui/styles"
)

// FocusMode represents the current state of the focus timer.
type FocusMode int

const (
	FocusModeIdle FocusMode = iota
	FocusModeRunning
	FocusModePaused
	FocusModeBreak
	FocusModeHistory
	FocusModeDuration // Duration picker
)

// Duration presets in minutes
var (
	WorkDurations  = []int{15, 25, 45, 60}
	BreakDurations = []int{5, 10, 15}
)

// tickMsg is sent every second when timer is running.
type tickMsg time.Time

// clearFeedbackMsg is sent to clear the "Saved" indicator after a delay.
type clearFeedbackMsg struct{}

// autoExitDurationMsg is sent to auto-exit the duration picker after selection.
type autoExitDurationMsg struct {
	sequence int // Used to cancel stale auto-exit timers
}

// FocusModel implements the focus session screen.
//
// Phase 5: Focus Sessions
//   - Pomodoro timer with configurable work/break durations
//   - Visual progress bar and large timer display
//   - Session history and statistics
//   - Streak tracking for consecutive focus days
//
// Keyboard Shortcuts:
//   - s: Start timer (or resume if paused)
//   - p: Pause timer
//   - c: Cancel current session
//   - h: Toggle history view
//   - d: Change duration (opens duration picker)
//   - b: Skip to break / Skip break
//   - Esc: Return to idle / Cancel action
type FocusModel struct {
	store          *sqlite.Store
	mode           FocusMode
	workDuration   int           // Work duration in minutes
	breakDuration  int           // Break duration in minutes
	remaining      time.Duration // Time remaining
	totalDuration  time.Duration // Total duration for progress calculation
	startTime      time.Time     // When current session started
	currentSession *models.FocusSession
	sessions       []models.FocusSession
	sessionList    list.Model
	stats          *sqlite.SessionStats
	header         components.Header
	helpBar        components.HelpBar
	width          int
	height         int
	// Duration picker state
	durationIndex       int    // Currently selected duration preset
	selectingWork       bool   // true = selecting work duration, false = break duration
	durationJustChanged bool   // Show "Saved" indicator briefly
	lastChangedField    string // "work" or "break" - which field was just changed
	autoExitSequence    int    // Sequence number for auto-exit timer cancellation
}

// NewFocusModel creates a new focus session screen.
func NewFocusModel(store *sqlite.Store) FocusModel {
	items := []list.Item{}
	delegate := list.NewDefaultDelegate()

	l := list.New(items, delegate, 0, 0)
	l.Title = ""
	l.SetShowHelp(false)
	l.SetShowTitle(false)
	l.SetFilteringEnabled(false)

	return FocusModel{
		store:         store,
		mode:          FocusModeIdle,
		workDuration:  25, // Default Pomodoro duration
		breakDuration: 5,
		remaining:     25 * time.Minute,
		totalDuration: 25 * time.Minute,
		sessionList:   l,
		header:        components.NewHeader("🍅", "Focus Sessions"),
		helpBar:       components.NewHelpBar(components.FocusIdleHints),
	}
}

// Init implements tea.Model.
func (m *FocusModel) Init() tea.Cmd {
	return nil
}

// SetSize updates the model dimensions.
func (m *FocusModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.sessionList.SetSize(width-4, height-14)
	m.header.SetWidth(width - 4)
	m.helpBar.SetWidth(width - 4)
}

// LoadHistory loads session history from the database.
func (m *FocusModel) LoadHistory() error {
	sessions, err := m.store.ListSessions()
	if err != nil {
		return err
	}
	m.sessions = sessions

	items := make([]list.Item, 0, len(sessions))
	for _, session := range sessions {
		items = append(items, SessionItem{session: session})
	}
	m.sessionList.SetItems(items)

	// Load stats
	stats, err := m.store.GetSessionStats()
	if err != nil {
		return err
	}
	m.stats = stats

	return nil
}

// tickCmd returns a command that sends a tick every second.
func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Update handles messages for the focus screen.
func (m *FocusModel) Update(msg tea.Msg) (FocusModel, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tickMsg:
		if m.mode == FocusModeRunning || m.mode == FocusModeBreak {
			m.remaining -= time.Second
			if m.remaining <= 0 {
				return m.handleTimerComplete()
			}
			cmds = append(cmds, tickCmd())
		}

	case clearFeedbackMsg:
		// Clear the "Saved" indicator
		m.durationJustChanged = false
		m.lastChangedField = ""
		return *m, nil

	case autoExitDurationMsg:
		// Auto-exit duration picker if sequence matches (not cancelled by new input)
		if m.mode == FocusModeDuration && msg.sequence == m.autoExitSequence {
			m.mode = FocusModeIdle
			m.durationJustChanged = false
			m.lastChangedField = ""
		}
		return *m, nil

	case tea.KeyMsg:
		switch m.mode {
		case FocusModeDuration:
			return m.handleDurationInput(msg)
		case FocusModeHistory:
			return m.handleHistoryInput(msg)
		default:
			return m.handleTimerInput(msg)
		}
	}

	return *m, tea.Batch(cmds...)
}

// handleTimerComplete handles when the timer reaches zero.
func (m *FocusModel) handleTimerComplete() (FocusModel, tea.Cmd) {
	if m.mode == FocusModeRunning {
		// Work session completed - NOW save to database
		now := time.Now()
		if m.currentSession != nil {
			m.currentSession.EndTime = &now
			m.currentSession.Status = models.SessionStatusCompleted
			// Create the session in DB only on completion
			if err := m.store.CreateSession(m.currentSession); err != nil {
				// Log error but continue (session tracking is best-effort)
			}
		}

		// Start break
		m.mode = FocusModeBreak
		m.remaining = time.Duration(m.breakDuration) * time.Minute
		m.totalDuration = m.remaining
		m.currentSession = nil

		return *m, tickCmd()
	} else if m.mode == FocusModeBreak {
		// Break completed - return to idle
		m.mode = FocusModeIdle
		m.remaining = time.Duration(m.workDuration) * time.Minute
		m.totalDuration = m.remaining
		m.LoadHistory() // Refresh stats

		return *m, nil
	}

	return *m, nil
}

// handleTimerInput handles keyboard input for timer modes (idle, running, paused, break).
func (m *FocusModel) handleTimerInput(msg tea.KeyMsg) (FocusModel, tea.Cmd) {
	switch msg.String() {
	case "s":
		if m.mode == FocusModeIdle || m.mode == FocusModePaused {
			// Start or resume timer
			if m.mode == FocusModeIdle {
				// Create in-memory session for tracking (NOT saved to DB yet)
				// Session will only be saved when completed successfully
				m.currentSession = &models.FocusSession{
					StartTime: time.Now(),
					Duration:  m.workDuration * 60, // Store in seconds
					Status:    models.SessionStatusRunning,
				}
				m.remaining = time.Duration(m.workDuration) * time.Minute
				m.totalDuration = m.remaining
				m.startTime = time.Now()
			}
			m.mode = FocusModeRunning
			return *m, tickCmd()
		}

	case "p":
		if m.mode == FocusModeRunning {
			m.mode = FocusModePaused
			return *m, nil
		}

	case "c":
		if m.mode == FocusModeRunning || m.mode == FocusModePaused || m.mode == FocusModeBreak {
			// Cancel current session - just discard, don't save to DB
			// (cancelled sessions are not worth tracking)
			m.currentSession = nil
			m.mode = FocusModeIdle
			m.remaining = time.Duration(m.workDuration) * time.Minute
			m.totalDuration = m.remaining
			m.LoadHistory()
			return *m, nil
		}

	case "h":
		if m.mode == FocusModeIdle {
			m.mode = FocusModeHistory
			m.LoadHistory()
			return *m, nil
		}

	case "d":
		if m.mode == FocusModeIdle {
			m.mode = FocusModeDuration
			m.selectingWork = true
			m.durationIndex = findDurationIndex(m.workDuration, WorkDurations)
			return *m, nil
		}

	case "b":
		if m.mode == FocusModeRunning {
			// Skip to break (complete current session early)
			now := time.Now()
			if m.currentSession != nil {
				m.currentSession.EndTime = &now
				m.currentSession.Status = models.SessionStatusCompleted
				// Save session to DB on early completion
				m.store.CreateSession(m.currentSession)
				m.currentSession = nil
			}
			m.mode = FocusModeBreak
			m.remaining = time.Duration(m.breakDuration) * time.Minute
			m.totalDuration = m.remaining
			return *m, tickCmd()
		} else if m.mode == FocusModeBreak {
			// Skip break
			m.mode = FocusModeIdle
			m.remaining = time.Duration(m.workDuration) * time.Minute
			m.totalDuration = m.remaining
			m.LoadHistory()
			return *m, nil
		}

	case "esc":
		if m.mode == FocusModeBreak {
			// Allow skipping break with Esc
			m.mode = FocusModeIdle
			m.remaining = time.Duration(m.workDuration) * time.Minute
			m.totalDuration = m.remaining
			m.LoadHistory()
			return *m, nil
		}
	}

	return *m, nil
}

// handleDurationInput handles keyboard input for duration picker.
// UX: Arrow keys update values immediately (live preview) with visual feedback,
// Tab switches fields, Enter confirms all and exits.
func (m *FocusModel) handleDurationInput(msg tea.KeyMsg) (FocusModel, tea.Cmd) {
	durations := WorkDurations
	if !m.selectingWork {
		durations = BreakDurations
	}

	switch msg.String() {
	case "left", "h":
		if m.durationIndex > 0 {
			m.durationIndex--
			// Live update: immediately apply and show feedback
			cmd := m.applySelectedDuration(durations)
			return *m, cmd
		}
	case "right", "l":
		if m.durationIndex < len(durations)-1 {
			m.durationIndex++
			// Live update: immediately apply and show feedback
			cmd := m.applySelectedDuration(durations)
			return *m, cmd
		}
	case "tab", "shift+tab":
		// Switch between work and break duration selection
		m.selectingWork = !m.selectingWork
		if m.selectingWork {
			m.durationIndex = findDurationIndex(m.workDuration, WorkDurations)
		} else {
			m.durationIndex = findDurationIndex(m.breakDuration, BreakDurations)
		}
	case "enter":
		// Confirm both values and exit to idle
		// Values are already applied via live update, just exit
		m.mode = FocusModeIdle
		m.durationJustChanged = false
	case "esc":
		// Cancel - restore original values would need tracking, for now just exit
		m.mode = FocusModeIdle
		m.durationJustChanged = false
	}

	return *m, nil
}

// applySelectedDuration applies the currently selected duration immediately.
// Returns commands to show feedback briefly and auto-exit after 500ms.
func (m *FocusModel) applySelectedDuration(durations []int) tea.Cmd {
	if m.selectingWork {
		m.workDuration = durations[m.durationIndex]
		m.remaining = time.Duration(m.workDuration) * time.Minute
		m.totalDuration = m.remaining
		m.lastChangedField = "work"
	} else {
		m.breakDuration = durations[m.durationIndex]
		m.lastChangedField = "break"
	}

	// Show "Saved" indicator
	m.durationJustChanged = true

	// Increment sequence to cancel any pending auto-exit timers
	m.autoExitSequence++
	currentSequence := m.autoExitSequence

	// Return commands: clear feedback after 300ms, auto-exit after 500ms
	return tea.Batch(
		tea.Tick(300*time.Millisecond, func(t time.Time) tea.Msg {
			return clearFeedbackMsg{}
		}),
		tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
			return autoExitDurationMsg{sequence: currentSequence}
		}),
	)
}

// handleHistoryInput handles keyboard input for history view.
func (m *FocusModel) handleHistoryInput(msg tea.KeyMsg) (FocusModel, tea.Cmd) {
	switch msg.String() {
	case "esc", "h":
		m.mode = FocusModeIdle
		return *m, nil
	case "d":
		// Delete selected session
		if len(m.sessionList.Items()) > 0 {
			if selected, ok := m.sessionList.SelectedItem().(SessionItem); ok {
				m.store.DeleteSession(selected.session.ID)
				m.LoadHistory()
			}
		}
		return *m, nil
	}

	// Pass navigation keys to list
	var cmd tea.Cmd
	m.sessionList, cmd = m.sessionList.Update(msg)
	return *m, cmd
}

// findDurationIndex finds the index of a duration in the preset list, or returns 0.
func findDurationIndex(duration int, presets []int) int {
	for i, d := range presets {
		if d == duration {
			return i
		}
	}
	return 0
}

// View renders the focus screen.
func (m *FocusModel) View() string {
	switch m.mode {
	case FocusModeHistory:
		return m.renderHistory()
	case FocusModeDuration:
		return m.renderDurationPicker()
	default:
		return m.renderTimer()
	}
}

// renderTimer renders the main timer view.
func (m *FocusModel) renderTimer() string {
	// Update help hints based on mode
	switch m.mode {
	case FocusModeIdle:
		m.helpBar.SetHints(components.FocusIdleHints)
	case FocusModeRunning:
		m.helpBar.SetHints(components.FocusRunningHints)
	case FocusModePaused:
		m.helpBar.SetHints(components.FocusPausedHints)
	case FocusModeBreak:
		m.helpBar.SetHints(components.FocusBreakHints)
	}

	// Phase indicator
	phaseStyle := lipgloss.NewStyle().
		Foreground(styles.PrimaryColor).
		Bold(true).
		Padding(0, 1)

	var phaseText string
	switch m.mode {
	case FocusModeIdle:
		phaseText = "Ready to Focus"
	case FocusModeRunning:
		phaseText = "🍅 Work Session"
	case FocusModePaused:
		phaseText = "⏸ Paused"
	case FocusModeBreak:
		phaseText = "☕ Break Time"
	}
	phase := phaseStyle.Render(phaseText)

	// Large ASCII timer display
	timer := m.renderLargeTimer()

	// Progress bar
	progress := m.renderProgressBar()

	// Stats summary
	stats := m.renderStatsSummary()

	// Build content
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		m.header.View(),
		"",
		phase,
		"",
		timer,
		"",
		progress,
		"",
		stats,
		"",
		m.helpBar.View(),
	)

	return styles.PanelStyle.Render(content)
}

// renderLargeTimer renders the timer in large ASCII-style digits.
func (m *FocusModel) renderLargeTimer() string {
	minutes := int(m.remaining.Minutes())
	seconds := int(m.remaining.Seconds()) % 60
	timeStr := fmt.Sprintf("%02d:%02d", minutes, seconds)

	timerStyle := styles.TimerStyle
	if m.mode == FocusModeRunning {
		timerStyle = styles.TimerActiveStyle
	} else if m.mode == FocusModeBreak {
		timerStyle = lipgloss.NewStyle().
			Foreground(styles.SecondaryColor).
			Bold(true).
			Padding(1, 4)
	} else if m.mode == FocusModePaused {
		timerStyle = lipgloss.NewStyle().
			Foreground(styles.WarningColor).
			Bold(true).
			Padding(1, 4)
	}

	// Create large timer display
	largeTimerStyle := lipgloss.NewStyle().
		Foreground(timerStyle.GetForeground()).
		Bold(true).
		Padding(1, 2)

	// Simple but large representation
	timerDisplay := fmt.Sprintf(`
    ╔═══════════════════╗
    ║                   ║
    ║      %s       ║
    ║                   ║
    ╚═══════════════════╝`, timeStr)

	return largeTimerStyle.Render(timerDisplay)
}

// renderProgressBar renders a visual progress bar with vaporwave gradient.
func (m *FocusModel) renderProgressBar() string {
	if m.totalDuration == 0 {
		return ""
	}

	elapsed := m.totalDuration - m.remaining
	progress := float64(elapsed) / float64(m.totalDuration)
	if progress < 0 {
		progress = 0
	}
	if progress > 1 {
		progress = 1
	}

	// Use vaporwave progress bar with gradient effect
	bar := styles.VaporwaveProgressBar(progress, 30)

	percentageStyle := lipgloss.NewStyle().Foreground(styles.SecondaryColor).Bold(true)
	percentage := percentageStyle.Render(fmt.Sprintf(" %d%%", int(progress*100)))

	return lipgloss.JoinHorizontal(
		lipgloss.Center,
		bar,
		percentage,
	)
}

// renderStatsSummary renders a brief statistics summary.
func (m *FocusModel) renderStatsSummary() string {
	if m.stats == nil {
		// Load stats if not loaded
		m.LoadHistory()
	}

	statsStyle := lipgloss.NewStyle().
		Foreground(styles.MutedColor).
		Padding(0, 1)

	statItemStyle := lipgloss.NewStyle().
		Foreground(styles.TextColor)

	statValueStyle := lipgloss.NewStyle().
		Foreground(styles.SecondaryColor).
		Bold(true)

	var todaySessions, streak, totalMinutes int
	if m.stats != nil {
		todaySessions = m.stats.TodaySessions
		streak = m.stats.CurrentStreak
		totalMinutes = m.stats.TotalFocusMinutes
	}

	statsContent := lipgloss.JoinHorizontal(
		lipgloss.Center,
		statItemStyle.Render("Today: ")+statValueStyle.Render(fmt.Sprintf("%d", todaySessions)),
		statsStyle.Render(" │ "),
		statItemStyle.Render("Streak: ")+statValueStyle.Render(fmt.Sprintf("%d days", streak)),
		statsStyle.Render(" │ "),
		statItemStyle.Render("Total: ")+statValueStyle.Render(fmt.Sprintf("%dh %dm", totalMinutes/60, totalMinutes%60)),
	)

	return statsContent
}

// renderHistory renders the session history view.
func (m *FocusModel) renderHistory() string {
	m.helpBar.SetHints(components.FocusHistoryHints)

	title := styles.TitleStyle.Render("📊 Session History")

	if len(m.sessionList.Items()) == 0 {
		emptyState := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			"",
			styles.SubtitleStyle.Render("No sessions yet. Start focusing!"),
			"",
			m.helpBar.View(),
		)
		return styles.PanelStyle.Render(emptyState)
	}

	// Stats header
	statsHeader := m.renderStatsSummary()

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		statsHeader,
		"",
		m.sessionList.View(),
		"",
		m.helpBar.View(),
	)

	return styles.PanelStyle.Render(content)
}

// renderDurationPicker renders the duration selection UI.
func (m *FocusModel) renderDurationPicker() string {
	m.helpBar.SetHints(components.FocusDurationHints)

	title := styles.TitleStyle.Render("⏱ Set Duration")

	// Saved indicator style
	savedStyle := lipgloss.NewStyle().
		Foreground(styles.SuccessColor).
		Bold(true)

	// Work duration selection
	workLabel := styles.SubtitleStyle.Render("Work Duration:")
	workSaved := ""
	if m.selectingWork {
		workLabel = styles.SelectedItemStyle.Render("▶ Work Duration:")
	}
	if m.durationJustChanged && m.lastChangedField == "work" {
		workSaved = savedStyle.Render(" ✓ Saved")
	}
	workOptions := m.renderDurationOptions(WorkDurations, m.workDuration, m.selectingWork)
	workRow := lipgloss.JoinHorizontal(lipgloss.Left, workLabel, workSaved)

	// Break duration selection
	breakLabel := styles.SubtitleStyle.Render("Break Duration:")
	breakSaved := ""
	if !m.selectingWork {
		breakLabel = styles.SelectedItemStyle.Render("▶ Break Duration:")
	}
	if m.durationJustChanged && m.lastChangedField == "break" {
		breakSaved = savedStyle.Render(" ✓ Saved")
	}
	breakOptions := m.renderDurationOptions(BreakDurations, m.breakDuration, !m.selectingWork)
	breakRow := lipgloss.JoinHorizontal(lipgloss.Left, breakLabel, breakSaved)

	// Current values summary
	summaryStyle := lipgloss.NewStyle().
		Foreground(styles.MutedColor).
		Italic(true)
	summary := summaryStyle.Render(fmt.Sprintf("Current: %d min work / %d min break", m.workDuration, m.breakDuration))

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		summary,
		"",
		workRow,
		workOptions,
		"",
		breakRow,
		breakOptions,
		"",
		m.helpBar.View(),
	)

	return styles.PanelStyle.Render(content)
}

// renderDurationOptions renders the duration preset options.
func (m *FocusModel) renderDurationOptions(durations []int, current int, isActive bool) string {
	normalStyle := lipgloss.NewStyle().
		Foreground(styles.MutedColor).
		Padding(0, 1)

	selectedStyle := lipgloss.NewStyle().
		Foreground(styles.SecondaryColor).
		Bold(true).
		Background(styles.SurfaceColor).
		Padding(0, 1)

	currentStyle := lipgloss.NewStyle().
		Foreground(styles.PrimaryColor).
		Padding(0, 1)

	var options []string
	for i, d := range durations {
		label := fmt.Sprintf("%d min", d)
		style := normalStyle

		if isActive && i == m.durationIndex {
			style = selectedStyle
		} else if d == current {
			style = currentStyle
		}

		options = append(options, style.Render(label))
	}

	return lipgloss.JoinHorizontal(lipgloss.Center, options...)
}

// SessionItem implements list.Item for displaying sessions in the history list.
type SessionItem struct {
	session models.FocusSession
}

func (s SessionItem) Title() string {
	date := s.session.StartTime.Format("2006-01-02 15:04")
	duration := s.session.Duration / 60 // Convert to minutes

	statusIcon := "✓"
	if s.session.Status == models.SessionStatusCancelled {
		statusIcon = "✗"
	} else if s.session.Status == models.SessionStatusRunning {
		statusIcon = "●"
	}

	return fmt.Sprintf("%s %s - %d min", statusIcon, date, duration)
}

func (s SessionItem) Description() string {
	if s.session.EndTime != nil {
		elapsed := s.session.EndTime.Sub(s.session.StartTime)
		return fmt.Sprintf("Actual: %d min", int(elapsed.Minutes()))
	}
	return "In progress"
}

func (s SessionItem) FilterValue() string {
	return s.session.StartTime.Format("2006-01-02")
}
