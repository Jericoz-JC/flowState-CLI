package components

import (
	"strings"
	"testing"
	"time"
)

func TestNewAnimatedSpinner(t *testing.T) {
	s := NewAnimatedSpinner()

	if s.IsActive() {
		t.Error("spinner should not be active initially")
	}

	if len(s.frames) == 0 {
		t.Error("spinner should have frames")
	}
}

func TestSpinnerFrameCycle(t *testing.T) {
	s := NewAnimatedSpinner()
	s.Start()

	initialFrame := s.current

	// Simulate tick
	s, _ = s.Update(SpinnerTickMsg{})

	if s.current == initialFrame && len(s.frames) > 1 {
		t.Error("spinner frame should advance on tick")
	}

	// Cycle through all frames
	for i := 0; i < len(s.frames); i++ {
		s, _ = s.Update(SpinnerTickMsg{})
	}

	// Should wrap around
	if s.current >= len(s.frames) {
		t.Error("spinner frame should wrap around")
	}
}

func TestSpinnerView(t *testing.T) {
	s := NewAnimatedSpinner()

	view := s.View()
	if view == "" {
		t.Error("spinner view should not be empty")
	}

	// View should contain current frame character
	found := false
	for _, frame := range s.frames {
		if strings.Contains(view, frame) {
			found = true
			break
		}
	}
	if !found {
		t.Error("spinner view should contain a frame character")
	}
}

func TestSpinnerViewWithLabel(t *testing.T) {
	s := NewAnimatedSpinner()
	s.SetLabel("Loading notes...")

	view := s.View()
	if !strings.Contains(view, "Loading notes...") {
		t.Error("spinner view should contain the label")
	}
}

func TestSpinnerStartStop(t *testing.T) {
	s := NewAnimatedSpinner()

	if s.IsActive() {
		t.Error("spinner should not be active initially")
	}

	cmd := s.Start()
	if cmd == nil {
		t.Error("Start should return a tick command")
	}
	if !s.IsActive() {
		t.Error("spinner should be active after Start")
	}

	s.Stop()
	if s.IsActive() {
		t.Error("spinner should not be active after Stop")
	}
}

func TestSpinnerTickCommand(t *testing.T) {
	s := NewAnimatedSpinner()
	s.Start()

	// Tick should return another tick command when active
	_, cmd := s.Update(SpinnerTickMsg{})
	if cmd == nil {
		t.Error("active spinner should return tick command")
	}

	// Stop and tick again
	s.Stop()
	_, cmd = s.Update(SpinnerTickMsg{})
	if cmd != nil {
		t.Error("stopped spinner should not return tick command")
	}
}

func TestSpinnerInterval(t *testing.T) {
	s := NewAnimatedSpinner()

	// Default interval should be reasonable (80-150ms typical)
	if s.interval < 50*time.Millisecond || s.interval > 200*time.Millisecond {
		t.Errorf("spinner interval %v outside reasonable range", s.interval)
	}
}

func TestVaporwaveFrames(t *testing.T) {
	// Verify our vaporwave frames are defined
	if len(VaporwaveSpinnerFrames) == 0 {
		t.Error("VaporwaveSpinnerFrames should be defined")
	}

	// Each frame should be non-empty
	for i, frame := range VaporwaveSpinnerFrames {
		if frame == "" {
			t.Errorf("frame %d is empty", i)
		}
	}
}

func TestSpinnerCustomFrames(t *testing.T) {
	customFrames := []string{"A", "B", "C"}
	s := NewAnimatedSpinnerWithFrames(customFrames)

	if len(s.frames) != len(customFrames) {
		t.Errorf("expected %d frames, got %d", len(customFrames), len(s.frames))
	}
}
