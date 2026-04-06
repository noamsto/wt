package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestDetectStale_MergedBranch(t *testing.T) {
	dir := t.TempDir()
	repo := filepath.Join(dir, "repo")

	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = repo
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=Test",
			"GIT_AUTHOR_EMAIL=test@test.com",
			"GIT_COMMITTER_NAME=Test",
			"GIT_COMMITTER_EMAIL=test@test.com",
		)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("command %v failed: %v\n%s", args, err, out)
		}
	}

	os.MkdirAll(repo, 0o755)
	run("git", "init", "-b", "main")
	run("git", "commit", "--allow-empty", "-m", "init")

	run("git", "checkout", "-b", "feature-merged")
	run("git", "commit", "--allow-empty", "-m", "feature work")
	run("git", "checkout", "main")
	run("git", "merge", "feature-merged")

	wtPath := filepath.Join(repo, ".worktrees", "feature-merged")
	run("git", "worktree", "add", wtPath, "feature-merged")

	run("git", "checkout", "-b", "feature-active")
	run("git", "commit", "--allow-empty", "-m", "active work")
	run("git", "checkout", "main")
	wtPath2 := filepath.Join(repo, ".worktrees", "feature-active")
	run("git", "worktree", "add", wtPath2, "feature-active")

	worktrees := []Worktree{
		{Branch: "feature-merged", Path: wtPath},
		{Branch: "feature-active", Path: wtPath2},
	}

	DetectStale(repo, "main", worktrees)

	if !worktrees[0].IsStale() {
		t.Error("expected feature-merged to be stale")
	}
	if worktrees[0].StaleReason != "merged into main" {
		t.Errorf("expected stale reason 'merged into main', got %q", worktrees[0].StaleReason)
	}
	if worktrees[1].IsStale() {
		t.Error("expected feature-active to not be stale")
	}
}
