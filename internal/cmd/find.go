package cmd

import (
	"fmt"

	"github.com/noamsto/wt/internal/runtime"
	"github.com/noamsto/wt/internal/tmux"
)

// Find fuzzy-finds a worktree and outputs its path.
func Find(repoRoot, query string, rt runtime.Runtime, tmuxClient *tmux.Client) error {
	return fmt.Errorf("not implemented: find")
}
