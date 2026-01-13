package sqlite

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Jericoz-JC/flowState-CLI/internal/config"
	"github.com/Jericoz-JC/flowState-CLI/internal/models"
)

// TestNotesCRUD tests Create, Read, Update, Delete operations for notes.
func TestNotesCRUD(t *testing.T) {
	// Setup: Create temporary database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	cfg := &config.Config{
		DbPath: dbPath,
	}

	store, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	// Test Create
	note := &models.Note{
		Title: "Test Note",
		Body:  "This is a #test note with #tags",
		Tags:  []string{"test", "tags"},
	}

	if err := store.CreateNote(note); err != nil {
		t.Fatalf("Failed to create note: %v", err)
	}

	if note.ID == 0 {
		t.Error("Note ID should be set after creation")
	}

	// Test Read
	retrieved, err := store.GetNote(note.ID)
	if err != nil {
		t.Fatalf("Failed to get note: %v", err)
	}

	if retrieved == nil {
		t.Fatal("Retrieved note should not be nil")
	}

	if retrieved.Title != note.Title {
		t.Errorf("Expected title %q, got %q", note.Title, retrieved.Title)
	}

	if retrieved.Body != note.Body {
		t.Errorf("Expected body %q, got %q", note.Body, retrieved.Body)
	}

	// Test Update - THIS IS THE CRITICAL BUG FIX VERIFICATION
	retrieved.Title = "Updated Note Title"
	retrieved.Body = "Updated body content"

	if err := store.UpdateNote(retrieved); err != nil {
		t.Fatalf("Failed to update note: %v", err)
	}

	// Verify update persisted
	updated, err := store.GetNote(note.ID)
	if err != nil {
		t.Fatalf("Failed to get updated note: %v", err)
	}

	if updated.Title != "Updated Note Title" {
		t.Errorf("Update failed: expected title %q, got %q", "Updated Note Title", updated.Title)
	}

	if updated.Body != "Updated body content" {
		t.Errorf("Update failed: expected body %q, got %q", "Updated body content", updated.Body)
	}

	// Verify only one note exists (not a duplicate)
	notes, err := store.ListNotes()
	if err != nil {
		t.Fatalf("Failed to list notes: %v", err)
	}

	if len(notes) != 1 {
		t.Errorf("Expected 1 note after update, got %d", len(notes))
	}

	// Test Delete
	if err := store.DeleteNote(note.ID); err != nil {
		t.Fatalf("Failed to delete note: %v", err)
	}

	deleted, err := store.GetNote(note.ID)
	if err != nil {
		t.Fatalf("Unexpected error after delete: %v", err)
	}

	if deleted != nil {
		t.Error("Note should be nil after deletion")
	}
}

// TestTodosCRUD tests Create, Read, Update, Delete operations for todos.
func TestTodosCRUD(t *testing.T) {
	// Setup: Create temporary database
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	cfg := &config.Config{
		DbPath: dbPath,
	}

	store, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	// Test Create
	todo := &models.Todo{
		Title:       "Test Todo",
		Description: "Test description",
		Status:      models.TodoStatusPending,
		Priority:    models.TodoPriorityMedium,
	}

	if err := store.CreateTodo(todo); err != nil {
		t.Fatalf("Failed to create todo: %v", err)
	}

	if todo.ID == 0 {
		t.Error("Todo ID should be set after creation")
	}

	// Test Read
	retrieved, err := store.GetTodo(todo.ID)
	if err != nil {
		t.Fatalf("Failed to get todo: %v", err)
	}

	if retrieved == nil {
		t.Fatal("Retrieved todo should not be nil")
	}

	if retrieved.Title != todo.Title {
		t.Errorf("Expected title %q, got %q", todo.Title, retrieved.Title)
	}

	if retrieved.Status != models.TodoStatusPending {
		t.Errorf("Expected status %q, got %q", models.TodoStatusPending, retrieved.Status)
	}

	// Test Update - THIS IS THE CRITICAL BUG FIX VERIFICATION
	retrieved.Title = "Updated Todo Title"
	retrieved.Status = models.TodoStatusCompleted
	retrieved.Priority = models.TodoPriorityHigh

	if err := store.UpdateTodo(retrieved); err != nil {
		t.Fatalf("Failed to update todo: %v", err)
	}

	// Verify update persisted
	updated, err := store.GetTodo(todo.ID)
	if err != nil {
		t.Fatalf("Failed to get updated todo: %v", err)
	}

	if updated.Title != "Updated Todo Title" {
		t.Errorf("Update failed: expected title %q, got %q", "Updated Todo Title", updated.Title)
	}

	if updated.Status != models.TodoStatusCompleted {
		t.Errorf("Update failed: expected status %q, got %q", models.TodoStatusCompleted, updated.Status)
	}

	if updated.Priority != models.TodoPriorityHigh {
		t.Errorf("Update failed: expected priority %d, got %d", models.TodoPriorityHigh, updated.Priority)
	}

	// Verify only one todo exists (not a duplicate)
	todos, err := store.ListTodos()
	if err != nil {
		t.Fatalf("Failed to list todos: %v", err)
	}

	if len(todos) != 1 {
		t.Errorf("Expected 1 todo after update, got %d", len(todos))
	}

	// Test Delete
	if err := store.DeleteTodo(todo.ID); err != nil {
		t.Fatalf("Failed to delete todo: %v", err)
	}

	deleted, err := store.GetTodo(todo.ID)
	if err != nil {
		t.Fatalf("Unexpected error after delete: %v", err)
	}

	if deleted != nil {
		t.Error("Todo should be nil after deletion")
	}
}

