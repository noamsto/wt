# Demo

Reproducible assets for the `wtc -i` tmux-popup demo.

## Files

| File | Purpose |
|------|---------|
| `setup.sh` | Builds a throwaway repo at `/tmp/wtc-demo` with four worktrees, one per stale-detection path (merged / remote-deleted / dirty / clean). |
| `demo.tape` | [VHS](https://github.com/charmbracelet/vhs) script that drives the flow and renders `wtc-demo.gif`. |
| `generate.sh` | Nix-shell wrapper that builds `wtc` from source, pulls `vhs` + deps, and renders the gif. |

## Regenerate the gif

```bash
./demo/generate.sh
```

Builds the current source, spins up an isolated tmux session, runs `wtc -i`,
selects all stale worktrees, deletes them, quits, and prints `git worktree list`.
The popup is triggered by `wtc` itself — no wrapper required.

## Inspect manually

```bash
bash demo/setup.sh
cd /tmp/wtc-demo/repo
tmux new-session -A -s wtc-demo
# inside tmux:
wtc -i              # opens as a floating popup automatically
WTC_NO_POPUP=1 wtc -i   # force inline (opt out of popup)
```
