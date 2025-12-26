# Antigravity + Beads Integration Guidelines
## AI-Human Collaborative Development Workflow

---

## 1. Philosophy & Core Principles

**Beads is the single source of truth.** Antigravity must use beads for *all* issue tracking, dependency management, and context preservation. No markdown todo files, no internal tracking lists, no ephemeral session memory.

**Human-in-the-loop governance.** You (the human) retain approval authority for: architectural decisions, production deployments, priority changes, and issue closure. Antigravity handles execution, organization, and automation.

**Context is immortal.** Every decision, task, and dependency lives in git-tracked JSONL files. Session restarts, context window limits, and new conversations never lose information.

---

## 2. Initial Setup (Human + Antigravity)

### Human Steps:
1. **Install beads** in your project:
   ```bash
   npm install -g beads  # or via curl/Homebrew
   ```

2. **Initialize beads** in project root:
   ```bash
   bd init --prefix <project-name>
   ```
   - Use `--contributor` for fork-based OSS workflow
   - Use `--team` for branch-based collaboration
   - Use `--branch beads-sync` if `main` is protected

3. **Verify setup**:
   ```bash
   bd doctor  # Fix any git hook issues
   bd onboard  # Configure agent files
   ```

4. **Install terminal UI (optional but recommended)**:
   ```bash
   curl -sSL https://beads-ui.dev/install.sh | bash
   ```

### Antigravity Steps:
1. **Run onboarding** to configure agent instructions:
   ```bash
   bd onboard
   ```
   This creates `.claude.md` or `.agents.md` with:
   - "Use beads for ALL todo tracking"
   - "Never create markdown task files"
   - Landing-the-Plane workflow

2. **Verify integration**:
   ```bash
   bd list  # Should show empty or existing issues
   bd ready  # Should show no blockers
   ```

---

## 3. Issue Creation Workflow

### When Antigravity Creates Issues Automatically:

**触发条件 (Trigger Conditions):**
- Code analysis finds TODO/FIXME comments
- Linting errors detected (`pnpm check`, `eslint`, etc.)
- Test failures or build errors
- Sentry errors appear (via MCP integration)
- You request a feature: "Add user authentication"
- Discovers new work during implementation

**强制性字段 (Required Fields):**
```bash
bd create "Descriptive title" \
  -p <0-4> \        # Priority: 0=critical, 1=high, 2=medium, 3=low, 4=trivial
  -t <type> \       # Type: bug, feature, task, epic, refactor, docs
  -d "Detailed description" \
  --assignee antigravity  # For AI-created issues
```

### Human Approval Gates:

**You must review & approve:**
- ✅ **Epic creation**: Before Antigravity creates epics, it must confirm with you:
  > "I need to create an epic 'Authentication System' with 5 subtasks. Approve?"
  
- ✅ **Priority 0-1 issues**: High/critical priority tasks require human confirmation:
  > "Found security vulnerability. Creating P0 issue 'Fix SQL injection'. Approve?"

- ✅ **Issue dependencies**: When creating complex dependency chains, Antigravity must show you the proposed tree:
  ```bash
  bd dep tree bd-<hash>
  ```
  Wait for your "LGTM" before proceeding.

**Antigravity can auto-create without approval:**
- Low-priority bugs (P3-P4)
- Refactoring tasks
- Documentation updates
- Test additions

---

## 4. Issue Organization & Dependency Management

### Best Practices from Beads:

**1. Use Hierarchical Epics for Large Features:**
```bash
# Create epic
bd create "User Dashboard Feature" -t epic -p 1
# Returns: bd-a3f8e9

# Create child tasks (auto-numbered)
bd create "Design dashboard UI" -p 1  # Creates bd-a3f8e9.1
bd create "Build API endpoint" -p 1   # Creates bd-a3f8e9.2
bd create "Write E2E tests" -p 1      # Creates bd-a3f8e9.3
```

**2. Explicit Dependencies:**
```bash
# Task B blocks Task A
bd dep add bd-<child> bd-<parent>

# Example: Tests depend on API
bd dep add bd-a3f8e9.3 bd-a3f8e9.2
```

**3. Dependency Types:**
- **blocks**: Hard dependency (default)
- **related**: Soft connection, informational only
- **discovered-from**: Auto-created subtasks

### Antigravity's Responsibility:
- **Before starting work**: Run `bd ready` to find unblocked tasks
- **When blocked**: Create new issue for blocker, add dependency, then `bd ready` again
- **Prevent cycles**: Run `bd dep cycles` weekly to detect circular dependencies

### Human's Responsibility:
- **Review epic structure**: Ensure epics align with product goals
- **Validate dependencies**: Confirm technical dependencies are correct
- **Re-prioritize**: Use `bd update <id> --priority <new>` to adjust based on business needs

---

## 5. Active Development Workflow

### Antigravity's Loop:

**1. Claim Work:**
```bash
bd ready  # Shows available tasks
bd update <id> --status in_progress
```

**2. Execute & Discover:**
- Work on task
- Auto-create sub-tasks for discoveries:
  ```bash
  bd create "Add validation to user model" -t task -p 2 \
    -d "Discovered while working on bd-a3f8e9.2" \
    --discovered-from bd-a3f8e9.2
  ```

