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
	"github.com/noamsto/wt/internal/zoxide"
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
  • GitHub PR squash-merged`

type flags struct {
	yes         bool
	quiet       bool
	interactive bool
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

	tmuxClient := tmux.New(rt.TmuxActive())
	zoxideClient := zoxide.New(rt.HasZoxide)

	if err := wcmd.Clean(repoRoot, f.interactive, rt, tmuxClient, zoxideClient); err != nil {
		prompt.LogError("%v", err)
		os.Exit(1)
	}
}
