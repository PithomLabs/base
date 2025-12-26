# Comprehensive Integration Guidelines: Beads Issue Tracker for Antigravity AI Coding Agents

## 1. Overview & Philosophy

Beads (bd) is a **dependency-aware issue tracker** designed specifically for AI-human collaborative coding workflows. Unlike ephemeral to-do lists that disappear with session resets, Beads provides:
- **Persistent issue tracking** with Git synchronization
- **Explicit dependency chains** preventing out-of-order work
- **Context-efficient queries** instead of loading entire markdown files
- **Team collaboration** through JSONL file commits

**Core Principle**: Beads serves as the **single source of truth** for work tracking, replacing Antigravity's internal to-do systems and ad-hoc markdown files.

## 2. Installation & Initialization

### For Antigravity (Automated Setup)
```bash
# Check if beads is initialized
if [ ! -d ".beads" ]; then
    echo "Initializing beads..."
    bd init
    bd doctor  # Fix any setup issues
fi

# Configure Antigravity to use beads
bd onboard  # Sets up agent instructions automatically
```

### For Human Review
- Verify `.beads/` directory creation
- Confirm `beads.jsonl` is in `.gitignore` (should NOT be ignored)
- Check git hooks are installed for auto-sync

## 3. AI Agent (Antigravity) Workflow Guidelines

### 3.1 Session Start Protocol
```bash
# ALWAYS run at session start
bd sync  # Pull latest issues from Git
bd ready --json  # Get next available work
```

**Antigravity MUST:**
1. Check for existing work assigned to itself
2. Review dependency graph before claiming work
3. Never start blocked tasks (use `bd blocked` to check)

### 3.2 Issue Creation & Management

**When discovering new work:**
```bash
# Instead of internal to-dos, create beads issues
bd create "Fix TypeScript compilation error in auth module" \
  -p 1 \
  -t bug \
  -d "Error occurs during production build, see logs in #123" \
  --label "backend,urgent"
```

**Issue Creation Rules:**
- **One issue per discrete problem/feature**
- **Always assign priority** (-p 0-4, 0=highest)
- **Use descriptive titles** that standalone
- **Include relevant context** in description
- **Add appropriate labels** for filtering

### 3.3 Dependency Management

**Explicit Dependencies Required:**
```bash
# When task B depends on task A
bd dep add bd-15 bd-14  # bd-15 blocks bd-14

# For epics/subtasks
bd create "Implement User Dashboard" -t epic -p 2
# Returns: bd-a3f8e9
bd create "Add dashboard UI components" -p 2  # Automatically becomes bd-a3f8e9.1
```

**Dependency Types to Use:**
1. **blocks** - Hard dependency (must complete first)
2. **related** - Soft connection for context
3. **parent-child** - Epic/subtask relationships
4. **discovered-from** - For issues found while working

### 3.4 Work Execution Flow

```bash
# 1. Claim work
bd update bd-14 --status in_progress --assignee antigravity

# 2. Work on it (Antigravity's normal coding process)

# 3. Create sub-issues for discovered work
#    Instead of context drift, file new issues
if [ "$DISCOVERED_NEW_PROBLEM" = true ]; then
    bd create "Fix memory leak in data processor" \
      -p 0 \
      -t bug \
      --parent bd-14  # Link to current work
fi

# 4. Complete with quality gates
run_tests && run_linter && run_build
if [ $? -eq 0 ]; then
    bd close bd-14 --reason "Implemented with tests passing"
else
    bd update bd-14 --status blocked -d "Tests failing, see output"
    # Create follow-up issue automatically
    bd create "Fix failing tests for bd-14" -p 0 -t bug
fi
```

### 3.5 Context Management During Coding

**Instead of:**
- Loading entire project markdown files
- Maintaining internal task lists
- Losing context on restart

**Antigravity Should:**
```bash
# Query for specific context
bd show bd-14 --json  # Get just this issue's context
bd dep tree bd-14  # See dependencies without full project load
bd list --status open --priority 0  # Get high-priority work only
```

### 3.6 "Landing the Plane" - Session Completion Protocol

**MANDATORY Before Session End:**
```bash
# 1. File issues for remaining work
for todo in "$REMAINING_TODOS"; do
    bd create "$todo" -p 3 -t task -d "Discovered during session, needs follow-up"
done

# 2. Run quality gates (if code changed)
pnpm check  # or equivalent
if [ $? -ne 0 ]; then
    bd create "Fix linting errors from session $(date)" -p 1 -t bug
fi

# 3. Update all issue statuses
bd update --status closed $(bd list --assignee antigravity --status in_progress --quiet)

# 4. PUSH TO REMOTE (CRITICAL)
git pull --rebase
bd sync  # Export beads changes to JSONL
git add .
git commit -m "Session complete: $(date)"
git push  # MUST SUCCEED

# 5. Verify push succeeded
if git status | grep -q "up to date with origin"; then
    echo "✓ Session properly landed"
else
    # Create issue for failed push
    bd create "Git push failed - manual intervention required" -p 0 -t blocker
    exit 1
fi
```

## 4. Human Review & Approval Workflow

### 4.1 Daily/Weekly Review Cycle

