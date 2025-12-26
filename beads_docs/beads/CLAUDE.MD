Should I check Sentry for recent errors?"
```

**After Approval:**
```bash
# Using Sentry MCP
# Antigravity fetches errors and creates issues:

bd create "Fix: structuredClone not defined in Chrome 79" \
  -t bug \
  -p 1 \
  -d "Sentry issue: https://sentry.io/issues/12345
  
Root cause (Sentry SEER AI): structuredClone lacks polyfill for Chrome 79, 
which is 19 versions too old for native support.

Solution: Add polyfill for older browsers or update browser requirements.

Affected users: 23 in last 7 days"
```

**Linting/Testing Integration:**

```bash
# Antigravity can batch-create issues from test failures

# Run tests
npm test 2>&1 | tee test-results.txt

# Antigravity proposes:
"I found 5 failing tests. Should I create issues for each?"

# After approval:
bd create "Fix: UserService.login() throws null pointer" -t bug -p 1
bd create "Fix: PaymentProcessor.charge() returns wrong status" -t bug -p 0
# ... etc
```

### 4.3 Dependency Management

**Dependency Types:**

1. **blocks**: Task B must complete before Task A
2. **related**: Soft connection, doesn't block
3. **parent-child**: Epic/subtask relationship
4. **discovered-from**: Auto-created when AI discovers related work

**Example:**
```bash
# API must be built before frontend
bd dep add bd-frontend bd-api --type blocks

# Documentation related to feature
bd dep add bd-docs bd-feature --type related

# Bug discovered during feature work
bd dep add bd-bug bd-feature --type discovered-from
```

**Antigravity Should:**
- Propose dependencies when discovered
- Explain why dependency exists
- Never assume dependencies without human approval

### 4.4 Priority Management

**Priority Levels (0-4):**
- **P0**: Critical/Hotfix - stops production
- **P1**: High - important for current sprint
- **P2**: Medium - planned work
- **P3**: Low - nice to have
- **P4**: Backlog - future consideration

**Antigravity Guidelines:**
```
Suggest P0 only for:
- Security vulnerabilities
- Data loss bugs
- Complete system outages

Suggest P1 for:
- Important features in current sprint
- Significant bugs affecting users
- Blocking dependencies

Suggest P2 for:
- Standard features
- Minor bugs
- Improvements

Always explain priority rationale and let human override
```

---

## 5. Best Practices

### 5.1 Issue Writing Guidelines

**Good Issue Title:**
- Clear and concise
- Action-oriented
- Specific

```
✓ "Fix authentication bypass in login endpoint"
✓ "Add pagination to user list API"
✗ "Login problem"
✗ "Fix bug"
```

**Good Issue Description:**
```markdown
# What
Brief description of the issue/feature

# Why
Context and motivation

# How (optional for bugs, required for features)
Proposed approach

# Acceptance Criteria
- [ ] Criterion 1
- [ ] Criterion 2

# Additional Context
- Related issues: bd-123
- External links: https://...
- Error logs/screenshots if applicable
```

### 5.2 Context Management

**Antigravity Should:**
- Only load issues currently being worked on
- Use `bd show <id>` to get issue details
- Never maintain separate task lists
- Clear context freely - issues persist
- Use `bd ready` to find next work

**Anti-patterns:**
```
✗ Maintaining TODO lists in chat
✗ Asking human to remember what's next
✗ Loading all issues into context
✗ Using markdown files for tracking
```

### 5.3 Session Continuity

**Between Sessions:**
```
Session 1 ends:
  → Antigravity lands the plane (pushes everything)
  → Issues bd-1, bd-2 marked in_progress

Session 2 starts (days later):
  → Antigravity runs: bd list --status in_progress
  → Reports: "bd-1 and bd-2 were in progress"
  → Asks: "Should I continue bd-1?"
  → No context lost!
```

### 5.4 Team Collaboration

**With Multiple Developers:**

```bash
# Developer A creates issue
bd create "Add search feature" -p 1

# Pushes to git
git push

# Developer B pulls
git pull
bd sync  # Automatically syncs issues

# Developer B sees the issue
bd list
# → bd-abc123: Add search feature [P1] (open)

# Developer B can claim it
bd update bd-abc123 --assignee "developer-b" --status in_progress
```

**Antigravity in Team Context:**
- Check assignees before claiming issues
- Report if issue already assigned
- Suggest unassigned issues from `bd ready`

---

## 6. Troubleshooting

### 6.1 Common Issues

**Problem: "bd: command not found"**
```bash
# Solution: Install beads
npm install -g @beads/cli
# or check PATH
echo $PATH
```

**Problem: "Database locked"**
```bash
# Solution: Check for zombie processes
bd daemons list
bd daemons kill <pid>

