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

const helpText = `Git Worktree Manager

Usage:
  wt <branch>           Smart switch/create (prompts before creating)
  wt -y <branch>        Skip prompts
  wt -q <branch>        Quiet mode (only output path)
  wt -n <branch>        No tmux (skip window creation/switching)
  wt -yqn <branch>      Combine flags (for Claude/scripts)
  wt z [query]          Fuzzy find worktree, output path (cd "$(wt z)")
  wt main               Switch to root repository window
  wt list               List all worktrees
  wt remove <branch>    Remove worktree + kill window
  wt clean              Remove stale worktrees (merged, squash-merged, deleted)
  wt clean -i           Interactive explorer: inspect worktrees, force-remove
  wt help               Show this help

Model: Session = Project, Window = Worktree

Smart mode:
  Worktree exists     → switch to window (unless -n)
  Branch exists       → prompt to create worktree
  Branch not found    → prompt to create new branch

Worktree location: .worktrees/<branch-name>`

type flags struct {
	yes         bool
	quiet       bool
	noSwitch    bool
	interactive bool
}

func parseArgs(args []string) (flags, []string) {
	var f flags
	var rest []string

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") && !strings.HasPrefix(arg, "--") && arg != "-" {
			chars := arg[1:]
			allShort := true
			for _, c := range chars {
				switch c {
				case 'y':
					f.yes = true
				case 'q':
					f.quiet = true
				case 'n':
					f.noSwitch = true
				case 'i':
					f.interactive = true
				default:
					allShort = false
				}
			}
			if !allShort {
				rest = append(rest, arg)
			}
			continue
		}

		switch arg {
		case "--yes":
			f.yes = true
		case "--quiet":
			f.quiet = true
		case "--no-switch":
			f.noSwitch = true
		case "--interactive":
			f.interactive = true
		default:
			rest = append(rest, arg)
		}
	}

	return f, rest
}

func main() {
	f, args := parseArgs(os.Args[1:])

	sub := ""
	if len(args) > 0 {
		sub = args[0]
	}

	// Help doesn't need a git repo
	switch sub {
	case "help", "-h", "--help", "":
		fmt.Println(helpText)
		return
	}

	// Everything else requires a git repo
	repoRoot, err := git.RepoRoot()
	if err != nil {
		prompt.LogError("Not in a git repository")
		os.Exit(1)
	}

	rt := runtime.Detect()
	rt.NoSwitch = f.noSwitch
	rt.Quiet = f.quiet
	rt.Yes = f.yes

	tmuxClient := tmux.New(rt.TmuxActive())
	zoxideClient := zoxide.New(rt.HasZoxide)

	switch sub {
	case "list", "ls":
		if err := wcmd.List(repoRoot); err != nil {
			prompt.LogError("%v", err)
			os.Exit(1)
		}

	case "remove", "rm":
		branch := ""
		if len(args) > 1 {
			branch = args[1]
		}
		if err := wcmd.Remove(repoRoot, branch, rt, tmuxClient, zoxideClient); err != nil {
			prompt.LogError("%v", err)
			os.Exit(1)
		}

	case "clean", "prune":
		if err := wcmd.Clean(repoRoot, f.interactive, rt, tmuxClient, zoxideClient); err != nil {
			prompt.LogError("%v", err)
			os.Exit(1)
		}

	case "z":
		query := ""
		if len(args) > 1 {
			query = args[1]
		}
		if err := wcmd.Find(repoRoot, query, rt, tmuxClient); err != nil {
			prompt.LogError("%v", err)
			os.Exit(1)
		}

	case "main":
		if err := wcmd.MainSwitch(repoRoot, rt, tmuxClient); err != nil {
			prompt.LogError("%v", err)
			os.Exit(1)
		}

	default:
		// Smart mode: treat as branch name
		if err := wcmd.Smart(repoRoot, sub, rt, tmuxClient, zoxideClient); err != nil {
			prompt.LogError("%v", err)
			os.Exit(1)
		}
	}
}
