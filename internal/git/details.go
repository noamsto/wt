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
			wt.DirtyFileNames = nil
		} else {
			names := make([]string, 0, len(lines))
			for _, l := range lines {
				if len(l) > 3 {
					names = append(names, l[3:])
				}
			}
			wt.DirtyFileNames = names
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
