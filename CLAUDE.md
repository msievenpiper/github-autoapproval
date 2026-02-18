# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Run
go run index.go --branch=<branch> <owner/repo> [<owner/repo>...]

# Build
go build

# Probe (dry-run, no actual approvals)
go run index.go --branch=<branch> --probe <owner/repo>

# Auto-merge after approval
go run index.go --branch=<branch> --merge <owner/repo>
```

There are no tests in this project.

## Prerequisites

Requires `gh` CLI authenticated via `gh auth login`.

## Architecture

**Entry point** (`index.go`): Validates `gh` auth, parses inputs, iterates repos → PRs → approve (and optionally merge).

**`internal/` package:**
- `inputs.go` — CLI flag parsing: `--branch` (required), `--probe` (dry-run), `--merge` (auto-merge); positional args are repo slugs (`owner/repo`)
- `gitAuth.go` — Checks `gh auth status` output for "Logged in" and parses the username via regex
- `pullRequests.go` — All PR operations via `gh.Exec()` from `github.com/cli/go-gh/v2`:
  - `GetPullRequests`: runs `gh pr list --search <branch>`, parses tab-separated output
  - `getPullRequestStatus`: runs `gh pr view --json latestReviews,state,author`
  - `ApprovePullRequest`: skips if already approved; runs `gh pr review --approve` unless probe mode
  - `MergePullRequest`: runs `gh pr merge`
- `helpers.go` — Generic `Filter[T]` and `Map[T, E]` utilities

**Key data flow:** `GetPullRequests` fetches and parses PR list → each PR gets its review state hydrated via `getPullRequestStatus` → `IsAppoved()` checks `latestReviews` array for any "approve" state before acting.

**`--probe` flag**: gates the actual `gh pr review --approve` call in `ApprovePullRequest` and prevents `MergePullRequest` from being called at all (checked in `index.go`).
