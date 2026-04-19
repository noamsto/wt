# Regenerate wtc-demo.gif
#   nix shell nixpkgs#vhs -- vhs demo/demo.tape
# Output lands next to the tape; move it wherever the README points to.

Output demo/wtc-demo.gif

Set Shell bash
Set FontSize 14
Set Width 1200
Set Height 720
Set Padding 20
Set TypingSpeed 50ms
Set Theme "Catppuccin Mocha"

# --- Setup (hidden) ---
Hide
Type "bash demo/setup.sh /tmp/wtc-demo"
Enter
Sleep 3s
Type "tmux -L wtcdemo kill-server 2>/dev/null; cd /tmp/wtc-demo/repo; clear"
Enter
Show

# --- Start tmux on an isolated socket with an empty config (no user tmux.conf, bash default shell) ---
Type "tmux -L wtcdemo -f /dev/null new-session -s demo bash"
Enter
Sleep 1500ms
Type "clear"
Enter
Sleep 500ms

# --- Launch wtc in a floating popup ---
Type "tmux display-popup -E -w 90% -h 90% wtc -i"
Sleep 500ms
Enter
Sleep 2500ms

# --- Explore: navigate and select all stale ---
Down
Sleep 800ms
Down
Sleep 800ms
Type "a"
Sleep 1500ms

# --- Delete selected ---
Type "d"
Sleep 1500ms
Type "y"
Sleep 2500ms

# --- Quit popup ---
Type "q"
Sleep 1s

# --- Show the result back in the shell ---
Type "git worktree list"
Enter
Sleep 2s