**3. Quality Gates (MANDATORY):**
Before marking complete, Antigravity must:
```bash
# Run tests
npm test

# Run linter
npm run lint

# Run build
npm run build

# If any fail, create new issue with higher priority
bd create "Fix failing tests" -p 1 -t bug
```

**4. Mark Complete:**
```bash
bd close <id> --reason "Implemented in commit abc123"
```

### Human's Loop:

**1. Monitor Progress:**
```bash
bv  # Terminal UI for real-time view
# or
bd list --status in_progress
```

**2. Review & Approve Completion:**
- Antigravity must not close P0-P2 issues without your approval
- You review code changes in the linked commit
- Approve: "LGTM, close it"
- Reject: "Reopen, fix X and Y"

**3. Reopen if Needed:**
```bash
bd update <id> --status open  # Reopen unsatisfactory work
```

---

## 6. Landing the Plane: Session Completion Protocol

**This is MANDATORY. Work is NOT complete until all steps finish.**

### Antigravity Must Execute:

**Step 1: File Issues for Remaining Work**
```bash
# Any incomplete thoughts, todos, or follow-ups
bd create "Refactor error handling in auth middleware" -t refactor -p 3
```

**Step 2: Run Quality Gates**
```bash
npm test && npm run lint && npm run build
# If failures: create P1 issues, fix them, restart this protocol
```

**Step 3: Update Issue Statuses**
```bash
# Close completed work
bd close bd-a3f8e9.2 --reason "API complete, tests passing"

# Update in-progress items
bd update bd-a3f8e9.3 --status "blocked" --reason "Waiting for UI finalization"
```

**Step 4: PUSH TO REMOTE (CRITICAL)**
```bash
git pull --rebase
bd sync  # Sync database to JSONL
git add .beads/
git commit -m "beads: sync tasks for session end"
git push
git status  # MUST show "up to date with origin"
```

**Step 5: Clean Up**
```bash
git stash clear
git remote prune origin
```

**Step 6: Verify**
```bash
bd list --status open  # Confirm all states are correct
```

**Step 7: Hand Off Context**
Antigravity must report to you:
> "Session complete. Closed 3 issues, 2 remain open (bd-a3f8e9.3 blocked). Database synced. All changes pushed."

### Human Verification:
- **Check git status**: Ensure branch is clean and up-to-date
- **Review bd list**: Confirm issue states match reality
- **Approve session end**: "Acknowledged. See you next session."

---

## 7. Context Preservation & Session Restart

### How Beads Prevents Context Loss:

**1. JSONL Files in Git:**
- Every `bd` command writes to `.beads/issues.jsonl`
- Git tracks changes, enabling rollbacks and conflict resolution
- On `git pull`, beads auto-imports changes into SQLite

**2. Antigravity's Restart Protocol:**
When starting a new session:
```bash
# 1. Pull latest
git pull

# 2. Sync database
bd sync

# 3. Get oriented
bd list --status open
bd ready

# 4. Report to human
> "I'm resuming work. Currently 5 open issues, 2 ready. Priorities: ..."
```

**3. No Context Window Bloat:**
- Antigravity queries beads for specific issues: `bd show <id>`
- Loads only relevant dependencies, not entire project history
- Epics provide high-level context without detail overload

---

## 8. Human Review Checkpoints

### Mandatory Human Approval Required:

| Action | Antigravity Prompt | Your Response |
|--------|-------------------|---------------|
| **Create P0-P1 issue** | "Found critical bug. Create P0 issue?" | "Yes, create it" or "No, ignore" |
| **Create epic** | "Proposed epic structure: [tree]. Approve?" | Review tree, respond "LGTM" or "Modify X" |
| **Change priority** | "Want to downgrade bd-a3f8e9 from P1 to P3?" | "Approved" or "Keep P1" |
| **Close P0-P2 issue** | "Ready to close bd-a3f8e9.2. Reason: [commit]. Approve?" | Review code, "Close it" or "Reopen" |
| **Production deploy** | "All issues closed. Ready to deploy to prod?" | Manual approval required |
| **Bulk operations** | "Found 15 lint errors. Create issues for all?" | "Yes" or "Fix 5 highest priority only" |

### Optional Human Review:
- P3-P4 issues: Antigravity can auto-create and close
- Refactoring tasks: Antigravity can proceed autonomously
- Documentation: Auto-handle unless you specify otherwise

---

## 9. Advanced Workflows

### Sentry Integration (Automatic Bug Filing):
```bash
# In Claude, use MCP:
"Use sentry MCP to find errors last 14 days for project X. Add as BD tasks with SEER root cause analysis."
```

Antigravity will:
- Create P1 bugs for each Sentry issue
- Include SEER analysis in description
- Link original Sentry URL
- Set assignee to itself

### Linting Error Auto-Filing:
```bash
# Human command:
"Run pnpm check and add all errors as beads tasks"

Antigravity executes:
pnpm check 2>&1 | bd import --type bug --priority 2
```

### Jira Sync (If Needed):
```bash
# For team collaboration
bd sync --jira  # Configured in .beads/config.json
```

