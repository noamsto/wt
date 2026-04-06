package tmux

import (
	"os/exec"
	"path/filepath"
	"strings"
)

// Client manages tmux operations. All methods are no-ops if tmux is unavailable.
type Client struct {
	active bool
}

// New creates a Client. If active is false, all operations are no-ops.
func New(active bool) *Client {
	return &Client{active: active}
}

// FindWindowByWorktree returns the window index for a worktree path, or empty string.
func (c *Client) FindWindowByWorktree(session, worktreePath string) string {
	if !c.active {
		return ""
	}

	out, err := exec.Command("tmux", "list-windows", "-t", session,
		"-F", "#{window_index}\t#{@worktree}\t#{pane_current_path}").Output()
	if err != nil {
		return ""
	}

	for line := range strings.SplitSeq(strings.TrimSpace(string(out)), "\n") {
		parts := strings.SplitN(line, "\t", 3)
		if len(parts) != 3 {
			continue
		}
		if parts[1] == worktreePath || parts[2] == worktreePath {
			return parts[0]
		}
	}
	return ""
}

// SwitchToWorktree handles the full tmux session/window lifecycle for a worktree.
// Returns true if a tmux operation was performed.
func (c *Client) SwitchToWorktree(repoRoot, branch, worktreePath string) bool {
	if !c.active {
		return false
	}

	repoName := filepath.Base(repoRoot)
	inTmux := isInsideTmux()

	if inTmux {
		return c.switchInsideTmux(repoName, branch, worktreePath)
	}
	return c.switchOutsideTmux(repoName, branch, worktreePath)
}

func (c *Client) switchInsideTmux(repoName, branch, worktreePath string) bool {
	currentSession := tmuxCmd("display-message", "-p", "#{session_name}")

	inRepoSession := currentSession == repoName || strings.HasPrefix(currentSession, repoName+"/")

	if inRepoSession {
		targetSession := currentSession
		windowIdx := c.FindWindowByWorktree(targetSession, worktreePath)
		if windowIdx != "" {
			_ = exec.Command("tmux", "select-window", "-t", targetSession+":"+windowIdx).Run()
		} else {
			_ = exec.Command("tmux", "new-window", "-a", "-t", targetSession, "-c", worktreePath).Run()
			c.setWindowOptions(targetSession, branch, worktreePath)
		}
		return true
	}

	// Different session — switch to repo session
	if !hasSession(repoName) {
		_ = exec.Command("tmux", "new-session", "-d", "-s", repoName, "-c", worktreePath).Run()
		c.setWindowOptions(repoName, branch, worktreePath)
	} else if c.FindWindowByWorktree(repoName, worktreePath) == "" {
		_ = exec.Command("tmux", "new-window", "-a", "-t", repoName, "-c", worktreePath).Run()
		c.setWindowOptions(repoName, branch, worktreePath)
	}

	windowIdx := c.FindWindowByWorktree(repoName, worktreePath)
	if windowIdx != "" {
		_ = exec.Command("tmux", "switch-client", "-t", repoName+":"+windowIdx).Run()
	} else {
		_ = exec.Command("tmux", "switch-client", "-t", repoName).Run()
	}
	return true
}

func (c *Client) switchOutsideTmux(repoName, branch, worktreePath string) bool {
	if !hasSession(repoName) {
		_ = exec.Command("tmux", "new-session", "-d", "-s", repoName, "-c", worktreePath).Run()
		c.setWindowOptions(repoName, branch, worktreePath)
	} else if c.FindWindowByWorktree(repoName, worktreePath) == "" {
		_ = exec.Command("tmux", "new-window", "-a", "-t", repoName, "-c", worktreePath).Run()
		c.setWindowOptions(repoName, branch, worktreePath)
	}

	windowIdx := c.FindWindowByWorktree(repoName, worktreePath)
	if windowIdx != "" {
		_ = exec.Command("tmux", "attach-session", "-t", repoName+":"+windowIdx).Run()
	} else {
		_ = exec.Command("tmux", "attach-session", "-t", repoName).Run()
	}
	return true
}

func (c *Client) setWindowOptions(session, branch, worktreePath string) {
	_ = exec.Command("tmux", "set-option", "-t", session, "-w", "@worktree", worktreePath).Run()
	_ = exec.Command("tmux", "set-option", "-t", session, "-w", "@branch", branch).Run()
}

// KillWindow kills the tmux window associated with a worktree path.
func (c *Client) KillWindow(repoRoot, worktreePath string) {
	if !c.active {
		return
	}
	sessionName := filepath.Base(repoRoot)
	if !hasSession(sessionName) {
		return
	}
	windowIdx := c.FindWindowByWorktree(sessionName, worktreePath)
	if windowIdx != "" {
		_ = exec.Command("tmux", "kill-window", "-t", sessionName+":"+windowIdx).Run()
	}
}

// UpdateWindowMetadata sets @worktree and @branch on the current window.
func (c *Client) UpdateWindowMetadata(worktreePath, branch string) {
	if !c.active {
		return
	}
	_ = exec.Command("tmux", "set-option", "-w", "@worktree", worktreePath).Run()
	_ = exec.Command("tmux", "set-option", "-w", "@branch", branch).Run()
}

func isInsideTmux() bool {
	return strings.TrimSpace(tmuxCmd("display-message", "-p", "#{pid}")) != ""
}

func hasSession(name string) bool {
	return exec.Command("tmux", "has-session", "-t", name).Run() == nil
}

func tmuxCmd(args ...string) string {
	out, err := exec.Command("tmux", args...).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
