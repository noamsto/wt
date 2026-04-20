package main

import (
	"fmt"
	"os"
	"strings"

	wcmd "github.com/noamsto/wt/internal/cmd"
	"github.com/noamsto/wt/internal/git"
	"github.com/noamsto/wt/internal/runtime"
	"github.com/noamsto/wt/internal/tmux"
	"github.com/noamsto/wt/internal/tui/prompt"
)

const helpText = `wtc - Worktree Cleanup

Usage:
  wtc              Non-interactive clean (remove stale worktrees)
  wtc -i           Interactive TUI explorer
  wtc -y           Skip confirmation prompts
  wtc -q           Suppress non-essential output
  wtc -h           Show this help

Stale detection:
  • Merged into default branch
  • Remote branch deleted
  • GitHub PR squash-merged

Environment:
  WTC_NO_POPUP=1                 Disable auto-popup; keep the TUI inline
  WTC_POPUP_WIDTH / _HEIGHT      Override popup size (default 90%)`

type flags struct {
	yes         bool
	quiet       bool
	interactive bool
	inPopup     bool
}

func parseArgs(args []string) (flags, bool) {
	var f flags
	showHelp := false

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") && !strings.HasPrefix(arg, "--") && arg != "-" {
			for _, c := range arg[1:] {
				switch c {
				case 'y':
					f.yes = true
				case 'q':
					f.quiet = true
				case 'i':
					f.interactive = true
				case 'h':
					showHelp = true
				}
			}
			continue
		}

		switch arg {
		case "--yes":
			f.yes = true
		case "--quiet":
			f.quiet = true
		case "--interactive":
			f.interactive = true
		case "--in-popup":
			f.inPopup = true
		case "--help", "help":
			showHelp = true
		}
	}

	return f, showHelp
}

func main() {
	f, showHelp := parseArgs(os.Args[1:])

	if showHelp {
		fmt.Println(helpText)
		return
	}

	repoRoot, err := git.RepoRoot()
	if err != nil {
		prompt.LogError("Not in a git repository")
		os.Exit(1)
	}

	rt := runtime.Detect()
	rt.Quiet = f.quiet
	rt.Yes = f.yes

	// Auto-popup: when -i is run inside tmux, re-exec ourselves inside a
	// display-popup so the TUI floats instead of taking over the parent pane.
	// WTC_NO_POPUP=1 disables this. Does not return on success.
	if f.interactive && !f.inPopup && rt.InTmux && rt.HasTmux && os.Getenv("WTC_NO_POPUP") == "" {
		popupArgs := append([]string{"--in-popup"}, os.Args[1:]...)
		if err := tmux.ReExecInPopup(popupArgs...); err != nil {
			prompt.LogError("tmux popup unavailable, running inline: %v", err)
		}
	}

	tmuxClient := tmux.New(rt.TmuxActive())

	if err := wcmd.Clean(repoRoot, f.interactive, rt, tmuxClient); err != nil {
		prompt.LogError("%v", err)
		os.Exit(1)
	}
}
