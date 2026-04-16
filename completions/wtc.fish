# Completions for wt (worktree manager)

# Flags
complete -c wt -f -s y -l yes -d 'Skip confirmation prompts'
complete -c wt -f -s q -l quiet -d 'Quiet mode (only output path)'
complete -c wt -f -s n -l no-switch -d 'Skip tmux window operations'
complete -c wt -f -s i -l interactive -d 'Interactive TUI mode (for clean)'

# Subcommands
complete -c wt -f -n '__fish_use_subcommand' -a 'list' -d 'List worktrees'
complete -c wt -f -n '__fish_use_subcommand' -a 'ls' -d 'List worktrees (alias)'
complete -c wt -f -n '__fish_use_subcommand' -a 'remove' -d 'Remove worktree + window'
complete -c wt -f -n '__fish_use_subcommand' -a 'rm' -d 'Remove worktree + window (alias)'
complete -c wt -f -n '__fish_use_subcommand' -a 'clean' -d 'Remove merged worktrees'
complete -c wt -f -n '__fish_use_subcommand' -a 'prune' -d 'Remove merged worktrees (alias)'
complete -c wt -f -n '__fish_use_subcommand' -a 'z' -d 'Fuzzy find worktree'
complete -c wt -f -n '__fish_use_subcommand' -a 'main' -d 'Switch to main worktree'
complete -c wt -f -n '__fish_use_subcommand' -a 'help' -d 'Show help'

# Complete existing worktree branches
function __wt_list_worktree_branches
    if git rev-parse --git-dir >/dev/null 2>&1
        set -l repo_root (git rev-parse --show-toplevel)
        git -C "$repo_root" worktree list 2>/dev/null | while read -l line
            if string match -rq '\[(.+)\]$' -- $line
                set -l branch (string match -r '\[(.+)\]$' -- $line)[2]
                echo $branch
            end
        end
    end
end

# Complete all branches (for smart mode)
function __wt_list_all_branches
    if git rev-parse --git-dir >/dev/null 2>&1
        set -l repo_root (git rev-parse --show-toplevel)

        __wt_list_worktree_branches

        set -l used_branches (__wt_list_worktree_branches)

        for branch in (git -C "$repo_root" branch --format='%(refname:short)' 2>/dev/null)
            if not contains $branch $used_branches
                echo $branch
            end
        end

        for branch in (git -C "$repo_root" branch -r --format='%(refname:short)' 2>/dev/null)
            set -l short_name (string replace -r '^origin/' "" -- $branch)
            if test "$short_name" != "HEAD"; and not contains $short_name $used_branches
                echo $short_name
            end
        end | sort -u
    end
end

# Branch completions for remove and z (only existing worktrees)
complete -c wt -f -n '__fish_seen_subcommand_from remove rm z' -a '(__wt_list_worktree_branches)'

# Branch completions for smart mode (all branches)
complete -c wt -f -n '__fish_use_subcommand' -a '(__wt_list_all_branches)'
