package components

import (
	"strings"
	"testing"
)

func TestQuickCaptureHints(t *testing.T) {
	t.Parallel()

	// Verify QuickCaptureHints shows Ctrl+S (not Enter) for save
	foundCtrlS := false
	foundEnterForSave := false

	for _, hint := range QuickCaptureHints {
		if hint.Key == "Ctrl+S" && hint.Description == "Save" {
			foundCtrlS = true
		}
		if hint.Key == "Enter" && hint.Description == "Save" {
			foundEnterForSave = true
		}
	}

	if !foundCtrlS {
		t.Errorf("expected QuickCaptureHints to have Ctrl+S for Save")
	}

	if foundEnterForSave {
		t.Errorf("QuickCaptureHints should NOT have Enter for Save (was incorrectly showing this before)")
	}
}

func TestQuickCaptureHintsHasEsc(t *testing.T) {
	t.Parallel()

	foundEsc := false
	for _, hint := range QuickCaptureHints {
		if hint.Key == "Esc" {
			foundEsc = true
			break
		}
	}

	if !foundEsc {
		t.Errorf("expected QuickCaptureHints to have Esc for cancel")
	}
}

func TestFocusDurationHintsShowLiveUpdate(t *testing.T) {
	t.Parallel()

	// Verify FocusDurationHints indicates live/auto-save update behavior
	foundLiveHint := false
	for _, hint := range FocusDurationHints {
		desc := strings.ToLower(hint.Description)
		if strings.Contains(desc, "live") || strings.Contains(desc, "auto") {
			foundLiveHint = true
			break
		}
	}

	if !foundLiveHint {
		t.Errorf("expected FocusDurationHints to indicate live/auto-save update behavior")
	}
}

func TestHelpBarRender(t *testing.T) {
	t.Parallel()

	hints := []HelpHint{
		{Key: "a", Description: "Action A", Primary: true},
		{Key: "b", Description: "Action B", Primary: false},
	}

	bar := NewHelpBar(hints)
	bar.SetWidth(80)

	view := bar.View()
	if view == "" {
		t.Errorf("expected non-empty view from HelpBar")
	}

	// Should contain the keys
	if !strings.Contains(view, "a") {
		t.Errorf("expected view to contain key 'a'")
	}
	if !strings.Contains(view, "b") {
		t.Errorf("expected view to contain key 'b'")
	}
}

func TestNotesEditHints(t *testing.T) {
	t.Parallel()

	// Verify NotesEditHints has proper save hint
	foundCtrlS := false
	for _, hint := range NotesEditHints {
		if hint.Key == "Ctrl+S" && hint.Primary {
			foundCtrlS = true
			break
		}
	}

	if !foundCtrlS {
		t.Errorf("expected NotesEditHints to have Ctrl+S as primary save action")
	}
}

func TestAllHintSetsHaveRequiredKeys(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		hints       []HelpHint
		requiredKey string
	}{
		{"NotesListHints has Ctrl+H", NotesListHints, "Ctrl+H"},
		{"TodosListHints has Ctrl+H", TodosListHints, "Ctrl+H"},
		{"FocusIdleHints has Ctrl+H", FocusIdleHints, "Ctrl+H"},
		{"SearchInputHints has Ctrl+H", SearchInputHints, "Ctrl+H"},
		{"MindMapHints has Ctrl+H", MindMapHints, "Ctrl+H"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			found := false
			for _, hint := range tc.hints {
				if hint.Key == tc.requiredKey {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("%s: expected to find key %q", tc.name, tc.requiredKey)
			}
		})
	}
}