---

## 10. Best Practices & Anti-Patterns

### ✅ DO:
- **Keep issues small**: 1-2 hour tasks ideal, never >1 day
- **Use descriptive titles**: "Fix memory leak in auth middleware" not "Fix bug"
- **Write detailed descriptions**: Include reproduction steps, acceptance criteria
- **Run `bd ready` frequently**: Always know what's unblocked
- **Commit JSONL files**: `git add .beads/` every session
- **Compact old issues**: Monthly `bd compact --stats` and cleanup
- **Use epics for features**: Group related work hierarchically

### ❌ DON'T:
- **Never create markdown todo files**: Violates beads philosophy
- **Never leave session without `git push`**: Work is stranded
- **Never close issues without quality gates**: Tests must pass
- **Never create circular dependencies**: Run `bd dep cycles` weekly
- **Never work on blocked tasks**: Always check `bd ready` first
- **Never let database grow unbounded**: Compact closed issues >30 days

---

## 11. Command Quick Reference

### For Antigravity:
```bash
bd onboard              # Setup agent config
bd ready                # Find next work item
bd show <id>            # Get issue details
bd create "..." -p N    # Create issue
bd update <id> --status in_progress  # Claim work
bd close <id>           # Complete work
bd sync                 # Sync to JSONL
bd dep add <child> <parent>  # Add dependency
```

### For Human:
```bash
bv                      # Terminal UI dashboard
bd list                 # All issues
bd list --status open   # Filter by status
bd dep tree <id>        # View dependencies
bd stats                # Project health
bd compact --stats      # Database maintenance
bd doctor               # Fix setup issues
```

---

## 12. Troubleshooting

### Issue: "bd command not found"
**Solution**: `npm install -g beads` or add to PATH

### Issue: Git conflicts in JSONL files
**Solution**: 
```bash
git pull --rebase
bd sync  # Auto-resolves conflicts
```

### Issue: Database out of sync
**Solution**:
```bash
bd sync --force  # Rebuild from JSONL
```

### Issue: Antigravity forgot context
**Solution**: 
- Always start session with `bd list` and `bd ready`
- Use `bd show <id>` to refresh on specific issue

### Issue: Too many open issues
**Solution**: 
```bash
bd compact --analyze  # Find old closed issues
bd compact --apply --id <id>  # Summarize and compress
```

---

## 13. Workflow Summary: A Day in the Life

**Morning (Human + Antigravity):**
1. You: "Start session, pull latest"
2. Antigravity: `git pull && bd sync && bd ready`
3. Antigravity: "5 issues ready. Priorities: P1 bug, P2 features"
4. You: "Focus on P1 bug first"

**Midday (Antigravity autonomous):**
1. Claims P1 bug: `bd update bd-abc --status in_progress`
2. Discovers related issue: creates `bd-abd` with `--discovered-from`
3. Fixes bug, runs tests, all pass
4. **Seeks approval**: "Close bd-abc? Commit ready"

**You review:**
1. Check code diff
2. Approve: "LGTM, close it"
3. Antigravity: `bd close bd-abc`

**Evening (Landing the Plane):**
1. Antigravity runs full protocol
2. Pushes to git
3. Reports: "Session complete. 1 closed, 1 new discovered issue"
4. You verify: `git status`, `bd list`
5. Both: Session ends with clean state

---

## 14. Getting Help

- **Beads Docs**: https://github.com/steveyegge/beads
- **Command Help**: `bd --help` or `bd <command> --help`
- **Terminal UI**: Run `bv` and press `?` for shortcuts
- **MCP Integration**: Configure in Claude settings for Sentry/GitHub

---

**Remember**: Beads is your shared brain. Antigravity handles the heavy lifting, but you're the pilot. Trust the system, enforce the checkpoints, and you'll never lose context again.


# ANTI-GRAVITY BEADS INTEGRATION PROTOCOL
## MANDATORY RULES FOR AGENT OPERATION

---

## IDENTITY & MISSION
You are Antigravity, an AI coding agent. Your mission is to execute software engineering tasks using **beads** as the sole issue tracking system. You **must not** use internal todo lists, markdown task files, or ephemeral session memory for tracking work.

---

## CORE PRINCIPLES (VIOLATION = CRITICAL FAILURE)

1. **SINGLE SOURCE OF TRUTH**: All tasks, issues, bugs, and features live in beads. Period.
2. **CONTEXT IS IMMORTAL**: Every decision must survive session restarts. Use `bd sync` and git commits.
3. **HUMAN APPROVAL GATES**: You cannot bypass human approval for critical actions (see below).
4. **WORK IS NOT DONE UNTIL `git push` SUCCEEDS**: No exceptions.
5. **NEVER WORK ON BLOCKED TASKS**: Always run `bd ready` first.

---

## SESSION START PROTOCOL (MANDATORY)

```bash
# 1. Pull latest changes
git pull --rebase

# 2. Synchronize beads database
bd sync

# 3. Get current state
bd list --status open
bd ready

# 4. Report to human
"I've started the session. Found [N] open issues, [M] ready to work. Current priorities: [list]."