# Or force direct mode
bd --no-daemon list
```

**Problem: "Merge conflict in issues.jsonl"**
```bash
# Solution: Beads handles this automatically
git pull  # Conflict appears
# Edit .beads/issues.jsonl to resolve
git add .beads/issues.jsonl
git commit
bd sync  # Beads re-imports
```

**Problem: "Issues not syncing across sessions"**
```bash
# Check daemon status
bd info | grep daemon

# Manually sync
bd sync

# Check git status
git status .beads/issues.jsonl
```

### 6.2 Database Maintenance

**Regular Maintenance:**
```bash
# Check database size
bd compact --stats

# Compact old closed issues (30+ days)
bd compact --analyze

# View detailed statistics
bd stats --json
```

**When Database Grows Large (>10MB):**
```bash
# Preview compaction
bd compact --analyze --json --no-daemon

# Apply compaction (with human approval)
bd compact --apply --id bd-42 --summary summary.txt
```

---

## 7. Configuration Files

### 7.1 .agents.md (Antigravity Instructions)

```markdown
# Beads Issue Tracking

## Core Rules
1. NEVER use markdown files for task tracking
2. ALL tasks must be Beads issues
3. ALWAYS get human approval before creating/closing issues
4. ALWAYS run quality gates before requesting issue closure
5. ALWAYS land the plane before ending session

## Commands
- Find work: `bd ready`
- Show issue: `bd show <id>`
- Update status: `bd update <id> --status <status>`
- Create issue: `bd create "<title>" -t <type> -p <priority>`
- Close issue: `bd close <id> --reason "<reason>"`

## Workflow
1. Start: Check `bd ready` and `bd list --status in_progress`
2. Work: Update to in_progress with approval
3. Complete: Run tests, get approval, close with reason
4. End: Land the plane (commit, push, sync)

## Context Management
- Load only current issue details
- Never maintain separate task lists
- Issues persist across session restarts
- Use `bd show <id>` for issue context
```

### 7.2 Project-Specific Configuration

Create `.beads/config.json`:
```json
{
  "prefix": "myapp",
  "priority_levels": {
    "0": "critical",
    "1": "high",
    "2": "medium",
    "3": "low",
    "4": "backlog"
  },
  "required_fields": ["description", "priority"],
  "auto_close_parent": false,
  "quality_gates": {
    "test": "npm test",
    "lint": "npm run lint",
    "build": "npm run build"
  }
}
```

---

## 8. Integration Examples

### 8.1 With Git Workflow

```bash
# Feature branch workflow
git checkout -b feature/user-auth

# Create feature epic
bd create "User Authentication System" -t epic -p 1

# Work on subtasks
bd create "Implement JWT tokens" -p 1
bd update bd-abc --status in_progress

# Commit incrementally
git add .
git commit -m "feat: add JWT token generation"

# When feature complete
bd close bd-abc --reason "JWT implementation complete with tests"
bd sync
git push origin feature/user-auth

# Create PR
# Merge to main
# bd issues automatically sync to main branch
```

### 8.2 With CI/CD

```yaml
# .github/workflows/beads-check.yml
name: Beads Check

on: [pull_request]

jobs:
  check-issues:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Install Beads
        run: npm install -g @beads/cli
      
      - name: Check for blockers
        run: |
          bd --no-daemon list --status in_progress
          bd --no-daemon dep cycles
      
      - name: Verify sync
        run: |
          bd --no-daemon sync
          git diff --exit-code .beads/issues.jsonl
```

### 8.3 With Project Management Tools

```bash
# Export to Jira (if configured)
bd export --format jira > issues.json

# Generate reports
bd stats --json > weekly-report.json

# Visualize in terminal
bd # Opens TUI
# Press 'b' for kanban board
# Press 'g' for dependency graph
```

---

## 9. Quick Reference Card

### Human Commands

```bash
# Status checks
bd ready                     # Show available work
bd list                      # List all issues
bd list --status open        # Filter by status
bd stats                     # Project statistics

# Issue management
bd show <id>                 # View details
bd update <id> --priority 0  # Change priority
bd update <id> --assignee me # Assign to self

# Visualization
bd                           # Terminal UI
bd dep tree <id>             # Dependency tree
bd dep cycles                # Find circular deps