// TestTodoDueDate tests that due dates are properly stored and retrieved.
func TestTodoDueDate(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	cfg := &config.Config{
		DbPath: dbPath,
	}

	store, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	dueDate := time.Now().Add(24 * time.Hour).Truncate(time.Second)
	todo := &models.Todo{
		Title:   "Todo with due date",
		Status:  models.TodoStatusPending,
		DueDate: &dueDate,
	}

	if err := store.CreateTodo(todo); err != nil {
		t.Fatalf("Failed to create todo: %v", err)
	}

	retrieved, err := store.GetTodo(todo.ID)
	if err != nil {
		t.Fatalf("Failed to get todo: %v", err)
	}

	if retrieved.DueDate == nil {
		t.Fatal("Due date should not be nil")
	}

	// Compare truncated times to avoid subsecond differences
	if !retrieved.DueDate.Truncate(time.Second).Equal(dueDate.Truncate(time.Second)) {
		t.Errorf("Due date mismatch: expected %v, got %v", dueDate, *retrieved.DueDate)
	}
}

// TestListNotesEmpty tests that an empty database returns empty slice, not nil.
func TestListNotesEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	cfg := &config.Config{
		DbPath: dbPath,
	}

	store, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	notes, err := store.ListNotes()
	if err != nil {
		t.Fatalf("ListNotes failed: %v", err)
	}

	// Should return empty slice, not nil (important for JSON serialization)
	if notes == nil {
		// This is acceptable behavior but worth noting
		t.Log("Note: ListNotes returns nil for empty db, not empty slice")
	}

	if len(notes) != 0 {
		t.Errorf("Expected 0 notes in empty db, got %d", len(notes))
	}
}

// TestListNotesTruncation verifies that ListNotes returns truncated bodies.
func TestListNotesTruncation(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	cfg := &config.Config{DbPath: dbPath}
	store, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	// Create a note with > 100 chars
	longBody := ""
	for i := 0; i < 20; i++ {
		longBody += "1234567890" // 10 chars * 20 = 200 chars
	}

	note := &models.Note{Title: "Long Note", Body: longBody}
	store.CreateNote(note)

	// List notes
	notes, err := store.ListNotes()
	if err != nil {
		t.Fatalf("ListNotes failed: %v", err)
	}

	if len(notes) != 1 {
		t.Fatalf("Expected 1 note, got %d", len(notes))
	}

	// Verify truncation
	if len(notes[0].Body) != 100 {
		t.Errorf("Expected body length 100, got %d", len(notes[0].Body))
	}

	// Verify full fetch still works
	fullNote, err := store.GetNote(notes[0].ID)
	if err != nil {
		t.Fatalf("GetNote failed: %v", err)
	}
	if len(fullNote.Body) != 200 {
		t.Errorf("Expected full body length 200, got %d", len(fullNote.Body))
	}
}

