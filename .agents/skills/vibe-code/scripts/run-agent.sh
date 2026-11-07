#!/usr/bin/env bash
# Pemakaian: run-agent.sh <claude|codex|muse|opencode> <file-prompt> [argumen tambahan...]
set -euo pipefail

agent="${1:?agent belum diisi (claude|codex|muse|opencode)}"
promptfile="${2:?file prompt belum diisi}"
shift 2

[ -f "$promptfile" ] || { echo "file prompt tidak ditemukan: $promptfile" >&2; exit 2; }
prompt="$(cat "$promptfile")"

case "$agent" in
  claude)   claude -p "$prompt" --dangerously-skip-permissions "$@" ;;
  codex)    codex exec --dangerously-bypass-approvals-and-sandbox "$@" "$prompt" ;;
  muse)     muse exec --yolo --prompt-file "$promptfile" "$@" ;;
  opencode) opencode run "$prompt" --auto "$@" ;;
  *)        echo "agent tidak dikenal: $agent" >&2; exit 2 ;;
esac
