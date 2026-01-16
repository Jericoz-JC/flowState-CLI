# flowState-cli Development Plan

> This file tracks the current development plan and progress. Updated after each phase completion.

## Current Status: Phase 6 Complete, v0.1.9 Released
**Last Updated:** January 16, 2026
**Current Version:** v0.1.9
**Next Target:** v0.1.10 (NPM Install Fixes)

---

## Phase Overview

| Version | Phase | Status | Description |
|---------|-------|--------|-------------|
| v0.1.4 | 1 | ✅ Complete | NPM Package Fix |
| v0.1.5 | 2 | ✅ Complete | Focus Timer UX Enhancement |
| v0.1.6 | 3 | ✅ Complete | Todos Notion-Inspired Overhaul |
| v0.1.7 | 4 | ✅ Complete | Bug Fixes & UX Polish |
| v0.1.8 | 5 | ✅ Complete | Critical Bug Fixes & Layout Issues |
| v0.1.9 | 6 | ✅ Complete | Notes System Overhaul |
| v0.1.10 | 7 | 🔴 Next | NPM Install Fixes (ia32, Linux PATH) |
| v0.1.11 | 8 | ⏳ Pending | Focus Screen Visual Overhaul |
| v0.1.12 | 9 | ⏳ Pending | Unified Theme & Design System |
| v0.1.13 | 10 | ⏳ Pending | Technical Debt Cleanup |
| v0.2.0 | 11 | ⏳ Pending | Final Polish & Documentation |

---

## Development Workflow (MANDATORY)

### Test-Driven Development (TDD)

**All new code MUST follow TDD. No exceptions.**

#### The RED-GREEN-REFACTOR Cycle

1. **RED**: Write a failing test first demonstrating desired behavior
2. **VERIFY RED**: Run test, confirm it fails for the RIGHT reason (missing feature, not syntax error)
3. **GREEN**: Implement the MINIMAL code to pass the test
4. **VERIFY GREEN**: Run test, confirm it passes. Run ALL tests, nothing else breaks.
5. **REFACTOR**: Clean up while maintaining green status

#### Go Test Patterns

```go
// Table-driven tests (preferred)
func TestParseWikilinks(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected []string
    }{
        {"empty string", "", nil},
        {"single wikilink", "see [[Note]]", []string{"Note"}},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := parseWikilinks(tt.input)
            if !reflect.DeepEqual(result, tt.expected) {
                t.Errorf("got %v, want %v", result, tt.expected)
            }
        })
    }
}
```

#### Running Tests
```bash
go test ./...                    # All tests
go test ./internal/tui/screens/  # Specific package
go test -v ./...                 # Verbose
go test -cover ./...             # With coverage
```

