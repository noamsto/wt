package cmd

import (
	"fmt"

	"github.com/noamsto/wt/internal/runtime"
	"github.com/noamsto/wt/internal/tmux"
	"github.com/noamsto/wt/internal/zoxide"
)

// Smart handles the default wt <branch> command.
func Smart(repoRoot, branch string, rt runtime.Runtime, tmuxClient *tmux.Client, zoxideClient *zoxide.Client) error {
	return fmt.Errorf("not implemented: smart")
}
