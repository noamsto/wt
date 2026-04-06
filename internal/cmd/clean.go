package cmd

import (
	"fmt"

	"github.com/noamsto/wt/internal/runtime"
	"github.com/noamsto/wt/internal/tmux"
	"github.com/noamsto/wt/internal/zoxide"
)

// Clean removes stale worktrees.
func Clean(repoRoot string, interactive bool, rt runtime.Runtime, tmuxClient *tmux.Client, zoxideClient *zoxide.Client) error {
	return fmt.Errorf("not implemented: clean")
}