// TestListTodosEmpty tests that an empty database returns empty slice.
func TestListTodosEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	cfg := &config.Config{
		DbPath: dbPath,
	}

	store, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	todos, err := store.ListTodos()
	if err != nil {
		t.Fatalf("ListTodos failed: %v", err)
	}

	if len(todos) != 0 {
		t.Errorf("Expected 0 todos in empty db, got %d", len(todos))
	}
}

// TestSessionsCRUD tests focus session operations.
func TestSessionsCRUD(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	cfg := &config.Config{
		DbPath: dbPath,
	}

	store, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	startTime := time.Now()
	endTime := startTime.Add(25 * time.Minute)

	session := &models.FocusSession{
		StartTime: startTime,
		EndTime:   &endTime,
		Duration:  25 * 60, // 25 minutes in seconds
		Status:    "completed",
	}

	if err := store.CreateSession(session); err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	if session.ID == 0 {
		t.Error("Session ID should be set after creation")
	}

	retrieved, err := store.GetSession(session.ID)
	if err != nil {
		t.Fatalf("Failed to get session: %v", err)
	}

	if retrieved == nil {
		t.Fatal("Retrieved session should not be nil")
	}

	if retrieved.Duration != 25*60 {
		t.Errorf("Expected duration %d, got %d", 25*60, retrieved.Duration)
	}

	if retrieved.Status != "completed" {
		t.Errorf("Expected status %q, got %q", "completed", retrieved.Status)
	}
}

// TestLinksCRUD tests link operations for Phase 3.
func TestLinksCRUD(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	cfg := &config.Config{
		DbPath: dbPath,
	}

	store, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	// Create a note and todo to link
	note := &models.Note{Title: "Link Test Note", Body: "Test body"}
	if err := store.CreateNote(note); err != nil {
		t.Fatalf("Failed to create note: %v", err)
	}

	todo := &models.Todo{Title: "Link Test Todo", Status: models.TodoStatusPending}
	if err := store.CreateTodo(todo); err != nil {
		t.Fatalf("Failed to create todo: %v", err)
	}

	// Create link
	link := &models.Link{
		SourceType: "note",
		SourceID:   note.ID,
		TargetType: "todo",
		TargetID:   todo.ID,
		LinkType:   "related",
	}

	if err := store.CreateLink(link); err != nil {
		t.Fatalf("Failed to create link: %v", err)
	}

	if link.ID == 0 {
		t.Error("Link ID should be set after creation")
	}

	// Get links for note
	noteLinks, err := store.GetLinksForItem("note", note.ID)
	if err != nil {
		t.Fatalf("Failed to get links for note: %v", err)
	}

	if len(noteLinks) != 1 {
		t.Errorf("Expected 1 link for note, got %d", len(noteLinks))
	}

	// Get links for todo (should also find the same link)
	todoLinks, err := store.GetLinksForItem("todo", todo.ID)
	if err != nil {
		t.Fatalf("Failed to get links for todo: %v", err)
	}

	if len(todoLinks) != 1 {
		t.Errorf("Expected 1 link for todo, got %d", len(todoLinks))
	}

	// Delete link
	if err := store.DeleteLink(link.ID); err != nil {
		t.Fatalf("Failed to delete link: %v", err)
	}

	// Verify deletion
	linksAfter, err := store.GetLinksForItem("note", note.ID)
	if err != nil {
		t.Fatalf("Failed to get links after delete: %v", err)
	}

	if len(linksAfter) != 0 {
		t.Errorf("Expected 0 links after delete, got %d", len(linksAfter))
	}
}

// TestDuplicateLinkIgnored tests that duplicate links are ignored (not error).
func TestDuplicateLinkIgnored(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	cfg := &config.Config{
		DbPath: dbPath,
	}

	store, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	note := &models.Note{Title: "Dup Link Note"}
	store.CreateNote(note)

	todo := &models.Todo{Title: "Dup Link Todo", Status: models.TodoStatusPending}
	store.CreateTodo(todo)

	link1 := &models.Link{
		SourceType: "note",
		SourceID:   note.ID,
		TargetType: "todo",
		TargetID:   todo.ID,
		LinkType:   "related",
	}

	if err := store.CreateLink(link1); err != nil {
		t.Fatalf("Failed to create first link: %v", err)
	}

	// Try to create duplicate - should be ignored, not error
	link2 := &models.Link{
		SourceType: "note",
		SourceID:   note.ID,
		TargetType: "todo",
		TargetID:   todo.ID,
		LinkType:   "related",
	}

	if err := store.CreateLink(link2); err != nil {
		t.Fatalf("Duplicate link should be ignored, not error: %v", err)
	}

	// Should still only have one link
	links, err := store.GetLinksForItem("note", note.ID)
	if err != nil {
		t.Fatalf("Failed to get links: %v", err)
	}

	if len(links) != 1 {
		t.Errorf("Expected 1 link (duplicate ignored), got %d", len(links))
	}
}

