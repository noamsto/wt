# Regenerate via demo/generate.sh (builds wtc, pre-creates the tmux session).
# NOTE: VHS requires all Set directives before Env/Output.
Set Shell bash
Set FontSize 22
Set Width 1600
Set Height 1000
Set Padding 20
Set TypingSpeed 50ms
Set Theme "Catppuccin Mocha"

Output demo/wtc-demo.gif

# --- Attach to the pre-created session on its isolated socket (see generate.sh) ---
Hide
Type "tmux -L wtcdemo attach-session -t demo"
Enter
Sleep 800ms
Type "clear"
Enter
Sleep 400ms
Show

# --- Baseline: what's in the repo before cleanup ---
Type "git worktree list"
Enter
Sleep 3500ms
Type "clear"
Enter
Sleep 500ms

# --- Step 1: non-interactive wtc handles obvious stale (merged + remote-deleted) ---
Type "wtc -y"
Enter
Sleep 4s
Type "clear"
Enter
Sleep 500ms

# --- Step 2: wtc -i for the dirty leftover; auto-pops up inside tmux ---
Type "wtc -i"
Sleep 500ms
Enter
Sleep 3500ms

# --- Expand to show why the delete failed (dirty files) ---
Type "e"
Sleep 3500ms
Type "e"
Sleep 1s

# --- Force-delete with `D` (capital) ---
Type "D"
Sleep 2500ms
Type "y"
Sleep 3s

# --- Quit popup; verify the cleanup ---
Type "q"
Sleep 1500ms
Type "git worktree list"
Enter
Sleep 4s
