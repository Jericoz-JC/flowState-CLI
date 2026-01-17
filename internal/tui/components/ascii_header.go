// Package components provides reusable TUI components for flowState-cli.
//
// ASCIIHeader provides screen headers with ASCII art decoration
// in the ARCHWAVE vaporwave style.
package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/Jericoz-JC/flowState-CLI/internal/tui/styles"
)

// HeaderStyle defines the visual style for the header
type HeaderStyle int

const (
	// HeaderStyleMinimal shows just icon + title
	HeaderStyleMinimal HeaderStyle = iota
	// HeaderStyleBoxed wraps in a double-border box
	HeaderStyleBoxed
	// HeaderStyleBanner creates a full-width banner
	HeaderStyleBanner
)

// ScreenASCII contains ASCII art headers for each screen
var ScreenASCII = map[string]string{
	"notes": `
╔═══════════════════════════════════╗
║  📝  N O T E S                    ║
╚═══════════════════════════════════╝`,
	"todos": `
╔═══════════════════════════════════╗
║  ✅  T O D O S                    ║
╚═══════════════════════════════════╝`,
	"focus": `
╔═══════════════════════════════════╗
║  🍅  F O C U S                    ║
╚═══════════════════════════════════╝`,
	"search": `
╔═══════════════════════════════════╗
║  🔍  S E A R C H                  ║
╚═══════════════════════════════════╝`,
	"mindmap": `
╔═══════════════════════════════════╗
║  🧠  M I N D   M A P              ║
╚═══════════════════════════════════╝`,
}

// ASCIIHeader provides screen headers with ASCII art decoration
type ASCIIHeader struct {
	icon      string
	title     string
	subtitle  string
	style     HeaderStyle
	width     int
	itemCount int
}

// NewASCIIHeader creates a new ASCII header with the given icon and title
func NewASCIIHeader(icon, title string) ASCIIHeader {
	return ASCIIHeader{
		icon:      icon,
		title:     title,
		subtitle:  "",
		style:     HeaderStyleBoxed,
		width:     80,
		itemCount: -1, // -1 means hidden
	}
}

// SetWidth updates the width constraint for rendering
func (h *ASCIIHeader) SetWidth(width int) {
	h.width = width
}

// SetStyle sets the header rendering style
func (h *ASCIIHeader) SetStyle(style HeaderStyle) {
	h.style = style
}

// SetItemCount sets the item count to display (-1 to hide)
func (h *ASCIIHeader) SetItemCount(count int) {
	h.itemCount = count
}

// SetSubtitle sets an optional subtitle below the main title
func (h *ASCIIHeader) SetSubtitle(subtitle string) {
	h.subtitle = subtitle
}

// View renders the header based on current style
func (h ASCIIHeader) View() string {
	switch h.style {
	case HeaderStyleMinimal:
		return h.renderMinimal()
	case HeaderStyleBanner:
		return h.renderBanner()
	case HeaderStyleBoxed:
		fallthrough
	default:
		return h.renderBoxed()
	}
}

// renderMinimal renders a simple icon + title header
func (h ASCIIHeader) renderMinimal() string {
	iconStyle := lipgloss.NewStyle().Foreground(styles.SecondaryColor)
	titleStyle := lipgloss.NewStyle().Foreground(styles.PrimaryColor).Bold(true)

	result := iconStyle.Render(h.icon) + " " + titleStyle.Render(h.title)

	if h.itemCount >= 0 {
		countStyle := lipgloss.NewStyle().Foreground(styles.MutedColor).Italic(true)
		countText := fmt.Sprintf("%d items", h.itemCount)
		if h.itemCount == 1 {
			countText = "1 item"
		}
		result += " " + countStyle.Render(countText)
	}

	if h.subtitle != "" {
		subtitleStyle := lipgloss.NewStyle().Foreground(styles.MutedColor)
		result += "\n" + subtitleStyle.Render(h.subtitle)
	}

	return result
}

