package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/noamsto/wt/internal/runtime"
	"github.com/noamsto/wt/internal/tmux"
)

// MainSwitch switches to the repo root worktree.
func MainSwitch(repoRoot string, rt runtime.Runtime, tmuxClient *tmux.Client) error {
	if rt.Quiet {
		fmt.Println(repoRoot)
		return nil
	}

	tmuxClient.SwitchToWorktree(repoRoot, filepath.Base(repoRoot), repoRoot)

	if !rt.Quiet {
		fmt.Println()
		fmt.Println("Worktree:", repoRoot)
	}
	return nil
}
