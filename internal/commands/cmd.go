package commands

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type TickMsg struct {
	Time time.Time
}

func Tick() tea.Cmd {
	return func() tea.Msg {
		return TickMsg{Time: time.Now()}
	}
}

type LoadingMsg struct {
	Message string
}

func Loading(message string) tea.Cmd {
	return func() tea.Msg {
		return LoadingMsg{Message: message}
	}
}
