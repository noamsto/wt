#!/usr/bin/env nix-shell
#!nix-shell -i bash -p vhs ttyd ffmpeg go git tmux
# shellcheck shell=bash
# Regenerate demo/wtc-demo.gif from demo/demo.tape using the local wtc source.
#
# Uses an isolated tmux socket (`-L wtcdemo`) so the recording never touches
# the caller's real tmux sessions. Pre-creates the session detached so the
# tape can attach instantly (no slow tmux startup inside the recording).
set -euo pipefail

TMUX_SOCKET="wtcdemo"
cd "$(git rev-parse --show-toplevel)"

BIN_DIR="$(mktemp -d)"
cleanup() {
    tmux -L "$TMUX_SOCKET" kill-server 2>/dev/null || true
    rm -rf "$BIN_DIR"
}
trap cleanup EXIT

go build -o "$BIN_DIR/wtc" ./cmd/wtc
export PATH="$BIN_DIR:$PATH"

# Build a fresh demo repo and pre-create the tmux session on an isolated socket.
bash demo/setup.sh /tmp/wtc-demo >/dev/null
tmux -L "$TMUX_SOCKET" kill-server 2>/dev/null || true
tmux -L "$TMUX_SOCKET" new-session -d -s demo -c /tmp/wtc-demo/repo

vhs demo/demo.tape
