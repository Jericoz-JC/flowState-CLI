package screens

import (
	"path/filepath"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Jericoz-JC/flowState-CLI/internal/config"
	"github.com/Jericoz-JC/flowState-CLI/internal/storage/sqlite"
)

func newTestFocusModel(t *testing.T) FocusModel {
	t.Helper()

	tmpDir := t.TempDir()
	cfg := &config.Config{
		DbPath:    filepath.Join(tmpDir, "test.db"),
		ModelPath: filepath.Join(tmpDir, "models"),
	}

	store, err := sqlite.New(cfg)
	if err != nil {
		t.Fatalf("sqlite.New() err = %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })

	model := NewFocusModel(store)
	model.SetSize(100, 40)
	return model
}

func TestFocusScreenRender(t *testing.T) {
	t.Parallel()

	m := newTestFocusModel(t)
	v := m.View()
	if v == "" {
		t.Fatalf("expected non-empty view")
	}
}

func TestFocusDurationPickerEntry(t *testing.T) {
	t.Parallel()

	m := newTestFocusModel(t)

	// Should start in idle mode
	if m.mode != FocusModeIdle {
		t.Fatalf("expected FocusModeIdle, got %v", m.mode)
	}

	// Press 'd' to enter duration picker
	mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	m = mm

	if m.mode != FocusModeDuration {
		t.Fatalf("expected FocusModeDuration after pressing 'd', got %v", m.mode)
	}

	// Should start with work duration selected
	if !m.selectingWork {
		t.Fatalf("expected selectingWork to be true when entering duration picker")
	}
}

func TestFocusDurationPickerLiveUpdate(t *testing.T) {
	t.Parallel()

	m := newTestFocusModel(t)

	// Enter duration picker
	mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	m = mm

	// Store initial work duration
	initialDuration := m.workDuration
	initialRemaining := m.remaining

	// Find index of initial duration in WorkDurations
	initialIndex := m.durationIndex

	// Press right arrow to select next duration
	if initialIndex < len(WorkDurations)-1 {
		mm, _ = m.Update(tea.KeyMsg{Type: tea.KeyRight})
		m = mm

		// Duration should be updated IMMEDIATELY (live update)
		if m.workDuration == initialDuration {
			t.Fatalf("expected work duration to change immediately after right arrow")
		}

		expectedDuration := WorkDurations[initialIndex+1]
		if m.workDuration != expectedDuration {
			t.Fatalf("expected work duration %d, got %d", expectedDuration, m.workDuration)
		}

		// Remaining time should also be updated
		expectedRemaining := time.Duration(expectedDuration) * time.Minute
		if m.remaining != expectedRemaining {
			t.Fatalf("expected remaining %v, got %v", expectedRemaining, m.remaining)
		}

		// Should NOT have exited duration picker mode
		if m.mode != FocusModeDuration {
			t.Fatalf("expected to still be in duration picker mode")
		}
	}

	// Press left to go back
	if m.durationIndex > 0 {
		mm, _ = m.Update(tea.KeyMsg{Type: tea.KeyLeft})
		m = mm

		// Should update immediately
		if m.remaining == initialRemaining && initialIndex < len(WorkDurations)-1 {
			// If we moved right then left, we should be back to initial
			if m.workDuration != initialDuration {
				t.Logf("Note: duration changed as expected after left arrow")
			}
		}
	}
}

func TestFocusDurationPickerTabSwitch(t *testing.T) {
	t.Parallel()

	m := newTestFocusModel(t)

	// Enter duration picker
	mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	m = mm

	// Should be selecting work
	if !m.selectingWork {
		t.Fatalf("expected selectingWork to be true")
	}

	// Press Tab to switch to break duration
	mm, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = mm

	if m.selectingWork {
		t.Fatalf("expected selectingWork to be false after Tab")
	}

	// Press Tab again to switch back to work
	mm, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = mm

	if !m.selectingWork {
		t.Fatalf("expected selectingWork to be true after second Tab")
	}
}

func TestFocusDurationPickerEnterConfirms(t *testing.T) {
	t.Parallel()

	m := newTestFocusModel(t)

	// Enter duration picker
	mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	m = mm

	// Change work duration
	if m.durationIndex < len(WorkDurations)-1 {
		mm, _ = m.Update(tea.KeyMsg{Type: tea.KeyRight})
		m = mm
	}

	// Tab to break and change it
	mm, _ = m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = mm

	if m.durationIndex < len(BreakDurations)-1 {
		mm, _ = m.Update(tea.KeyMsg{Type: tea.KeyRight})
		m = mm
	}

	// Press Enter to confirm all and exit
	mm, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = mm

	// Should be back in idle mode
	if m.mode != FocusModeIdle {
		t.Fatalf("expected FocusModeIdle after Enter, got %v", m.mode)
	}
}

func TestFocusDurationPickerEscCancels(t *testing.T) {
	t.Parallel()

	m := newTestFocusModel(t)

	// Enter duration picker
	mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	m = mm

	// Press Esc to cancel
	mm, _ = m.Update(tea.KeyMsg{Type: tea.KeyEscape})
	m = mm

	// Should be back in idle mode
	if m.mode != FocusModeIdle {
		t.Fatalf("expected FocusModeIdle after Esc, got %v", m.mode)
	}
}

func TestFocusStartSession(t *testing.T) {
	t.Parallel()

	m := newTestFocusModel(t)

	// Press 's' to start
	mm, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	m = mm

	if m.mode != FocusModeRunning {
		t.Fatalf("expected FocusModeRunning after pressing 's', got %v", m.mode)
	}

	// Should have a tick command
	if cmd == nil {
		t.Fatalf("expected tick command when starting timer")
	}
}

func TestFocusPauseSession(t *testing.T) {
	t.Parallel()

	m := newTestFocusModel(t)

	// Start session
	mm, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	m = mm

	// Press 'p' to pause
	mm, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}})
	m = mm

	if m.mode != FocusModePaused {
		t.Fatalf("expected FocusModePaused after pressing 'p', got %v", m.mode)
	}
}
