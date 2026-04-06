package explorer

import (
	"fmt"

	"github.com/noamsto/wt/internal/git"
	"github.com/noamsto/wt/internal/tmux"
	"github.com/noamsto/wt/internal/zoxide"
)

// Run launches the interactive TUI explorer for worktree management.
func Run(repoRoot string, worktrees []git.Worktree, tmuxClient *tmux.Client, zoxideClient *zoxide.Client) error {
	return fmt.Errorf("not implemented: explorer TUI")
}