# Maintenance
bd sync                      # Manual sync
bd compact --stats           # Check DB size
bd doctor                    # Diagnose issues
```

### Antigravity Commands

```bash
# Session start
bd ready
bd list --status in_progress

# During work
bd show <id>
bd update <id> --status in_progress
bd create "<title>" -t <type> -p <priority>
bd dep add <blocked-id> <blocker-id>

# Session end (landing the plane)
bd sync
git add . && git commit && git pull --rebase && git push
bd list --status open
```

---

## 10. Success Metrics

### Tracking Effectiveness

**Metrics to Monitor:**
```bash
# Issue velocity
bd stats --json | jq '.closed_last_7_days'

# Cycle time (creation to closure)
bd list --status closed --json | jq '.[] | .cycle_time'

# Blocked vs ready ratio
bd blocked --json | jq 'length'
bd ready --json | jq 'length'

# Context usage (how often issues are queried)
# Lower is better - means context engineering is working
```

**Good Indicators:**
- Issues rarely get lost
- Sessions resume smoothly
- Minimal "what was I doing?" moments
- Clear prioritization
- Smooth team handoffs

**Red Flags:**
- Duplicate issues
- Issues with no descriptions
- Many stale in_progress issues
- Circular dependencies
- Large uncommitted changes

---

## 11. Appendix: Command Reference

### Complete BD CLI Commands

```bash
# Initialization
bd init [--prefix <name>] [--contributor] [--team]
bd onboard
bd doctor

# Issue CRUD
bd create "<title>" [-t type] [-p priority] [-d description]
bd show <id>
bd update <id> [--status status] [--priority P] [--assignee name]
bd close <id> [--reason "<reason>"]
bd delete <id>

# Listing & Filtering
bd list [--status status] [--priority P] [--assignee name]
bd ready [--priority P]
bd blocked

# Dependencies
bd dep add <blocked-id> <blocker-id> [--type type]
bd dep remove <issue-id> <dep-id>
bd dep tree <id>
bd dep cycles

# Synchronization
bd sync
bd export [--format format]
bd import <file>

# Database Management
bd info [--schema] [--json]
bd migrate [--inspect] [--dry-run] [--cleanup]
bd compact [--stats] [--analyze] [--apply]

# Daemon Management
bd daemons list
bd daemons kill <pid>

# Utilities
bd stats [--json]
bd quickstart
bd --version
bd --help
```

### Status Values

- `open` - Not started
- `in_progress` - Currently working
- `blocked` - Waiting on dependency
- `review` - Code review pending
- `closed` - Completed
- `cancelled` - Won't do

### Type Values

- `bug` - Defect or error
- `feature` - New functionality
- `task` - General work item
- `epic` - Large feature with subtasks
- `chore` - Maintenance work
- `docs` - Documentation

---

## Summary

This guide establishes a robust workflow where:

1. **Beads manages all task state** - no markdown files, no lost context
2. **Antigravity proposes, human approves** - maintains control
3. **Issues persist across sessions** - seamless continuity
4. **Dependencies are explicit** - clear work ordering
5. **Quality gates are mandatory** - maintains code quality
6. **Git integration is automatic** - team collaboration works
7. **Landing the plane is enforced** - no abandoned work

The key insight: Beads provides the persistent memory layer that allows AI agents to work effectively across multiple sessions while keeping humans in control of all critical decisions.Should I create this structure?"
```

**After Approval:**
```bash
# Create epic
bd create "User Authentication System" -t epic -p 1
# Returns: bd-a3f8e9

# Create child tasks (automatically numbered .1, .2, .3...)
bd create "Design database schema" -p 1
bd create "Implement backend API" -p 1
bd create "Create frontend UI" -p 1
bd create "Write integration tests" -p 1
bd create "Add documentation" -p 2

# Add dependencies
bd dep add bd-a3f8e9.2 bd-a3f8e9.1  # API depends on schema
bd dep add bd-a3f8e9.3 bd-a3f8e9.2  # UI depends on API
bd dep add bd-a3f8e9.4 bd-a3f8e9.3  # Tests depend on UI

# View hierarchy
bd dep tree bd-a3f8e9
```

### 4.2 Integration with External Tools

**Sentry Integration (Bug Tracking):**

```
Antigravity: "I can integrate with Sentry to automatically file bugs. 
Should I check Sentry for recent errors?"Should I create this structure?"Should I create this issue?"
```

**Human Response:**
```bash
# Option 1: Approve
"Yes, create it"

# Option 2: Modify
"Yes, but make it P1 instead of P0, we can fix it in the next sprint"

# Option 3: Reject
"No, that's actually expected behavior per the security audit"
```

