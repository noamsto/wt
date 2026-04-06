package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/noamsto/wt/internal/git"
	"github.com/noamsto/wt/internal/runtime"
	"github.com/noamsto/wt/internal/tmux"
	"github.com/noamsto/wt/internal/tui/prompt"
)

// Find fuzzy-finds a worktree and outputs its path.
func Find(repoRoot, query string, rt runtime.Runtime, tmuxClient *tmux.Client) error {
	defaultBranch, err := git.DefaultBranch(repoRoot)
	if err != nil {
		return err
	}

	worktrees, err := git.ListWorktrees(repoRoot, defaultBranch)
	if err != nil {
		return err
	}

	if len(worktrees) == 0 {
		return fmt.Errorf("no worktrees found (besides main)")
	}

	var paths []string
	for _, wt := range worktrees {
		paths = append(paths, wt.Path)
	}

	var result string
	if query != "" {
		// Filter paths matching query
		var matches []string
		for _, p := range paths {
			if strings.Contains(p, query) {
				matches = append(matches, p)
			}
		}
		if len(matches) == 0 {
			return fmt.Errorf("no worktree matching '%s'", query)
		}
		if len(matches) == 1 {
			result = matches[0]
		} else {
			result, err = prompt.Filter(matches, "Select worktree...", query)
			if err != nil {
				return err
			}
		}
	} else {
		result, err = prompt.Filter(paths, "Select worktree...", "")
		if err != nil {
			return err
		}
	}

	if result == "" {
		return fmt.Errorf("no worktree selected")
	}

	// Update tmux window metadata
	branch := filepath.Base(result)
	tmuxClient.UpdateWindowMetadata(result, branch)

	fmt.Println(result)
	return nil
}
