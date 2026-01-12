package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	PrimaryColor    = lipgloss.Color("#7D56C4")
	SecondaryColor  = lipgloss.Color("#4ECDC4")
	SuccessColor    = lipgloss.Color("#10B981")
	WarningColor    = lipgloss.Color("#F59E0B")
	ErrorColor      = lipgloss.Color("#EF4444")
	BackgroundColor = lipgloss.Color("#1A1B26")
	SurfaceColor    = lipgloss.Color("#24283B")
	TextColor       = lipgloss.Color("#C0CAF5")
	MutedColor      = lipgloss.Color("#565F89")

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(PrimaryColor).
			MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(MutedColor).
			MarginBottom(2)

	MenuItemStyle = lipgloss.NewStyle().
			Foreground(TextColor).
			Padding(0, 1)

	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(SecondaryColor).
				Bold(true)

	StatusBarStyle = lipgloss.NewStyle().
			Background(SurfaceColor).
			Foreground(MutedColor).
			Padding(0, 2)

	TimerStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true)

	ContainerStyle = lipgloss.NewStyle().
			Background(BackgroundColor).
			Padding(1, 2)

	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(SurfaceColor)

	TagStyle = lipgloss.NewStyle().
			Foreground(SecondaryColor).
			Background(SurfaceColor).
			Padding(0, 1).
			MarginRight(1)

	ProgressBarStyle = lipgloss.NewStyle().
				Foreground(PrimaryColor)
)

func Screen(width, height int) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Background(BackgroundColor)
}
