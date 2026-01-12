package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// ASCII Art banner for the home screen
const LogoASCII = `
  ╭─────────────────────────────────────────────────────────────╮
  │  ███████╗██╗      ██████╗ ██╗    ██╗                        │
  │  ██╔════╝██║     ██╔═══██╗██║    ██║                        │
  │  █████╗  ██║     ██║   ██║██║ █╗ ██║                        │
  │  ██╔══╝  ██║     ██║   ██║██║███╗██║                        │
  │  ██║     ███████╗╚██████╔╝╚███╔███╔╝                        │
  │  ╚═╝     ╚══════╝ ╚═════╝  ╚══╝╚══╝  S T A T E             │
  ╰─────────────────────────────────────────────────────────────╯`

// Enhanced color palette for a more vibrant, modern look
var (
	// Primary colors - brighter purple spectrum
	PrimaryColor   = lipgloss.Color("#A78BFA") // Bright violet
	SecondaryColor = lipgloss.Color("#22D3EE") // Cyan glow
	AccentColor    = lipgloss.Color("#F472B6") // Pink accent

	// Semantic colors
	SuccessColor = lipgloss.Color("#34D399") // Mint green
	WarningColor = lipgloss.Color("#FBBF24") // Amber
	ErrorColor   = lipgloss.Color("#F87171") // Soft red
	TimerColor   = lipgloss.Color("#FF6B6B") // Warm red for timer

	// Background colors - deeper darks
	BackgroundColor = lipgloss.Color("#0F0F1A") // Deep dark
	SurfaceColor    = lipgloss.Color("#1E1E2E") // Elevated surface
	BorderColor     = lipgloss.Color("#313244") // Subtle border

	// Text colors
	TextColor      = lipgloss.Color("#CDD6F4") // Bright text
	MutedColor     = lipgloss.Color("#6C7086") // Muted text
	HighlightColor = lipgloss.Color("#F5F5F5") // Pure white for highlights

	// Logo style with gradient effect
	LogoStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true)

	// Title style - larger, more prominent
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(PrimaryColor).
			MarginBottom(1).
			Padding(0, 1)

	// Subtitle style
	SubtitleStyle = lipgloss.NewStyle().
			Foreground(MutedColor).
			Italic(true).
			MarginBottom(2)

	// Menu item styles
	MenuItemStyle = lipgloss.NewStyle().
			Foreground(TextColor).
			Padding(0, 2).
			MarginLeft(2)

	MenuItemActiveStyle = lipgloss.NewStyle().
				Foreground(SecondaryColor).
				Bold(true).
				Padding(0, 2).
				MarginLeft(2)

	// Selected/active item style
	SelectedItemStyle = lipgloss.NewStyle().
				Foreground(SecondaryColor).
				Bold(true).
				Background(SurfaceColor).
				Padding(0, 1)

	// Status bar - more prominent with accent
	StatusBarStyle = lipgloss.NewStyle().
			Background(SurfaceColor).
			Foreground(MutedColor).
			Padding(0, 2).
			MarginTop(1)

	// Timer display style
	TimerStyle = lipgloss.NewStyle().
			Foreground(TimerColor).
			Bold(true).
			Padding(1, 4)

	TimerActiveStyle = lipgloss.NewStyle().
				Foreground(SuccessColor).
				Bold(true).
				Padding(1, 4)

	// Container with border
	ContainerStyle = lipgloss.NewStyle().
			Background(BackgroundColor).
			Padding(1, 2)

	// Panel with rounded border
	PanelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(BorderColor).
			Padding(1, 2)

	// Highlighted panel (for focused elements)
	PanelActiveStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(PrimaryColor).
				Padding(1, 2)

	// Border style for sections
	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(BorderColor)

	// Tag style - pill-like appearance
	TagStyle = lipgloss.NewStyle().
			Foreground(SecondaryColor).
			Background(SurfaceColor).
			Padding(0, 1).
			MarginRight(1)

	// Progress bar style
	ProgressBarStyle = lipgloss.NewStyle().
				Foreground(PrimaryColor)

	ProgressBarFilledStyle = lipgloss.NewStyle().
				Foreground(SuccessColor)

	// Input field styles
	InputStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(BorderColor).
			Padding(0, 1)

	InputFocusedStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(PrimaryColor).
				Padding(0, 1)

	// Help text style
	HelpStyle = lipgloss.NewStyle().
			Foreground(MutedColor).
			Italic(true).
			MarginTop(1)

	// Keyboard shortcut style
	KeyStyle = lipgloss.NewStyle().
			Foreground(AccentColor).
			Bold(true)

	// Description/label in help
	DescStyle = lipgloss.NewStyle().
			Foreground(MutedColor)

	// Success message style
	SuccessStyle = lipgloss.NewStyle().
			Foreground(SuccessColor).
			Bold(true)

	// Error message style
	ErrorStyle = lipgloss.NewStyle().
			Foreground(ErrorColor).
			Bold(true)

	// Warning message style
	WarningStyle = lipgloss.NewStyle().
			Foreground(WarningColor)

	// Divider line
	DividerStyle = lipgloss.NewStyle().
			Foreground(BorderColor)
)

// Helper function to create a full-screen container
func Screen(width, height int) lipgloss.Style {
	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Background(BackgroundColor)
}

// Helper to render a keyboard shortcut hint
func KeyHint(key, description string) string {
	return KeyStyle.Render("["+key+"]") + " " + DescStyle.Render(description)
}

// Helper to create a horizontal divider
func Divider(width int) string {
	line := ""
	for i := 0; i < width; i++ {
		line += "─"
	}
	return DividerStyle.Render(line)
}
