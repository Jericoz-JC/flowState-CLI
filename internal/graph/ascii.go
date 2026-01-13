package graph

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// RenderSingleNode renders a single boxed node.
func RenderSingleNode(label string, boxWidth int, colorHex string) string {
	if boxWidth < 6 {
		boxWidth = 6
	}
	inner := boxWidth - 2
	lbl := truncate(label, inner)
	lbl = padRight(lbl, inner)

	top := "┌" + strings.Repeat("─", inner) + "┐"
	mid := "│" + lbl + "│"
	bot := "└" + strings.Repeat("─", inner) + "┘"

	if colorHex != "" {
		st := lipgloss.NewStyle().Foreground(lipgloss.Color(colorHex))
		top = st.Render(top)
		mid = st.Render(mid)
		bot = st.Render(bot)
	}

	return strings.Join([]string{top, mid, bot}, "\n")
}

// RenderGraphASCII renders a small graph into an ASCII canvas using provided positions.
// positions are interpreted as top-left coordinates for each node box.
func RenderGraphASCII(g Graph, labels map[string]string, positions map[string]Point, width, height int, nodeColorHex map[string]string) string {
	if width <= 0 || height <= 0 {
		return ""
	}

	grid := make([][]rune, height)
	for y := range grid {
		grid[y] = make([]rune, width)
		for x := range grid[y] {
			grid[y][x] = ' '
		}
	}

	// Draw nodes
	const boxW = 18
	for key := range g.Nodes {
		p, ok := positions[key]
		if !ok {
			continue
		}
		label := key
		if labels != nil && labels[key] != "" {
			label = labels[key]
		}

		color := ""
		if nodeColorHex != nil {
			color = nodeColorHex[key]
		}

		box := RenderSingleNode(label, boxW, color)
		lines := strings.Split(box, "\n")
		placeLines(grid, p.X, p.Y, lines)
	}

	// Draw simple horizontal edges (only when nodes share same Y for deterministic output).
	for from, m := range g.Adj {
		for to := range m {
			// ensure we only draw one direction
			if from >= to {
				continue
			}
			p1, ok1 := positions[from]
			p2, ok2 := positions[to]
			if !ok1 || !ok2 || p1.Y != p2.Y {
				continue
			}
			y := p1.Y + 1 // middle row of box
			x1 := p1.X + boxW
			x2 := p2.X - 1
			if x2 <= x1 || y < 0 || y >= height {
				continue
			}
			drawHorizontalArrow(grid, x1, x2, y)
		}
	}

	// Convert grid to lines, trimming trailing spaces but keeping width constraint.
	out := make([]string, 0, height)
	for _, row := range grid {
		s := string(row)
		s = strings.TrimRight(s, " ")
		if len(s) > width {
			s = s[:width]
		}
		out = append(out, s)
	}
	return strings.Join(out, "\n")
}

func drawHorizontalArrow(grid [][]rune, x1, x2, y int) {
	if y < 0 || y >= len(grid) {
		return
	}
	row := grid[y]
	if x1 < 0 {
		x1 = 0
	}
	if x2 >= len(row) {
		x2 = len(row) - 1
	}
	for x := x1; x < x2; x++ {
		row[x] = '─'
	}
	row[x2] = '▶'
}

func placeLines(grid [][]rune, x, y int, lines []string) {
	for dy, line := range lines {
		yy := y + dy
		if yy < 0 || yy >= len(grid) {
			continue
		}
		row := grid[yy]
		for dx, ch := range []rune(line) {
			xx := x + dx
			if xx < 0 || xx >= len(row) {
				continue
			}
			// Skip ANSI sequences (lipgloss output); for now, place only printable runes.
			// This keeps our ASCII canvas stable; color codes will remain in the string output
			// for tests that check for ANSI presence.
			row[xx] = ch
		}
	}
}

func truncate(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n])
}

func padRight(s string, n int) string {
	r := []rune(s)
	if len(r) >= n {
		return string(r[:n])
	}
	return s + strings.Repeat(" ", n-len(r))
}


