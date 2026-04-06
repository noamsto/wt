package gh

import "os/exec"

// Available returns true if the gh CLI is installed.
func Available() bool {
	_, err := exec.LookPath("gh")
	return err == nil
}