#### Red Flags - STOP and RESTART
- You wrote production code before the test
- The test passed immediately (didn't test anything new)
- You can't explain why the test should fail

### Systematic Debugging

When fixing bugs:
1. **Root Cause Investigation** - Read error, reproduce consistently, trace data flow
2. **Pattern Analysis** - Find working examples, identify differences
3. **Hypothesis Testing** - Single change at a time, verify results
4. **Implementation** - Write failing test, implement fix, verify

**Never propose fixes without understanding the root cause first.**

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

## Phase 3: Todos Notion-Inspired Overhaul ✅
**Version:** v0.1.6 | **Status:** Complete

### Features Implemented

| Feature | Implementation |
|---------|----------------|
| **Sort Modes** | 's' key cycles: Date↓ → Priority → Date↑ → A-Z → Due Date |
| **Date Display** | Show due date with relative time in list items |
| **Priority Filter** | 'p' key cycles through All → High → Medium → Low |
| **Tag Support** | Extract #hashtags from title/description, 't' key to filter |
| **Preview Mode** | 'v' key shows full todo details with status/priority badges |
| **Status Badges** | Colored badges: Pending (yellow), In Progress (cyan), Done (green) |
| **Reset Filters** | Ctrl/Cmd+R resets all filters |

### UI Enhancements
- Tags displayed as colored badges in list and preview
- Due date with relative time ("Due in 2 days", "Overdue", "Due today")
- Priority indicators with emoji (🔴 high, 🟢 low)
- Due date indicators (⚠️ overdue, 📅 today, ⏰ soon)
- Sort indicator showing current sort mode
- Active filter status line showing all applied filters

### Files Modified
- `internal/tui/screens/todos.go` - Major overhaul with new features

### Key Implementation Details
- `TodoSortMode` enum with 5 sort modes
- `extractTagsFromTodo()` function extracts #hashtags
- `priorityFilter` uses -1 for "all" to avoid conflict with `TodoPriorityLow = 0`
- `renderPreview()` renders full todo details view
- Filter status line shows active search, status, priority, and tag filters

### Tests
- Existing tests pass
- New features integrated with existing test coverage

---

## Phase 4: Bug Fixes & UX Polish ✅
**Version:** v0.1.7 | **Status:** Complete

### Changes Made

#### Focus Timer Duration Picker Auto-Exit
- Added `autoExitDurationMsg` message type with sequence number for cancellation
- Added `autoExitSequence` field to FocusModel to track timer cancellation
- Modified `applySelectedDuration()` to schedule auto-exit after 500ms
- Duration picker now auto-exits after arrow key selection

#### Notes Preview Edit Shortcut Fix
- Added explicit `m.bodyInput.Blur()` call when entering edit from preview
- Ensures proper focus state when transitioning from preview to edit mode

#### Notes Edit Mode Label Cleanup
- When body is focused, title is now displayed as styled header text
- Title input and label are hidden, showing just the title value
- Shows "(Untitled)" if title is empty

#### Notes Body Enter Key Fix
- Restructured key handling from switch/case to if statements
- Enter now only triggers save when title is focused
- When body is focused, Enter passes through to textarea for newlines

#### Help Modal for Links & Mind Map (? Shortcut)
- Added `LinkModeHelp` mode to links.go with full help content
- Added `showHelp bool` to MindMapModel with help view
- '?' key opens contextual help in both screens
- Any key closes the help modal
- Updated `LinksHints` and `MindMapHints` to include '?' hint

### Files Modified
- `internal/tui/screens/focus.go` - Auto-exit timer implementation
- `internal/tui/screens/notes.go` - Preview 'e' fix, title label cleanup, Enter key fix
- `internal/tui/screens/links.go` - Help modal for Links
- `internal/tui/screens/mindmap.go` - Help modal for Mind Map
- `internal/tui/components/helpbar.go` - Added '?' hints

### Tests
- All existing tests pass
- Build succeeds with no errors

### Release
- Tag: v0.1.7
- Release: https://github.com/Jericoz-JC/flowState-CLI/releases/tag/v0.1.7

---

## Phase 5: Critical Bug Fixes & Layout Issues ✅
**Version:** v0.1.8 | **Status:** Complete

### Functional Bug Fixes

#### Todo Module - Title Bug
- [x] Added comprehensive test suite for todo create flow (todos_test.go)
- [x] Verified title input captures keystrokes correctly
- [x] Verified title saves properly on Enter
- Note: Core functionality verified working - issue may be display-specific

#### Focus Session - Auto-Save & Navigation
- [x] Disabled auto-save on Focus session start
- [x] Sessions now only save to DB when completed (not on start)
- [x] Cancelled sessions are discarded without saving
- [x] Users can freely navigate/adjust before starting timer
- [x] Added tests: `TestFocusStartDoesNotSaveSession`, `TestFocusCancelDoesNotSaveSession`

### Visual & Layout Fixes

#### Home Page ASCII Art
- [x] Added responsive logo system - small logo for terminals < 72 chars
- [x] Added `LogoASCIISmall` constant for narrow terminals
- [x] Added `LogoMinWidth` constant (72) for responsive switching
- [x] Removed "productivity reimagined" tagline from logo

#### TUI Layout & Alignment
- [x] Responsive logo prevents overflow on narrow terminals
- [x] Home view now checks terminal width before rendering logo
- Note: Further alignment issues may need visual testing

#### Copy Cleanup
- [x] Removed "✦ productivity reimagined ✦" line from ASCII art

### Files Modified
- `internal/tui/screens/focus.go` - Auto-save disable, sessions save on completion only
- `internal/tui/screens/focus_test.go` - Added 2 new tests for auto-save behavior
- `internal/tui/screens/todos_test.go` - NEW: 6 tests for todo create flow
- `internal/tui/app.go` - Responsive logo rendering based on terminal width
- `internal/tui/styles/theme.go` - Added small logo, removed tagline

### Tests Added
- `TestFocusStartDoesNotSaveSession` - verifies no DB save on timer start
- `TestFocusCancelDoesNotSaveSession` - verifies cancelled sessions not saved
- `TestTodosScreenRender` - basic render test
- `TestTodosCreateModeEntry` - verifies 'c' enters create mode
- `TestTodosTitleInputCapture` - verifies keystrokes captured in title
- `TestTodosCreateAndSave` - verifies full create and save flow
- `TestTodosEscCancelsCreate` - verifies Esc cancels create
- `TestTodosTabSwitchesFocus` - verifies Tab toggles title/description focus

### Release
- Tag: v0.1.8
- Release: https://github.com/Jericoz-JC/flowState-CLI/releases/tag/v0.1.8

---

## Phase 6: Notes System Overhaul ✅
**Version:** v0.1.9 | **Status:** Complete

### Tagging Engine Overhaul
- [x] Fixed @tag and #tag functionality - both now extract as tags
- [x] Tags extracted from both title AND body (not just body)
- [x] Added `cleanTag()` function for tag normalization
- [x] Implemented "Quick-Tag" picker modal (Ctrl+G in edit mode)
  - Multi-select tags from existing tags
  - Appends selected tags to note body
- [x] Added Tag Filter picker ('t' key in list view)
  - Shows all available tags
  - Multi-select for filtering by multiple tags
  - Pre-selects currently active filters

### Features Added
- **@mention Support**: Both `#hashtag` and `@mention` syntax now create tags
- **Quick-Tag Picker**: Press Ctrl+G while editing to add existing tags
- **Tag Filter Picker**: Press 't' to filter notes by tags (multi-select)
- **Tags from Title**: Tags now extracted from both title and body

### Files Modified
- `internal/tui/screens/notes.go` - Major overhaul with tag picker
- `internal/tui/screens/notes_test.go` - Added tag extraction tests

### Tests Added
- `TestExtractTagsHashtag` - verifies #hashtag extraction
- `TestExtractTagsAtSign` - verifies @mention extraction
- `TestExtractTagsFromTitle` - verifies combined title+body extraction

### Release
- Tag: v0.1.9
- Release: https://github.com/Jericoz-JC/flowState-CLI/releases/tag/v0.1.9

---

## Phase 7: NPM Install Fixes 🔴 NEXT
**Version:** v0.1.10 | **Status:** In Progress

### Issues to Fix

#### 1. ia32/x86 CPU Error on Windows
**Problem**: Users with 32-bit Node.js on 64-bit Windows get unhelpful error:
```
Unsupported platform: win32-ia32
```

**Root Cause**: `npm/install.js` only maps `x64` and `arm64` in `archMap`. When 32-bit Node.js is installed, `process.arch` returns `ia32`.

**Fix**: Add clear error message with instructions to install 64-bit Node.js.

**Location**: `npm/install.js:23-34`

#### 2. Linux SSH - Can't Run App After Install
**Problem**: Users install via npm on Linux but `flowstate` command not found.

**Root Cause**: npm global bin directory not in PATH.

**Fix**: Show PATH setup instructions after successful install.

**Location**: `npm/install.js` (success message)

#### 3. Binary Not Found Error Handling
**Problem**: If binary download fails silently, wrapper gives unhelpful error.

**Fix**: Add existence check in `npm/bin/flowstate` with clear reinstall instructions.

**Location**: `npm/bin/flowstate`

#### 4. Dev Machine Old Version Conflict
**Problem**: Running `flowstate` opens old version, `go build` gives latest.

**Root Cause**: Old global npm install or PATH ordering.

**Documentation**: Add troubleshooting section for PATH conflicts.

### Checklist
- [x] Update `npm/install.js` with ia32 error handling
- [x] Add success message with PATH instructions for Linux
- [x] Update `npm/bin/flowstate` with binary existence check
- [x] Update `npm/package.json` version to 0.1.10
- [ ] Test on Windows (64-bit Node)
- [ ] Test on Linux
- [ ] Test on macOS
- [ ] Create release tag v0.1.10

---

## Technical Debt Analysis

### Priority: High (Address before v0.2.0)

#### 1. Placeholder Embedding System
**Location**: `internal/embeddings/embedder.go:101-144`
```go
// embedSimple creates simple hash-based embeddings.
// Future: Replace with ONNX inference
func (e *Embedder) embedSimple(texts []string) ([][]float32, error)
```
**Issue**: Using character-weighted hash instead of real ML embeddings. Semantic search produces poor results.

**Resolution**: Phase 10+ ONNX model integration (~90MB model).

#### 2. Missing Test Coverage
**Packages without tests**:
- `internal/commands` - Command wrappers
- `internal/config` - Configuration loading
- `internal/models` - Data models
- `internal/storage/qdrant` - Vector store (dead code)
- `internal/tui/keymap` - Key bindings
- `internal/tui/styles` - Theme

**Current coverage**: 17 test files

### Priority: Medium

#### 3. Large Screen Files
**Files**:
- `notes.go` (~1,400 lines)
- `todos.go` (~1,119 lines)

**Issue**: Model, update logic, and view rendering all in one file.

**Potential Refactor**: Extract common patterns (tag filtering, sort modes) into shared utilities.

#### 4. Duplicate Vector Storage
**Locations**:
- `internal/storage/qdrant/vector_store.go` - In-memory (data lost on restart)
- `internal/storage/sqlite/store.go` - SQLite-backed `note_vectors` table

**Issue**: Qdrant package is effectively dead code.

**Resolution**: Remove `qdrant/` package or document as development stub.

### Priority: Low

#### 5. TODO Comments in Code
```
notes.go:472:  // TODO: proper cursor handling
notes.go:562:  // TODO: Show error message
```
Only 2 TODOs - relatively clean.

#### 6. Error Handling in Main
**Location**: `cmd/flowState/main.go`
```go
log.Fatalf("Failed to load config: %v", err)  // line 56
log.Fatalf("Failed to create app: %v", err)   // line 62
```
Uses `log.Fatal` - prevents graceful shutdown but acceptable for CLI.

#### 7. Hardcoded Strings
Screen titles, help text, and UI labels scattered rather than centralized.

---

## Phase 8: Focus Screen Visual Overhaul
**Version:** v0.1.11 | **Status:** Pending

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

## Phase 9: Unified Theme & Design System
**Version:** v0.1.12 | **Status:** Pending

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

## Phase 10: Technical Debt Cleanup
**Version:** v0.1.13 | **Status:** Pending

### High Priority Items
- [ ] Add test coverage for `internal/commands`
- [ ] Add test coverage for `internal/config`
- [ ] Add test coverage for `internal/models`
- [ ] Add test coverage for `internal/tui/keymap`
- [ ] Add test coverage for `internal/tui/styles`
- [ ] Remove or document `internal/storage/qdrant` package

### Medium Priority Items
- [ ] Extract common filter logic from `notes.go` and `todos.go`
- [ ] Create shared `internal/tui/filters/` package
- [ ] Fix TODO comments in `notes.go`

### Optional
- [ ] Centralize hardcoded strings
- [ ] Improve error handling in `main.go`

---

## Phase 11: Final Polish & Documentation
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
| `internal/tui/app.go` | Main app, layout, navigation |
| `internal/tui/screens/home.go` | Home screen with ASCII art |
| `internal/tui/screens/focus.go` | Focus/Pomodoro screen |
| `internal/tui/screens/todos.go` | Todos screen |
| `internal/tui/screens/notes.go` | Notes screen (Phase 6 overhaul) |
| `internal/tui/styles/theme.go` | Theme system |
| `internal/tui/components/helpbar.go` | Help bar hints |
| `internal/models/note.go` | Note model with tags |

---

## Release Workflow

For each phase:
1. Implement features
2. Write/update unit tests (TDD)
3. Run full test suite: `go test ./...`
4. Update CLAUDE.md with progress
5. Update README.md with new features
6. **Update `npm/package.json` version** to match the new version
7. Commit with descriptive message
8. Push to main AND create/push tag to trigger release:
   ```bash
   git push origin main
   git tag v0.1.X
   git push origin v0.1.X
   ```
9. CI/CD automatically:
   - Builds binaries for all 6 platforms
   - Creates GitHub release with binaries
   - Publishes to npm
10. Verify installation on all platforms

**IMPORTANT:**
- Always tag and push after completing a phase. The CI/CD pipeline only triggers on tag push (`v*`), not on main branch push.
- **Always update `npm/package.json` version** before tagging! npm publish will fail if the version already exists.

---

## Notes for Future Sessions

If credits run out mid-phase:
1. Check this file for current progress
2. Look at the phase status and checklist
3. Continue from where left off
4. Update this file after completing work
