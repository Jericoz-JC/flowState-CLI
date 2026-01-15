package screens

import (
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Jericoz-JC/flowState-CLI/internal/graph"
	"github.com/Jericoz-JC/flowState-CLI/internal/storage/sqlite"
	"github.com/Jericoz-JC/flowState-CLI/internal/tui/components"
	"github.com/Jericoz-JC/flowState-CLI/internal/tui/styles"
)

type MindMapModel struct {
	store *sqlite.Store

	g         graph.Graph
	labels    map[string]string
	positions map[string]graph.Point
	nodeOrder []string

	selected int
	zoom     int
	showHelp bool // Help modal state

	header  components.Header
	helpBar components.HelpBar
	width   int
	height  int
}

func NewMindMapModel(store *sqlite.Store) MindMapModel {
	return MindMapModel{
		store:     store,
		g:         graph.New(),
		labels:    map[string]string{},
		positions: map[string]graph.Point{},
		nodeOrder: nil,
		selected:  0,
		zoom:      1,
		header:    components.NewHeader("🧠", "Mind Map"),
		helpBar:   components.NewHelpBar(components.MindMapHints),
	}
}

func (m *MindMapModel) Init() tea.Cmd { return nil }

func (m *MindMapModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.header.SetWidth(width - 4)
	m.helpBar.SetWidth(width - 4)
}

func (m *MindMapModel) LoadGraph() error {
	links, err := m.store.ListLinks()
	if err != nil {
		return err
	}

	notes, err := m.store.ListNotes()
	if err != nil {
		return err
	}

	nodeTags := make(map[string][]string)
	labels := make(map[string]string)
	for _, n := range notes {
		key := graph.NodeKey("note", n.ID)
		nodeTags[key] = n.Tags
		labels[key] = n.Title
	}

	m.g = graph.BuildGraphFromLinks(links, nodeTags)
	m.labels = labels

	m.nodeOrder = make([]string, 0, len(m.g.Nodes))
	for k := range m.g.Nodes {
		m.nodeOrder = append(m.nodeOrder, k)
	}
	sortStrings(m.nodeOrder)
	m.selected = 0

	canvasW, canvasH := m.canvasSize()
	m.positions = graph.CalculateNodePositions(m.g, canvasW, canvasH)
	return nil
}

func (m *MindMapModel) Update(msg tea.Msg) (MindMapModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle help modal
		if m.showHelp {
			// Any key closes help
			m.showHelp = false
			return *m, nil
		}

		switch msg.String() {
		case "?":
			m.showHelp = true
			return *m, nil
		case "+", "=":
			if m.zoom < 3 {
				m.zoom++
				canvasW, canvasH := m.canvasSize()
				m.positions = graph.CalculateNodePositions(m.g, canvasW, canvasH)
			}
			return *m, nil
		case "-":
			if m.zoom > 1 {
				m.zoom--
				canvasW, canvasH := m.canvasSize()
				m.positions = graph.CalculateNodePositions(m.g, canvasW, canvasH)
			}
			return *m, nil
		case "h":
			m.moveSelection(-1, 0)
			return *m, nil
		case "l":
			m.moveSelection(1, 0)
			return *m, nil
		case "k":
			m.moveSelection(0, -1)
			return *m, nil
		case "j":
			m.moveSelection(0, 1)
			return *m, nil
		case "enter":
			if len(m.nodeOrder) == 0 {
				return *m, nil
			}
			key := m.nodeOrder[m.selected]
			if strings.HasPrefix(key, "note:") {
				id, _ := strconv.ParseInt(strings.TrimPrefix(key, "note:"), 10, 64)
				return *m, func() tea.Msg { return OpenNoteMsg{NoteID: id} }
			}
			return *m, nil
		}
	}

	return *m, nil
}

func (m *MindMapModel) View() string {
	panel := lipgloss.NewStyle().Padding(1, 2).Width(m.width).Height(m.height)

	// Show help modal if active
	if m.showHelp {
		return panel.Render(m.helpView())
	}

	canvasW, canvasH := m.canvasSize()

	// Mark selected node with a distinct color (ARCHWAVE neon cyan).
	colors := map[string]string{}
	if len(m.nodeOrder) > 0 {
		colors[m.nodeOrder[m.selected]] = "#5ffbf1"
	}

	art := graph.RenderGraphASCII(m.g, m.labels, m.positions, canvasW, canvasH, colors)
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		m.header.View(),
		"",
		art,
		"",
		m.helpBar.View(),
	)
	return panel.Render(content)
}

func (m *MindMapModel) helpView() string {
	title := styles.TitleStyle.Render("🧠 MIND MAP - Help")

	helpText := `The Mind Map visualizes your notes and their connections as an interactive graph.

` + styles.SelectedItemStyle.Render("Navigation:") + `
• ` + styles.NeonStyle.Render("h/j/k/l") + ` or Arrow Keys: Pan the view
• ` + styles.NeonStyle.Render("+/-") + ` or Scroll: Zoom in/out
• ` + styles.NeonStyle.Render("Enter") + `: Open the selected note
• ` + styles.NeonStyle.Render("Esc") + `: Return to notes list

` + styles.SelectedItemStyle.Render("Visual Elements:") + `
• ` + styles.NeonStyle.Render("Nodes") + `: Each note appears as a node
• ` + styles.NeonStyle.Render("Edges") + `: Lines connect linked notes
• ` + styles.NeonStyle.Render("Colors") + `: Nodes are colored by tag
• ` + styles.NeonStyle.Render("Size") + `: Node size reflects connection count

` + styles.SelectedItemStyle.Render("Tips:") + `
• Notes with more links appear larger
• Cyan highlight shows current selection
• Clusters indicate related topics
• Use zoom to see more detail`

	help := styles.HelpStyle.Render("Press any key to close")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		helpText,
		"",
		help,
	)
}

func (m *MindMapModel) canvasSize() (int, int) {
	// Leave room for header + help bar; keep inside the panel padding (2 on each side).
	w := m.width - 4
	h := m.height - 10
	if w < 20 {
		w = 20
	}
	if h < 6 {
		h = 6
	}
	// Zoom increases the layout space, but we still clamp to terminal size.
	return w, h
}

func (m *MindMapModel) moveSelection(dx, dy int) {
	if len(m.nodeOrder) == 0 {
		return
	}
	curKey := m.nodeOrder[m.selected]
	curPos, ok := m.positions[curKey]
	if !ok {
		return
	}

	best := m.selected
	bestDist := int(^uint(0) >> 1) // max int

	for i, k := range m.nodeOrder {
		if i == m.selected {
			continue
		}
		p, ok := m.positions[k]
		if !ok {
			continue
		}
		ddx := p.X - curPos.X
		ddy := p.Y - curPos.Y

		// Direction filter.
		if dx < 0 && ddx >= 0 {
			continue
		}
		if dx > 0 && ddx <= 0 {
			continue
		}
		if dy < 0 && ddy >= 0 {
			continue
		}
		if dy > 0 && ddy <= 0 {
			continue
		}

		dist := abs(ddx) + abs(ddy)
		if dist < bestDist {
			bestDist = dist
			best = i
		}
	}

	m.selected = best
}

func sortStrings(s []string) {
	// Tiny in-place sort for small slices.
	for i := 0; i < len(s); i++ {
		best := i
		for j := i + 1; j < len(s); j++ {
			if s[j] < s[best] {
				best = j
			}
		}
		s[i], s[best] = s[best], s[i]
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
