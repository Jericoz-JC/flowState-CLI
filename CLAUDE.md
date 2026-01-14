# flowState-cli Development Plan

> This file tracks the current development plan and progress. Updated after each phase completion.

## Current Status: Phase 2 Complete, Phase 3 Next
**Last Updated:** January 14, 2025
**Current Version:** v0.1.5
**Next Target:** v0.1.6

---

## Phase Overview

| Version | Phase | Status | Description |
|---------|-------|--------|-------------|
| v0.1.4 | 1 | ✅ Complete | NPM Package Fix |
| v0.1.5 | 2 | ✅ Complete | Focus Timer UX Enhancement |
| v0.1.6 | 3 | 🔄 Next | Todos Notion-Inspired Overhaul |
| v0.1.7 | 4 | ⏳ Pending | Focus Screen Visual Overhaul |
| v0.1.8 | 5 | ⏳ Pending | Unified Theme & Design System |
| v0.2.0 | 6 | ⏳ Pending | Final Polish & Documentation |

---

## Phase 1: NPM Package Fix ✅
**Version:** v0.1.4 | **Status:** Complete

### Changes Made
- Fixed `npm/install.js` to read version from package.json dynamically
- Improved error messages with platform-specific troubleshooting guidance
- Built and uploaded binaries for all 6 platforms
- Release: https://github.com/Jericoz-JC/flowState-CLI/releases/tag/v0.1.4

---

## Phase 2: Focus Timer UX Enhancement ✅
**Version:** v0.1.5 | **Status:** Complete

### Changes Made
- Added `durationJustChanged` and `lastChangedField` to FocusModel for visual feedback
- Added `clearFeedbackMsg` message type for auto-clearing feedback after 800ms
- Updated `applySelectedDuration()` to trigger visual feedback
- Updated `renderDurationPicker()` to show "✓ Saved" indicator when duration changes
- Added "Current: X min work / Y min break" summary line showing both values
- Updated help hints to say "Adjust (auto-saves)" instead of "Adjust (live)"

### Files Modified
- `internal/tui/screens/focus.go` - Visual feedback implementation
- `internal/tui/screens/focus_test.go` - Added new tests
- `internal/tui/components/helpbar.go` - Updated help hints
- `internal/tui/components/helpbar_test.go` - Updated test for new hint wording

### Tests Added
- `TestFocusDurationPickerVisualFeedback` - verifies feedback flag and command
- `TestFocusDurationPickerShowsBothValues` - verifies both durations displayed

---

## Phase 3: Todos Notion-Inspired Overhaul
**Version:** v0.1.6 | **Status:** Pending

### Features to Add

| Feature | Implementation |
|---------|----------------|
| **Sort Modes** | 's' key cycles: Date↓ → Priority → Date↑ → Alphabetical |
| **Date Display** | Show due date and creation date in list items |
| **Priority Filter** | 'p' key cycles through priority levels |
| **Tag Support** | Extract #hashtags from description, 't' key to filter |
| **Preview Mode** | 'p' key shows todo details with markdown rendering |
| **Markdown Description** | Support for formatted descriptions |
| **Status Badges** | Colored badges: Pending (yellow), In Progress (cyan), Done (green) |

### UI Enhancements
- Visual cards for each todo item
- Tags displayed as colored badges
- Due date with relative time ("Due in 2 days")
- Priority indicators with colors (High=red, Medium=yellow, Low=green)

### Files to Modify
- `internal/tui/screens/todos.go` - Major overhaul
- `internal/tui/styles/theme.go` - Add todo-specific styles

### Tests
- [ ] `TestTodosSortModes`
- [ ] `TestTodosTagExtraction`
- [ ] `TestTodosFilterByPriority`
- [ ] `TestTodosPreviewMode`

---

## Phase 4: Focus Screen Visual Overhaul
**Version:** v0.1.7 | **Status:** Pending

### Enhanced Timer Display
- Large ASCII art numbers for countdown
- Animated progress ring/bar
- Session counter with visual indicators
- Break vs Work mode visual distinction

### Session History Improvements
- Graph/chart of recent sessions
- Statistics panel (total focus time, streaks)
- Calendar heatmap of activity

### Visual Polish
- Gradient backgrounds for different modes
- Smooth transitions between states
- Sound notification indicators (visual bell)

### Files to Modify
- `internal/tui/screens/focus.go`
- `internal/tui/styles/theme.go`

---

## Phase 5: Unified Theme & Design System
**Version:** v0.1.8 | **Status:** Pending

### Charmbracelet Library Integration
```go
// Add to go.mod
github.com/charmbracelet/glamour   // Markdown rendering
github.com/charmbracelet/harmonica // Animations
```

### Design System Enhancements

#### ASCII Art Headers
- Consistent ASCII art headers for each screen
- flowState logo variations for different contexts

#### Animation System
- Page transition animations
- Loading spinners for async operations
- Smooth scroll animations
- Fade effects for modals

#### ARCHWAVE Theme Polish
- Consistent color usage across all screens
- Unified border styles (double borders for primary, rounded for secondary)
- Gradient text for titles
- Neon glow effects for focused elements

### Component Library
- [ ] `AnimatedSpinner` - for loading states
- [ ] `GradientProgress` - animated progress bars
- [ ] `GlowBorder` - neon border effect
- [ ] `FadeModal` - animated modal component
- [ ] `ASCIIHeader` - screen headers with art

---

## Phase 6: Final Polish & Documentation
**Version:** v0.2.0 | **Status:** Pending

### Cross-Screen Consistency Audit
- Verify all screens use same patterns
- Consistent keybindings
- Unified help bar format
- Matching visual hierarchy

### Performance Optimization
- Benchmark all screens
- Optimize render loops
- Lazy load heavy components

### Documentation
- Update README with all new features
- Add keyboard shortcut reference
- Screenshot gallery
- Demo GIF

---

## Key Files Reference

| File | Purpose |
|------|---------|
| `npm/install.js` | Binary download script |
| `npm/package.json` | NPM package config |
| `.goreleaser.yaml` | Release build config |
| `internal/tui/screens/focus.go` | Focus/Pomodoro screen |
| `internal/tui/screens/todos.go` | Todos screen |
| `internal/tui/screens/notes.go` | Notes screen (reference) |
| `internal/tui/styles/theme.go` | Theme system |
| `internal/tui/components/helpbar.go` | Help bar hints |

---

## Release Workflow

For each phase:
1. Implement features
2. Write/update unit tests (TDD)
3. Run full test suite: `go test ./...`
4. Update CLAUDE.md with progress
5. Update README.md with new features
6. Commit with descriptive message
7. Push to main
8. Create tag (e.g., v0.1.5)
9. Build binaries and create GitHub release
10. Verify installation on all platforms

---

## Notes for Future Sessions

If credits run out mid-phase:
1. Check this file for current progress
2. Look at the phase status and checklist
3. Continue from where left off
4. Update this file after completing work
