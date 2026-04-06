package cmd

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/noamsto/wt/internal/git"
	"github.com/noamsto/wt/internal/runtime"
	"github.com/noamsto/wt/internal/tmux"
	"github.com/noamsto/wt/internal/zoxide"
)

// testRepo creates a temp git repo with an initial commit on main.
func testRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	repo := filepath.Join(dir, "repo")
	os.MkdirAll(repo, 0o755)
	gitRun(t, repo, "init", "-b", "main")
	gitRun(t, repo, "commit", "--allow-empty", "-m", "init")
	return repo
}

func gitRun(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=Test",
		"GIT_AUTHOR_EMAIL=test@test.com",
		"GIT_COMMITTER_NAME=Test",
		"GIT_COMMITTER_EMAIL=test@test.com",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, out)
	}
}

func testRuntime() runtime.Runtime {
	return runtime.Runtime{Yes: true, Quiet: true, NoSwitch: true}
}

func testTmux() *tmux.Client {
	return tmux.New(false)
}

func testZoxide() *zoxide.Client {
	return zoxide.New(false)
}

func TestSmart_CreateNewBranch(t *testing.T) {
	repo := testRepo(t)
	rt := testRuntime()

	err := Smart(repo, "feat-new", rt, testTmux(), testZoxide())
	if err != nil {
		t.Fatalf("Smart() error: %v", err)
	}

	// Verify worktree was created
	wtPath := filepath.Join(repo, ".worktrees", "feat-new")
	if _, err := os.Stat(wtPath); os.IsNotExist(err) {
		t.Error("expected worktree directory to exist")
	}

	// Verify branch was created
	if !git.BranchExists(repo, "feat-new") {
		t.Error("expected branch feat-new to exist")
	}
}

func TestSmart_SwitchExistingWorktree(t *testing.T) {
	repo := testRepo(t)
	rt := testRuntime()

	// Create worktree first
	wtPath := filepath.Join(repo, ".worktrees", "feat-existing")
	gitRun(t, repo, "worktree", "add", "-b", "feat-existing", wtPath)

	// Smart should switch to it (no error)
	err := Smart(repo, "feat-existing", rt, testTmux(), testZoxide())
	if err != nil {
		t.Fatalf("Smart() error: %v", err)
	}
}

func TestSmart_ExistingLocalBranch(t *testing.T) {
	repo := testRepo(t)
	rt := testRuntime()

	// Create a branch without a worktree
	gitRun(t, repo, "branch", "feat-local")

	err := Smart(repo, "feat-local", rt, testTmux(), testZoxide())
	if err != nil {
		t.Fatalf("Smart() error: %v", err)
	}

	// Verify worktree was created
	wtPath := filepath.Join(repo, ".worktrees", "feat-local")
	if _, err := os.Stat(wtPath); os.IsNotExist(err) {
		t.Error("expected worktree directory to exist")
	}
}

func TestRemove_Success(t *testing.T) {
	repo := testRepo(t)
	rt := testRuntime()

	// Create worktree
	wtPath := filepath.Join(repo, ".worktrees", "feat-remove")
	gitRun(t, repo, "worktree", "add", "-b", "feat-remove", wtPath)

	err := Remove(repo, "feat-remove", rt, testTmux(), testZoxide())
	if err != nil {
		t.Fatalf("Remove() error: %v", err)
	}

	// Verify worktree directory is gone
	if _, err := os.Stat(wtPath); !os.IsNotExist(err) {
		t.Error("expected worktree directory to be removed")
	}
}

func TestRemove_EmptyBranch(t *testing.T) {
	repo := testRepo(t)
	rt := testRuntime()

	err := Remove(repo, "", rt, testTmux(), testZoxide())
	if err == nil {
		t.Fatal("expected error for empty branch")
	}
	if !strings.Contains(err.Error(), "branch name required") {
		t.Errorf("expected 'branch name required' error, got: %v", err)
	}
}

func TestRemove_NonexistentBranch(t *testing.T) {
	repo := testRepo(t)
	rt := testRuntime()

	err := Remove(repo, "nonexistent", rt, testTmux(), testZoxide())
	if err == nil {
		t.Fatal("expected error for nonexistent branch")
	}
	if !strings.Contains(err.Error(), "no worktree found") {
		t.Errorf("expected 'no worktree found' error, got: %v", err)
	}
}

func TestList_ShowsWorktrees(t *testing.T) {
	repo := testRepo(t)

	// Capture stdout — List prints to stdout
	// We just verify it doesn't error
	err := List(repo)
	if err != nil {
		t.Fatalf("List() error: %v", err)
	}
}

func TestFind_NoWorktrees(t *testing.T) {
	repo := testRepo(t)
	rt := testRuntime()

	err := Find(repo, "", rt, testTmux())
	if err == nil {
		t.Fatal("expected error when no worktrees")
	}
	if !strings.Contains(err.Error(), "no worktrees found") {
		t.Errorf("expected 'no worktrees found' error, got: %v", err)
	}
}

func TestFind_SingleMatchByQuery(t *testing.T) {
	repo := testRepo(t)
	rt := testRuntime()

	// Create a worktree
	wtPath := filepath.Join(repo, ".worktrees", "feat-find")
	gitRun(t, repo, "worktree", "add", "-b", "feat-find", wtPath)

	// Find with exact query — should succeed (prints to stdout)
	err := Find(repo, "feat-find", rt, testTmux())
	if err != nil {
		t.Fatalf("Find() error: %v", err)
	}
}

func TestFind_NoMatchingQuery(t *testing.T) {
	repo := testRepo(t)
	rt := testRuntime()

	wtPath := filepath.Join(repo, ".worktrees", "feat-find")
	gitRun(t, repo, "worktree", "add", "-b", "feat-find", wtPath)

	err := Find(repo, "nonexistent", rt, testTmux())
	if err == nil {
		t.Fatal("expected error for non-matching query")
	}
	if !strings.Contains(err.Error(), "no worktree matching") {
		t.Errorf("expected 'no worktree matching' error, got: %v", err)
	}
}
