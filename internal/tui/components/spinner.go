// Package components provides reusable TUI components for flowState-cli.
//
// AnimatedSpinner provides a loading indicator with ARCHWAVE vaporwave styling.
package components

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Jericoz-JC/flowState-CLI/internal/tui/styles"
)

// VaporwaveSpinnerFrames are the default spinner animation frames
// Uses circular quarter-block characters for smooth animation
var VaporwaveSpinnerFrames = []string{"◐", "◓", "◑", "◒"}

// NeonSpinnerFrames alternative frames with horizontal bars
var NeonSpinnerFrames = []string{"▰▱▱", "▰▰▱", "▰▰▰", "▱▰▰", "▱▱▰", "▱▱▱"}

// DotsSpinnerFrames braille-style dots animation
var DotsSpinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// SpinnerTickMsg signals that the spinner should advance to the next frame
type SpinnerTickMsg struct{}

// AnimatedSpinner provides a loading indicator with ARCHWAVE styling
type AnimatedSpinner struct {
	frames   []string
	current  int
	style    lipgloss.Style
	label    string
	isActive bool
	interval time.Duration
}

// NewAnimatedSpinner creates a new spinner with default vaporwave frames
func NewAnimatedSpinner() AnimatedSpinner {
	return AnimatedSpinner{
		frames:   VaporwaveSpinnerFrames,
		current:  0,
		style:    lipgloss.NewStyle().Foreground(styles.SecondaryColor),
		label:    "",
		isActive: false,
		interval: 100 * time.Millisecond,
	}
}

// NewAnimatedSpinnerWithFrames creates a spinner with custom frames
func NewAnimatedSpinnerWithFrames(frames []string) AnimatedSpinner {
	s := NewAnimatedSpinner()
	if len(frames) > 0 {
		s.frames = frames
	}
	return s
}

// Start begins the spinner animation and returns the initial tick command
func (s *AnimatedSpinner) Start() tea.Cmd {
	s.isActive = true
	return s.tick()
}

// Stop halts the spinner animation
func (s *AnimatedSpinner) Stop() {
	s.isActive = false
}

// IsActive returns whether the spinner is currently animating
func (s AnimatedSpinner) IsActive() bool {
	return s.isActive
}

// SetLabel sets the text displayed next to the spinner
func (s *AnimatedSpinner) SetLabel(label string) {
	s.label = label
}

// SetStyle sets the lipgloss style for the spinner frame
func (s *AnimatedSpinner) SetStyle(style lipgloss.Style) {
	s.style = style
}

// SetInterval sets the animation speed (time between frames)
func (s *AnimatedSpinner) SetInterval(d time.Duration) {
	s.interval = d
}

// Update handles spinner tick messages
func (s AnimatedSpinner) Update(msg tea.Msg) (AnimatedSpinner, tea.Cmd) {
	switch msg.(type) {
	case SpinnerTickMsg:
		if !s.isActive {
			return s, nil
		}

		// Advance to next frame
		s.current = (s.current + 1) % len(s.frames)

		// Return command for next tick
		return s, s.tick()
	}

	return s, nil
}

// View renders the spinner with its current frame and optional label
func (s AnimatedSpinner) View() string {
	frame := s.style.Render(s.frames[s.current])

	if s.label == "" {
		return frame
	}

	labelStyle := lipgloss.NewStyle().Foreground(styles.MutedColor)
	return frame + " " + labelStyle.Render(s.label)
}

// tick returns a command that sends a SpinnerTickMsg after the interval
func (s AnimatedSpinner) tick() tea.Cmd {
	return tea.Tick(s.interval, func(t time.Time) tea.Msg {
		return SpinnerTickMsg{}
	})
}
