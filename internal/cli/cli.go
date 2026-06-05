// Package cli implements flowstate's non-TUI command-line surface:
// the --version and --paths support flags used for install/run debugging.
//
// The string builders are kept pure so they can be unit tested without
// touching the filesystem or starting the TUI.
package cli

import (
	"fmt"
	"runtime"
	"strings"
)

// Flag identifies a recognized top-level command-line flag.
type Flag int

const (
	FlagNone Flag = iota
	FlagVersion
	FlagPaths
)

// ParseFlag inspects the program arguments (os.Args[1:]) and returns the
// recognized support flag, if any. Only the first argument is considered.
func ParseFlag(args []string) Flag {
	if len(args) == 0 {
		return FlagNone
	}
	switch args[0] {
	case "--version", "-v", "version":
		return FlagVersion
	case "--paths", "paths":
		return FlagPaths
	default:
		return FlagNone
	}
}

// VersionString renders the output of `flowstate --version`. It includes the
// running executable's path so PATH conflicts (an old binary shadowing a new
// one) are immediately visible.
func VersionString(version, commit, execPath string) string {
	return fmt.Sprintf(
		"flowstate %s\n  commit:     %s\n  platform:   %s/%s\n  executable: %s",
		version, commit, runtime.GOOS, runtime.GOARCH, execPath,
	)
}

// Paths bundles the resolved filesystem locations for `flowstate --paths`.
type Paths struct {
	ExecPath  string
	ConfigDir string
	DataDir   string
	DbPath    string
	ModelPath string
	LogPath   string
}

// PathsString renders the output of `flowstate --paths`, listing every
// location the app reads from or writes to.
func PathsString(p Paths) string {
	var b strings.Builder
	b.WriteString("flowstate paths:\n")
	fmt.Fprintf(&b, "  executable: %s\n", p.ExecPath)
	fmt.Fprintf(&b, "  config dir: %s\n", p.ConfigDir)
	fmt.Fprintf(&b, "  data dir:   %s\n", p.DataDir)
	fmt.Fprintf(&b, "  database:   %s\n", p.DbPath)
	fmt.Fprintf(&b, "  models:     %s\n", p.ModelPath)
	fmt.Fprintf(&b, "  log file:   %s", p.LogPath)
	return b.String()
}
