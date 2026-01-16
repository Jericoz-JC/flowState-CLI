package styles

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ASCII Art banner for the home screen - ARCHWAVE vaporwave style
// Full logo requires ~72 chars width
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
╚══════════════════════════════════════════════════════════════════╝`

// LogoASCIISmall is a compact logo for narrow terminals (< 72 chars)
const LogoASCIISmall = `
╔════════════════════════════════╗
║  ╔═╗╦  ╔═╗╦ ╦  ╔═╗╔╦╗╔═╗╔╦╗╔═╗ ║
║  ╠╣ ║  ║ ║║║║  ╚═╗ ║ ╠═╣ ║ ║╣  ║
║  ╚  ╩═╝╚═╝╚╩╝  ╚═╝ ╩ ╩ ╩ ╩ ╚═╝ ║
╚════════════════════════════════╝`

// LogoMinWidth is the minimum terminal width for full ASCII logo
const LogoMinWidth = 72

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

	// Card styles for list items (enhanced visual hierarchy)
	CardStyle = lipgloss.NewStyle().
			Background(SurfaceColor).
			Padding(0, 1).
			MarginBottom(1)

	CardActiveStyle = lipgloss.NewStyle().
			Background(SurfaceColor).
			BorderLeft(true).
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(AccentColor).
			Padding(0, 1).
			MarginBottom(1)

	// Badge styles for status indicators
	BadgeStyle = lipgloss.NewStyle().
			Foreground(TextColor).
			Background(BorderColor).
			Padding(0, 1)

	BadgeSuccessStyle = lipgloss.NewStyle().
				Foreground(BackgroundColor).
				Background(SuccessColor).
				Bold(true).
				Padding(0, 1)

	BadgeWarningStyle = lipgloss.NewStyle().
				Foreground(BackgroundColor).
				Background(WarningColor).
				Bold(true).
				Padding(0, 1)

	BadgeErrorStyle = lipgloss.NewStyle().
			Foreground(BackgroundColor).
			Background(ErrorColor).
			Bold(true).
			Padding(0, 1)

	BadgeInfoStyle = lipgloss.NewStyle().
			Foreground(BackgroundColor).
			Background(SecondaryColor).
			Bold(true).
			Padding(0, 1)

	// Section header style with decorative line
	SectionHeaderStyle = lipgloss.NewStyle().
				Foreground(PrimaryColor).
				Bold(true).
				BorderBottom(true).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(BorderColor).
				MarginBottom(1).
				PaddingBottom(0)

	// Muted card for completed/inactive items
	CardMutedStyle = lipgloss.NewStyle().
			Foreground(MutedColor).
			Background(BackgroundColor).
			Padding(0, 1).
			MarginBottom(1)

	// Highlight box for important messages
	HighlightBoxStyle = lipgloss.NewStyle().
				Foreground(TextColor).
				Background(SurfaceColor).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(AccentColor).
				Padding(1, 2)

	// Empty state style
	EmptyStateStyle = lipgloss.NewStyle().
			Foreground(MutedColor).
			Align(lipgloss.Center).
			Italic(true).
			Padding(2, 4)

	// Count badge (for item counts in headers)
	CountBadgeStyle = lipgloss.NewStyle().
			Foreground(SecondaryColor).
			Background(SurfaceColor).
			Padding(0, 1).
			Bold(true)

	// Inline link style
	LinkStyle = lipgloss.NewStyle().
			Foreground(SecondaryColor).
			Underline(true)

	// Code/monospace style
	CodeStyle = lipgloss.NewStyle().
			Foreground(PaleAqua).
			Background(SurfaceColor).
			Padding(0, 1)
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

// StatusBadge returns a styled badge based on status type
func StatusBadge(status string) string {
	switch status {
	case "completed", "done", "success":
		return BadgeSuccessStyle.Render(" ✓ " + status + " ")
	case "pending", "todo":
		return BadgeStyle.Render(" ○ " + status + " ")
	case "in_progress", "active":
		return BadgeWarningStyle.Render(" ▶ " + status + " ")
	case "error", "failed":
		return BadgeErrorStyle.Render(" ✕ " + status + " ")
	default:
		return BadgeInfoStyle.Render(" " + status + " ")
	}
}

// SectionHeader renders a styled section header with optional count
func SectionHeader(title string, count int) string {
	countStr := ""
	if count >= 0 {
		countStr = " " + CountBadgeStyle.Render(strconv.Itoa(count))
	}
	return SectionHeaderStyle.Render(title + countStr)
}

// HighlightBox renders text in a highlighted box
func HighlightBox(text string) string {
	return HighlightBoxStyle.Render(text)
}

// EmptyState renders a centered empty state message
func EmptyState(message string) string {
	return EmptyStateStyle.Render(DecoStar + " " + message + " " + DecoStar)
}

// FormatTag renders a tag with proper styling
func FormatTag(tag string) string {
	return TagStyle.Render("#" + tag)
}

// FormatTags renders multiple tags with proper styling
func FormatTags(tags []string) string {
	if len(tags) == 0 {
		return ""
	}
	var result strings.Builder
	for i, tag := range tags {
		if i > 0 {
			result.WriteString(" ")
		}
		result.WriteString(FormatTag(tag))
	}
	return result.String()
}

// ASCII art digit definitions - 5 lines tall, 6 chars wide
// Each digit is represented as 5 strings (lines)
var asciiDigits = map[rune][]string{
	'0': {
		"█████╗ ",
		"██╔██║ ",
		"█╔╝██║ ",
		"██╔██║ ",
		"╚████╝ ",
	},
	'1': {
		" ██╗   ",
		"███║   ",
		"╚██║   ",
		" ██║   ",
		"███╝   ",
	},
	'2': {
		"█████╗ ",
		"╚═══██╗",
		" ████╔╝",
		"██╔══╝ ",
		"██████╗",
	},
	'3': {
		"█████╗ ",
		"╚═══██╗",
		" ████╔╝",
		"╚═══██╗",
		"█████╔╝",
	},
	'4': {
		"██╗██╗ ",
		"██║██║ ",
		"██████╗",
		"╚══██║ ",
		"   ██╝ ",
	},
	'5': {
		"██████╗",
		"██╔═══╝",
		"█████╗ ",
		"╚═══██╗",
		"█████╔╝",
	},
	'6': {
		"█████╗ ",
		"██╔═══╝",
		"█████╗ ",
		"██╔═██╗",
		"╚████╔╝",
	},
	'7': {
		"██████╗",
		"╚═══██║",
		"   ██╔╝",
		"  ██╔╝ ",
		"  ██║  ",
	},
	'8': {
		"█████╗ ",
		"██╔═██╗",
		"╚███╔╝ ",
		"██╔═██╗",
		"╚████╔╝",
	},
	'9': {
		"█████╗ ",
		"██╔═██╗",
		"╚████║ ",
		"╚══██║ ",
		"█████╔╝",
	},
	':': {
		"   ",
		" █ ",
		"   ",
		" █ ",
		"   ",
	},
}

// RenderASCIITime renders a time string (e.g., "25:00") as large ASCII art
// Returns a slice of strings, one per line
func RenderASCIITime(timeStr string, color lipgloss.Color) string {
	lines := make([]string, 5)

	style := lipgloss.NewStyle().Foreground(color).Bold(true)

	for _, char := range timeStr {
		digit, ok := asciiDigits[char]
		if !ok {
			continue
		}
		for i := 0; i < 5; i++ {
			lines[i] += digit[i]
		}
	}

	// Apply styling to each line
	var result strings.Builder
	for i, line := range lines {
		result.WriteString(style.Render(line))
		if i < 4 {
			result.WriteString("\n")
		}
	}

	return result.String()
}

// ASCII art for session mode indicators
const (
	WorkModeASCII = `
