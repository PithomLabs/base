# Agent Instructions

This project uses **bd** (beads) for issue tracking. Run `bd onboard` to get started.

## Quick Reference

```bash
bd ready              # Find available work
bd show <id>          # View issue details
bd update <id> --status in_progress  # Claim work
bd close <id>         # Complete work
bd sync               # Sync with git
```

## Landing the Plane (Session Completion)

**When ending a work session**, you MUST complete ALL steps below. Work is NOT complete until `git push` succeeds.

**MANDATORY WORKFLOW:**

1. **File issues for remaining work** - Create issues for anything that needs follow-up
2. **Run quality gates** (if code changed) - Tests, linters, builds
3. **Update issue status** - Close finished work, update in-progress items
4. **PUSH TO REMOTE** - This is MANDATORY:
   ```bash
   git pull --rebase
   bd sync
   git push
   git status  # MUST show "up to date with origin"
   ```
5. **Clean up** - Clear stashes, prune remote branches
6. **Verify** - All changes committed AND pushed
7. **Hand off** - Provide context for next session

**CRITICAL RULES:**
- Work is NOT complete until `git push` succeeds
- NEVER stop before pushing - that leaves work stranded locally
- NEVER say "ready to push when you are" - YOU must push
- If push fails, resolve and retry until it succeeds


# AGENTS.MD - Antigravity × Beads Integration Guide

## Prime Directive

> **If it is not in Beads, it does not exist.**

Beads (`bd`) is the **authoritative issue tracker, memory system, and context engine**. Chat history, markdown to-do lists, and agent memory are **NOT** sources of truth.

---

## Table of Contents

