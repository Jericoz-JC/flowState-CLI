// Package sqlite provides SQLite-based persistent storage for flowState-cli.
//
// Phase 1: Core Infrastructure
// - Uses modernc.org/sqlite (pure Go SQLite implementation)
// - Stores notes, todos, sessions, and links
// - Automatic schema migration on startup
// - Indexed fields for efficient querying
//
// Database Schema:
//   - notes: id, title, body, tags (JSON), created_at, updated_at
//   - todos: id, title, description, status, priority, due_date, note_id, created_at, updated_at
//   - sessions: id, start_time, end_time, duration, status, created_at
//   - links: id, source_type, source_id, target_type, target_id, link_type, created_at
//
// Phase 2: Notes & Todos
// - Note CRUD operations with tag handling
// - Todo CRUD operations with status/priority handling
// - Automatic timestamp management
//
// Usage:
//
//	store, err := sqlite.New(cfg)
//	if err != nil { ... }
//	note := &models.Note{Title: "Test", Body: "Content"}
//	store.CreateNote(note)
//	notes, _ := store.ListNotes()
package sqlite

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "modernc.org/sqlite"

	"flowState-cli/internal/config"
	"flowState-cli/internal/models"
)

// Store manages SQLite database operations for flowState.
//
// Phase 1: Core Infrastructure
//   - db: SQLite database connection
//   - New(): Creates store and runs migrations
//   - Close(): Properly closes database connection
//
// Phase 2: Notes & Todos
//   - CreateNote/UpdateNote/DeleteNote/GetNote/ListNotes
//   - CreateTodo/UpdateTodo/DeleteTodo/GetTodo/ListTodos
//   - CreateSession/GetSession/ListSessions/UpdateSession
//   - CreateLink/GetLinksForItem/DeleteLink
type Store struct {
	db *sql.DB
}

