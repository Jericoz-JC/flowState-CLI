package graph

import (
	"fmt"
	"math"
	"sort"

	"flowState-cli/internal/models"
)

// NodeKey creates a stable node identifier for graph nodes.
func NodeKey(itemType string, itemID int64) string {
	return fmt.Sprintf("%s:%d", itemType, itemID)
}

type Node struct {
	Key      string
	ItemType string
	ItemID   int64
	Tags     []string
}

type Edge struct {
	From     string
	To       string
	LinkType models.LinkType
}

type Graph struct {
	Nodes map[string]*Node
	Adj   map[string]map[string]Edge // Adj[from][to] = edge
}

func New() Graph {
	return Graph{
		Nodes: map[string]*Node{},
		Adj:   map[string]map[string]Edge{},
	}
}

// BuildGraphFromLinks builds an undirected graph from link records.
// Each link creates edges in both directions for navigation purposes.
func BuildGraphFromLinks(links []models.Link, nodeTags map[string][]string) Graph {
	g := New()

	ensureNode := func(itemType string, itemID int64) string {
		key := NodeKey(itemType, itemID)
		if _, ok := g.Nodes[key]; !ok {
			g.Nodes[key] = &Node{
				Key:      key,
				ItemType: itemType,
				ItemID:   itemID,
				Tags:     nodeTags[key],
			}
		}
		if _, ok := g.Adj[key]; !ok {
			g.Adj[key] = map[string]Edge{}
		}
		return key
	}

	for _, l := range links {
		a := ensureNode(l.SourceType, l.SourceID)
		b := ensureNode(l.TargetType, l.TargetID)

		edgeAB := Edge{From: a, To: b, LinkType: l.LinkType}
		edgeBA := Edge{From: b, To: a, LinkType: l.LinkType}
		g.Adj[a][b] = edgeAB
		g.Adj[b][a] = edgeBA
	}

	return g
}

// ConnectedComponents returns connected components as slices of node keys.
func ConnectedComponents(g Graph) [][]string {
	visited := map[string]bool{}
	var comps [][]string

	var keys []string
	for k := range g.Nodes {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, start := range keys {
		if visited[start] {
			continue
		}
		// BFS
		queue := []string{start}
		visited[start] = true
		var comp []string
		for len(queue) > 0 {
			cur := queue[0]
			queue = queue[1:]
			comp = append(comp, cur)
			for nb := range g.Adj[cur] {
				if !visited[nb] {
					visited[nb] = true
					queue = append(queue, nb)
				}
			}
		}
		comps = append(comps, comp)
	}

	return comps
}

type Point struct {
	X int
	Y int
}

// CalculateNodePositions assigns deterministic positions within a bounding box.
// This is a simple circular layout to keep behavior stable for tests.
func CalculateNodePositions(g Graph, width, height int) map[string]Point {
	pos := make(map[string]Point, len(g.Nodes))
	if width <= 0 || height <= 0 || len(g.Nodes) == 0 {
		return pos
	}

	var keys []string
	for k := range g.Nodes {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	cx := float64(width-1) / 2.0
	cy := float64(height-1) / 2.0
	r := math.Min(float64(width), float64(height)) * 0.35
	if r < 1 {
		r = 1
	}

	n := float64(len(keys))
	for i, k := range keys {
		theta := 2 * math.Pi * (float64(i) / n)
		x := int(math.Round(cx + r*math.Cos(theta)))
		y := int(math.Round(cy + r*math.Sin(theta)))
		// Clamp
		if x < 0 {
			x = 0
		} else if x >= width {
			x = width - 1
		}
		if y < 0 {
			y = 0
		} else if y >= height {
			y = height - 1
		}
		pos[k] = Point{X: x, Y: y}
	}

	return pos
}

// TagColors assigns a stable color name (string token) for each unique tag.
func TagColors(tags []string) map[string]string {
	palette := []string{"cyan", "magenta", "green", "yellow", "blue", "red"}
	out := make(map[string]string, len(tags))
	seen := map[string]bool{}
	var uniq []string
	for _, t := range tags {
		if !seen[t] {
			seen[t] = true
			uniq = append(uniq, t)
		}
	}
	sort.Strings(uniq)
	for i, t := range uniq {
		out[t] = palette[i%len(palette)]
	}
	return out
}


