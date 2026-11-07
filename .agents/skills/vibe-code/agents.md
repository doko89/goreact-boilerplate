# agents.md — CLI agent reference for vibe-code

Headless commands run via `scripts/run-agent.sh`, which forwards
extra arguments (`"$@"`) to each CLI. If flags change or a CLI
behaves differently, check `<cli> --help` — that wins.

## Headless commands (step 5)

| Agent | Command the script runs | Meaning of the yolo flag |
|---|---|---|
| `claude` | `claude -p "$prompt" --dangerously-skip-permissions` | Skip all permission checks |
| `codex` | `codex exec --dangerously-bypass-approvals-and-sandbox "$prompt"` | No approvals and no sandbox |
| `muse` | `muse exec --yolo --prompt-file "$promptfile"` | No approvals and no sandbox, workspace trusted |
| `opencode` | `opencode run "$prompt" --auto` | Auto-approve everything not explicitly denied |

## Session continuation (step 7)

Write follow-up instructions to a **new** prompt file (unique name), then
continue the session with extra arguments to the script:

| Agent | Example | Notes |
|---|---|---|
| `claude` | `run-agent.sh claude /tmp/vibe-followup-....md --continue` | `--continue` = latest session; or `--resume <id>` |
| `codex` | `run-agent.sh codex /tmp/vibe-followup-....md resume --last` | `resume` subcommand, prompt is sent after resuming |
| `muse` | `run-agent.sh muse /tmp/vibe-followup-....md` (no resume argument) | `muse resume` is interactive-only (TUI), so re-run `exec` with a self-contained follow-up prompt: summarize context + new instructions |
| `opencode` | `run-agent.sh opencode /tmp/vibe-followup-....md --continue` | `--continue` = latest session; or `--session <id>` |

Note on muse: since follow-ups run as a new session, the follow-up prompt
must include the latest status (latest commit/hash, changed files)
so the agent doesn't start over from scratch.
