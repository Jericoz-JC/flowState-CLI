// Package models defines the core data structures for flowState-cli.
//
// Phase 1: Core Infrastructure
// - Note: Represents a note with title, body, tags, and timestamps
// - Todo: Represents a task with status, priority, and optional due date
// - FocusSession: Represents a focus timer session
// - Link: Represents relationships between notes and todos
//
// Phase 2: Notes & Todos
// - Note.Tags: Auto-extracted from #hashtag syntax in note body
// - Todo.Status: Tracks task progress (pending, in_progress, completed)
// - Todo.Priority: Task importance level (low, medium, high)
// - Todo.NoteID: Optional link to a related note
//
// All models include CreatedAt and UpdatedAt timestamps for auditing.
package models

import (
	"time"
)

// Note represents a note in the flowState system.
//
// Phase 1: Core Infrastructure
//   - ID: Unique identifier (auto-incremented by SQLite)
//   - Title: Short descriptive title
//   - Body: Full note content
//   - Tags: Slice of tags extracted from #hashtag syntax
//   - CreatedAt/UpdatedAt: Timestamps for auditing
//
// Phase 2: Notes
//   - Tags are automatically extracted when note is saved
//   - Supports filtering by tags in the UI
type Note struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TodoStatus represents the status of a todo item.
//
// Phase 2: Todos
//   - Pending: Task not started
//   - InProgress: Task currently being worked on
//   - Completed: Task finished
type TodoStatus string

const (
	TodoStatusPending    TodoStatus = "pending"
	TodoStatusInProgress TodoStatus = "in_progress"
	TodoStatusCompleted  TodoStatus = "completed"
)

// TodoPriority represents the priority level of a todo item.
//
// Phase 2: Todos
//   - Low: Minor tasks
//   - Medium: Normal tasks
//   - High: Urgent/important tasks
type TodoPriority int

const (
	TodoPriorityLow    TodoPriority = 0
	TodoPriorityMedium TodoPriority = 1
	TodoPriorityHigh   TodoPriority = 2
)

// Todo represents a task in the flowState system.
//
// Phase 1: Core Infrastructure
//   - ID: Unique identifier (auto-incremented by SQLite)
//   - Title: Short task title
//   - Description: Detailed task description
//   - Status: Current task status
//   - Priority: Task importance level
//   - DueDate: Optional due date
//   - NoteID: Optional link to related note
//   - CreatedAt/UpdatedAt: Timestamps for auditing
//
// Phase 2: Todos
//   - Press SPACE to toggle status between pending/completed
//   - Visual indicators: [ ] pending, [~] in progress, [x] completed
//   - Priority shown as 🔴 (high), 🟢 (low), nothing (medium)
type Todo struct {
	ID          int64        `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Status      TodoStatus   `json:"status"`
	Priority    TodoPriority `json:"priority"`
	DueDate     *time.Time   `json:"due_date,omitempty"`
	NoteID      *int64       `json:"note_id,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// SessionStatus represents the status of a focus session.
//
// Phase 4: Focus Sessions (upcoming)
type SessionStatus string

const (
	SessionStatusRunning   SessionStatus = "running"
	SessionStatusCompleted SessionStatus = "completed"
	SessionStatusCancelled SessionStatus = "cancelled"
)

// FocusSession represents a focus timer session.
//
// Phase 1: Core Infrastructure
//   - ID: Unique identifier
//   - StartTime: When the session began
//   - EndTime: When the session ended (nil if running)
//   - Duration: Planned duration in seconds
//   - Status: Session status
//   - CreatedAt: Timestamp
//
// Phase 4: Focus Sessions
//   - Pomodoro-style timer (25 min work, 5 min break)
//   - Session history tracking
//   - Daily/weekly statistics
type FocusSession struct {
	ID        int64         `json:"id"`
	StartTime time.Time     `json:"start_time"`
	EndTime   *time.Time    `json:"end_time,omitempty"`
	Duration  int           `json:"duration"`
	Status    SessionStatus `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
}

// LinkType represents the type of relationship between items.
//
// Phase 3: Linking System (upcoming)
type LinkType string

const (
	LinkTypeRelated    LinkType = "related"
	LinkTypeContains   LinkType = "contains"
	LinkTypeReferences LinkType = "references"
)

// Link represents a relationship between two items.
//
// Phase 1: Core Infrastructure
//   - ID: Unique identifier
//   - SourceType/SourceID: The source item (e.g., "note", 5)
//   - TargetType/TargetID: The target item
//   - LinkType: Type of relationship
//   - CreatedAt: Timestamp
//
// Phase 3: Linking System
//   - Bidirectional links between notes and todos
//   - Visual indicators showing linked items
//   - Press Ctrl+L to create links
type Link struct {
	ID         int64     `json:"id"`
	SourceType string    `json:"source_type"`
	SourceID   int64     `json:"source_id"`
	TargetType string    `json:"target_type"`
	TargetID   int64     `json:"target_id"`
	LinkType   LinkType  `json:"link_type"`
	CreatedAt  time.Time `json:"created_at"`
}

// SearchableItem is an interface for items that can be indexed for search.
//
// Phase 5: Semantic Search (upcoming)
type SearchableItem interface {
	GetID() int64
	GetContent() string
	GetType() string
}