**After Approval, Antigravity Executes:**
```bash
bd create "Fix password reset authentication bug" \
  -t bug \
  -p 0 \
  -d "Password reset endpoint bypasses authentication check, allowing unauthenticated requests. Same root cause as bd-a3f2. Found during work on login bug fix."

# If dependencies mentioned:
bd dep add bd-b1c4 bd-x7y9  # bd-x7y9 is the new issue ID
```

### 3.3 Working on Issues

**Antigravity Protocol:**

```
Before starting work:
1. Ask human which issue to work on (or suggest from bd ready)
2. Get explicit approval
3. Update status to in_progress
4. Query issue details: bd show <issue-id>
5. Check dependencies: bd dep tree <issue-id>
6. Load ONLY this issue's context
```

**Example Commands:**
```bash
# Human approves bd-a3f2
bd update bd-a3f2 --status in_progress
bd show bd-a3f2
bd dep tree bd-a3f2

# Output shows:
# → bd-a3f2: Fix critical login bug [P0] (in_progress)
#   No dependencies (ready to work!)
```

**During Work:**
- Provide progress updates
- If discovering blocking issues, propose new issues with dependencies
- If solution changes, propose description updates
- Never make assumptions about architecture

### 3.4 Code Review & Quality Gates

**Before Closing Issues:**

```
Antigravity must run quality gates:
1. Run tests: npm test / pytest / cargo test
2. Run linters: eslint / flake8 / clippy
3. Run type checker: tsc / mypy
4. Run formatter check: prettier --check / black --check
5. Build: npm run build / cargo build
```

**Example Dialogue:**
```
Antigravity: "I've completed bd-a3f2 (Fix critical login bug). Running quality gates:

✓ Tests passed (142 passed, 0 failed)
✓ Linter passed (0 errors, 2 warnings in unrelated files)
✓ TypeScript passed (0 errors)
✓ Build successful

Changes made:
- Added authentication check in /src/auth/login.ts
- Updated tests in /tests/auth.test.ts
- Fixed edge case for expired tokens

Ready for your review. Should I close this issue?"
```

**Human Actions:**
```bash
# Review the code
git diff

# Test locally
npm test

# Approve or request changes
"Looks good, close it"
# or
"Please add JSDoc comments to the new auth function first"
```

### 3.5 Closing Issues

**Antigravity Protocol:**
```
Only close issues after:
1. Human reviews and approves code
2. All quality gates pass
3. Changes are committed
4. Human explicitly says to close
```

**Closing Command:**
```bash
bd close bd-a3f2 --reason "Fixed authentication bypass. Added check in login endpoint and comprehensive tests. Reviewed and approved by human."
```

### 3.6 Landing the Plane (Session End)

**CRITICAL: This is mandatory before ending any session**

**Antigravity Protocol:**
```
At session end, you MUST complete ALL these steps:

1. File issues for remaining work
   - Any TODOs discovered
   - Any follow-up tasks
   - Any technical debt identified

2. Run quality gates (if code changed)
   - Tests, linters, builds
   - Report any failures as issues

3. Update all issue statuses
   - Close finished issues
   - Update in_progress issues with status notes
   - Set blocked issues if dependencies discovered

4. GIT PUSH (MANDATORY - this is non-negotiable)
   git add .
   git commit -m "Session work: <summary>"
   git pull --rebase
   bd sync
   git push
   git status  # Must show "up to date with origin"

5. Clean up
   git stash clear (if safe)
   
6. Verify
   - All changes committed AND pushed
   - No uncommitted changes
   - Beads state synced

7. Handoff report
   - What was completed
   - What's in progress
   - What's blocked
   - Suggested next steps
```

**Example Handoff:**
```
Antigravity: "Session complete. Landing the plane:

✓ Issues filed: 2 new issues for follow-up work
  - bd-f8a1: Add rate limiting to login endpoint
  - bd-g2b5: Update API documentation

✓ Quality gates: All passed

✓ Issue updates:
  - bd-a3f2: CLOSED (login bug fixed)
  - bd-d7f3: UPDATED (OAuth 80% complete, blocked on API key)

✓ Git push: SUCCESSFUL
  - 3 commits pushed to origin/main
  - Repository up to date
  - Beads state synced

✓ Cleanup: Stashes cleared

✓ Verification: git status clean

Next session recommendations:
1. Continue bd-d7f3 (need API key from DevOps)
2. Start bd-f8a1 (unblocked, P1 priority)
3. Review bd-g2b5 with technical writer

Session summary: Fixed critical authentication bug, discovered 2 related issues, 
all work committed and pushed. Project in good state for next session."
```

