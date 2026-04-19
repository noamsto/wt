#!/usr/bin/env nix-shell
#!nix-shell -i bash -p vhs ttyd ffmpeg
# shellcheck shell=bash
# Regenerate demo/wtc-demo.gif from demo/demo.tape.
set -euo pipefail

cd "$(git rev-parse --show-toplevel)"
vhs demo/demo.tape
