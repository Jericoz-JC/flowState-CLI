// Package main implements the entry point for flowState-cli.
// flowState-cli is a unified terminal productivity system for notes, todos,
// and focus sessions.
//
// Usage:
//
//	flowstate              # Run the application (npm global install)
//	./flowstate            # Run the downloaded binary (macOS/Linux)
//	.\flowstate.exe        # Run the downloaded binary (Windows)
//	flowstate --version    # Print version, commit, and the running binary path
//	flowstate --paths      # Print config/data/db/model/log locations
package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Jericoz-JC/flowState-CLI/internal/cli"
	"github.com/Jericoz-JC/flowState-CLI/internal/config"
	app "github.com/Jericoz-JC/flowState-CLI/internal/tui"
)

// Build metadata, injected at release time via -ldflags (see .goreleaser.yaml).
var (
	version = "dev"
	commit  = "none"
)

func main() {
	// Resolve configuration/paths FIRST so support flags and the log file can
	// point at the platform-native locations before any TUI init.
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	execPath, _ := os.Executable()

	// Handle non-TUI support flags used for install/run debugging.
	switch cli.ParseFlag(os.Args[1:]) {
	case cli.FlagVersion:
		fmt.Println(cli.VersionString(version, commit, execPath))
		return
	case cli.FlagPaths:
		fmt.Println(cli.PathsString(cli.Paths{
			ExecPath:  execPath,
			ConfigDir: cfg.DataDir,
			DataDir:   cfg.DataDir,
			DbPath:    cfg.DbPath,
			ModelPath: cfg.ModelPath,
			LogPath:   cfg.LogPath,
		}))
		return
	}

	// File logging to the platform-native log path (not the CWD).
	f, err := tea.LogToFile(cfg.LogPath, "debug")
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: could not open log file %s: %v\n", cfg.LogPath, err)
		os.Exit(1)
	}
	defer f.Close()

	// Global panic recovery: surface the absolute log path so users can find it.
	defer func() {
		if r := recover(); r != nil {
			log.Printf("CRITICAL PANIC: %v", r)
			fmt.Printf("\n\nEncountered a critical error: %v\nCheck the log for details:\n  %s\n", r, cfg.LogPath)
			os.Exit(1)
		}
	}()

	// Initialize TUI application with storage connections.
	application, err := app.New(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create app: %v\nLog: %s\n", err, cfg.LogPath)
		os.Exit(1)
	}
	defer application.Close()

	// Start the Bubble Tea event loop with alternate screen.
	p := tea.NewProgram(application, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running app: %v\nLog: %s\n", err, cfg.LogPath)
		os.Exit(1)
	}
}
