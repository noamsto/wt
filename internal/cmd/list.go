package cmd

import (
	"fmt"

	"github.com/noamsto/wt/internal/git"
)

// List prints all worktrees.
func List(repoRoot string) error {
	output, err := git.ListWorktreesRaw(repoRoot)
	if err != nil {
		return err
	}
	fmt.Println(output)
	return nil
}
