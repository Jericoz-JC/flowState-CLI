// Package keymap provides cross-platform keyboard shortcut handling.
//
// Phase 4: UX Overhaul - Cross-platform support
//   - Detects OS and maps Ctrl (Windows/Linux) vs Cmd (macOS)
//   - Provides consistent key binding checks across platforms
package keymap

import (
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// IsMacOS returns true if running on macOS.
func IsMacOS() bool {
	return runtime.GOOS == "darwin"
}

// IsWindows returns true if running on Windows.
func IsWindows() bool {
	return runtime.GOOS == "windows"
}

// IsLinux returns true if running on Linux.
func IsLinux() bool {
	return runtime.GOOS == "linux"
}

// ModKey returns the modifier key name for the current platform.
// Returns "cmd" on macOS, "ctrl" on Windows/Linux.
func ModKey() string {
	if IsMacOS() {
		return "cmd"
	}
	return "ctrl"
}

// IsModN checks if the key message is Ctrl+N (or Cmd+N on macOS).
func IsModN(msg tea.KeyMsg) bool {
	key := strings.ToLower(msg.String())
	if IsMacOS() {
		return key == "cmd+n" || key == "ctrl+n"
	}
	return key == "ctrl+n"
}

// IsModT checks if the key message is Ctrl+T (or Cmd+T on macOS).
func IsModT(msg tea.KeyMsg) bool {
	key := strings.ToLower(msg.String())
	if IsMacOS() {
		return key == "cmd+t" || key == "ctrl+t"
	}
	return key == "ctrl+t"
}

// IsModF checks if the key message is Ctrl+F (or Cmd+F on macOS).
func IsModF(msg tea.KeyMsg) bool {
	key := strings.ToLower(msg.String())
	if IsMacOS() {
		return key == "cmd+f" || key == "ctrl+f"
	}
	return key == "ctrl+f"
}

// IsModH checks if the key message is Ctrl+H (or Cmd+H on macOS).
func IsModH(msg tea.KeyMsg) bool {
	key := strings.ToLower(msg.String())
	if IsMacOS() {
		// Note: Cmd+H hides windows on macOS, so we allow both
		return key == "cmd+h" || key == "ctrl+h"
	}
	return key == "ctrl+h"
}

// IsModX checks if the key message is Ctrl+X (or Cmd+X on macOS).
func IsModX(msg tea.KeyMsg) bool {
	key := strings.ToLower(msg.String())
	if IsMacOS() {
		return key == "cmd+x" || key == "ctrl+x"
	}
	return key == "ctrl+x"
}

// IsModL checks if the key message is Ctrl+L (or Cmd+L on macOS).
func IsModL(msg tea.KeyMsg) bool {
	key := strings.ToLower(msg.String())
	if IsMacOS() {
		return key == "cmd+l" || key == "ctrl+l"
	}
	return key == "ctrl+l"
}

// IsModS checks if the key message is Ctrl+S (or Cmd+S on macOS).
func IsModS(msg tea.KeyMsg) bool {
	key := strings.ToLower(msg.String())
	if IsMacOS() {
		return key == "cmd+s" || key == "ctrl+s"
	}
	return key == "ctrl+s"
}

// IsModR checks if the key message is Ctrl+R (or Cmd+R on macOS).
func IsModR(msg tea.KeyMsg) bool {
	key := strings.ToLower(msg.String())
	if IsMacOS() {
		return key == "cmd+r" || key == "ctrl+r"
	}
	return key == "ctrl+r"
}

// IsModSlash checks if the key message is Ctrl+/ (or Cmd+/ on macOS).
func IsModSlash(msg tea.KeyMsg) bool {
	key := strings.ToLower(msg.String())
	if IsMacOS() {
		return key == "cmd+/" || key == "ctrl+/"
	}
	return key == "ctrl+/"
}

// IsModG checks if the key message is Ctrl+G (or Cmd+G on macOS).
func IsModG(msg tea.KeyMsg) bool {
	key := strings.ToLower(msg.String())
	if IsMacOS() {
		return key == "cmd+g" || key == "ctrl+g"
	}
	return key == "ctrl+g"
}

// ModKeyDisplay returns the display string for the modifier key.
// Returns "⌘" on macOS, "Ctrl" on Windows/Linux.
func ModKeyDisplay() string {
	if IsMacOS() {
		return "⌘"
	}
	return "Ctrl"
}

// FormatShortcut formats a keyboard shortcut for display.
// Example: FormatShortcut("N") returns "Ctrl+N" on Windows/Linux, "⌘N" on macOS.
func FormatShortcut(key string) string {
	if IsMacOS() {
		return "⌘" + key
	}
	return "Ctrl+" + key
}

