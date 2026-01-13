// Package search provides semantic search functionality for flowState-cli.
//
// Phase 5: Semantic Search (upcoming)
//   - Natural language query support
//   - Semantic similarity ranking
//   - Incremental indexing of notes
//
// Architecture:
//   - Embedder: Converts text to vectors
//   - VectorStore: Stores and searches vectors
//   - SemanticSearch: Orchestrates the search pipeline
//
// Usage:
//
// Phase 5: Enable semantic search in your notes
//
//	searcher.IndexNote(noteID, noteContent)
//	results, _ := searcher.Search("your query", 10)
//	for _, r := range results {
//	    fmt.Printf("Score: %.2f - %s\n", r.Score, r.NoteText)
//	}
package search

import (
	embeddings "github.com/Jericoz-JC/flowState-CLI/internal/embeddings"
	"github.com/Jericoz-JC/flowState-CLI/internal/storage/sqlite"
)

type SemanticSearch struct {
	embedder *embeddings.Embedder
	store    *sqlite.Store
}

func New(embedder *embeddings.Embedder, store *sqlite.Store) *SemanticSearch {
	return &SemanticSearch{
		embedder: embedder,
		store:    store,
	}
}

// Search performs semantic similarity search.
//
// Phase 5: Natural language query support
//   - Converts query to embedding vector
//   - Searches for similar note embeddings
//   - Returns ranked results by similarity score
//
// Example:
//
//	results, err := searcher.Search("project ideas", 10)
//	if err != nil { ... }
//	for _, r := range results {
//	    fmt.Printf("Match: %s (score: %.2f)\n", r.NoteText, r.Score)
//	}
func (s *SemanticSearch) Search(query string, limit int) ([]SearchResult, error) {
	if len(query) == 0 {
		return []SearchResult{}, nil
	}

	queryEmbedding, err := s.embedder.EmbedSingle(query)
	if err != nil {
		return nil, err
	}

	results, err := s.store.SearchNoteEmbeddings(queryEmbedding, limit)
	if err != nil {
		return nil, err
	}

	searchResults := make([]SearchResult, 0, len(results))
	for _, r := range results {
		note, err := s.store.GetNote(r.NoteID)
		if err != nil {
			return nil, err
		}
		if note == nil {
			// Note was deleted, vector should get cleaned up eventually.
			continue
		}

		preview := note.Title
		if note.Body != "" {
			preview = note.Title + "\n" + note.Body
		}
		if len(preview) > 300 {
			preview = preview[:300]
		}

		searchResults = append(searchResults, SearchResult{
			NoteID:   r.NoteID,
			Score:    r.Score,
			NoteText: preview,
		})
	}
	return searchResults, nil
}

// IndexNote adds a note to the search index.
//
// Phase 5: Incremental indexing
//   - Generates embedding for note content
//   - Stores in vector database
//   - Available for search immediately
func (s *SemanticSearch) IndexNote(noteID int64, text string) error {
	embeddings, err := s.embedder.Embed([]string{text})
	if err != nil {
		return err
	}

	return s.store.UpsertNoteEmbedding(noteID, embeddings[0])
}

// RemoveNote removes a note from the search index.
func (s *SemanticSearch) RemoveNote(noteID int64) error {
	return s.store.DeleteNoteEmbedding(noteID)
}

// IndexAllNotes bulk-indexes all notes currently in the database.
func (s *SemanticSearch) IndexAllNotes() error {
	notes, err := s.store.ListNotes()
	if err != nil {
		return err
	}

	for _, n := range notes {
		// ListNotes truncates body for performance; embed the full note.
		full, err := s.store.GetNote(n.ID)
		if err != nil {
			return err
		}
		if full == nil {
			continue
		}
		text := full.Title
		if full.Body != "" {
			text += "\n" + full.Body
		}
		if err := s.IndexNote(full.ID, text); err != nil {
			return err
		}
	}
	return nil
}

// SearchWithTagFilter searches and filters results to notes containing all requiredTags.
func (s *SemanticSearch) SearchWithTagFilter(query string, limit int, requiredTags []string) ([]SearchResult, error) {
	results, err := s.Search(query, limit)
	if err != nil {
		return nil, err
	}
	if len(requiredTags) == 0 {
		return results, nil
	}

	out := make([]SearchResult, 0, len(results))
	for _, r := range results {
		note, err := s.store.GetNote(r.NoteID)
		if err != nil {
			return nil, err
		}
		if note == nil {
			continue
		}
		if hasAllTags(note.Tags, requiredTags) {
			out = append(out, r)
		}
	}
	return out, nil
}

func hasAllTags(noteTags []string, required []string) bool {
	if len(required) == 0 {
		return true
	}
	set := make(map[string]struct{}, len(noteTags))
	for _, t := range noteTags {
		set[t] = struct{}{}
	}
	for _, req := range required {
		if _, ok := set[req]; !ok {
			return false
		}
	}
	return true
}

// SearchResult represents a single search result.
//
// Phase 5: Semantic Search Results
//   - NoteID: ID of the matching note
//   - Score: Cosine similarity (0.0 to 1.0)
//   - NoteText: Original note text for display
type SearchResult struct {
	NoteID   int64
	Score    float32
	NoteText string
}