---

## 4. Advanced Workflows

### 4.1 Epics & Hierarchical Issues

**For Large Features:**

```
Antigravity: "This feature is large. I recommend breaking it into an epic:

Epic: 'User Authentication System'
Subtasks:
  1. Design database schema
  2. Implement backend API
  3. Create frontend UI
  4. Write integration tests
  5. Add documentation

Should I create this structure?"Should I create this issue?"Should I create this issue?"Which would you like me to work on? I recommend bd-a3f2 due to P0 priority."
```

### 3.2 Issue Discovery & Creation

**When Antigravity Discovers New Work:**

```
Antigravity Protocol:
1. Identify the new work item
2. Propose issue creation with:
   - Clear, descriptive title
   - Detailed description with context
   - Suggested priority (0-4, 0=highest)
   - Suggested type (bug/feature/task/epic)
   - Suggested dependencies (if any)
3. Wait for human approval
4. Only create issue after explicit approval
```

**Example Dialogue:**
```
Antigravity: "While working on bd-a3f2, I discovered that the password reset 
functionality also has the same authentication issue. 

I recommend creating a new issue:
- Title: 'Fix password reset authentication bug'
- Type: bug
- Priority: 0 (critical security issue)
- Dependency: Should block bd-b1c4 (user profile validation)
- Description: Password reset endpoint bypasses authentication check, allowing 
  unauthenticated requests. Same root cause as bd-a3f2.

Should I create this issue?"Which would you like me to work on? I recommend bd-a3f2 due to P0 priority."# Beads Integration Guide for Antigravity AI Agent

## Overview

This guide establishes a comprehensive workflow for integrating Beads issue tracker with Antigravity (AI coding agent) while maintaining human oversight and following software engineering best practices.

---

## 1. Initial Setup

### 1.1 Installation & Configuration

```bash
# Install Beads
npm install -g @beads/cli
# or
curl -fsSL https://beads.dev/install.sh | sh

# Initialize in project
cd your-project
bd init

# Run doctor to verify setup
bd doctor

# Configure agent integration
bd onboard
```

### 1.2 Project Structure

After initialization, your project will have:
```
your-project/
├── .beads/
│   ├── beads.db           # SQLite database (local only)
│   ├── issues.jsonl       # Git-tracked issue state
│   ├── daemon.sock        # Background sync service
│   └── config.json        # Project configuration
├── .claude.md             # Agent instructions (auto-generated)
└── .agents.md             # Multi-agent instructions
```

---

## 2. Core Principles

### 2.1 Context Engineering Philosophy

1. **Issues Over Markdown**: Never use markdown files for task tracking. All tasks, todos, and work items MUST be Beads issues.
2. **Atomic Context**: Load only the specific issue being worked on, not entire task lists.
3. **Persistent Memory**: Issues survive session restarts, context clears, and agent switches.
4. **Explicit Dependencies**: Use dependency graphs instead of implicit ordering in lists.
5. **Human-in-the-Loop**: All issue creation, status changes, and closures require human approval.

### 2.2 Agent Responsibilities vs Human Responsibilities

**Antigravity (Agent) Should:**
- Propose new issues when discovering work
- Suggest dependencies between issues
- Recommend priority levels
- Query for ready work
- Update issue status with human approval
- Provide detailed descriptions with context

**Human Should:**
- Approve/reject issue creation
- Approve/reject dependency additions
- Approve/reject status changes
- Set final priorities
- Make architectural decisions
- Review code before closing issues

---

## 3. Development Workflow

### 3.1 Starting a New Session

**Human Actions:**
```bash
# Check current state
bd ready                    # See available work
bd list --status open       # Review all open issues
bd stats                    # Get project overview
```

**Antigravity Actions:**
```
When session starts:
1. Run `bd ready` to see unblocked issues
2. Run `bd list --status in_progress` to check ongoing work
3. Report findings to human
4. Ask human which issue to work on OR suggest highest priority ready issue
```

**Example Dialogue:**
```
Antigravity: "I've checked the issue tracker. Here's the current state:

Ready to work on (3 issues):
- [P0] bd-a3f2: Fix critical login bug
- [P1] bd-b1c4: Add user profile validation
- [P2] bd-c5e9: Update documentation

In progress (1 issue):
- [P1] bd-d7f3: Implement OAuth flow (started yesterday)

Which would you like me to work on? I recommend bd-a3f2 due to P0 priority."