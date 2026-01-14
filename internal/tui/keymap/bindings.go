// Package keymap provides cross-platform keyboard shortcut handling.
//
// Phase 4: UX Overhaul - Unified Keybinding Definitions
//   - Centralized keybinding constants for consistency
//   - Standard action descriptions
//   - Used by helpbar hints and documentation
package keymap

// Standard keybinding constants
// These define the canonical keybindings used throughout the app
const (
	// Global Navigation (work from any screen)
	KeyHome        = "Ctrl+H" // Navigate to Home screen
	KeyNotes       = "Ctrl+N" // Navigate to Notes screen
	KeyTodos       = "Ctrl+T" // Navigate to Todos screen
	KeyFocus       = "Ctrl+F" // Navigate to Focus screen
	KeySearch      = "Ctrl+/" // Navigate to Search screen
	KeyMindMap     = "Ctrl+G" // Navigate to Mind Map screen
	KeyQuickCap    = "Ctrl+X" // Open Quick Capture modal
	KeyLinks       = "Ctrl+L" // Open Links modal
	KeyHelp        = "?"      // Toggle help modal
	KeyQuit        = "q"      // Quit application

	// Edit Mode Actions
	KeySave   = "Ctrl+S" // Save current item
	KeyCancel = "Esc"    // Cancel/back
	KeyTab    = "Tab"    // Next field
	KeyShiftTab = "Shift+Tab" // Previous field

	// List Actions
	KeyCreate = "c"     // Create new item
	KeyEdit   = "e"     // Edit selected item
	KeyDelete = "d"     // Delete selected item
	KeyFilter = "/"     // Open filter/search
	KeyReset  = "Ctrl+R" // Reset filters
	KeyToggle = "Space" // Toggle (checkbox, etc.)
	KeyPreview = "p"    // Preview item

	// Navigation
	KeyUp    = "k"       // Move up (vim)
	KeyDown  = "j"       // Move down (vim)
	KeyLeft  = "h"       // Move left (vim)
	KeyRight = "l"       // Move right (vim)
	KeyEnter = "Enter"   // Confirm/select
)

// Binding represents a keyboard shortcut with its description
type Binding struct {
	Key         string
	Description string
	Primary     bool // Whether this is a primary action
}

// GlobalBindings are shortcuts that work from any screen
var GlobalBindings = []Binding{
	{Key: KeyHome, Description: "Home", Primary: false},
	{Key: KeyNotes, Description: "Notes", Primary: true},
	{Key: KeyTodos, Description: "Todos", Primary: true},
	{Key: KeyFocus, Description: "Focus", Primary: false},
	{Key: KeySearch, Description: "Search", Primary: false},
	{Key: KeyMindMap, Description: "Mind Map", Primary: false},
	{Key: KeyQuickCap, Description: "Quick Capture", Primary: true},
	{Key: KeyHelp, Description: "Help", Primary: false},
	{Key: KeyQuit, Description: "Quit", Primary: false},
}

// EditBindings are shortcuts for edit/create forms
var EditBindings = []Binding{
	{Key: KeySave, Description: "Save", Primary: true},
	{Key: KeyTab, Description: "Next Field", Primary: false},
	{Key: KeyCancel, Description: "Cancel", Primary: false},
}

// ListBindings are shortcuts for list views
var ListBindings = []Binding{
	{Key: KeyCreate, Description: "Create", Primary: true},
	{Key: KeyEdit, Description: "Edit", Primary: false},
	{Key: KeyDelete, Description: "Delete", Primary: false},
	{Key: KeyFilter, Description: "Filter", Primary: false},
	{Key: KeyPreview, Description: "Preview", Primary: false},
}

// GetDisplayKey returns the platform-appropriate display key
// For modifier keys, uses ModKeyDisplay() to show ⌘ on macOS, Ctrl on others
func GetDisplayKey(key string) string {
	// Replace "Ctrl+" prefix with platform-appropriate modifier
	if len(key) > 5 && key[:5] == "Ctrl+" {
		return ModKeyDisplay() + "+" + key[5:]
	}
	return key
}
