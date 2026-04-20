package tmux

import (
	"os"
	"os/exec"
	"syscall"
)

// ReExecInPopup replaces the current process with `tmux display-popup -E <self> <args...>`.
// It does not return on success. Callers should fall through to normal execution
// when it returns a non-nil error (e.g. tmux not on PATH).
//
// Popup dimensions can be overridden via WTC_POPUP_WIDTH and WTC_POPUP_HEIGHT
// (accepted formats are any tmux display-popup size, e.g. "80%", "100", "40").
func ReExecInPopup(args ...string) error {
	tmuxPath, err := exec.LookPath("tmux")
	if err != nil {
		return err
	}
	self, err := os.Executable()
	if err != nil {
		return err
	}

	width := envOr("WTC_POPUP_WIDTH", "90%")
	height := envOr("WTC_POPUP_HEIGHT", "90%")

	popupArgs := []string{"tmux", "display-popup", "-E", "-w", width, "-h", height, self}
	popupArgs = append(popupArgs, args...)

	return syscall.Exec(tmuxPath, popupArgs, os.Environ())
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
