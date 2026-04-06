package zoxide

import "os/exec"

// Client manages zoxide operations. All methods are no-ops if zoxide is unavailable.
type Client struct {
	active bool
}

// New creates a Client. If active is false, all operations are no-ops.
func New(active bool) *Client {
	return &Client{active: active}
}

// Add registers a path with zoxide.
func (c *Client) Add(path string) {
	if !c.active {
		return
	}
	_ = exec.Command("zoxide", "add", path).Run()
}

// Remove unregisters a path from zoxide.
func (c *Client) Remove(path string) {
	if !c.active {
		return
	}
	_ = exec.Command("zoxide", "remove", path).Run()
}
