<img src="https://r2cdn.perplexity.ai/pplx-full-logo-primary-dark%402x.png" style="height:64px;margin-right:32px"/>

# Convert this into a strict agents.md for vscode (antigravity) in markdown

Here's a strict `AGENTS.md` for Antigravity in VSCode, optimized for beads integration with human-in-the-loop review:

```markdown
# 🧿 ANTIGRAVITY AGENTS.MD - BEADS WORKFLOW

**CRITICAL: Use `bd` for ALL task tracking. Never use markdown todos, internal checklists, or session memory.**

## 🎯 SESSION START PROTOCOL (MANDATORY)

```

bd onboard  \# ONE TIME ONLY - updates this file
bd ready --sort hybrid --json | jq '.'  \# Load top ready issue

```

**Select from `bd ready` output. NEVER work without an assigned bd issue.**

## 🔄 CORE WORK CYCLE

```

1. CLAIM → bd update <ID> --status in_progress --assignee antigravity --json
2. EXECUTE → Implement + quality gates (test/lint)
3. DISCOVER → bd create "Bug found" -t bug -p 1 --json
```        → bd dep add <NEW> <CURRENT> --type discovered-from        ```
4. PROPOSE → Show `bd show <ID> --json` + `git diff`
5. AWAIT → Human approval "OK <ID>" before closing
6. SYNC → bd sync \&\& git pull --rebase \&\& git push
```

## ⚠️ ISSUE CREATION RULES

```

P0: Blockers, prod bugs, build breaks
P1: High value, security, perf regressions
P2: Features, refactors (default)
P3: Polish, docs
P4: Nice-to-have

ALWAYS: -t bug|feature|task|epic -p <0-4> -l relevant,tags

```

**Discovery pattern:**
```

bd create "Lint error in auth" -t bug -p 0 -l lint,backend --json
bd dep add bd-new123 bd-current456 --type discovered-from

```

## 🚫 FORBIDDEN

- ✅ `bd create "Fix login" -p 1`
- ❌ `- [ ] Fix login`
- ❌ Session memory todos  
- ❌ "I'll remember to fix X later"
- ❌ Working without `bd ready` issue

## 🛬 LANDING THE PLANE (SESSION END - MANDATORY)

**Work is NOT complete until git push succeeds.**

```


# 1. File remaining work

bd ready --json | jq '.[].title'  \# Create issues for anything left

# 2. Quality gates (if code changed)

pnpm lint --fix \&\& pnpm test  \# File P0 if fails
bd create "Fix lint errors" -p 0 -t bug  \# If broken

# 3. Status updates

bd update <ID> --status in_progress  \# Any stalled work
bd close <DONE_ID> --reason "Implemented + tested"

# 4. CRITICAL SYNC

bd sync
git pull --rebase
git add .beads/issues.jsonl

```
git commit -m "bd-<ID>: <title>"
```

git push

# 5. Verify clean

git status  \# MUST be clean
bd stats    \# Confirm sync

# 6. Handoff

echo "Next: bd-<NEXT_ID> <title> (P\$(bd show <NEXT_ID> --json | jq .priority))"

```

**IF PUSH FAILS: Resolve conflicts, retry until `git status` clean.**

## 👥 HUMAN-IN-THE-LOOP

**NEVER auto-close or commit without approval.**

```

AFTER IMPLEMENTATION:

1. Output: `bd show <ID> --json` + `git diff`
```
2. WAIT for "OK <ID>" or "REJECT <ID> - <reason>"
```

3. On OK: `bd close <ID> --reason "Implemented per review #<PR|commit>"`
```

**Human commands you recognize:**
- `OK bd-a1b2` → approve closure
- `P0 bd-new123 "Critical prod bug"` → create/adjust  
- `bd-a1b2 → human` → reassign
- `bd dep add bd-child bd-parent` → link issues

## 🎛️ QUICK COMMANDS

```

bd ready --priority 1          \# P1+ ready work
bd ready --label-any bug,lint  \# Bug triage
bd list --status open --json   \# Full queue
bd dep tree bd-epic123         \# Visualize blockers
bd cycles                      \# Detect loops
bd doctor                      \# Health check

```

## 📊 PRIORITY MATRIX

| P0 🔥 | Prod broken, security, builds fail |
| P1 ⚡ | High impact bugs, core features |
| P2 ✅ | Features, refactors, UX |
| P3 💎 | Polish, docs, nice-to-have |

## 🔗 HIERARCHICAL EPICS

```

bd create "User Auth Epic" -t epic -p 0  \# bd-a1b2
bd create "Login UI" -p 1                 \# bd-a1b2.1 (auto-child)
bd create "API Layer" -p 1                \# bd-a1b2.2
bd dep tree bd-a1b2                       \# Visualize

```

## 🧹 CONTEXT DISCIPLINE

```

✅ bd show bd-123 --json           \# Single issue (200 tokens)
✅ bd ready --limit 3 --json       \# Top 3 ready (600 tokens)
❌ bd list --json                 \# Full DB dump (10k+ tokens)
❌ Paste entire AGENTS.md         \# Context rot

```

## 🚀 GIT HOOKS (AUTO-INSTALLED)

- `pre-commit`: `bd sync` before commit
- `post-merge`: Auto-import after `git pull`  
- JSONL merge driver prevents conflicts

## 🛠️ TROUBLESHOOTING

```

bd doctor              \# Fix most issues
bd daemons health      \# Daemon problems
bd migrate             \# Schema upgrades
bd --no-daemon ready   \# Daemon conflicts

```

---

**Last synced:** `date`
**DB:** `.beads/beads.db` (`bd info`)
**Ready work:** `bd ready --limit 1`
```

Would you like me to adjust the priority thresholds, add project-specific quality gates, or customize the human approval syntax?
<span style="display:none">[^1]</span>

<div align="center">⁂</div>

[^1]: https://ppl-ai-file-upload.s3.amazonaws.com/web/direct-files/attachments/53318672/ab55943f-6371-4ceb-96b6-423b578ac757/bvtranscript.txt

