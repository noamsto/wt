#!/usr/bin/env bash
# Build a throwaway git repo with worktrees covering each stale-detection path
# that wtc handles: merged-into-main, remote-deleted, dirty, and clean.
set -euo pipefail

DEMO_DIR="${1:-/tmp/wtc-demo}"
rm -rf "$DEMO_DIR"
mkdir -p "$DEMO_DIR"
cd "$DEMO_DIR"

git init -q --bare remote.git
git clone -q remote.git repo
cd repo
git config user.email demo@example.com
git config user.name Demo

echo "# Demo" > README.md
git add README.md
git commit -qm "initial"
git branch -qM main
git push -q origin main

mkdir -p .worktrees

# feature-a: merged into main → stale
git worktree add -q .worktrees/feature-a -b feature-a
(cd .worktrees/feature-a && echo a > a.txt && git add a.txt && git commit -qm "add a" && git push -qu origin feature-a)
git merge -q --no-ff feature-a -m "merge feature-a"
git push -q

# feature-b: clean, not merged → not stale
git worktree add -q .worktrees/feature-b -b feature-b
(cd .worktrees/feature-b && echo b > b.txt && git add b.txt && git commit -qm "add b" && git push -qu origin feature-b)

# feature-c: remote branch deleted → stale
git worktree add -q .worktrees/feature-c -b feature-c
(cd .worktrees/feature-c && echo c > c.txt && git add c.txt && git commit -qm "add c" && git push -qu origin feature-c)
git push -q origin --delete feature-c

# feature-d: dirty, not merged → flagged
git worktree add -q .worktrees/feature-d -b feature-d
(cd .worktrees/feature-d && echo "wip" > wip.txt)

echo
echo "Demo repo: $PWD"
git worktree list