1. [Core Principles](#core-principles)
2. [Initial Setup](#initial-setup)
3. [Issue Creation Protocol](#issue-creation-protocol)
4. [Work Execution Flow](#work-execution-flow)
5. [Landing the Plane](#landing-the-plane)
6. [Human Review Workflow](#human-review-workflow)
7. [Best Practices](#best-practices)
8. [Troubleshooting](#troubleshooting)
9. [Quick Reference](#quick-reference)

---

## Core Principles

### Single Source of Truth (SSOT)

| Artifact | Status |
|----------|--------|
| Beads issues (SQLite + JSONL) | ✅ Canonical |
| Git history | ✅ Canonical |
| Chat history | ❌ Ephemeral |
| Markdown task lists | ❌ Forbidden |
| Agent memory | ❌ Forbidden |

### Separation of Concerns

- **Beads**: Persistence, ordering, dependency graph, auditability
- **Antigravity**: Execution, discovery, automation, implementation
- **Human**: Intent, approval, prioritization, architectural decisions

### Context Engineering Rule

**Never carry more context than the current issue requires.**

- Load only: Issue description, dependencies, relevant files
- Never: Full backlog, reconstructed state from chat, entire project history

---

## Initial Setup

### Prerequisites

```bash
# Install Beads
npm install -g @beads/cli
# or
curl -fsSL https://beads.dev/install.sh | sh
```

### Human: One-Time Repository Setup

```bash
# Navigate to project root
cd your-project

# Initialize Beads
bd init --prefix <project-name>

# Run diagnostics
bd doctor

# Install git hooks (recommended)
# Hooks auto-sync on git pull/push

# For protected main branch
bd init --branch beads-sync

# Verify setup
bd list
bd ready
```

**Expected directory structure:**
```
your-project/
├── .beads/
│   ├── beads.db           # SQLite cache (local only)
│   ├── issues.jsonl       # Git-tracked issue state
│   ├── daemon.sock        # Background sync service
│   └── config.json        # Project configuration
├── .agents.md             # This file (auto-generated)
└── .claude.md             # Agent instructions (auto-generated)
```

### Antigravity: Session Bootstrap

```bash
# ALWAYS run at session start
bd onboard  # Configures agent instructions (first time only)
bd sync     # Pull latest issues from git
bd ready    # Show available work
```

---

## Issue Creation Protocol

### When to Create Issues
Create a beads issue for EVERY discrete unit of work:
- Every bug found
- Every feature request
- Every refactoring task
- Every dependency update
- Every test that needs writing

### Issue Creation Template
```bash
# BASIC TEMPLATE
bd create "[BRIEF DESCRIPTIVE TITLE]" \
  -p [0-4] \
  -t [bug|feature|task|epic] \
  -d "[Detailed description including:
      - Steps to reproduce (for bugs)
      - Acceptance criteria (for features)
      - Current vs desired behavior
      - Relevant code locations
      - Any error messages]" \
  --label "[comma,separated,labels]"

# EXAMPLES
# Bug:
bd create "TypeError in user authentication middleware" \
  -p 0 \
  -t bug \
  -d "Server crashes when JWT token is malformed. 
      Steps:
      1. Send POST to /api/login with invalid token
      2. Observe 500 error
      
      Location: src/middleware/auth.js:45
      Error: 'Cannot read property userId of undefined'" \
  --label "backend,security,blocker"

# Feature:
bd create "Add user profile image upload" \
  -p 2 \
  -t feature \
  -d "Users should be able to upload profile pictures.
      
      Acceptance Criteria:
      - [ ] Upload via drag-and-drop or file picker
      - [ ] Support JPG, PNG, WebP up to 5MB
      - [ ] Auto-generate thumbnails
      - [ ] Show in user profile and comments
      
      Dependencies: Cloud storage service configured" \
  --label "frontend,feature,ui"

# Technical Debt:
bd create "Refactor legacy payment processor module" \
  -p 3 \
  -t task \
  -d "Module is tightly coupled and hard to test.
      
      Goals:
      - Extract interfaces for payment providers
      - Add unit tests (currently 0% coverage)
      - Implement retry logic with exponential backoff
      
      Files: src/payment/legacy-processor.js" \
  --label "refactor,backend,tech-debt"
```

### Priority Guidelines
- **P0**: Blockers, production outages, security issues
- **P1**: Critical features, major bugs affecting users
- **P2**: Important features, significant improvements
- **P3**: Nice-to-have features, minor improvements
- **P4**: Backlog, future enhancements

### When Issues Are Created

**Antigravity MUST create Beads issues for:**
- Every bug discovered (compilation errors, runtime errors, test failures)
- Every TODO/FIXME comment found in code
- Every feature request from human
- Every refactoring opportunity identified
- Every test that needs writing
- Every discovered dependency or blocker
- Every error from external tools (linter, Sentry, CI/CD)

**Antigravity MUST NOT:**
- Use internal to-do lists
- Create markdown task files
- Store work items in chat memory
- Say "I'll remember to fix X later"

---

### Issue Creation Workflow: Step by Step

#### Step 1: Discovery (Antigravity)

**Trigger conditions:**
```bash
# Code analysis
grep -r "TODO\|FIXME\|HACK" src/

# Linting errors
npm run lint
pnpm check

# Test failures
npm test

# Build errors
npm run build

# External integrations
# - Sentry MCP: production errors
# - GitHub Issues: imported bugs
# - Manual human request
```

#### Step 2: Classification (Antigravity)

**Analyze and categorize the issue:**

```bash
# Required metadata
TYPE: bug | feature | task | epic | chore | docs | investigation
PRIORITY: 0 (critical) | 1 (high) | 2 (medium) | 3 (low) | 4 (backlog)
DEPENDENCIES: List of blocking issues (if any)
LABELS: Tags for filtering (e.g., backend, frontend, security, test)
```

**Priority Guidelines:**

| Priority | Use Case | Examples |
|----------|----------|----------|
| P0 | Production broken, security vulnerability, data loss | SQL injection, server crash, build broken |
| P1 | High-impact bug, critical feature, blocking dependency | Authentication broken, major UI bug |
| P2 | Standard feature, significant improvement | New API endpoint, performance optimization |
| P3 | Nice-to-have, minor bug, tech debt | UI polish, code refactoring |
| P4 | Backlog, future consideration | Ideas, long-term improvements |

#### Step 3: Proposal (Antigravity → Human)

**Antigravity MUST propose before creating:**

```
I discovered new work that should be tracked:

Title: "Fix TypeScript compilation error in auth middleware"
Type: bug
Priority: 1 (blocks development)
Dependencies: None
Labels: backend, typescript, blocker

Description:
Server fails to build with error: "Property 'userId' does not exist on type 'JWTPayload'"

Location: src/middleware/auth.ts:45
Error message: [paste error]
Steps to reproduce:
1. Run `npm run build`
2. Observe compilation failure

Suggested fix: Add userId to JWTPayload interface

Should I create this issue?
```

**Human responses:**
- `"Yes, create it"` → Proceed to Step 4
- `"Yes, but P0"` → Adjust priority, then create
- `"No, expected behavior"` → Discard
- `"File as P2, not a blocker"` → Adjust priority, then create

#### Step 4: Creation (Antigravity, after approval)

```bash
# Basic creation
bd create "Fix TypeScript compilation error in auth middleware" \
  -t bug \
  -p 1 \
  -d "Server fails to build with error: Property 'userId' does not exist on type 'JWTPayload'. Location: src/middleware/auth.ts:45" \
  --label backend,typescript,blocker

# Returns: bd-a3f8e9 (issue ID)
```

**Good vs Bad Issue Titles:**

✅ **Good:**
- "Fix SQL injection vulnerability in login endpoint"
- "Add pagination to user list API"
- "Refactor database connection pooling for performance"
- "Update API documentation for v2 endpoints"

❌ **Bad:**
- "Login problem" (too vague)
- "Fix bug" (no context)
- "TODO" (not actionable)
- "Update docs" (not specific)

#### Step 5: Linking (Antigravity, automatic)

**If discovered during work on another issue:**

```bash
# Link to parent issue
bd dep add bd-a3f8e9 bd-parent123 --type discovered-from

# For blocking dependencies
bd dep add bd-blocked-issue bd-blocker-issue --type blocks

# For epic hierarchies
bd create "User Authentication System" -t epic -p 1
# Returns: bd-parent123
bd create "Implement JWT tokens" -p 1  # Auto-becomes bd-parent123.1
bd create "Add OAuth integration" -p 1  # Auto-becomes bd-parent123.2
```

#### Step 6: Confirmation (Antigravity → Human)

```
✓ Created issue bd-a3f8e9
  Title: Fix TypeScript compilation error in auth middleware
  Type: bug
  Priority: 1
  Labels: backend, typescript, blocker

  View details: bd show bd-a3f8e9
  
Moving on to next task.
```

---

### Special Case: Epic Creation

**When to suggest epics:**
- Feature requires 5+ subtasks
- Work spans multiple areas (frontend, backend, tests, docs)
- Implementation takes multiple sessions
- Clear hierarchical structure exists

**Antigravity proposal:**

```
This feature is large enough for an epic. I recommend:

Epic: "User Authentication System"

Subtasks:
1. Design database schema for users/sessions (task, P1)
2. Implement JWT token generation (feature, P1)
3. Create login/logout API endpoints (feature, P1)
4. Build frontend login UI (feature, P2)
5. Write integration tests (task, P1)
6. Add API documentation (docs, P2)

Dependencies:
- Task 2 depends on Task 1 (schema required)
- Task 3 depends on Task 2 (needs JWT)
- Task 4 depends on Task 3 (needs API)
- Task 5 depends on Task 4 (needs full flow)
- Task 6 can run in parallel

Should I create this epic structure?
```

**After human approval:**

```bash
# Create epic
bd create "User Authentication System" -t epic -p 1
# Returns: bd-abc123

# Create children (auto-numbered)
bd create "Design database schema" -p 1          # bd-abc123.1
bd create "Implement JWT tokens" -p 1            # bd-abc123.2
bd create "Create login/logout endpoints" -p 1   # bd-abc123.3
bd create "Build frontend login UI" -p 2         # bd-abc123.4
bd create "Write integration tests" -p 1         # bd-abc123.5
bd create "Add API documentation" -p 2           # bd-abc123.6

# Add dependencies
bd dep add bd-abc123.2 bd-abc123.1  # Task 2 blocks Task 1
bd dep add bd-abc123.3 bd-abc123.2
bd dep add bd-abc123.4 bd-abc123.3
bd dep add bd-abc123.5 bd-abc123.4

# Visualize
bd dep tree bd-abc123
```

---

### Special Case: External Tool Integration

#### Sentry Error Integration

**Antigravity workflow:**

```bash
# Human request: "Check Sentry for recent errors"

# Using Sentry MCP
# 1. Fetch errors
sentry_errors=$(sentry-cli issues list --status unresolved)

# 2. For each error, propose issue
"Sentry Error #12345: Uncaught TypeError in payment processor

Title: Fix TypeError: Cannot read property 'amount' of undefined
Type: bug
Priority: 0 (affects 45 users)

Description:
Sentry Issue: https://sentry.io/issues/12345
Error: TypeError: Cannot read property 'amount' of undefined
Stack trace: [truncated stack trace]

Root Cause (Sentry SEER AI): Null payment object passed when user cancels checkout

Frequency: 45 occurrences in last 24 hours
First seen: 2024-01-15 14:23:00
Last seen: 2024-01-16 09:15:00

Should I create this issue?"
```

**After approval:**

```bash
bd create "Fix TypeError: Cannot read property 'amount' of undefined" \
  -t bug \
  -p 0 \
  -d "Sentry issue #12345. Root cause: Null payment object when user cancels checkout. Affects 45 users." \
  --label production,sentry,payment
```

#### Linting Error Batch Creation

**Antigravity workflow:**

```bash
# Run linter
npm run lint > lint-output.txt

# Found 15 errors

# Propose to human:
"Found 15 linting errors. Options:
1. Create individual issues for each (15 issues)
2. Create one issue for all related errors (1 issue)
3. Fix critical P0/P1 now, file P2 for the rest

What would you like me to do?"
```

**Human chooses option 3:**

```bash
# Antigravity creates issues for P0/P1 only
bd create "Fix ESLint error: no-unused-vars in auth.ts" -t bug -p 1
bd create "Fix ESLint error: no-console in logger.ts" -t bug -p 1
bd create "Resolve remaining 13 lint warnings" -t chore -p 2
```

---

### Auto-Creation Rules (No Approval Needed)

**Antigravity CAN auto-create these without human approval:**

1. **P3-P4 issues** (low priority, non-blocking)
2. **Discovered-from subtasks** (during active work)
3. **Test-related tasks** (missing test coverage)
4. **Documentation updates** (outdated API docs)

**Example auto-creation:**

```bash
# During work on bd-abc123, Antigravity discovers missing tests
bd create "Add unit tests for JWT validation" \
  -t task \
  -p 3 \
  -d "Discovered while working on bd-abc123. Current coverage: 45%" \
  --label test,backend

bd dep add bd-new-id bd-abc123 --type discovered-from

# Report to human:
"✓ Auto-created bd-new-id for missing tests (P3). Continuing work on bd-abc123."
```

---

## Work Execution Flow

### Step 1: Session Start (Antigravity)

```bash
# Sync with git
git pull --rebase
bd sync

# Get current state
bd list --status in_progress --assignee antigravity
bd list --status blocked
bd ready --priority 1

# Report to human
"Status Report:

Ready to work on (3 issues):
- [P0] bd-xyz789: Fix production login bug
- [P1] bd-abc123: Implement user dashboard
- [P1] bd-def456: Add API rate limiting

In progress (1 issue):
- [P1] bd-prev111: OAuth integration (started yesterday)

Blocked (2 issues):
- [P2] bd-block222: Requires API key from DevOps
- [P2] bd-block333: Depends on bd-abc123

Recommendation: Complete bd-prev111 (80% done) or start bd-xyz789 (P0 critical)

Which issue would you like me to work on?"
```

### Step 2: Claim Work (Antigravity, after human approval)

```bash
# Human says: "Finish bd-prev111"

bd update bd-prev111 --status in_progress --assignee antigravity

# Load context
bd show bd-prev111
bd dep tree bd-prev111

# Check for blockers
bd blocked | grep bd-prev111
# If blocked, report and switch to bd ready
```

### Step 3: Execute (Antigravity)

**During implementation:**

- Stay strictly within issue scope
- Do NOT bundle unrelated changes
- If discovering new work → Create issue → Link it → Continue

**Example:**

```bash
# While working on bd-prev111, discovers memory leak

# 1. Propose issue
"While implementing OAuth, I discovered a memory leak in session storage.

Title: Fix memory leak in session storage cleanup
Type: bug
Priority: 1 (performance issue)
Dependencies: None (can be fixed independently)

Should I create this issue?"

# Human: "Yes, but P2, not urgent"

# 2. Create and link
bd create "Fix memory leak in session storage cleanup" -t bug -p 2
bd dep add bd-new-leak bd-prev111 --type discovered-from

# 3. Continue original work
"✓ Created bd-new-leak (P2). Continuing work on bd-prev111."
```

### Step 4: Quality Gates (Antigravity, MANDATORY)

**Before requesting closure, MUST run:**

```bash
# Run all quality checks
npm test
npm run lint
npm run build
npm run type-check  # if TypeScript

# Report results
"Quality Gates Report for bd-prev111:

✓ Tests: 142 passed, 0 failed
✓ Linter: 0 errors, 2 warnings (in unrelated files)
✓ Type checker: 0 errors
✓ Build: successful

Changes summary:
- src/auth/oauth.ts: Added OAuth2 flow implementation
- src/auth/token.ts: Added token refresh logic
- tests/auth/oauth.test.ts: Added 15 test cases

Ready for your review. Should I close this issue?"
```

**If quality gates fail:**

```bash
"Quality gate failures for bd-prev111:

✗ Tests: 3 failed
  - test_oauth_callback: Expected 200, got 401
  - test_token_refresh: Timeout after 5s
  - test_session_cleanup: Null pointer exception

Options:
1. Fix these failures now (recommend)
2. Create P0 blocker issues and mark bd-prev111 blocked
3. Investigate root cause before deciding

What would you like me to do?"
```

### Step 5: Request Review (Antigravity → Human)

```bash
# After quality gates pass, Antigravity proposes closure

"bd-prev111 ready for review:

Issue: Implement OAuth integration
Status: Implementation complete
Quality: All tests passing

Changes:
- 3 files modified
- 15 tests added
- 0 linter errors
- Documentation updated

Git diff available. Please review.

Approve closure with: 'OK bd-prev111'
Request changes with: 'REJECT bd-prev111 - <reason>'"
```

### Step 6: Human Review

**Human reviews:**

```bash
# Check code changes
git diff

# Test locally
npm test

# Review issue details
bd show bd-prev111
```

**Human response options:**

1. **Approve:** `"OK bd-prev111"`
2. **Request changes:** `"REJECT bd-prev111 - Add error handling for network timeouts"`
3. **Request more tests:** `"Add tests for edge case X before closing"`

### Step 7: Close (Antigravity, after approval)

```bash
# Only after human approval
bd close bd-prev111 --reason "OAuth integration implemented and tested. All quality gates passed. Reviewed and approved by human in commit abc123."

# Sync to git
bd sync
git add .beads/issues.jsonl
git commit -m "bd-prev111: Implement OAuth integration"
git push
```

---

## Landing the Plane

**THIS IS MANDATORY. NEVER END A SESSION WITHOUT COMPLETING ALL STEPS.**

### Step 1: File Remaining Work (Antigravity)

```bash
# Identify leftover work
- Any TODO comments in code
- Any follow-up tasks discovered
- Any technical debt identified
- Any bugs found but not fixed
- Any improvements noted

# Propose each issue
"Found 3 items that need tracking:

1. TODO in auth.ts: Implement rate limiting
   → Suggest: bd create "Add rate limiting to auth endpoints" -t task -p 2

2. FIXME in payment.ts: Validate currency format
   → Suggest: bd create "Fix currency validation logic" -t bug -p 1

3. Discovered: API documentation outdated
   → Suggest: bd create "Update API docs for OAuth flow" -t docs -p 3

Should I create these issues?"
```

### Step 2: Run Quality Gates (Antigravity)

```bash
# If ANY code changed during session
npm test
npm run lint
npm run build

# If failures, create issues
if [ $? -ne 0 ]; then
    bd create "Fix post-session test failures" -t bug -p 0
    bd create "Resolve linting errors from session" -t chore -p 1
fi
```

### Step 3: Update Issue Statuses (Antigravity)

```bash
# Close completed issues (already done in Step 7 above)

# Update in-progress issues
bd update bd-incomplete123 --status open \
  -d "Session ended at 60% complete. Next: Implement error handling."

# Mark blocked issues
bd update bd-blocked456 --status blocked \
  -d "Waiting for API key from DevOps team."
```

### Step 4: GIT PUSH (Antigravity, NON-NEGOTIABLE)

```bash
# Stage all changes
git add .

# Commit with descriptive message
git commit -m "Session $(date): Completed bd-prev111, created 3 follow-up issues"

# Pull with rebase
git pull --rebase

# Sync Beads state
bd sync

# PUSH TO REMOTE (CRITICAL)
git push

# Verify push succeeded
git status
# MUST show "up to date with origin"

# If not, resolve conflicts and retry
while ! git status | grep -q "up to date with origin"; do
    echo "Push failed. Resolving..."
    git pull --rebase
    git push
done
```

**If git push fails repeatedly:**

```bash
# Create blocker issue
bd create "Git push failed - manual intervention required" \
  -t blocker \
  -p 0 \
  -d "Failed at $(date). Repository may be in inconsistent state. Needs human review."

# Alert human
"⚠️ CRITICAL: Unable to push to remote after 3 attempts. Created blocker bd-emergency123. Please investigate."
```

### Step 5: Clean Up (Antigravity)

```bash
# Clear stashes (if safe)
git stash list
# Only if empty or human approves
git stash clear

# Prune remote branches
git remote prune origin
```

### Step 6: Verification Checklist (Antigravity)

```bash
# Verify ALL of these:
[ ] All changes committed
[ ] Git push successful
[ ] git status shows clean working tree
[ ] git status shows "up to date with origin"
[ ] Beads state synced
[ ] No uncommitted changes
[ ] No unpushed commits
[ ] All in-progress issues updated
[ ] All discovered work filed

# Report to human
"✓ Verification complete. All checks passed."
```

### Step 7: Handoff Report (Antigravity → Human)

```
SESSION COMPLETE - LANDING THE PLANE

✓ Issues Filed: 3 new issues
  - bd-rate123: Add rate limiting to auth endpoints (P2, task)
  - bd-curr456: Fix currency validation logic (P1, bug)
  - bd-docs789: Update API docs for OAuth flow (P3, docs)

✓ Quality Gates: ALL PASSED
  - Tests: 142/142 ✓
  - Linter: 0 errors ✓
  - Build: success ✓

✓ Issue Updates:
  - bd-prev111: CLOSED - OAuth integration complete
  - bd-incomplete123: UPDATED - 60% complete, needs error handling
  - bd-blocked456: BLOCKED - Waiting for API key

✓ Git Push: SUCCESSFUL
  - 4 commits pushed to origin/main
  - Repository up to date
  - Beads state synced

✓ Cleanup: Complete

✓ Verification: All checks passed

NEXT SESSION RECOMMENDATIONS:
1. Continue bd-incomplete123 (P1, 60% done)
2. Start bd-curr456 (P1, high priority bug)
3. Unblock bd-blocked456 (obtain API key from DevOps)

SESSION SUMMARY:
Completed OAuth integration with full test coverage. Discovered and filed 3 follow-up tasks. All changes committed and pushed. Project in good state for next session.
```

---

## Human Review Workflow

### Daily/Weekly Oversight

**Check project health:**

```bash
# View all open issues
bd list --status open

# Check ready work
bd ready

# View statistics
bd stats

# Use terminal UI
bv
```

**Prioritize work:**

```bash
# Elevate critical issues
bd update bd-urgent123 --priority 0

# Deprioritize less important
bd update bd-minor456 --priority 3

# Assign to Antigravity
bd update bd-next789 --assignee antigravity
```

### Code Review Integration

**When Antigravity completes work:**

1. **Review code changes** (normal PR process)
2. **Test locally** if needed
3. **Approve or request changes**

```bash
# If approved
"OK bd-prev111"
# Antigravity will close issue

# If changes needed
"REJECT bd-prev111 - Add error handling for network timeouts"
# Antigravity will reopen and continue
```

### Dependency Management

**Review dependency graphs:**

```bash
# Visualize project structure
bd dep tree --all

# Check for circular dependencies
bd dep cycles

# If cycles found, break them
bd dep remove bd-issue1 bd-issue2
```

### Priority & Scope Management

**Adjust priorities based on business needs:**

```bash
# Critical business requirement
bd update bd-feature123 --priority 0

# Technical debt
bd update bd-refactor456 --priority 3
```

**Break down epics:**

```bash
# Create epic
bd create "Q1 Feature Launch" -t epic -p 0

# Add subtasks
bd create "Phase 1: Backend API" -p 0  # Auto-numbered
bd create "Phase 2: Frontend UI" -p 0
bd create "Phase 3: Testing" -p 0
bd create "Phase 4: Documentation" -p 1
```

---

## Best Practices

### For Antigravity

**DO:**
- ✅ Use `bd ready` to find next work
- ✅ Create issues for ALL discovered work
- ✅ Run quality gates before requesting closure
- ✅ Link related issues with dependencies
- ✅ Provide clear, detailed commit messages
- ✅ Sync with `bd sync` after git operations
- ✅ Report status clearly and concisely
- ✅ Ask questions when uncertain

**DON'T:**
- ❌ Use internal to-do lists
- ❌ Create markdown task files
- ❌ Work on blocked issues
- ❌ Close issues without quality gates
- ❌ End session without pushing
- ❌ Assume priorities or dependencies
- ❌ Make architectural decisions without approval
- ❌ Load entire project backlog into context

### For Humans

**DO:**
- ✅ Review and approve issue closures
- ✅ Set clear priorities based on business needs
- ✅ Provide feedback on proposed issues
- ✅ Verify git push succeeded at session end
- ✅ Use `bd stats` for weekly health checks
- ✅ Break down large epics into manageable tasks
- ✅ Resolve blockers promptly

**DON'T:**
- ❌ Let issues accumulate without review
- ❌ Skip code review before approving closure
- ❌ Bypass quality gates
- ❌ Create parallel task tracking systems
- ❌ Ignore dependency graphs
- ❌ Let blocked issues pile up

### Software Engineering Best Practices

**Atomicity:**
- One issue = one discrete unit of work
- Issues should be completable in 1-2 hours
- Break large features into epics with subtasks

**Traceability:**
- Link issues to commits: `git commit -m "bd-abc123: Fix login bug"`
- Reference issues in PRs
- Use dependency links to show relationships

**Reproducibility:**
- Include steps to reproduce bugs
- Document acceptance criteria for features
- Add relevant error messages and stack traces

**Maintainability:**
- Keep issue descriptions up to date
- Close stale issues promptly
- Compact old closed issues periodically

**Collaboration:**
- Use assignees to indicate ownership
- Add labels for filtering and organization
- Provide clear handoffs between sessions

---

## Troubleshooting

### Common Issues

**"bd: command not found"**
```bash
# Solution: Install Beads
npm install -g @beads/cli
# or
curl -fsSL https://beads.dev/install.sh | sh
```

**"Database locked"**
```bash
# Solution: Use no-daemon mode
bd --no-daemon ready

# Or check for zombie processes
bd daemons list
bd daemons kill <pid>
```

**"Merge conflict in issues.jsonl"**
```bash
# Solution: Beads auto-resolves on sync
git pull
# Manually resolve if needed (JSONL is human-readable)
bd sync
```

**"Issue not found"**
```bash
# Solution: Verify issue ID
bd list | grep <partial-id>
bd show <full-id>
```

**"Git push failed"**
```bash
# Solution: Pull and rebase
git pull --rebase
git push

# If conflicts, resolve manually
git status
# Fix conflicts
git add .
git rebase --continue
git push
```

### Recovery Procedures

**Lost local changes:**
```bash
# Recover from git history
git log .beads/issues.jsonl
git checkout <commit> -- .beads/issues.jsonl
bd sync
```

**Corrupt database:**
```bash
# Rebuild from JSONL
rm .beads/beads.db
bd sync --import
```

**Agent confusion:**
```bash
# Reset instructions
bd onboard
```

**Circular dependencies:**
```bash
# Detect cycles
bd dep cycles

# Break cycle
bd dep remove bd-issue1 bd-issue2
```

---

## Quick Reference

### Essential Commands

```bash
# Session Management
bd sync                                    # Sync with git
bd ready                                   # Show available work
bd ready --priority 1                      # Show high-priority work
bd list --status open                      # List open issues
bd show <id>                               # View issue details

# Issue Creation
bd create "title" -t bug -p 1 -d "description"
bd create "title" -t feature -p 2 --label frontend,ui

# Issue Updates
bd update <id> --status in_progress
bd update <id> --priority 0
bd update <id> --assignee antigravity

# Issue Closure
bd close <id> --reason "Implemented and tested"

# Dependencies
bd dep add <blocked-id> <blocker-id>       # Add dependency
bd dep tree <id>                           # View dependency tree
bd dep cycles                              # Check for circular deps

# Quality & Maintenance
bd stats                                   # Project statistics
bd compact --stats                         # Check database size
bd doctor                                  # Diagnose issues

# Session End (MANDATORY)
git pull --rebase && bd sync && git push
git status  # Verify "up to date with origin"
```

### Priority Guide

- **P0**: Drop everything, fix now (prod broken, security, data loss)
- **P1**: Next available slot (critical features, major bugs)
- **P2**: This week (standard features, improvements)
- **P3**: This month (polish, tech debt, minor bugs)
- **P4**: Backlog (nice-to-have, future ideas)

### Issue Type Guide

- **bug**: Something is broken
- **feature**: New functionality
- **task**: Technical work
- **epic**: Large feature (auto-creates subtasks)
- **chore**: Refactoring, cleanup, infrastructure
- **docs**: Documentation updates
- **investigation**: Unknown solution space

---

## Mental Model

> **Beads is the brain.**
> **Antigravity is the hands.**
> **The human is the mind.**

Violating this model is a correctness error, not a stylistic choice.

---

## Success Criteria

A successful Antigravity session:

1. ✅ All work tracked in Beads
2. ✅ No blocked issues ignored
3. ✅ Quality gates passed before closing
4. ✅ Git repository pushed and up to date
5. ✅ All discovered work filed as new issues
6. ✅ Human review points addressed
7. ✅ Clear handoff for next session

---

## Final Reminders

- **Context is king**: Issues persist, context doesn't. Use Beads as your memory.
- **Human approves everything**: Antigravity proposes, human decides.
- **Quality gates are mandatory**: No shortcuts, ever.
- **Land the plane**: Every session ends with a complete handoff.
- **Git push is non-negotiable**: Work isn't done until it's pushed.
- **When in doubt, ask**: Never guess, always clarify.

Your role is to be a reliable, thorough coding assistant that never loses track of work, always maintains quality, and keeps the human in full control of all decisions.

---

*Last updated: $(date)*  
*Version: 1.0*  
*Maintained by: Antigravity + Human collaboration*