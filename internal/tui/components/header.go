// Package components provides reusable TUI components for flowState-cli.
//
// Phase 4: UX Overhaul
//   - Header: Screen title with breadcrumb navigation
//   - Provides consistent navigation context across all screens
package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/Jericoz-JC/flowState-CLI/internal/tui/styles"
)

// Breadcrumb represents a navigation path element.
type Breadcrumb struct {
	Icon  string
	Title string
}

// Header provides a screen title with breadcrumb navigation.
type Header struct {
	icon       string
	title      string
	breadcrumb []Breadcrumb
	itemCount  int
	width      int
}

// NewHeader creates a new header with title.
func NewHeader(icon, title string) Header {
	return Header{
		icon:       icon,
		title:      title,
		breadcrumb: []Breadcrumb{},
		itemCount:  -1, // -1 means don't show count
		width:      80,
	}
}

// SetWidth updates the width constraint for rendering.
func (h *Header) SetWidth(width int) {
	h.width = width
}

// SetBreadcrumb sets the navigation path.
func (h *Header) SetBreadcrumb(crumbs []Breadcrumb) {
	h.breadcrumb = crumbs
}

// SetItemCount sets the item count to display (set to -1 to hide).
func (h *Header) SetItemCount(count int) {
	h.itemCount = count
}

// SetTitle updates the current title.
func (h *Header) SetTitle(icon, title string) {
	h.icon = icon
	h.title = title
}

// View renders the header.
func (h *Header) View() string {
	// Styles using centralized ARCHWAVE theme
	titleStyle := lipgloss.NewStyle().
		Foreground(styles.PrimaryColor).
		Bold(true)

	iconStyle := lipgloss.NewStyle().
		Foreground(styles.SecondaryColor)

	breadcrumbStyle := lipgloss.NewStyle().
		Foreground(styles.MutedColor)

	countStyle := lipgloss.NewStyle().
		Foreground(styles.MutedColor).
		Italic(true)

	// Build title line
	titleLine := iconStyle.Render(h.icon) + " " + titleStyle.Render(h.title)

	// Add item count if set
	if h.itemCount >= 0 {
		countText := fmt.Sprintf("%d items", h.itemCount)
		if h.itemCount == 1 {
			countText = "1 item"
		}
		// Calculate padding to right-align count
		titleLen := len(h.icon) + 1 + len(h.title)
		countLen := len(countText)
		padding := h.width - titleLen - countLen - 6 // Account for border/padding
		if padding > 0 {
			titleLine += strings.Repeat(" ", padding) + countStyle.Render(countText)
		}
	}

	// Build breadcrumb line if we have crumbs
	var breadcrumbLine string
	if len(h.breadcrumb) > 0 {
		var parts []string
		for _, crumb := range h.breadcrumb {
			parts = append(parts, crumb.Icon+" "+crumb.Title)
		}
		parts = append(parts, h.icon+" "+h.title)
		breadcrumbLine = breadcrumbStyle.Render(strings.Join(parts, " " + styles.DecoArrow + " "))
	}

	// Build vaporwave divider
	dividerLen := h.width - 4 // Account for padding
	if dividerLen < 10 {
		dividerLen = 10
	}
	divider := styles.VaporwaveDivider(dividerLen)

	// Combine parts
	if breadcrumbLine != "" {
		return lipgloss.JoinVertical(
			lipgloss.Left,
			titleLine,
			breadcrumbLine,
			divider,
		)
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		titleLine,
		divider,
	)
}
