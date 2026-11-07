#!/usr/bin/env bash
# Usage: run-agent.sh <claude|codex|muse|opencode> <prompt-file> [extra args...]
set -euo pipefail

agent="${1:?agent is required (claude|codex|muse|opencode)}"
promptfile="${2:?prompt file is required}"
shift 2

[ -f "$promptfile" ] || { echo "prompt file not found: $promptfile" >&2; exit 2; }
prompt="$(cat "$promptfile")"

case "$agent" in
  claude)   claude -p "$prompt" --dangerously-skip-permissions "$@" ;;
  codex)    codex exec --dangerously-bypass-approvals-and-sandbox "$@" "$prompt" ;;
  muse)     muse exec --yolo --prompt-file "$promptfile" "$@" ;;
  opencode) opencode run "$prompt" --auto "$@" ;;
  *)        echo "unknown agent: $agent" >&2; exit 2 ;;
esac
