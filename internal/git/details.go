package git

import (
	"os/exec"
	"strings"
)

// LoadDetails populates a worktree's detail fields (dirty files, unpushed commits, last commit).
func LoadDetails(wt *Worktree) {
	out, err := exec.Command("git", "-C", wt.Path, "status", "--porcelain").Output()
	if err == nil {
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		if len(lines) == 1 && lines[0] == "" {
			wt.DirtyFiles = 0
		} else {
			wt.DirtyFiles = len(lines)
		}
	}

	out, err = exec.Command("git", "-C", wt.Path, "log", "--oneline", "@{upstream}..HEAD").Output()
	if err == nil {
		text := strings.TrimSpace(string(out))
		if text != "" {
			wt.UnpushedLog = strings.Split(text, "\n")
		}
	}

	out, err = exec.Command("git", "-C", wt.Path, "log", "-1",
		"--format=%h %s (%cr)", "--date=relative").Output()
	if err == nil {
		wt.LastCommit = strings.TrimSpace(string(out))
	}

	wt.DetailsLoaded = true
}
