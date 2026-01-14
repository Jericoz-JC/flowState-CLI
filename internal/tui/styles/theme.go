package styles

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ASCII Art banner for the home screen - ARCHWAVE vaporwave style
const LogoASCII = `
  ╔══════════════════════════════════════════════════════════════════╗
  ║                                                                  ║
  ║   ███████╗██╗      ██████╗ ██╗    ██╗                            ║
  ║   ██╔════╝██║     ██╔═══██╗██║    ██║                            ║
  ║   █████╗  ██║     ██║   ██║██║ █╗ ██║                            ║
  ║   ██╔══╝  ██║     ██║   ██║██║███╗██║                            ║
  ║   ██║     ███████╗╚██████╔╝╚███╔███╔╝                            ║
  ║   ╚═╝     ╚══════╝ ╚═════╝  ╚══╝╚══╝                             ║
  ║                                                                  ║
  ║     ███████╗████████╗ █████╗ ████████╗███████╗                   ║
  ║     ██╔════╝╚══██╔══╝██╔══██╗╚══██╔══╝██╔════╝                   ║
  ║     ███████╗   ██║   ███████║   ██║   █████╗                     ║
  ║     ╚════██║   ██║   ██╔══██║   ██║   ██╔══╝                     ║
  ║     ███████║   ██║   ██║  ██║   ██║   ███████╗                   ║
  ║     ╚══════╝   ╚═╝   ╚═╝  ╚═╝   ╚═╝   ╚══════╝                   ║
  ║                                                                  ║
  ║            ✦  productivity reimagined  ✦                        ║
  ╚══════════════════════════════════════════════════════════════════╝`

// ARCHWAVE vaporwave color palette - pastel pinks, purples, and cyans
var (
	// Primary colors - vaporwave spectrum
	PrimaryColor   = lipgloss.Color("#d4a5ff") // Soft lavender
	SecondaryColor = lipgloss.Color("#5ffbf1") // Neon cyan
	AccentColor    = lipgloss.Color("#ff6ec7") // Hot pink

	// Semantic colors
	SuccessColor = lipgloss.Color("#8ffef4") // Light cyan
	WarningColor = lipgloss.Color("#f9f871") // Pale yellow
	ErrorColor   = lipgloss.Color("#ff9adc") // Soft pink
	TimerColor   = lipgloss.Color("#ff6ec7") // Hot pink for timer

	// Background colors - deep purple-black
	BackgroundColor = lipgloss.Color("#1a0d2e") // Deep purple-black
	SurfaceColor    = lipgloss.Color("#2d1b4e") // Dark purple
	BorderColor     = lipgloss.Color("#543a6e") // Muted purple

	// Text colors
	TextColor      = lipgloss.Color("#fef6ff") // Off-white pink
	MutedColor     = lipgloss.Color("#b8c1ff") // Pale blue
	HighlightColor = lipgloss.Color("#ffffff") // Pure white for highlights

	// Additional ARCHWAVE colors
	NeonPink    = lipgloss.Color("#f4a5ff") // Soft neon pink
	PaleAqua    = lipgloss.Color("#8ffef4") // Pale aqua
	CreamYellow = lipgloss.Color("#fbf9a5") // Cream yellow
	Periwinkle  = lipgloss.Color("#8b9aff") // Periwinkle blue
	PalePink    = lipgloss.Color("#ffc8ff") // Pale pink

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

	// Panel with double border (vaporwave aesthetic)
	PanelStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(BorderColor).
			Padding(1, 2)

	// Highlighted panel (for focused elements)
	PanelActiveStyle = lipgloss.NewStyle().
				Border(lipgloss.DoubleBorder()).
				BorderForeground(AccentColor).
				Padding(1, 2)

	// Border style for sections
	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(BorderColor)

	// Neon glow style for important elements
	NeonStyle = lipgloss.NewStyle().
			Foreground(SecondaryColor).
			Bold(true)

	// Retro box style with hot pink border
	RetroBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(AccentColor).
			Padding(1, 2)

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
				BorderForeground(AccentColor).
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

// Vaporwave decorative elements
const (
	DecoStar    = "✦"
	DecoWave    = "～"
	DecoNeonBar = "▓"
	DecoSoftBar = "▒"
	DecoFadeBar = "░"
	DecoArrow   = "▸"
	DecoBullet  = "◈"
)

// VaporwaveProgressBar renders a gradient progress bar with vaporwave styling
func VaporwaveProgressBar(progress float64, width int) string {
	if width <= 0 {
		return ""
	}
	filled := int(float64(width) * progress)
	if filled > width {
		filled = width
	}
	empty := width - filled

	filledStyle := lipgloss.NewStyle().Foreground(SecondaryColor)
	emptyStyle := lipgloss.NewStyle().Foreground(BorderColor)

	return filledStyle.Render(strings.Repeat(DecoNeonBar, filled)) +
		emptyStyle.Render(strings.Repeat(DecoFadeBar, empty))
}

// GradientText applies alternating colors to text for a gradient-like effect
func GradientText(text string, colors ...lipgloss.Color) string {
	if len(colors) == 0 || len(text) == 0 {
		return text
	}

	var result strings.Builder
	runes := []rune(text)
	for i, char := range runes {
		if char == ' ' {
			result.WriteRune(char)
			continue
		}
		colorIndex := i % len(colors)
		style := lipgloss.NewStyle().Foreground(colors[colorIndex])
		result.WriteString(style.Render(string(char)))
	}
	return result.String()
}

// VaporwaveDivider creates a decorative divider with vaporwave styling
func VaporwaveDivider(width int) string {
	if width <= 4 {
		return DividerStyle.Render("══")
	}
	starStyle := lipgloss.NewStyle().Foreground(AccentColor)
	lineStyle := lipgloss.NewStyle().Foreground(BorderColor)

	lineWidth := (width - 4) / 2
	return lineStyle.Render(strings.Repeat("═", lineWidth)) +
		starStyle.Render(" "+DecoStar+" ") +
		lineStyle.Render(strings.Repeat("═", lineWidth))
}

// VaporwaveSeparator creates a small separator for inline use
func VaporwaveSeparator() string {
	sepStyle := lipgloss.NewStyle().Foreground(BorderColor)
	return sepStyle.Render(" " + DecoBullet + " ")
}
