// Package embedder provides text embedding generation for semantic search.
//
// Phase 1: Core Infrastructure
//   - Text-to-vector embedding conversion
//   - 384-dimensional vectors (MiniLM-L6 architecture)
//   - Ready for ONNX model integration
//
// Phase 5: Semantic Search (upcoming)
//   - Generates embeddings for semantic similarity
//   - Supports batch processing for efficiency
//   - Configurable model path and settings
//
// Current Implementation:
//   - Hash-based embedding for development
//   - Computes simple character-weighted vectors
//   - Ready for ONNX runtime integration
//
// Future Enhancement:
//   - Integrate onnxruntime_go
//   - Use all-MiniLM-L6-v2-onnx model
//   - ~90MB model size
//   - Download from HuggingFace automatically
package embedder

import (
	"fmt"
	"math"
	"os"
	"path/filepath"

	"flowState-cli/internal/config"
)

// Embedder generates vector embeddings from text.
//
// Phase 1: Core Infrastructure
//   - modelPath: Directory for model files
//   - New(): Creates embedder and ensures model directory exists
//   - Embed(): Generates embeddings for multiple texts
//   - EmbedSingle(): Generates embedding for one text
//
// Phase 5: Semantic Search
//   - Text is converted to 384-dimensional vector
//   - Similar texts have similar vectors
//   - Enables semantic similarity search
type Embedder struct {
	modelPath string
}

// New creates a new Embedder instance.
//
// Phase 1: Creates model directory at ~/.config/flowState/models/
//   - Ready for model file storage
//   - Future: Download ONNX model automatically
func New(cfg *config.Config) (*Embedder, error) {
	modelPath := filepath.Join(cfg.ModelPath, "all-MiniLM-L6-v2")

	if err := os.MkdirAll(modelPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create model directory: %w", err)
	}

	return &Embedder{
		modelPath: modelPath,
	}, nil
}

// Embed generates embeddings for multiple texts.
//
// Phase 1: Returns 2D slice of float32 vectors
//   - Each input text gets one 384-dimensional vector
//   - Vectors are normalized (unit length)
//   - Ready for cosine similarity comparison
func (e *Embedder) Embed(texts []string) ([][]float32, error) {
	return e.embedSimple(texts)
}

// EmbedSingle generates an embedding for one text.
func (e *Embedder) EmbedSingle(text string) ([]float32, error) {
	embeddings, err := e.Embed([]string{text})
	if err != nil {
		return nil, err
	}
	return embeddings[0], nil
}

// embedSimple creates simple hash-based embeddings.
//
// Phase 1: Development implementation
//   - Character-weighted hashing
//   - Normalized to unit length
//   - Produces deterministic results
//
// Future: Replace with ONNX inference
//   - Use all-MiniLM-L6-v2-onnx model
//   - True semantic embeddings
//   - Better similarity detection
func (e *Embedder) embedSimple(texts []string) ([][]float32, error) {
	embeddings := make([][]float32, len(texts))

	for i, text := range texts {
		embeddings[i] = e.simpleHashEmbedding(text)
	}

	return embeddings, nil
}

// simpleHashEmbedding creates a 384-dimensional vector from text.
func (e *Embedder) simpleHashEmbedding(text string) []float32 {
	dim := 384
	embedding := make([]float32, dim)

	for i, ch := range text {
		idx := i % dim
		embedding[idx] += float32(ch) * 0.01
	}

	norm := float32(0)
	for _, v := range embedding {
		norm += v * v
	}
	norm = float32(math.Sqrt(float64(norm)))
	if norm > 0 {
		for i := range embedding {
			embedding[i] /= norm
		}
	}

	return embedding
}

// GetModelInfo returns information about the embedding model.
//
// Phase 1: Core Infrastructure
//   - Model name and dimensions
//   - Download URL for ONNX model
//   - Expected model size (~90MB)
func (e *Embedder) GetModelInfo() ModelInfo {
	return ModelInfo{
		Name:        "all-MiniLM-L6-v2",
		Dimensions:  384,
		ModelPath:   e.modelPath,
		DownloadURL: "https://huggingface.co/sentence-transformers/all-MiniLM-L6-v2-onnx",
		ModelSize:   "90MB",
	}
}

// IsModelLoaded always returns true for current implementation.
func (e *Embedder) IsModelLoaded() bool {
	return true
}

// ModelInfo contains metadata about the embedding model.
type ModelInfo struct {
	Name        string
	Dimensions  int
	ModelPath   string
	DownloadURL string
	ModelSize   string
}

func (e *Embedder) Close() error {
	return nil
}
