package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/noamsto/wt/internal/git"
	"github.com/noamsto/wt/internal/runtime"
	"github.com/noamsto/wt/internal/tmux"
	"github.com/noamsto/wt/internal/tui/prompt"
	"github.com/noamsto/wt/internal/zoxide"
)

// Smart handles the default wt <branch> command.
// It detects whether the worktree/branch exists and takes the appropriate action.
func Smart(repoRoot, branch string, rt runtime.Runtime, tmuxClient *tmux.Client, zoxideClient *zoxide.Client) error {
	worktreePath := filepath.Join(repoRoot, ".worktrees", branch)
	existingPath := git.FindWorktreeByBranch(repoRoot, branch)

	// Case 1: Worktree already exists
	if existingPath != "" {
		if _, err := os.Stat(existingPath); os.IsNotExist(err) {
			// Directory missing — offer to prune
			prompt.Log(rt.Quiet, "Worktree directory missing: %s", existingPath)
			prompt.Log(rt.Quiet, "Git thinks branch '%s' has a worktree, but directory doesn't exist.", branch)
			if prompt.Confirm("Run 'git worktree prune' to fix stale references?", rt.Yes) {
				if err := git.PruneWorktrees(repoRoot); err != nil {
					return fmt.Errorf("git worktree prune: %w", err)
				}
				prompt.Log(rt.Quiet, "✓ Pruned stale worktree references")
				// Fall through to create
			} else {
				return fmt.Errorf("cancelled")
			}
		} else {
			// Worktree exists and directory is present — switch to it
			return switchTo(existingPath, branch, repoRoot, rt, tmuxClient)
		}
	}

	// Case 2: Branch exists but no worktree
	isRemote := false
	branchExists := git.BranchExists(repoRoot, branch)
	if !branchExists {
		if git.RemoteBranchExists(repoRoot, branch) {
			branchExists = true
			isRemote = true
		}
	}

	if branchExists {
		sourceDesc := "local branch"
		if isRemote {
			sourceDesc = "remote branch origin/" + branch
		}
		prompt.Log(rt.Quiet, "Branch '%s' exists (%s) but has no worktree.", branch, sourceDesc)
		if !prompt.Confirm(fmt.Sprintf("Create worktree at .worktrees/%s?", branch), rt.Yes) {
			return fmt.Errorf("cancelled")
		}
		return createAndSwitch(repoRoot, branch, worktreePath, false, isRemote, rt, tmuxClient, zoxideClient)
	}

	// Case 3: Branch doesn't exist — create new
	prompt.Log(rt.Quiet, "Branch '%s' does not exist.", branch)
	if !prompt.Confirm("Create new branch + worktree?", rt.Yes) {
		return fmt.Errorf("cancelled")
	}
	return createAndSwitch(repoRoot, branch, worktreePath, true, false, rt, tmuxClient, zoxideClient)
}

func createAndSwitch(repoRoot, branch, worktreePath string, createBranch, trackRemote bool, rt runtime.Runtime, tmuxClient *tmux.Client, zoxideClient *zoxide.Client) error {
	prompt.Log(rt.Quiet, "Creating worktree...")
	if err := git.CreateWorktree(repoRoot, branch, worktreePath, createBranch, trackRemote); err != nil {
		return fmt.Errorf("failed to create worktree: %w", err)
	}
	prompt.Log(rt.Quiet, "✓ Worktree created: %s", worktreePath)
	zoxideClient.Add(worktreePath)
	return switchTo(worktreePath, branch, repoRoot, rt, tmuxClient)
}

func switchTo(worktreePath, branch, repoRoot string, rt runtime.Runtime, tmuxClient *tmux.Client) error {
	tmuxClient.SwitchToWorktree(repoRoot, branch, worktreePath)

	if rt.Quiet {
		fmt.Println(worktreePath)
	} else {
		fmt.Println()
		fmt.Println("Worktree:", worktreePath)
	}
	return nil
}
