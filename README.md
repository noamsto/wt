# wt

Git worktree manager. Single binary, no required dependencies beyond `git`.

Optional integrations are auto-detected and silently no-op when unavailable:

- **tmux** -- one window per worktree, one session per project
- **zoxide** -- frecency tracking for worktree directories
- **gh** -- squash-merge detection for stale branch cleanup

## Install

### Nix flake

```nix
# flake input
inputs.wt.url = "github:noamsto/wt";

# use the package directly
wt.packages.${system}.default

# or via home-manager module
imports = [ wt.homeManagerModules.default ];
programs.wt.enable = true;  # also installs fish completions
```

### Go

```bash
go install github.com/noamsto/wt/cmd/wt@latest
```

## Usage

```
wt <branch>           Smart switch/create (prompts before creating)
wt -y <branch>        Skip prompts
wt -q <branch>        Quiet mode (only output path)
wt -n <branch>        No tmux (skip window creation/switching)
wt -yqn <branch>      Combine flags (for scripts)
wt z [query]          Fuzzy find worktree, output path
wt main               Switch to root repository window
wt list               List all worktrees
wt remove <branch>    Remove worktree + kill window
wt clean              Remove stale worktrees (merged, squash-merged, deleted)
wt clean -i           Interactive TUI explorer
wt help               Show this help
```

### Smart mode

`wt <branch>` figures out what to do:

| Condition | Action |
|-----------|--------|
| Worktree exists | Switch to its tmux window |
| Branch exists | Prompt to create worktree |
| Branch not found | Prompt to create new branch |

### Scripting

```fish
# cd into a worktree (create if needed, no tmux, no prompts)
cd (wt -yqn my-feature)
```

### Interactive explorer

`wt clean -i` opens a TUI for inspecting and cleaning up worktrees:

```
j/k  navigate       space  select        a  select all stale
e    expand dirty    d      delete        D  force delete
/    search          q      quit
```

The preview pane shows branch details, dirty files, unpushed commits, and last commit info. Press `e` to expand the list of dirty files for a worktree.

### Stale detection

`wt clean` identifies worktrees whose branches are:

- Merged into the default branch
- Squash-merged via GitHub PR (requires `gh`)
- Deleted on the remote

## Worktree layout

Worktrees are created under `.worktrees/<branch-name>` relative to the repo root.

## License

MIT
