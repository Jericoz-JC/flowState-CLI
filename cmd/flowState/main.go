// Package main implements the entry point for flowState-cli.
// flowState-cli is a unified terminal productivity system for notes, todos, and focus sessions.
//
// Phase 1: Core Infrastructure
// - Project initialization with Go modules
// - Configuration management via config.Load()
// - Bubble Tea TUI framework initialization
// - Proper cleanup with deferred Close()
//
// Phase 2: Notes & Todos
// - Full CRUD for notes (Create, Read, Update, Delete)
// - Full CRUD for todos (Create, Read, Update, Delete)
// - Auto-tagging from #hashtag syntax in note body
// - Status tracking for todos (pending, in_progress, completed)
// - Priority levels for todos (low, medium, high)
//
// Usage:
//
//	./flowState           # Run the application
//	./flowState.exe       # Windows executable
package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Jericoz-JC/flowState-CLI/internal/config"
	app "github.com/Jericoz-JC/flowState-CLI/internal/tui"
)

func main() {
	// Phase 4: Robustness - File logging
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	// Phase 4: Robustness - Global Panic Recovery
	defer func() {
		if r := recover(); r != nil {
			log.Printf("CRITICAL PANIC: %v", r)
			// Ensure terminal is reset even if p.Run() didn't exit cleanly
			fmt.Printf("\n\nEncoutered a critical error: %v\nCheck debug.log for details.\n", r)
			os.Exit(1)
		}
	}()

	// Phase 1: Load configuration from ~/.config/flowState/
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Phase 1: Initialize TUI application with storage connections
	app, err := app.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create app: %v", err)
	}
	defer app.Close()

	// Phase 1: Start Bubble Tea event loop with alternate screen
	p := tea.NewProgram(app, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running app: %v\n", err)
		os.Exit(1)
	}
}
