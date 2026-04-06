package cmd

import (
	"fmt"
	"os"

	"github.com/noamsto/wt/internal/git"
	"github.com/noamsto/wt/internal/runtime"
	"github.com/noamsto/wt/internal/tmux"
	"github.com/noamsto/wt/internal/tui/prompt"
	"github.com/noamsto/wt/internal/zoxide"
)

// Remove removes a worktree and its associated tmux window.
func Remove(repoRoot, branch string, rt runtime.Runtime, tmuxClient *tmux.Client, zoxideClient *zoxide.Client) error {
	if branch == "" {
		return fmt.Errorf("branch name required\nUsage: wt remove <branch-name>")
	}

	worktreePath := git.FindWorktreeByBranch(repoRoot, branch)
	if worktreePath == "" {
		fmt.Fprintf(os.Stderr, "No worktree found for branch '%s'\n\nAvailable worktrees:\n", branch)
		output, err := git.ListWorktreesRaw(repoRoot)
		if err == nil {
			fmt.Fprintln(os.Stderr, output)
		}
		return fmt.Errorf("no worktree found for branch '%s'", branch)
	}

	tmuxClient.KillWindow(repoRoot, worktreePath)

	prompt.Log(rt.Quiet, "Removing worktree: %s", worktreePath)
	if err := git.RemoveWorktree(repoRoot, worktreePath, false); err != nil {
		return fmt.Errorf("failed to remove worktree: %w", err)
	}

	zoxideClient.Remove(worktreePath)
	prompt.Log(rt.Quiet, "✓ Worktree removed")

	return nil
}
