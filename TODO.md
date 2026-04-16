# TODO

## `wt clean -i` Explorer Improvements

### Checkmark spacing too dense

The list line packs `✓●` right next to each other with no breathing room.
Current: `> ✓● branch-name [merged into main]`
Fix: add a space between the checkmark and stale indicator columns.

### Expand mode for dirty worktrees

When a worktree shows "3 dirty file(s)" in the preview, there's no way to see
which files are dirty. Add an expand/collapse toggle (e.g. `e` or `enter`) that
shows the actual dirty filenames inline.

This requires:
- Storing dirty file names (not just count) in `git.Worktree`
- Loading them lazily alongside `DirtyFiles` in `git.LoadDetails`
- Rendering the file list in the preview pane when expanded
