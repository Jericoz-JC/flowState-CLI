package cli

import (
	"runtime"
	"strings"
	"testing"
)

func TestVersionStringIncludesVersionCommitAndPath(t *testing.T) {
	out := VersionString("0.1.14", "abc1234", "/usr/local/bin/flowstate")

	for _, want := range []string{"0.1.14", "abc1234", "/usr/local/bin/flowstate", runtime.GOOS, runtime.GOARCH} {
		if !strings.Contains(out, want) {
			t.Errorf("VersionString output missing %q\ngot:\n%s", want, out)
		}
	}
}

func TestPathsStringListsAllResolvedLocations(t *testing.T) {
	p := Paths{
		ExecPath:  "/usr/local/bin/flowstate",
		ConfigDir: "/home/u/.config/flowState",
		DataDir:   "/home/u/.config/flowState",
		DbPath:    "/home/u/.config/flowState/flowState.db",
		ModelPath: "/home/u/.config/flowState/models",
		LogPath:   "/home/u/.config/flowState/logs/debug.log",
	}
	out := PathsString(p)

	for _, want := range []string{p.ExecPath, p.DataDir, p.DbPath, p.ModelPath, p.LogPath} {
		if !strings.Contains(out, want) {
			t.Errorf("PathsString output missing %q\ngot:\n%s", want, out)
		}
	}
}

func TestParseFlagsRecognizesVersion(t *testing.T) {
	for _, arg := range []string{"--version", "-v", "version"} {
		if got := ParseFlag([]string{arg}); got != FlagVersion {
			t.Errorf("ParseFlag([%q]) = %v, want FlagVersion", arg, got)
		}
	}
}

func TestParseFlagsRecognizesPaths(t *testing.T) {
	for _, arg := range []string{"--paths", "paths"} {
		if got := ParseFlag([]string{arg}); got != FlagPaths {
			t.Errorf("ParseFlag([%q]) = %v, want FlagPaths", arg, got)
		}
	}
}

func TestParseFlagsDefaultsToNone(t *testing.T) {
	if got := ParseFlag(nil); got != FlagNone {
		t.Errorf("ParseFlag(nil) = %v, want FlagNone", got)
	}
	if got := ParseFlag([]string{"--bogus"}); got != FlagNone {
		t.Errorf("ParseFlag([--bogus]) = %v, want FlagNone", got)
	}
}
