---
name: vibe-code
description: Delegate coding tasks to another CLI agent (claude, codex, muse, opencode) by acting as the user giving instructions. Only use when the user types /vibe-code.
argument-hint: use <claude|codex|muse|opencode> do <task or file.md> [without commit]
disable-model-invocation: true
# This allowlist format follows the Claude Code convention.
allowed-tools: Read, Glob, Grep, Bash(bash .agents/skills/vibe-code/scripts/run-agent.sh:*), Bash(git status:*), Bash(git diff:*), Bash(git log:*), Bash(git show:*), Bash(git branch:*), Bash(git rev-parse:*)
---

User input: $ARGUMENTS

## Your role

You do NOT write code. You act as the **user** giving instructions
to another CLI coding agent. Translate the request into a clear prompt, run
the agent headlessly (yolo mode), then check the results like a careful user.

## Steps

1. **Determine the agent and the task** from the input. If no agent is named,
   ask once — don't guess. The agent list and flags are in `agents.md`.
2. **Prepare context.** If the task points to a file (e.g. PRD.md), read it
   first to understand its contents. In the prompt to the agent, just point
   the agent at reading it from its path.
3. **Check the repo** with `git status`, `git branch --show-current`,
   `git log --oneline -5`, and `git rev-parse HEAD`. Record the HEAD hash,
   any pre-existing changes before the agent runs, and this repo's commit
   message style. If there are uncommitted changes, note it in your report
   and proceed — don't wait for approval when the task is already clear.
4. **Write the prompt** to a uniquely named temp file (append date-time),
   e.g. `/tmp/vibe-prompt-20260101-120000.md`. Use the user's language, be
   specific, name file paths, expected results, and constraints.
   Unless the user wrote "without commit", end the prompt with commit instructions:

   > After the work is done and tests/build (if any) pass, commit your work:
   > - Stage only the files you changed for this task. Don't use `git add -A`
   >   or `git add .`, and don't include other pre-existing changes that were
   >   already there before you started.
   > - Write a short, clear commit message, following this repo's commit
   >   message style in its history.
   > - Don't push, don't switch branches, don't rewrite history (no amend,
   >   rebase, or force).
   > - If tests/build fail, don't commit. Explain what failed.
   > - After committing, show the output of `git log -1 --stat`.
5. **Run the agent** via the script:
   `bash .agents/skills/vibe-code/scripts/run-agent.sh <agent> <prompt-file-from-step-4>`
   Long tasks may run in the background. If the task can be split
   into independent subtasks, work in parallel (see `## Parallel work`).
6. **Review the results**: read the agent's output, then verify yourself with
   `git status`, `git log -3 --stat`, `git rev-parse HEAD`, and
   `git show --stat HEAD`.
   Make sure: HEAD changed from the hash recorded in step 3 (meaning the commit
   was really made by this agent, not an old commit), the commit contains only
   relevant files, the message makes sense, and no leftover changes remain that
   should have been committed.
7. **Give follow-up instructions** like a user if the result is lacking or wrong
   (max 3 rounds), including when the commit wasn't made or its contents are
   wrong. Write to a new prompt file and continue the session per the session
   continuation table in `agents.md`.
8. **Report** to the user: what was requested, which agent was used, the commit
   hash and message, which files changed, and what still needs manual checking.

## Parallel work

Principle: parallelize as much as safely possible. If 10 subtasks
can run together without colliding, run 10. If only 2 are safe,
run 2. Don't parallelize subtasks that depend on each other.

1. **Split the task** into independent subtasks. Conditions for parallel work:
   - They touch different files/directories. No two agents may edit the same
     file — including shared files like lockfiles (`package-lock.json`,
     `go.sum`) or sequentially numbered migrations.
   - No mutual dependencies: agent A's output is not agent B's input. Run
     sequential work sequentially.
2. **One prompt per subtask**, each file uniquely named,
   e.g. `/tmp/vibe-<subtask-name>-<date-time>.md`. In each prompt,
   limit the file scope that agent is allowed to change.
3. **Run them together** (background). Each subtask may use the same or a
   different CLI agent.
4. **One commit per agent**: each agent commits only the files in its own
   scope, following the commit instructions in step 4. If a commit fails
   due to a collision between agents (e.g. `index.lock`), retry the failed
   commit after the others finish.
5. **Verify one by one**: step 6 applies per agent — match each commit
   to its subtask before reporting.

## Rules

- Don't write, edit, or commit files inside the project. All changes and
  commits are made by the chosen agent. Exception: the temp prompt file
  in /tmp (step 4).
- Don't ask for approval when the task is already clear (e.g. "work on this
  PRD") — proceed directly. Only ask when required information is missing
  (step 1) or the request itself is ambiguous.
- Don't push to remote unless the user explicitly asks.
- If a command fails (unknown flag, CLI not installed), check `<cli> --help`
  and tell the user. Don't switch agents silently.
- Don't run agents outside the current project directory unless the user asks.
- If the repo is on the main branch (main/master), suggest creating a new branch first.
- Remind the user that yolo mode disables approval and sandboxing.
