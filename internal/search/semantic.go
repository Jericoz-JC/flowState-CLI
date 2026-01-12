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
	embeddings "flowState-cli/internal/embeddings"
	"flowState-cli/internal/storage/qdrant"
)

type SemanticSearch struct {
	embedder    *embeddings.Embedder
	vectorStore *qdrant.VectorStore
}

func New(embedder *embeddings.Embedder, vectorStore *qdrant.VectorStore) *SemanticSearch {
	return &SemanticSearch{
		embedder:    embedder,
		vectorStore: vectorStore,
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
	queryEmbedding, err := s.embedder.EmbedSingle(query)
	if err != nil {
		return nil, err
	}

	results := s.vectorStore.Search(queryEmbedding, limit)
	searchResults := make([]SearchResult, len(results))
	for i, r := range results {
		searchResults[i] = SearchResult{
			NoteID:   r.NoteID,
			Score:    r.Score,
			NoteText: r.NoteText,
		}
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

	s.vectorStore.AddEmbedding(noteID, embeddings[0], text)
	return nil
}

// RemoveNote removes a note from the search index.
func (s *SemanticSearch) RemoveNote(noteID int64) error {
	s.vectorStore.DeleteEmbedding(noteID)
	return nil
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