╔══════════════════════════════════════╗
║  🍅  W O R K   S E S S I O N  🍅    ║
╚══════════════════════════════════════╝`

	BreakModeASCII = `
╔══════════════════════════════════════╗
║  ☕  B R E A K   T I M E  ☕         ║
╚══════════════════════════════════════╝`

	IdleModeASCII = `
╔══════════════════════════════════════╗
║  ✦  R E A D Y   T O   F O C U S  ✦  ║
╚══════════════════════════════════════╝`

	PausedModeASCII = `
╔══════════════════════════════════════╗
║  ⏸  P A U S E D  ⏸                  ║
╚══════════════════════════════════════╝`
)

// RenderProgressRing renders a circular-style progress indicator
// using ASCII characters for terminal display
func RenderProgressRing(progress float64, width int) string {
	if width <= 0 {
		return ""
	}

	// Characters for different fill levels
	chars := []string{"░", "▒", "▓", "█"}

	filled := int(float64(width) * progress)
	if filled > width {
		filled = width
	}

	// Create gradient effect
	var result strings.Builder

	// Opening bracket
	bracketStyle := lipgloss.NewStyle().Foreground(BorderColor)
	result.WriteString(bracketStyle.Render("【"))

	for i := 0; i < width; i++ {
		var char string
		var style lipgloss.Style

		if i < filled {
			// Filled section with gradient colors
			gradientPos := float64(i) / float64(width)
			if gradientPos < 0.5 {
				style = lipgloss.NewStyle().Foreground(SecondaryColor) // Cyan
			} else {
				style = lipgloss.NewStyle().Foreground(AccentColor) // Pink
			}
			char = chars[3] // Full block
		} else {
			// Empty section
			style = lipgloss.NewStyle().Foreground(SurfaceColor)
			char = chars[0]
		}
		result.WriteString(style.Render(char))
	}

	// Closing bracket
	result.WriteString(bracketStyle.Render("】"))

	return result.String()
}

// RenderMiniBarChart renders a small bar chart for the last 7 days
func RenderMiniBarChart(values []int, maxHeight int, width int) string {
	if len(values) == 0 || maxHeight <= 0 {
		return ""
	}

	// Find max value for scaling
	maxVal := 1
	for _, v := range values {
		if v > maxVal {
			maxVal = v
		}
	}

	// Calculate bar width
	barWidth := width / len(values)
	if barWidth < 1 {
		barWidth = 1
	}

	// Build chart from top to bottom
	var lines []string
	for row := maxHeight; row >= 1; row-- {
		var line strings.Builder
		threshold := float64(row) / float64(maxHeight) * float64(maxVal)

		for i, v := range values {
			barStyle := lipgloss.NewStyle().Foreground(SecondaryColor)
			if i == len(values)-1 {
				// Today's bar in accent color
				barStyle = lipgloss.NewStyle().Foreground(AccentColor)
			}

			emptyStyle := lipgloss.NewStyle().Foreground(SurfaceColor)

			if float64(v) >= threshold {
				line.WriteString(barStyle.Render(strings.Repeat("█", barWidth)))
			} else {
				line.WriteString(emptyStyle.Render(strings.Repeat(" ", barWidth)))
			}

			// Add space between bars
			if i < len(values)-1 {
				line.WriteString(" ")
			}
		}
		lines = append(lines, line.String())
	}

	// Add day labels
	dayLabels := []string{"M", "T", "W", "T", "F", "S", "S"}
	if len(values) <= len(dayLabels) {
		var labelLine strings.Builder
		labelStyle := lipgloss.NewStyle().Foreground(MutedColor)
		for i := 0; i < len(values); i++ {
			dayIndex := (7 - len(values) + i) % 7
			padding := (barWidth - 1) / 2
			labelLine.WriteString(strings.Repeat(" ", padding))
			labelLine.WriteString(labelStyle.Render(dayLabels[dayIndex]))
			labelLine.WriteString(strings.Repeat(" ", barWidth-padding-1))
			if i < len(values)-1 {
				labelLine.WriteString(" ")
			}
		}
		lines = append(lines, labelLine.String())
	}

	return strings.Join(lines, "\n")
}

// SessionCountIndicator renders visual dots for completed sessions today
func SessionCountIndicator(count, max int) string {
	if max <= 0 {
		max = 8
	}

	var result strings.Builder
	completedStyle := lipgloss.NewStyle().Foreground(SuccessColor)
	emptyStyle := lipgloss.NewStyle().Foreground(SurfaceColor)

	for i := 0; i < max; i++ {
		if i < count {
			result.WriteString(completedStyle.Render("●"))
		} else {
			result.WriteString(emptyStyle.Render("○"))
		}
		if i < max-1 {
			result.WriteString(" ")
		}
	}

	return result.String()
}
