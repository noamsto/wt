package cmd

import (
	"fmt"
	"sort"

	"github.com/noamsto/wt/internal/git"
	"github.com/noamsto/wt/internal/runtime"
	"github.com/noamsto/wt/internal/tmux"
	"github.com/noamsto/wt/internal/tui/explorer"
	"github.com/noamsto/wt/internal/tui/prompt"
	"github.com/noamsto/wt/internal/zoxide"
)

// Clean removes stale worktrees. If interactive is true, launches the TUI explorer.
func Clean(repoRoot string, interactive bool, rt runtime.Runtime, tmuxClient *tmux.Client, zoxideClient *zoxide.Client) error {
	defaultBranch, err := git.DefaultBranch(repoRoot)
	if err != nil {
		return err
	}

	prompt.Log(rt.Quiet, "Fetching latest remote state...")
	_ = git.FetchPrune(repoRoot)

	worktrees, err := git.ListWorktrees(repoRoot, defaultBranch)
	if err != nil {
		return err
	}

	if len(worktrees) == 0 {
		prompt.Log(rt.Quiet, "No worktrees to clean (besides main).")
		return nil
	}

	git.DetectStale(repoRoot, defaultBranch, worktrees)
	if rt.HasGh {
		git.DetectStaleGh(repoRoot, worktrees)
	}

	if interactive {
		sort.Slice(worktrees, func(i, j int) bool {
			si := worktrees[i].IsStale()
			sj := worktrees[j].IsStale()
			if si != sj {
				return si
			}
			return worktrees[i].Branch < worktrees[j].Branch
		})
		return explorer.Run(repoRoot, worktrees, tmuxClient, zoxideClient)
	}

	// Non-interactive: find stale and prompt
	var stale []git.Worktree
	for _, wt := range worktrees {
		if wt.IsStale() {
			stale = append(stale, wt)
		}
	}

	if len(stale) == 0 {
		prompt.Log(rt.Quiet, "No stale worktrees found.")
		return nil
	}

	prompt.Log(rt.Quiet, "Found %d stale worktree(s):", len(stale))
	for _, wt := range stale {
		prompt.Log(rt.Quiet, "  • %s (%s)", wt.Branch, wt.StaleReason)
		prompt.Log(rt.Quiet, "    %s", wt.Path)
	}

	if !prompt.Confirm(fmt.Sprintf("Remove all %d stale worktrees?", len(stale)), rt.Yes) {
		return fmt.Errorf("cancelled")
	}

	var failed int
	for _, wt := range stale {
		prompt.Log(rt.Quiet, "Removing: %s", wt.Branch)
		tmuxClient.KillWindow(repoRoot, wt.Path)
		if err := git.RemoveWorktree(repoRoot, wt.Path, false); err != nil {
			prompt.Log(rt.Quiet, "  ❌ Failed: %v", err)
			failed++
		} else {
			zoxideClient.Remove(wt.Path)
			prompt.Log(rt.Quiet, "  ✓ Removed worktree")
		}
	}

	cleaned := len(stale) - failed
	if failed == 0 {
		prompt.Log(rt.Quiet, "✓ Cleaned %d worktree(s)", cleaned)
	} else {
		prompt.Log(rt.Quiet, "⚠ Cleaned %d worktree(s), %d failed", cleaned, failed)
	}
	return nil
}