// renderBoxed renders the header in a bordered box
func (h ASCIIHeader) renderBoxed() string {
	// Spaced title for vaporwave aesthetic
	spacedTitle := spacedText(h.title)

	// Calculate content width
	contentWidth := h.width - 6 // Account for borders and padding
	if contentWidth < 20 {
		contentWidth = 20
	}

	// Build title line with icon
	titleContent := fmt.Sprintf("  %s  %s", h.icon, spacedTitle)

	// Add item count if present
	if h.itemCount >= 0 {
		countText := fmt.Sprintf("%d", h.itemCount)
		padding := contentWidth - len(titleContent) - len(countText) - 2
		if padding > 0 {
			titleContent += strings.Repeat(" ", padding) + countText
		}
	}

	// Pad to content width
	if len(titleContent) < contentWidth {
		titleContent += strings.Repeat(" ", contentWidth-len(titleContent))
	}

	// Build the box
	borderStyle := lipgloss.NewStyle().Foreground(styles.BorderColor)
	contentStyle := lipgloss.NewStyle().Foreground(styles.PrimaryColor).Bold(true)

	topBorder := borderStyle.Render("╔" + strings.Repeat("═", contentWidth) + "╗")
	middleLine := borderStyle.Render("║") + contentStyle.Render(titleContent) + borderStyle.Render("║")
	bottomBorder := borderStyle.Render("╚" + strings.Repeat("═", contentWidth) + "╝")

	result := topBorder + "\n" + middleLine + "\n" + bottomBorder

	// Add subtitle if present
	if h.subtitle != "" {
		subtitleStyle := lipgloss.NewStyle().Foreground(styles.MutedColor)
		result += "\n" + subtitleStyle.Render(h.subtitle)
	}

	return result
}

// renderBanner renders a full-width banner with decorations
func (h ASCIIHeader) renderBanner() string {
	// Spaced title
	spacedTitle := spacedText(h.title)

	// Calculate widths
	contentWidth := h.width - 4
	if contentWidth < 20 {
		contentWidth = 20
	}

	// Create decorative line
	decoStyle := lipgloss.NewStyle().Foreground(styles.BorderColor)
	accentStyle := lipgloss.NewStyle().Foreground(styles.AccentColor)
	titleStyle := lipgloss.NewStyle().Foreground(styles.PrimaryColor).Bold(true)
	iconStyle := lipgloss.NewStyle().Foreground(styles.SecondaryColor)

	// Build banner
	lineWidth := (contentWidth - len(spacedTitle) - 10) / 2
	if lineWidth < 2 {
		lineWidth = 2
	}

	topLine := decoStyle.Render(strings.Repeat("═", lineWidth)) +
		accentStyle.Render(" ✦ ") +
		iconStyle.Render(h.icon) + " " +
		titleStyle.Render(spacedTitle) + " " +
		iconStyle.Render(h.icon) +
		accentStyle.Render(" ✦ ") +
		decoStyle.Render(strings.Repeat("═", lineWidth))

	result := topLine

	// Add item count if present
	if h.itemCount >= 0 {
		countStyle := lipgloss.NewStyle().Foreground(styles.MutedColor)
		countText := fmt.Sprintf("%d items", h.itemCount)
		if h.itemCount == 1 {
			countText = "1 item"
		}
		result += "\n" + countStyle.Render(countText)
	}

	// Add subtitle if present
	if h.subtitle != "" {
		subtitleStyle := lipgloss.NewStyle().Foreground(styles.MutedColor)
		result += "\n" + subtitleStyle.Render(h.subtitle)
	}

	return result
}

// spacedText converts "Notes" to "N O T E S" for vaporwave aesthetic
func spacedText(text string) string {
	if text == "" {
		return ""
	}

	upper := strings.ToUpper(text)
	var result strings.Builder

	runes := []rune(upper)
	for i, r := range runes {
		result.WriteRune(r)
		if i < len(runes)-1 {
			result.WriteString(" ")
		}
	}

	return result.String()
}