// New creates a new SQLite store and runs migrations.
//
// Phase 1: Creates ~/.config/flowState/flowState.db
//   - Initializes all required tables
//   - Creates indexes for performance
//   - Handles existing databases gracefully
func New(cfg *config.Config) (*Store, error) {
	db, err := sql.Open("sqlite", cfg.DbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	store := &Store{db: db}
	if err := store.migrate(); err != nil {
		return nil, fmt.Errorf("failed to migrate: %w", err)
	}

	return store, nil
}

// migrate creates all required tables and indexes.
//
// Phase 1: Core Infrastructure
//   - Creates notes, todos, sessions, links tables
//   - Creates indexes for tags, status, foreign keys
func (s *Store) migrate() error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS notes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			body TEXT,
			tags TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS todos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			description TEXT,
			status TEXT DEFAULT 'pending',
			priority INTEGER DEFAULT 0,
			due_date DATETIME,
			note_id INTEGER REFERENCES notes(id),
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS sessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			start_time DATETIME,
			end_time DATETIME,
			duration INTEGER,
			status TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS links (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			source_type TEXT NOT NULL,
			source_id INTEGER NOT NULL,
			target_type TEXT NOT NULL,
			target_id INTEGER NOT NULL,
			link_type TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(source_type, source_id, target_type, target_id)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_notes_tags ON notes(tags)`,
		`CREATE INDEX IF NOT EXISTS idx_todos_status ON todos(status)`,
		`CREATE INDEX IF NOT EXISTS idx_todos_note_id ON todos(note_id)`,
		`CREATE INDEX IF NOT EXISTS idx_links_source ON links(source_type, source_id)`,
		`CREATE INDEX IF NOT EXISTS idx_links_target ON links(target_type, target_id)`,
	}

	for _, m := range migrations {
		if _, err := s.db.Exec(m); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	return nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

// Note Operations (Phase 2: Notes)

// CreateNote inserts a new note into the database.
// Automatically sets CreatedAt and UpdatedAt timestamps.
// Extracts and stores tags from body if present.
func (s *Store) CreateNote(note *models.Note) error {
	tagsJSON, _ := json.Marshal(note.Tags)
	now := time.Now()
	note.CreatedAt = now
	note.UpdatedAt = now

	result, err := s.db.Exec(
		"INSERT INTO notes (title, body, tags, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		note.Title, note.Body, string(tagsJSON), note.CreatedAt, note.UpdatedAt,
	)
	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	note.ID = id
	return nil
}

// GetNote retrieves a note by ID. Returns nil if not found.
func (s *Store) GetNote(id int64) (*models.Note, error) {
	var note models.Note
	var tagsStr string

	err := s.db.QueryRow(
		"SELECT id, title, body, tags, created_at, updated_at FROM notes WHERE id = ?",
		id,
	).Scan(&note.ID, &note.Title, &note.Body, &tagsStr, &note.CreatedAt, &note.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	json.Unmarshal([]byte(tagsStr), &note.Tags)
	return &note, nil
}

// ListNotes returns all notes ordered by updated_at descending.
func (s *Store) ListNotes() ([]models.Note, error) {
	rows, err := s.db.Query(
		"SELECT id, title, body, tags, created_at, updated_at FROM notes ORDER BY updated_at DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var note models.Note
		var tagsStr string
		if err := rows.Scan(&note.ID, &note.Title, &note.Body, &tagsStr, &note.CreatedAt, &note.UpdatedAt); err != nil {
			return nil, err
		}
		json.Unmarshal([]byte(tagsStr), &note.Tags)
		notes = append(notes, note)
	}
	return notes, nil
}

// UpdateNote modifies an existing note. Updates UpdatedAt timestamp.
func (s *Store) UpdateNote(note *models.Note) error {
	tagsJSON, _ := json.Marshal(note.Tags)
	note.UpdatedAt = time.Now()

	_, err := s.db.Exec(
		"UPDATE notes SET title = ?, body = ?, tags = ?, updated_at = ? WHERE id = ?",
		note.Title, note.Body, string(tagsJSON), note.UpdatedAt, note.ID,
	)
	return err
}

// DeleteNote removes a note by ID.
func (s *Store) DeleteNote(id int64) error {
	_, err := s.db.Exec("DELETE FROM notes WHERE id = ?", id)
	return err
}

// Todo Operations (Phase 2: Todos)

// CreateTodo inserts a new todo into the database.
func (s *Store) CreateTodo(todo *models.Todo) error {
	now := time.Now()
	todo.CreatedAt = now
	todo.UpdatedAt = now

	var dueDate interface{}
	if todo.DueDate != nil {
		dueDate = *todo.DueDate
	}

	var noteID interface{}
	if todo.NoteID != nil {
		noteID = *todo.NoteID
	}

	result, err := s.db.Exec(
		"INSERT INTO todos (title, description, status, priority, due_date, note_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		todo.Title, todo.Description, todo.Status, todo.Priority, dueDate, noteID, todo.CreatedAt, todo.UpdatedAt,
	)
	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	todo.ID = id
	return nil
}

// GetTodo retrieves a todo by ID.
func (s *Store) GetTodo(id int64) (*models.Todo, error) {
	var todo models.Todo
	var dueDate, noteID interface{}

	err := s.db.QueryRow(
		"SELECT id, title, description, status, priority, due_date, note_id, created_at, updated_at FROM todos WHERE id = ?",
		id,
	).Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Status, &todo.Priority, &dueDate, &noteID, &todo.CreatedAt, &todo.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if dueDate != nil {
		t := dueDate.(time.Time)
		todo.DueDate = &t
	}
	if noteID != nil {
		nid := noteID.(int64)
		todo.NoteID = &nid
	}

	return &todo, nil
}

// ListTodos returns all todos ordered by created_at descending.
func (s *Store) ListTodos() ([]models.Todo, error) {
	rows, err := s.db.Query(
		"SELECT id, title, description, status, priority, due_date, note_id, created_at, updated_at FROM todos ORDER BY created_at DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var todo models.Todo
		var dueDate, noteID interface{}
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Description, &todo.Status, &todo.Priority, &dueDate, &noteID, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
			return nil, err
		}
		if dueDate != nil {
			t := dueDate.(time.Time)
			todo.DueDate = &t
		}
		if noteID != nil {
			nid := noteID.(int64)
			todo.NoteID = &nid
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

// UpdateTodo modifies an existing todo.
func (s *Store) UpdateTodo(todo *models.Todo) error {
	todo.UpdatedAt = time.Now()

	var dueDate interface{}
	if todo.DueDate != nil {
		dueDate = *todo.DueDate
	}

	var noteID interface{}
	if todo.NoteID != nil {
		noteID = *todo.NoteID
	}

	_, err := s.db.Exec(
		"UPDATE todos SET title = ?, description = ?, status = ?, priority = ?, due_date = ?, note_id = ?, updated_at = ? WHERE id = ?",
		todo.Title, todo.Description, todo.Status, todo.Priority, dueDate, noteID, todo.UpdatedAt, todo.ID,
	)
	return err
}

// DeleteTodo removes a todo by ID.
func (s *Store) DeleteTodo(id int64) error {
	_, err := s.db.Exec("DELETE FROM todos WHERE id = ?", id)
	return err
}

// Session Operations (Phase 4: Focus Sessions - upcoming)

// CreateSession inserts a new focus session.
func (s *Store) CreateSession(session *models.FocusSession) error {
	session.CreatedAt = time.Now()

	result, err := s.db.Exec(
		"INSERT INTO sessions (start_time, end_time, duration, status, created_at) VALUES (?, ?, ?, ?, ?)",
		session.StartTime, session.EndTime, session.Duration, session.Status, session.CreatedAt,
	)
	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	session.ID = id
	return nil
}

// GetSession retrieves a session by ID.
func (s *Store) GetSession(id int64) (*models.FocusSession, error) {
	var session models.FocusSession

	err := s.db.QueryRow(
		"SELECT id, start_time, end_time, duration, status, created_at FROM sessions WHERE id = ?",
		id,
	).Scan(&session.ID, &session.StartTime, &session.EndTime, &session.Duration, &session.Status, &session.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &session, nil
}

// ListSessions returns all sessions ordered by created_at descending.
func (s *Store) ListSessions() ([]models.FocusSession, error) {
	rows, err := s.db.Query(
		"SELECT id, start_time, end_time, duration, status, created_at FROM sessions ORDER BY created_at DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []models.FocusSession
	for rows.Next() {
		var session models.FocusSession
		if err := rows.Scan(&session.ID, &session.StartTime, &session.EndTime, &session.Duration, &session.Status, &session.CreatedAt); err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}
	return sessions, nil
}

// UpdateSession modifies an existing session.
func (s *Store) UpdateSession(session *models.FocusSession) error {
	_, err := s.db.Exec(
		"UPDATE sessions SET start_time = ?, end_time = ?, duration = ?, status = ? WHERE id = ?",
		session.StartTime, session.EndTime, session.Duration, session.Status, session.ID,
	)
	return err
}

// Link Operations (Phase 3: Linking System - upcoming)

// CreateLink creates a relationship between two items.
// Source and target can be "note" or "todo".
func (s *Store) CreateLink(link *models.Link) error {
	link.CreatedAt = time.Now()

	result, err := s.db.Exec(
		"INSERT OR IGNORE INTO links (source_type, source_id, target_type, target_id, link_type, created_at) VALUES (?, ?, ?, ?, ?, ?)",
		link.SourceType, link.SourceID, link.TargetType, link.TargetID, link.LinkType, link.CreatedAt,
	)
	if err != nil {
		return err
	}

	id, _ := result.LastInsertId()
	link.ID = id
	return nil
}

// GetLinksForItem returns all links associated with an item.
func (s *Store) GetLinksForItem(itemType string, itemID int64) ([]models.Link, error) {
	rows, err := s.db.Query(
		"SELECT id, source_type, source_id, target_type, target_id, link_type, created_at FROM links WHERE (source_type = ? AND source_id = ?) OR (target_type = ? AND target_id = ?)",
		itemType, itemID, itemType, itemID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []models.Link
	for rows.Next() {
		var link models.Link
		if err := rows.Scan(&link.ID, &link.SourceType, &link.SourceID, &link.TargetType, &link.TargetID, &link.LinkType, &link.CreatedAt); err != nil {
			return nil, err
		}
		links = append(links, link)
	}
	return links, nil
}

// DeleteLink removes a link by ID.
func (s *Store) DeleteLink(id int64) error {
	_, err := s.db.Exec("DELETE FROM links WHERE id = ?", id)
	return err
}