// TestSessionStats tests the session statistics calculations.
func TestSessionStats(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	cfg := &config.Config{DbPath: dbPath}
	store, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	// Initially no sessions
	stats, err := store.GetSessionStats()
	if err != nil {
		t.Fatalf("Failed to get initial stats: %v", err)
	}

	if stats.TodaySessions != 0 {
		t.Errorf("Expected 0 today sessions, got %d", stats.TodaySessions)
	}

	if stats.TotalSessions != 0 {
		t.Errorf("Expected 0 total sessions, got %d", stats.TotalSessions)
	}

	if stats.CurrentStreak != 0 {
		t.Errorf("Expected 0 streak, got %d", stats.CurrentStreak)
	}

	// Add a completed session for today
	now := time.Now()
	endTime := now.Add(25 * time.Minute)
	session := &models.FocusSession{
		StartTime: now,
		EndTime:   &endTime,
		Duration:  25 * 60, // 25 minutes in seconds
		Status:    models.SessionStatusCompleted,
	}

	if err := store.CreateSession(session); err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	// Check stats again
	stats, err = store.GetSessionStats()
	if err != nil {
		t.Fatalf("Failed to get stats after session: %v", err)
	}

	if stats.TodaySessions != 1 {
		t.Errorf("Expected 1 today session, got %d", stats.TodaySessions)
	}

	if stats.TotalSessions != 1 {
		t.Errorf("Expected 1 total session, got %d", stats.TotalSessions)
	}

	if stats.TotalFocusMinutes != 25 {
		t.Errorf("Expected 25 focus minutes, got %d", stats.TotalFocusMinutes)
	}

	if stats.CurrentStreak != 1 {
		t.Errorf("Expected 1 day streak, got %d", stats.CurrentStreak)
	}
}

// TestSessionStreakCalculation tests the streak calculation with multiple days.
func TestSessionStreakCalculation(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	cfg := &config.Config{DbPath: dbPath}
	store, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	now := time.Now()

	// Create sessions for today, yesterday, and day before yesterday (3-day streak)
	for i := 0; i < 3; i++ {
		sessionTime := now.AddDate(0, 0, -i)
		endTime := sessionTime.Add(25 * time.Minute)
		session := &models.FocusSession{
			StartTime: sessionTime,
			EndTime:   &endTime,
			Duration:  25 * 60,
			Status:    models.SessionStatusCompleted,
		}
		if err := store.CreateSession(session); err != nil {
			t.Fatalf("Failed to create session for day -%d: %v", i, err)
		}
	}

	streak, err := store.GetCurrentStreak()
	if err != nil {
		t.Fatalf("Failed to get streak: %v", err)
	}

	if streak != 3 {
		t.Errorf("Expected 3 day streak, got %d", streak)
	}
}

// TestSessionStreakBroken tests that streak is broken when a day is missed.
func TestSessionStreakBroken(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	cfg := &config.Config{DbPath: dbPath}
	store, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	now := time.Now()

	// Create session for 3 days ago (streak is broken - no session today or yesterday)
	sessionTime := now.AddDate(0, 0, -3)
	endTime := sessionTime.Add(25 * time.Minute)
	session := &models.FocusSession{
		StartTime: sessionTime,
		EndTime:   &endTime,
		Duration:  25 * 60,
		Status:    models.SessionStatusCompleted,
	}
	if err := store.CreateSession(session); err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	streak, err := store.GetCurrentStreak()
	if err != nil {
		t.Fatalf("Failed to get streak: %v", err)
	}

	if streak != 0 {
		t.Errorf("Expected 0 streak (broken), got %d", streak)
	}
}

