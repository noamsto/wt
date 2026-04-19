# Demo

Reproducible assets for the `wtc -i` tmux-popup demo.

## Files

| File | Purpose |
|------|---------|
| `setup.sh` | Builds a throwaway repo at `/tmp/wtc-demo` with four worktrees, one per stale-detection path (merged / remote-deleted / dirty / clean). |
| `demo.tape` | [VHS](https://github.com/charmbracelet/vhs) script that drives the popup flow and renders `wtc-demo.gif`. |
| `generate.sh` | Nix-shell wrapper that pulls `vhs` + deps and renders the gif. |

## Regenerate the gif

```bash
./demo/generate.sh
```

Runs `setup.sh`, starts a tmux session, launches `wtc -i` as a `display-popup`,
navigates, selects all stale worktrees, deletes them, quits, and prints
`git worktree list` so the result is visible.

## Inspect manually

```bash
bash demo/setup.sh
cd /tmp/wtc-demo/repo
tmux new-session -A -s wtc-demo
# inside tmux:
tmux display-popup -E -w 90% -h 90% wtc -i
```

## Make the popup a keybinding

Add to `~/.config/tmux/tmux.conf`:

```tmux
bind-key W display-popup -E -w 90% -h 90% -d '#{pane_current_path}' "wtc -i"
```

Or, for fish users, wrap the binary so bare `wtc -i` auto-pops up when inside tmux:

```fish
function wtc
  set -l wtc_bin (command -s wtc)
  if contains -- -i $argv; and set -q TMUX; and command -q tmux
    command tmux display-popup -E -w 90% -h 90% -d "$PWD" "$wtc_bin $argv"
  else
    $wtc_bin $argv
  end
end
```
