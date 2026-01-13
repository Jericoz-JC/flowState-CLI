package graph

import (
	"fmt"
	"testing"

	"flowState-cli/internal/models"
)

func BenchmarkGraphRender(b *testing.B) {
	// Build a moderately sized graph.
	var links []models.Link
	nodeTags := map[string][]string{}
	labels := map[string]string{}

	for i := 0; i < 100; i++ {
		key := NodeKey("note", int64(i))
		nodeTags[key] = []string{"bench"}
		labels[key] = fmt.Sprintf("Note %d", i)
	}
	for i := 0; i < 200; i++ {
		a := int64(i % 100)
		bb := int64((i*7 + 3) % 100)
		if a == bb {
			bb = (bb + 1) % 100
		}
		links = append(links, models.Link{
			SourceType: "note",
			SourceID:   a,
			TargetType: "note",
			TargetID:   bb,
			LinkType:   models.LinkTypeRelated,
		})
	}

	g := BuildGraphFromLinks(links, nodeTags)
	pos := CalculateNodePositions(g, 120, 40)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = RenderGraphASCII(g, labels, pos, 120, 40, nil)
	}
}