**Human Should:**
1. **Review `bd ready` output** - See what's available
2. **Adjust priorities** based on business needs:
   ```bash
   bd update bd-17 --priority 0  # Elevate important work
   bd update bd-23 --priority 3  # Deprioritize less critical
   ```
3. **Assign work to Antigravity**:
   ```bash
   bd update bd-14 --assignee antigravity --status open
   ```
4. **Review dependency graphs**:
   ```bash
   bd dep cycles  # Check for circular dependencies
   bd dep tree --all  # Visualize entire project
   ```

### 4.2 Code Review Integration

**When Antigravity completes work:**
1. **Human reviews code changes** (normal PR process)
2. **If changes approved**:
   ```bash
   bd close bd-14 --reason "PR #42 merged"
   bd sync  # Ensure closure propagates
   ```
3. **If changes need revisions**:
   ```bash
   bd update bd-14 --status blocked -d "Revisions requested in PR #42"
   bd create "Address PR feedback for bd-14" -p 1 -t task --parent bd-14
   ```

### 4.3 Priority & Scope Management

**Human Responsibilities:**
- Set business priorities (P0-P4)
- Break down epics into manageable tasks
- Identify and resolve blockers
- Manage stakeholder expectations through Beads visibility

**Tools for Humans:**
```bash
# Web UI (optional)
npx beads-ui

# Terminal UI
bv  # After installing beads-tui

# Custom reports
bd list --status open --json | jq '.[] | select(.priority == 0)'
```

## 5. Advanced Integration Patterns

### 5.1 Sentry Error Automation (from transcript)
```bash
# Automatically create beads issues from Sentry errors
use sentry-mcp to find errors | while read error; do
    bd create "$error" -t bug -p 1 \
      -d "Automated from Sentry. Root cause: $(sentry-seer $error)"
done
```

### 5.2 CI/CD Integration
```yaml
# GitHub Actions example
jobs:
  beads-sync:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Sync Beads
        run: |
          bd sync --no-daemon
          git diff --quiet beads.jsonl || \
            (git add beads.jsonl && git commit -m "ci: beads sync")
```

### 5.3 Multi-Agent Collaboration
```bash
# For teams with multiple AI agents
bd create "Feature X" -t epic -p 2
bd update bd-xx --assignee antigravity-1
bd create "Subtask A" --parent bd-xx --assignee antigravity-2
```

## 6. Best Practices for Software Engineering

### 6.1 Issue Quality Standards
- **Atomicity**: One issue = one complete unit of work
- **Reproducibility**: Include steps to reproduce bugs
- **Traceability**: Link to related code, PRs, commits
- **Estimability**: Use priority levels consistently

### 6.2 Dependency Hygiene
- **Minimize blocking dependencies** - Parallelize when possible
- **Regularly check for cycles** - Run `bd dep cycles` weekly
- **Use epics for large features** - Maintain hierarchy
- **Review dependency graphs** before sprint planning

### 6.3 Git & Beads Synchronization
- **Always `bd sync` before push/pull**
- **Resolve JSONL merge conflicts** manually (they're human-readable)
- **Verify auto-export** is working: check `.beads/beads.log`
- **Backup database** before major migrations

### 6.4 Performance Optimization
```bash
# Compact old issues periodically
bd compact --stats
# When >10MB with old closed issues:
bd compact --apply --id bd-xx --summary "Archived old sprint"
```

## 7. Troubleshooting & Recovery

### 7.1 Common Issues
```bash
# Database locked
bd --no-daemon ready

# Sync conflicts
git checkout --theirs beads.jsonl  # Or --ours
bd sync --force

# Missing issues after git pull
bd sync --import
```

### 7.2 Recovery Procedures
1. **Lost local changes**: `git log beads.jsonl` to recover
2. **Corrupt database**: Delete `.beads/*.db`, run `bd sync --import`
3. **Agent confusion**: Run `bd onboard` to reset instructions

## 8. Success Metrics & Continuous Improvement

**Track:**
- Issues completed vs. created
- Average time in progress
- Dependency graph complexity
- Context switching frequency

**Review weekly:**
```bash
bd stats --json > weekly_report.json
# Analyze with human to improve workflow
```

---

## Quick Reference Card for Antigravity

```bash
# START SESSION
bd sync
WORK=$(bd ready --json | jq -r '.[0].id')
bd update $WORK --status in_progress --assignee antigravity

# DURING SESSION
# Found new work? bd create "description" -p X -t type
# Need context? bd show $CURRENT_ID --json
# Blocked? bd update $CURRENT_ID --status blocked

# END SESSION
# 1. File follow-ups: bd create "..."
# 2. Quality gates: run_tests && run_lint
# 3. Close work: bd close $WORK
# 4. PUSH: git pull --rebase && bd sync && git push
# 5. VERIFY: git status shows "up to date"
```

**Golden Rule**: If it's not in Beads, it doesn't exist. If it's not pushed to Git, the session isn't complete.

---

*This workflow ensures continuous progress tracking, maintains human oversight, and leverages Beads' strengths for AI-human collaborative development. The system scales from solo projects to team environments while preserving engineering best practices.*
