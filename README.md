# flowState-cli

A unified terminal productivity system for notes, todos, and focus sessions with semantic search capabilities.

## Overview

flowState-cli keeps you in the flow by making knowledge capture, task management, and focus timing instant and interconnected. Built with a clean Go codebase and a modern TUI using Bubble Tea.

## Features

- **Notes**: Quick capture, tagging, and organization
- **Todos**: Task management with priorities and due dates
- **Focus Sessions**: Pomodoro-style timer with session tracking
- **Semantic Search**: Local ONNX-powered semantic search with embeddings
- **Linking System**: Connect notes and todos through relationships
- **Tagging**: Simple tagging and filtering across all content

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      flowState CLI                          │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │  Notes UI   │  │  Todos UI   │  │  Focus Session UI   │  │
│  └──────┬──────┘  └──────┬──────┘  └──────────┬──────────┘  │
│         │                │                     │             │
│         └────────────────┼─────────────────────┘             │
│                          │                                   │
│                    ┌─────▼─────┐                            │
│                    │  Bubble   │                            │
│                    │   Tea     │                            │
│                    └─────┬─────┘                            │
│                          │                                   │
│         ┌────────────────┼────────────────┐                 │
│         ▼                ▼                ▼                 │
│  ┌─────────────┐  ┌─────────────┐  ┌──────────────────┐    │
│  │   SQLite    │  │   Qdrant    │  │  Embedding Model │    │
│  │  (Metadata) │  │  (Vectors)  │  │  (all-MiniLM)    │    │
│  └─────────────┘  └─────────────┘  └──────────────────┘    │
└─────────────────────────────────────────────────────────────┘
```

## Tech Stack

- **TUI Framework**: Bubble Tea + Lip Gloss + Bubbles
- **Structured Storage**: SQLite (pure Go via modernc.org/sqlite)
- **Vector Storage**: Qdrant for semantic search
- **Embeddings**: ONNX Runtime Go with all-MiniLM-L6-v2 model
- **Local Only**: No cloud dependencies, privacy-first

## Quick Start

### Installation

```bash
git clone https://github.com/yourusername/flowState-cli
cd flowState-cli
go build -o flowState ./cmd/flowState/
./flowState
```

### First Run

On first run, the application will:
1. Initialize SQLite database
2. Download the embedding model (~90MB)
3. Start the TUI interface

### Keyboard Shortcuts

#### Global Navigation
| Key | Action |
|-----|--------|
| `Ctrl+N` | Notes screen |
| `Ctrl+T` | Todos screen |
| `Ctrl+F` | Focus session screen |
| `Ctrl+L` | Link selected item |
| `Ctrl+H` | Home screen / Help |
| `Esc` | Go back / Cancel |
| `q` | Quit application |

#### Notes & Todos Screens
| Key | Action |
|-----|--------|
| `c` | Create new item |
| `e` | Edit selected item |
| `d` | Delete selected item (with confirmation) |
| `Space` | Toggle todo completion |
| `Tab` | Switch between form fields |
| `Ctrl+S` | Save item |
| `Ctrl+L` | Create link to another item |
| `j/↓` | Move selection down |
| `k/↑` | Move selection up |

#### Linking Modal
| Key | Action |
|-----|--------|
| `c` | Create new link |
| `d` | Delete selected link |
| `↑/↓` | Navigate link types or targets |
| `Enter` | Select / Confirm |
| `Esc` | Close modal / Go back |

## Project Structure

```
flowState-cli/
├── cmd/
│   └── flowState/
│       └── main.go                    # Entry point
├── internal/
│   ├── config/
│   │   └── config.go                  # Configuration management
│   ├── models/
│   │   ├── note.go                    # Note data structure
│   │   ├── todo.go                    # Todo data structure
│   │   ├── session.go                 # Focus session structure
│   │   └── link.go                    # Linking relationships
│   ├── storage/
│   │   ├── sqlite/
│   │   │   └── store.go               # SQLite operations
│   │   └── qdrant/
│   │       └── vector_store.go        # Qdrant vector operations
│   ├── embeddings/
│   │   └── embedder.go                # ONNX embedding service
│   ├── search/
│   │   └── semantic.go                # Semantic search logic
│   ├── tui/
│   │   ├── app.go                     # Main TUI application
│   │   ├── screens/
│   │   │   ├── notes.go               # Notes screen
│   │   │   ├── todos.go               # Todos screen
│   │   │   ├── focus.go               # Focus session screen
│   │   │   └── search.go              # Search results screen
│   │   ├── components/
│   │   │   ├── list.go                # Reusable list component
│   │   │   ├── editor.go              # Text editor component
│   │   │   ├── tag_input.go           # Tag input component
│   │   │   └── timer.go               # Focus timer component
│   │   └── styles/
│   │       └── theme.go               # Lip Gloss styling
│   └── commands/
│       └── cmd.go                     # Bubble Tea command wrappers
├── embeddings/
│   └── models/
│       └── README.md                  # Model download instructions
├── migrations/
│   └── 001_initial.sql                # SQLite schema
├── go.mod
├── go.sum
└── README.md
```

## Implementation Phases

### Phase 1: Core Infrastructure
- Project initialization
- SQLite database layer
- Qdrant integration
- Basic TUI shell with navigation

### Phase 2: Notes & Todos
- Full CRUD for notes
- Inline tag editing
- Full CRUD for todos
- Priority and due dates

### Phase 3: Linking System
- Bidirectional note-todo links
- Link creation UI
- Visual link indicators
- Link queries

### Phase 4: Focus Sessions
- Pomodoro timer (25/5)
- Custom durations
- Session tracking
- Streak statistics

### Phase 5: Semantic Search
- ONNX embedding model
- Vector indexing
- Natural language queries
- Filtered results

### Phase 6: Context & Polish
- Context-aware suggestions
- Quick capture
- UI polish
- Testing

## Database Schema

```sql
-- Notes table
CREATE TABLE notes (
    id INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    body TEXT,
    tags TEXT, -- JSON array
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Todos table
CREATE TABLE todos (
    id INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,
    status TEXT DEFAULT 'pending', -- pending, in_progress, completed
    priority INTEGER DEFAULT 0,
    due_date DATETIME,
    note_id INTEGER REFERENCES notes(id),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Focus sessions table
CREATE TABLE sessions (
    id INTEGER PRIMARY KEY,
    start_time DATETIME,
    end_time DATETIME,
    duration INTEGER, -- in seconds
    status TEXT, -- running, completed, cancelled
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Links table
CREATE TABLE links (
    id INTEGER PRIMARY KEY,
    source_type TEXT NOT NULL, -- 'note' or 'todo'
    source_id INTEGER NOT NULL,
    target_type TEXT NOT NULL,
    target_id INTEGER NOT NULL,
    link_type TEXT, -- 'related', 'contains', 'references'
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(source_type, source_id, target_type, target_id)
);

-- Indexes
CREATE INDEX idx_notes_tags ON notes(tags);
CREATE INDEX idx_todos_status ON todos(status);
CREATE INDEX idx_todos_note_id ON todos(note_id);
CREATE INDEX idx_links_source ON links(source_type, source_id);
CREATE INDEX idx_links_target ON links(target_type, target_id);
```

## Semantic Search

The application uses `all-MiniLM-L6-v2` embedding model for semantic search:

- **Model size**: ~90MB
- **Dimensions**: 384
- **Storage**: Qdrant vector database
- **Features**: Natural language queries, tag filtering, incremental indexing

## Requirements

- **OS**: Windows 10+, Linux, macOS 11+
- **RAM**: 512MB free (embedding model ~100MB)
- **Storage**: 200MB for app + models
- **Go**: 1.21+

## Development

```bash
# Install dependencies
go mod download

# Run tests
go test ./...

# Build
go build -o flowState ./cmd/flowState/

# Run with hot reload (requires air)
air
```

## License

MIT License - see LICENSE file for details.

## Acknowledgments

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) for the excellent TUI framework
- [ONNX Runtime](https://onnxruntime.ai/) for cross-platform ML inference
- [Qdrant](https://qdrant.tech/) for vector search infrastructure
