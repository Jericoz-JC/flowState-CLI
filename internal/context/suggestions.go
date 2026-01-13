package context

import (
	"sort"
	"strings"

	"flowState-cli/internal/models"
	"flowState-cli/internal/search"
	"flowState-cli/internal/storage/sqlite"
)

// SuggestRelatedNotes returns notes semantically similar to the given query, excluding excludeNoteID.
func SuggestRelatedNotes(searcher *search.SemanticSearch, store *sqlite.Store, excludeNoteID int64, query string, limit int) ([]models.Note, error) {
	if limit <= 0 {
		return []models.Note{}, nil
	}
	results, err := searcher.Search(query, limit+3)
	if err != nil {
		return nil, err
	}

	out := make([]models.Note, 0, limit)
	for _, r := range results {
		if r.NoteID == excludeNoteID {
			continue
		}
		n, err := store.GetNote(r.NoteID)
		if err != nil {
			return nil, err
		}
		if n == nil {
			continue
		}
		out = append(out, *n)
		if len(out) >= limit {
			break
		}
	}
	return out, nil
}

// SuggestTagsFromContent extracts #hashtags from content as lowercase unique tags.
func SuggestTagsFromContent(content string) []string {
	tags := make(map[string]struct{})
	for _, word := range strings.Fields(content) {
		if !strings.HasPrefix(word, "#") {
			continue
		}
		tag := strings.TrimPrefix(word, "#")
		tag = strings.TrimSpace(tag)
		tag = strings.Trim(tag, ".,;:!?()[]{}\"'")
		tag = strings.ToLower(tag)
		if tag != "" {
			tags[tag] = struct{}{}
		}
	}
	out := make([]string, 0, len(tags))
	for t := range tags {
		out = append(out, t)
	}
	sort.Strings(out)
	return out
}

// SuggestLinksFromWikilinks extracts [[Note Name]] patterns from content.
func SuggestLinksFromWikilinks(content string) []string {
	links := []string{}
	inLink := false
	start := 0
	for i := 0; i < len(content); i++ {
		if i < len(content)-1 && content[i] == '[' && content[i+1] == '[' {
			if !inLink {
				inLink = true
				start = i + 2
				i++
			}
			continue
		}
		if i < len(content)-1 && content[i] == ']' && content[i+1] == ']' && inLink {
			txt := strings.TrimSpace(content[start:i])
			if txt != "" {
				links = append(links, txt)
			}
			inLink = false
			i++
		}
	}
	return links
}


