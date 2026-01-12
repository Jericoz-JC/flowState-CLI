// Package qdrant provides in-memory vector storage for semantic search.
//
// Phase 1: Core Infrastructure
//   - In-memory vector store with cosine similarity search
//   - Compatible interface with Qdrant for future migration
//   - Thread-safe with RWMutex for concurrent access
//
// Phase 5: Semantic Search (upcoming)
//   - Stores embeddings for notes
//   - Supports semantic similarity queries
//   - Cosine similarity scoring
//
// Current Implementation:
//   - Uses hash-based embeddings for development
//   - Ready for ONNX model integration
//   - In-memory storage (data lost on restart)
//
// Future Migration:
//   - Replace with actual Qdrant client
//   - Persistent storage across restarts
//   - GPU acceleration support
package qdrant

import (
	"math"
	"sort"
	"sync"

	"flowState-cli/internal/config"
)

// VectorStore manages vector embeddings for semantic search.
//
// Phase 1: Core Infrastructure
//   - mu: RWMutex for thread-safe operations
//   - vectors: Map of noteID to embedding vector
//   - noteTexts: Map of noteID to original text for display
//
// Phase 5: Semantic Search
//   - AddEmbedding: Store a new embedding
//   - Search: Find similar notes by cosine similarity
//   - DeleteEmbedding: Remove an embedding
type VectorStore struct {
	mu        sync.RWMutex
	vectors   map[int64][]float32
	noteTexts map[int64]string
}

// New creates a new in-memory vector store.
//
// Phase 1: Initializes empty storage maps
//   - Ready for embedding additions
//   - No external dependencies
func New(cfg *config.Config) (*VectorStore, error) {
	return &VectorStore{
		vectors:   make(map[int64][]float32),
		noteTexts: make(map[int64]string),
	}, nil
}

func (v *VectorStore) Close() error {
	return nil
}

// AddEmbedding stores an embedding for a note.
//
// Phase 1: Core Infrastructure
//   - Thread-safe insertion
//   - Stores both vector and original text
//   - Ready for Phase 5 semantic search
func (v *VectorStore) AddEmbedding(noteID int64, embedding []float32, text string) {
	v.mu.Lock()
	defer v.mu.Unlock()

	v.vectors[noteID] = embedding
	v.noteTexts[noteID] = text
}

// DeleteEmbedding removes an embedding by note ID.
func (v *VectorStore) DeleteEmbedding(noteID int64) {
	v.mu.Lock()
	defer v.mu.Unlock()

	delete(v.vectors, noteID)
	delete(v.noteTexts, noteID)
}

// Search finds the most similar notes to the query embedding.
//
// Phase 5: Semantic Search
//   - Computes cosine similarity with all stored vectors
//   - Returns results sorted by score (highest first)
//   - Limits results to specified count
//
// Cosine Similarity:
//   - Measures angle between vectors
//   - 1.0 = identical, 0.0 = orthogonal, -1.0 = opposite
func (v *VectorStore) Search(queryEmbedding []float32, limit int) []SearchResult {
	v.mu.RLock()
	defer v.mu.RUnlock()

	type scoredResult struct {
		noteID int64
		score  float32
	}

	results := make([]scoredResult, 0)

	for noteID, emb := range v.vectors {
		score := cosineSimilarity(queryEmbedding, emb)
		results = append(results, scoredResult{noteID: noteID, score: score})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})

	if len(results) > limit {
		results = results[:limit]
	}

	searchResults := make([]SearchResult, len(results))
	for i, r := range results {
		searchResults[i] = SearchResult{
			NoteID:   r.noteID,
			Score:    r.score,
			NoteText: v.noteTexts[r.noteID],
		}
	}

	return searchResults
}

// GetNoteIDs returns all stored note IDs.
func (v *VectorStore) GetNoteIDs() []int64 {
	v.mu.RLock()
	defer v.mu.RUnlock()

	ids := make([]int64, 0, len(v.vectors))
	for id := range v.vectors {
		ids = append(ids, id)
	}
	return ids
}

// GetEmbedding retrieves an embedding by note ID.
func (v *VectorStore) GetEmbedding(noteID int64) ([]float32, bool) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	emb, ok := v.vectors[noteID]
	return emb, ok
}

// GetStats returns the number of stored embeddings.
func (v *VectorStore) GetStats() (int, error) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	return len(v.vectors), nil
}

// IsAvailable always returns true for in-memory store.
func (v *VectorStore) IsAvailable() bool {
	return true
}

// cosineSimilarity computes the cosine similarity between two vectors.
//
// Formula: dot(a, b) / (||a|| * ||b||)
func cosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) {
		return 0
	}

	var dotProduct float32
	var normA, normB float32

	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / float32(math.Sqrt(float64(normA))*math.Sqrt(float64(normB)))
}

// SearchResult represents a single search result.
//
// Phase 5: Semantic Search
//   - NoteID: ID of the matching note
//   - Score: Cosine similarity score (0.0 to 1.0)
//   - NoteText: Original note text for display
type SearchResult struct {
	NoteID   int64
	Score    float32
	NoteText string
}
