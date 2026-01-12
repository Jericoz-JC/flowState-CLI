package components

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"

	"flowState-cli/internal/tui/styles"
)

func DefaultListDelegate() list.DefaultDelegate {
	return list.NewDefaultDelegate()
}

func CreateList(title string) list.Model {
	items := []list.Item{}
	delegate := DefaultListDelegate()

	l := list.New(items, delegate, 0, 0)
	l.Title = title
	l.SetShowHelp(true)

	return l
}

func FormatTimer(duration time.Duration) string {
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	return lipgloss.NewStyle().
		Foreground(styles.PrimaryColor).
		Bold(true).
		Render(fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds))
}
