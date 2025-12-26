<img src="https://r2cdn.perplexity.ai/pplx-full-logo-primary-dark%402x.png" style="height:64px;margin-right:32px"/>

# based on attached file, write comprehensive guidelines for AI coding agents specifically Antigravity on how to integrate beads [https://github.com/steveyegge/beads](https://github.com/steveyegge/beads) as issue tracker and context engineering when working with a human prompt engineering on Antigravity

the goal is for Antigravity to call beads to automatically file issues, relate them per attached file, follow guidelines there for best practices so issues don't get lost even with coding session restart

however, I still want humans in the loop for review and approval

write the development workflow for both coding agent (Antigravity) and human (me), include best practices in software engineering

Antigravity coding agents integrate Beads as a persistent, git-backed issue tracker to maintain task context across sessions, using dependency graphs for ordered execution and preventing lost work. Humans retain approval authority through review gates before commits or closures. This workflow ensures software engineering best practices like traceability, incremental delivery, and verifiable progress.[^1][^2][^3]

## Setup Steps

Run these once per repository to initialize Beads for Antigravity workflows.

**Human (one-time setup):**

```
bd init  # Creates .beads/ with SQLite cache and JSONL (git-committed)
bd doctor  # Validates git hooks, merge drivers for conflict-free sync
echo "BEFORE ANY WORK: Run 'bd onboard' to integrate with Antigravity" >> AGENTS.md
```

Install git hooks during init for auto-sync on pull/push. Use `bd init --branch beads-metadata` for protected branches.[^1]

**Antigravity Agent (per session start):**

```
bd onboard  # Auto-updates AGENTS.md and CLAUDE.md with Beads protocols
bd ready --json  # Loads ready tasks (no blockers) into context
```

This bootstraps Antigravity with Beads CLI integration via JSON output.[^4][^1]

## Antigravity Workflow

Antigravity follows this cycle, filing issues proactively and querying ready work to avoid context bloat.

- **Session Start:** Run `bd ready --sort hybrid --json` to fetch top unblocked issues by priority/age. Select highest P0/P1 or human-specified ID.[^4][^1]
- **Claim Work:** `bd update <id> --status in_progress --assignee antigravity --json`. Add labels like `ai-working,backend`.[^1]
- **During Execution:** For discoveries (bugs, TODOs), create linked issues: `bd create "New bug" -t bug -p 1 --json` then `bd dep add <new-id> <current-id> --type discovered-from`.[^2][^1]
- **Quality Gates:** Run tests/lints post-changes. File P0 blockers if failures: `bd create "Fix lint errors" -p 0 -t bug`.[^4]

```
- **Propose Closure:** Never auto-close. Output `bd show <id> --json` diff and `git diff` for human review, then await approval before `bd close <id> --reason "Implemented and tested" --json`.[^1]
```

- **Sync:** `bd sync; git pull --rebase; git push` at session end to persist via git.[^3][^1]

Use hierarchical IDs for epics: Parent `bd-a1b2` auto-generates children like `bd-a1b2.1`. Detect cycles with `bd dep cycles`.[^1]

## Human Review Workflow

Maintain oversight without blocking velocity; approve via explicit gates.

- **Monitor Progress:** `bd list --status open --json` or `bv` TUI for real-time view. Check `bd ready` for next tasks.[^2][^1]

```
- **Approve Changes:** Review Antigravity's `git diff` and issue details (`bd show <id>`). Approve with "OK to commit <id>" or reject/add blockers.
```

```
- **Prioritize:** Create/adjust via CLI: `bd create "Epic: User Auth" -t epic -p 0`; `bd dep add <child> <epic> --type parent-child`. Reassign: `bd update <id> --assignee human`.[^1]
```

- **Session Handoff:** After Antigravity syncs, verify `git status` clean and `bd stats`. Add human issues for strategic shifts.[^4]
- **Compaction:** Manually trigger `bd compact --analyze --json` for old closed issues; approve summaries before `bd compact --apply`.[^1]

Delegate routine fixes but own epics and merges.

## Best Practices Table

| Aspect | Antigravity Agent | Human Reviewer | Rationale [^1][^4] |
| :-- | :-- | :-- | :-- |
| **Priorities** | P0=blockers/critical, P1=high-value, default P2 | Set P0 for business risks | Ensures focus on impact; `bd ready --priority 1` surfaces urgently. |
| **Dependencies** | Always link discoveries; use `blocks` for hard gates | Validate trees: `bd dep tree <epic>` | Prevents parallel rework; enables long-horizon planning across restarts. |
| **Labels** | `ai-generated,lint,tests-passed` | `human-review,deploy-ready` | Filters workflows; `bd list --label-any ai-generated,urgent`. |
| **Context Mgmt** | Query single issues via `bd show <id> --json`; avoid full dumps | Use `bd list --status open` for oversight | Saves tokens; git persistence eliminates markdown rot. |
| **Git Hygiene** | End with `bd sync; git commit -m "bd-<id>: <title>"; git push` | Merge PRs after review | Distributed DB via JSONL; hooks auto-reconcile multi-agent conflicts. |
| **Error Handling** | File P0 on failures; never ignore | `bd doctor` post-session | Audit trail in git; `bd migrate` for upgrades. |

<span style="display:none">[^10][^11][^5][^6][^7][^8][^9]</span>

<div align="center">⁂</div>

[^1]: https://github.com/steveyegge/beads

[^2]: bvtranscript.txt

[^3]: https://github.com/steveyegge/beads/releases

[^4]: https://betterstack.com/community/guides/ai/beads-issue-tracker-ai-agents/

[^5]: https://github.com/steveyegge/beads/blob/main/docs/ADVANCED.md

[^6]: https://www.linkedin.com/posts/steveyegge_github-steveyeggebeads-beads-a-memory-activity-7383408928665042944-tkcj

[^7]: https://github.com/steveyegge/beads/blob/main/RELEASING.md

[^8]: https://www.antigravityai.online

[^9]: https://stevespangler.com/experiments/newtons-inertia-beads-chain-fountain/

[^10]: https://ampcode.com/threads/T-adc03ba9-db60-49e6-bae9-e5f9749f4312

[^11]: https://antigravityai.pro

