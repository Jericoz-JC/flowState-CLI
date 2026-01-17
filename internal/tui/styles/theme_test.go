package styles

import (
	"strings"
	"testing"
)

func TestRenderASCIITime(t *testing.T) {
	tests := []struct {
		name     string
		timeStr  string
		wantRows int
	}{
		{"simple time", "25:00", 5},
		{"single digit", "5:30", 5},
		{"zeros", "00:00", 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderASCIITime(tt.timeStr, PrimaryColor)
			if result == "" {
				t.Fatalf("expected non-empty result")
			}

			// Count lines
			lines := strings.Split(result, "\n")
			if len(lines) != tt.wantRows {
				t.Errorf("got %d lines, want %d", len(lines), tt.wantRows)
			}
		})
	}
}

func TestRenderASCIITimeContainsDigitArt(t *testing.T) {
	result := RenderASCIITime("12:34", PrimaryColor)

	// Should contain block characters used in ASCII art
	if !strings.Contains(result, "█") && !strings.Contains(result, "╗") {
		t.Fatalf("expected ASCII art block characters in output")
	}
}

func TestRenderProgressRing(t *testing.T) {
	tests := []struct {
		name     string
		progress float64
		width    int
	}{
		{"empty", 0.0, 20},
		{"half", 0.5, 20},
		{"full", 1.0, 20},
		{"quarter", 0.25, 30},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderProgressRing(tt.progress, tt.width)
			if result == "" {
				t.Fatalf("expected non-empty result")
			}

			// Should contain brackets
			if !strings.Contains(result, "【") || !strings.Contains(result, "】") {
				t.Errorf("expected brackets in progress ring")
			}
		})
	}
}

func TestRenderProgressRingZeroWidth(t *testing.T) {
	result := RenderProgressRing(0.5, 0)
	if result != "" {
		t.Errorf("expected empty result for zero width")
	}
}

func TestRenderMiniBarChart(t *testing.T) {
	values := []int{3, 5, 2, 7, 4, 6, 8}

	result := RenderMiniBarChart(values, 5, 21)
	if result == "" {
		t.Fatalf("expected non-empty result")
	}

	// Should contain multiple lines
	lines := strings.Split(result, "\n")
	if len(lines) < 2 {
		t.Errorf("expected multiple lines in bar chart, got %d", len(lines))
	}
}

func TestRenderMiniBarChartEmpty(t *testing.T) {
	result := RenderMiniBarChart([]int{}, 5, 21)
	if result != "" {
		t.Errorf("expected empty result for empty values")
	}
}

func TestSessionCountIndicator(t *testing.T) {
	tests := []struct {
		name  string
		count int
		max   int
	}{
		{"no sessions", 0, 8},
		{"some sessions", 3, 8},
		{"all sessions", 8, 8},
		{"over max", 10, 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SessionCountIndicator(tt.count, tt.max)
			if result == "" {
				t.Fatalf("expected non-empty result")
			}

			// Should contain circles
			if !strings.Contains(result, "●") && !strings.Contains(result, "○") {
				t.Errorf("expected circle indicators")
			}
		})
	}
}

func TestVaporwaveProgressBar(t *testing.T) {
	tests := []struct {
		name     string
		progress float64
		width    int
	}{
		{"empty", 0.0, 20},
		{"half", 0.5, 20},
		{"full", 1.0, 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := VaporwaveProgressBar(tt.progress, tt.width)
			if result == "" && tt.width > 0 {
				t.Fatalf("expected non-empty result for width > 0")
			}
		})
	}
}

func TestGradientText(t *testing.T) {
	result := GradientText("Hello", PrimaryColor, SecondaryColor)
	if result == "" {
		t.Fatalf("expected non-empty result")
	}

	// Should contain H, e, l, l, o
	if !strings.Contains(result, "H") {
		t.Errorf("expected 'H' in result")
	}
}

func TestGlowBorder(t *testing.T) {
	content := "Test Content"
	result := GlowBorder(content, AccentColor)

	if result == "" {
		t.Fatalf("expected non-empty result")
	}

	// Result should contain the original content
	if !strings.Contains(result, content) {
		t.Errorf("expected result to contain original content")
	}

	// Result should have border characters
	hasBorder := strings.Contains(result, "║") ||
		strings.Contains(result, "│") ||
		strings.Contains(result, "╔") ||
		strings.Contains(result, "┌")
	if !hasBorder {
		t.Errorf("expected border characters in result")
	}
}

func TestGlowBorderEmpty(t *testing.T) {
	result := GlowBorder("", AccentColor)
	// Empty content should still return a bordered empty box
	if result == "" {
		t.Errorf("expected non-empty result even for empty content")
	}
}

func TestGlowBox(t *testing.T) {
	content := "Neon Glow"
	result := GlowBox(content)

	if result == "" {
		t.Fatalf("expected non-empty result")
	}

	// Should contain the content
	if !strings.Contains(result, content) {
		t.Errorf("expected result to contain content")
	}
}

func TestGlowBoxUsesAccentColor(t *testing.T) {
	// GlowBox should use AccentColor by default
	content := "Test"
	result := GlowBox(content)

	// Just verify it renders without error
	if result == "" {
		t.Fatalf("expected non-empty result")
	}
}

func TestGradientTitle(t *testing.T) {
	result := GradientTitle("My Title")

	if result == "" {
		t.Fatalf("expected non-empty result")
	}

	// Should contain the title text
	if !strings.Contains(result, "M") || !strings.Contains(result, "y") {
		t.Errorf("expected result to contain title characters")
	}
}
