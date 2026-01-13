package embedder

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"flowState-cli/internal/config"
)

func TestEmbedderModelDownload(t *testing.T) {
	t.Parallel()

	// Serve a tiny fake model file.
	const modelBytes = "fake-onnx-model"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/model.onnx":
			w.Header().Set("Content-Type", "application/octet-stream")
			_, _ = io.WriteString(w, modelBytes)
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)

	tmpDir := t.TempDir()
	cfg := &config.Config{ModelPath: tmpDir}

	e, err := NewWithHTTPClient(cfg, srv.Client())
	if err != nil {
		t.Fatalf("NewWithHTTPClient() err = %v", err)
	}

	modelPath := e.ModelFilePath()
	if _, err := os.Stat(modelPath); err == nil {
		t.Fatalf("expected model to be missing at start, but exists: %s", modelPath)
	}

	// First download creates the file.
	if err := e.EnsureModel(context.Background(), srv.URL+"/model.onnx"); err != nil {
		t.Fatalf("EnsureModel() err = %v", err)
	}

	got, err := os.ReadFile(modelPath)
	if err != nil {
		t.Fatalf("ReadFile(%s) err = %v", modelPath, err)
	}
	if string(got) != modelBytes {
		t.Fatalf("downloaded model mismatch: got=%q want=%q", string(got), modelBytes)
	}

	// Second call should be idempotent (no error, file still there).
	if err := e.EnsureModel(context.Background(), srv.URL+"/model.onnx"); err != nil {
		t.Fatalf("EnsureModel() second call err = %v", err)
	}
}

func TestEmbedSingleText(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()
	cfg := &config.Config{ModelPath: tmpDir}

	e, err := New(cfg)
	if err != nil {
		t.Fatalf("New() err = %v", err)
	}

	v1, err := e.EmbedSingle("hello world")
	if err != nil {
		t.Fatalf("EmbedSingle() err = %v", err)
	}
	if len(v1) != 384 {
		t.Fatalf("expected 384-dim embedding, got %d", len(v1))
	}

	// Deterministic for current provider: same text -> same vector.
	v2, err := e.EmbedSingle("hello world")
	if err != nil {
		t.Fatalf("EmbedSingle() 2nd err = %v", err)
	}
	for i := range v1 {
		if v1[i] != v2[i] {
			t.Fatalf("embedding not deterministic at idx=%d: %v != %v", i, v1[i], v2[i])
		}
	}

	// Ensure model path is under cfg.ModelPath (directory exists).
	wantDir := filepath.Join(tmpDir, "all-MiniLM-L6-v2")
	if e.modelPath != wantDir {
		t.Fatalf("modelPath mismatch: got=%s want=%s", e.modelPath, wantDir)
	}
}