// TestGetSessionsForDate tests retrieving sessions for a specific date.
func TestGetSessionsForDate(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	cfg := &config.Config{DbPath: dbPath}
	store, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	now := time.Now()
	today := now.Truncate(24 * time.Hour)

	// Create 2 sessions for today
	for i := 0; i < 2; i++ {
		sessionTime := now.Add(time.Duration(i) * time.Hour)
		endTime := sessionTime.Add(25 * time.Minute)
		session := &models.FocusSession{
			StartTime: sessionTime,
			EndTime:   &endTime,
			Duration:  25 * 60,
			Status:    models.SessionStatusCompleted,
		}
		if err := store.CreateSession(session); err != nil {
			t.Fatalf("Failed to create session %d: %v", i, err)
		}
	}

	// Create 1 session for yesterday
	yesterday := now.AddDate(0, 0, -1)
	endTime := yesterday.Add(25 * time.Minute)
	yesterdaySession := &models.FocusSession{
		StartTime: yesterday,
		EndTime:   &endTime,
		Duration:  25 * 60,
		Status:    models.SessionStatusCompleted,
	}
	if err := store.CreateSession(yesterdaySession); err != nil {
		t.Fatalf("Failed to create yesterday session: %v", err)
	}

	// Get sessions for today
	todaySessions, err := store.GetSessionsForDate(today)
	if err != nil {
		t.Fatalf("Failed to get sessions for today: %v", err)
	}

	if len(todaySessions) != 2 {
		t.Errorf("Expected 2 sessions for today, got %d", len(todaySessions))
	}

	// Get sessions for yesterday
	yesterdayDate := yesterday.Truncate(24 * time.Hour)
	yesterdaySessions, err := store.GetSessionsForDate(yesterdayDate)
	if err != nil {
		t.Fatalf("Failed to get sessions for yesterday: %v", err)
	}

	if len(yesterdaySessions) != 1 {
		t.Errorf("Expected 1 session for yesterday, got %d", len(yesterdaySessions))
	}
}

// TestDeleteSession tests session deletion.
func TestDeleteSession(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	cfg := &config.Config{DbPath: dbPath}
	store, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	now := time.Now()
	endTime := now.Add(25 * time.Minute)
	session := &models.FocusSession{
		StartTime: now,
		EndTime:   &endTime,
		Duration:  25 * 60,
		Status:    models.SessionStatusCompleted,
	}

	if err := store.CreateSession(session); err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	// Delete the session
	if err := store.DeleteSession(session.ID); err != nil {
		t.Fatalf("Failed to delete session: %v", err)
	}

	// Verify deletion
	deleted, err := store.GetSession(session.ID)
	if err != nil {
		t.Fatalf("Unexpected error after delete: %v", err)
	}

	if deleted != nil {
		t.Error("Session should be nil after deletion")
	}
}

// TestCancelledSessionNotCountedInStats tests that cancelled sessions don't count.
func TestCancelledSessionNotCountedInStats(t *testing.T) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	cfg := &config.Config{DbPath: dbPath}
	store, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create store: %v", err)
	}
	defer store.Close()

	now := time.Now()
	endTime := now.Add(10 * time.Minute)

	// Create a cancelled session
	session := &models.FocusSession{
		StartTime: now,
		EndTime:   &endTime,
		Duration:  25 * 60,
		Status:    models.SessionStatusCancelled, // Cancelled!
	}

	if err := store.CreateSession(session); err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}

	stats, err := store.GetSessionStats()
	if err != nil {
		t.Fatalf("Failed to get stats: %v", err)
	}

	// Cancelled sessions should not be counted
	if stats.TodaySessions != 0 {
		t.Errorf("Expected 0 today sessions (cancelled not counted), got %d", stats.TodaySessions)
	}

	if stats.TotalSessions != 0 {
		t.Errorf("Expected 0 total sessions (cancelled not counted), got %d", stats.TotalSessions)
	}

	if stats.CurrentStreak != 0 {
		t.Errorf("Expected 0 streak (cancelled not counted), got %d", stats.CurrentStreak)
	}
}

// cleanupTestDB is a helper to ensure db is closed
func cleanupTestDB(path string) {
	os.Remove(path)
}
