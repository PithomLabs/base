# Agent Instructions for Beads Issue Tracking

## CRITICAL RULES - READ FIRST

1. **NEVER create markdown files for task tracking**
2. **ALL tasks, todos, and work items MUST be Beads issues**
3. **ALWAYS get explicit human approval before:**
   - Creating new issues
   - Adding dependencies
   - Updating issue status
   - Closing issues
4. **ALWAYS run quality gates before requesting issue closure**
5. **ALWAYS land the plane before ending any session**
6. **NEVER make assumptions - always propose and wait for approval**

---

## Session Startup Protocol

### Step 1: Check Issue State
```bash
bd ready
bd list --status in_progress
bd list --status blocked
```

### Step 2: Report to Human
Present findings in this format:
```
Status Report:

Ready to work on (X issues):
- [P0] bd-xxxxx: <title>
- [P1] bd-xxxxx: <title>

In progress (X issues):
- [P1] bd-xxxxx: <title> (started: <date>)

Blocked (X issues):
- [P2] bd-xxxxx: <title> (blocked by: bd-xxxxx)

Recommendation: Start with <issue-id> because <reason>

Which issue would you like me to work on?
```

### Step 3: Wait for Human Decision
Do NOT start work until human explicitly chooses an issue.

---

## Issue Discovery Protocol

### When You Discover New Work

**DO NOT create issues automatically. Follow this exact protocol:**

1. **Identify** the work item
2. **Analyze** and categorize:
   - Type: bug/feature/task/epic/chore/docs
   - Priority: 0-4 (0=critical, 4=backlog)
   - Dependencies: what blocks this or is blocked by this
   - Description: clear context and details

3. **Propose** to human:
```
I discovered new work that should be tracked:

Title: "<clear, action-oriented title>"
Type: <type>
Priority: <number> - <reason for priority>
Dependencies: <if any, explain why>

Description:
<detailed description with context>

Should I create this issue?
```

4. **Wait** for explicit approval
5. **Execute** only after approval:
```bash
bd create "<title>" -t <type> -p <priority> -d "<description>"
```

6. **Add dependencies** if approved:
```bash
bd dep add <blocked-issue-id> <blocker-issue-id>
```

### Priority Guidelines

**P0 (Critical)** - Suggest only for:
- Security vulnerabilities
- Data loss bugs
- Complete system outages
- Production hotfixes

**P1 (High)** - Suggest for:
- Important sprint features
- Significant user-facing bugs
- Blocking dependencies

**P2 (Medium)** - Suggest for:
- Standard features
- Minor bugs
- Improvements

**P3 (Low)** - Suggest for:
- Nice-to-have features
- Technical debt
- Optimizations

**P4 (Backlog)** - Suggest for:
- Future considerations
- Ideas
- Long-term improvements

### Good vs Bad Issue Titles

**Good Examples:**
- ✓ "Fix authentication bypass in login endpoint"
- ✓ "Add pagination to user list API"
- ✓ "Refactor database connection pooling"
- ✓ "Update API documentation for v2 endpoints"

**Bad Examples:**
- ✗ "Login problem" (too vague)
- ✗ "Fix bug" (no context)
- ✗ "Update docs" (not specific)
- ✗ "TODO" (not actionable)

---

## Working on Issues Protocol

### Before Starting Work

1. **Get human approval** for which issue to work on
2. **Request status update approval:**
```
Ready to start work on bd-xxxxx. Should I update its status to in_progress?
```

3. **After approval, update status:**
```bash
bd update bd-xxxxx --status in_progress
```

4. **Load issue context:**
```bash
bd show bd-xxxxx
bd dep tree bd-xxxxx
```

5. **Check for blockers** - if issue has dependencies, verify they're closed

### During Work

- Provide progress updates periodically
- If you discover blocking issues, follow Issue Discovery Protocol
- If solution approach changes significantly, inform human
- Never make architectural decisions without approval
- If you encounter ambiguity, ask questions immediately

### Context Management Rules

**DO:**
- Query issue details with `bd show <id>` when needed
- Load only the current issue's context
- Use `bd ready` to find next available work
- Clear context freely - issues persist in database

**DO NOT:**
- Maintain separate TODO lists in chat
- Create markdown files with task lists
- Load all issues into context at once
- Ask human to remember what's next
- Assume what needs to be done without checking Beads

---

## Quality Gates Protocol

### Before Requesting Issue Closure

**YOU MUST run all applicable quality gates:**

```bash
# Run tests
npm test
# or: pytest / cargo test / go test / etc.

# Run linter
npm run lint
# or: eslint . / flake8 / clippy / golangci-lint

# Run type checker (if applicable)
tsc --noEmit
# or: mypy . / cargo check

# Run formatter check
prettier --check .
# or: black --check . / cargo fmt --check

# Run build
npm run build
# or: cargo build / go build / make
```

### Report Results

```
Quality Gates Report for bd-xxxxx:

✓ Tests: X passed, 0 failed
✓ Linter: 0 errors, Y warnings (in unrelated files)
✓ Type checker: 0 errors
✓ Build: successful

Changes summary:
- <file1>: <what changed>
- <file2>: <what changed>

Ready for your review. Should I close this issue?
```

### If Quality Gates Fail

**DO NOT request closure. Instead:**

```
Quality gate failures for bd-xxxxx:

✗ Tests: 3 failed
  - test_user_login: expected 200, got 401
  - test_password_reset: timeout
  - test_profile_update: null pointer exception

Should I:
1. Fix these failures and continue
2. Create separate issues for the failures
3. Investigate the root cause first

What would you like me to do?
```

---

## Issue Closure Protocol

### Requirements for Closure

Issue can ONLY be closed after ALL of:
1. ✓ All quality gates pass
2. ✓ Human reviews and approves code changes
3. ✓ Changes are committed to git
4. ✓ Human explicitly authorizes closure

### Closing Command

```bash
bd close bd-xxxxx --reason "<detailed reason including what was done and outcome>"
```

### Good Closure Reasons

**Good Examples:**
- ✓ "Fixed authentication bypass by adding JWT validation in middleware. Added 12 test cases covering edge cases. Reviewed and approved."
- ✓ "Implemented pagination with limit/offset parameters. Updated API docs. All tests passing."
- ✓ "Refactored connection pool to use singleton pattern. Reduced memory usage by 40%. Benchmarks included."

**Bad Examples:**
- ✗ "Done" (no details)
- ✗ "Fixed" (no explanation)
- ✗ "Complete" (no context)

---

## Landing the Plane - END OF SESSION PROTOCOL

**THIS IS MANDATORY. NEVER END A SESSION WITHOUT COMPLETING ALL STEPS.**

### Step 1: File Remaining Work

Identify and propose issues for:
- Any TODO comments discovered
- Any follow-up tasks needed
- Any technical debt identified
- Any bugs discovered but not fixed
- Any improvements noted

Propose each issue using Issue Discovery Protocol.

### Step 2: Run Quality Gates (if code changed)

Run all quality gates one final time:
```bash
npm test && npm run lint && npm run build
```

If failures exist, create issues for them.

### Step 3: Update All Issue Statuses

```bash
# Close completed issues (with human approval)
bd close bd-xxxxx --reason "<reason>"

# Update in-progress issues with status notes
bd update bd-yyyyy --status in_progress
# Add comment about current state

# Mark blocked issues if dependencies discovered
bd update bd-zzzzz --status blocked
```

### Step 4: GIT PUSH (MANDATORY - NON-NEGOTIABLE)

**This is the most critical step. Work is NOT complete until pushed.**

```bash
# Stage all changes
git add .

# Commit with descriptive message
git commit -m "Session work: <summary of what was accomplished>"

# Pull with rebase to get latest changes
git pull --rebase

# Sync Beads state
bd sync

# PUSH TO REMOTE (this is mandatory)
git push

# Verify push succeeded
git status
```

**If `git status` does NOT show "up to date with origin", the push FAILED. Resolve and retry.**

### Step 5: Clean Up

```bash
# Clear stashes (if safe to do so)
git stash list
# Only clear if you're certain nothing important is stashed
git stash clear
```

### Step 6: Verification Checklist

Verify ALL of these:
- [ ] All changes committed
- [ ] Git push successful
- [ ] `git status` shows clean working tree
- [ ] `git status` shows "up to date with origin"
- [ ] Beads state synced
- [ ] No uncommitted changes
- [ ] No unpushed commits

### Step 7: Handoff Report

Provide comprehensive handoff in this format:

```
SESSION COMPLETE - LANDING THE PLANE

✓ Issues Filed: X new issues
  - bd-xxxxx: <title>
  - bd-yyyyy: <title>

✓ Quality Gates: <status>

✓ Issue Updates:
  - bd-aaaaa: CLOSED - <reason>
  - bd-bbbbb: UPDATED - <status note>
  - bd-ccccc: BLOCKED - <blocker>

✓ Git Push: SUCCESSFUL
  - X commits pushed to origin/<branch>
  - Repository up to date
  - Beads state synced

✓ Cleanup: Complete

✓ Verification: All checks passed

NEXT SESSION RECOMMENDATIONS:
1. <recommendation 1 with issue ID>
2. <recommendation 2 with issue ID>
3. <recommendation 3 with issue ID>

SESSION SUMMARY:
<1-2 sentence summary of what was accomplished, what's in progress, and project state>
```

---

## Dependency Management

### Dependency Types

- **blocks**: Task B must complete before Task A can start
- **related**: Soft connection, doesn't block progress
- **parent-child**: Epic/subtask hierarchical relationship
- **discovered-from**: Auto-created when you discover related work

### Adding Dependencies

**Always propose dependencies, never add automatically:**

```
I noticed that issue bd-yyyyy depends on bd-xxxxx because <reason>.

Should I add this dependency?
Command: bd dep add bd-yyyyy bd-xxxxx
```

After approval:
```bash
bd dep add <blocked-issue-id> <blocker-issue-id>
```

### Checking Dependencies

```bash
# View dependency tree
bd dep tree bd-xxxxx

# Check for circular dependencies (bad!)
bd dep cycles
```

If circular dependencies detected, report immediately:
```
⚠️ WARNING: Circular dependency detected!
<output from bd dep cycles>

This needs to be resolved. Recommended actions:
1. <suggestion 1>
2. <suggestion 2>

How would you like to proceed?
```

---

## Epic Management (Large Features)

### When to Suggest an Epic

Suggest creating an epic when:
- Feature requires 5+ subtasks
- Work spans multiple areas (frontend, backend, docs, tests)
- Feature takes more than one session
- Work has clear hierarchical structure

### Epic Creation Proposal

```
This feature is large enough for an epic. I recommend:

Epic: "<epic title>"

Subtasks:
1. <subtask 1> - <brief description>
2. <subtask 2> - <brief description>
3. <subtask 3> - <brief description>
4. <subtask 4> - <brief description>
5. <subtask 5> - <brief description>

Dependencies:
- Task 2 depends on Task 1 because <reason>
- Task 3 depends on Task 2 because <reason>
- Tasks 4 and 5 can be done in parallel

Should I create this epic structure?
```

### After Approval

```bash
# Create epic (generates parent hash ID)
bd create "<epic title>" -t epic -p <priority>
# Returns: bd-a3f8e9

# Create child tasks (automatically numbered .1, .2, .3...)
bd create "<subtask 1>" -p <priority>  # bd-a3f8e9.1
bd create "<subtask 2>" -p <priority>  # bd-a3f8e9.2
bd create "<subtask 3>" -p <priority>  # bd-a3f8e9.3

# Add dependencies
bd dep add bd-a3f8e9.2 bd-a3f8e9.1
bd dep add bd-a3f8e9.3 bd-a3f8e9.2

# Show hierarchy
bd dep tree bd-a3f8e9
```

---

## External Tool Integration

### Sentry Integration (Bug Tracking)

If Sentry MCP is available and human requests it:

```bash
# Fetch recent errors
# Use Sentry MCP to get error details

# For each error, propose issue creation:
```

Propose format:
```
Sentry Error #<number>: <error title>

Title: "Fix: <error summary>"
Type: bug
Priority: <0-2 based on frequency/severity>

Description:
Sentry Issue: <URL>
Error: <error message>
Stack trace: <relevant portion>

Root Cause (Sentry SEER AI): <analysis>

Affected: X users in last Y days
First seen: <date>
Last seen: <date>

Should I create this issue?
```

### Linting/Testing Batch Issues

When running linters/tests that produce multiple errors:

```
I found X failing tests / linting errors.

Options:
1. Create individual issues for each (X issues)
2. Create one issue for all related errors
3. Fix the most critical ones now and file issues for the rest

What would you like me to do?
```

If human chooses option 1, propose each issue individually.

---

## Error Handling

### When Commands Fail

If any `bd` command fails:

```
⚠️ Beads command failed:
Command: <command>
Error: <error message>

Troubleshooting:
1. <suggestion 1>
2. <suggestion 2>

Should I:
A) Retry the command
B) Try alternative approach
C) Skip this step and continue
D) Stop and investigate

What would you like me to do?
```

### Common Errors and Solutions

**"bd: command not found"**
- Check if Beads is installed: `which bd`
- Installation needed: suggest installation method

**"Database locked"**
- Check for zombie processes: `bd daemons list`
- Suggest: `bd --no-daemon <command>`

**"Merge conflict in issues.jsonl"**
- This is normal after git pull
- Beads will auto-resolve on next `bd sync`
- Just continue normally

**"Issue not found"**
- Double-check issue ID
- Run `bd list` to verify ID exists
- May have been closed/deleted

---

## Team Collaboration

### Before Claiming Work

Always check if issue is assigned:
```bash
bd show bd-xxxxx
```

If assigned to someone else:
```
⚠️ Issue bd-xxxxx is assigned to <person>.

Should I:
1. Find a different issue to work on
2. Check with <person> if I can take over
3. Work on it anyway (not recommended)

What would you like me to do?
```

### Syncing with Team Changes

After `git pull`:
```bash
bd sync
```

Report any new issues or changes:
```
After syncing with remote, I found:
- X new issues created by <person>
- Y issues updated
- Z issues closed

Notable changes:
- bd-xxxxx: Now marked as blocked by bd-yyyyy
- bd-zzzzz: Priority changed from P2 to P0

Would you like me to review any of these?
```

---

## Troubleshooting Protocol

### When Uncertain

**NEVER guess. Always ask.**

If you're uncertain about:
- Whether to create an issue
- What priority to assign
- Whether something is a bug or expected behavior
- How to implement a solution
- Whether to add a dependency

Always present options and ask:
```
I'm uncertain about <situation>.

Options:
1. <option 1> - <pros/cons>
2. <option 2> - <pros/cons>
3. <option 3> - <pros/cons>

What would you like me to do?
```

### When Blocked

If you encounter a blocker:

```
⚠️ BLOCKED on bd-xxxxx

Blocker: <clear description of what's blocking progress>

This issue depends on:
- <dependency 1>: <status> - <why it blocks>
- <dependency 2>: <status> - <why it blocks>

Recommendations:
1. Switch to another issue from `bd ready`
2. Work around the blocker by <workaround>
3. Mark issue as blocked and investigate dependencies

What would you like me to do?
```

---

## Database Maintenance

### When to Check Database Size

Check database size if:
- Project is several months old
- Hundreds of closed issues
- Performance seems slow

```bash
bd compact --stats
```

Report findings:
```
Database Statistics:
- Total issues: X
- Closed issues: Y (Z% of total)
- Database size: N MB
- Old closed issues (30+ days): M

<if size > 10MB and many old issues>
Recommendation: Consider compacting database to improve performance.
Should I analyze compaction candidates?
```

### Compaction Protocol

**NEVER compact without explicit human approval.**

```bash
# Only after human approval
bd compact --analyze --json
```

Present analysis and wait for approval before applying.

---

## Quick Command Reference

### Most Common Commands

```bash
# Status checks
bd ready                          # Find available work
bd list                           # List all issues
bd list --status open             # Filter by status
bd show <id>                      # View issue details

# Issue lifecycle
bd create "<title>" -t <type> -p <priority>
bd update <id> --status in_progress
bd close <id> --reason "<reason>"

# Dependencies
bd dep add <blocked-id> <blocker-id>
bd dep tree <id>

# Sync
bd sync                           # Sync with git

# Utilities
bd stats                          # Project statistics
bd dep cycles                     # Check for circular deps
```

---

## Anti-Patterns - NEVER DO THESE

1. ✗ Creating markdown files for task tracking
2. ✗ Maintaining TODO lists in chat
3. ✗ Creating issues without human approval
4. ✗ Closing issues without running quality gates
5. ✗ Closing issues without human approval
6. ✗ Ending session without landing the plane
7. ✗ Pushing to git without verifying success
8. ✗ Loading all issues into context
9. ✗ Assuming priorities or dependencies
10. ✗ Making architectural decisions without approval

---

## Success Patterns - ALWAYS DO THESE

1. ✓ Get explicit approval for all issue operations
2. ✓ Run quality gates before requesting closure
3. ✓ Land the plane at end of every session
4. ✓ Provide detailed closure reasons
5. ✓ Report status clearly and concisely
6. ✓ Load only current issue context
7. ✓ Propose dependencies with explanations
8. ✓ Verify git push success
9. ✓ Ask questions when uncertain
10. ✓ Keep human in control at all times

---

## Final Reminders

- **Context is king**: Issues persist, context doesn't. Use Beads as your memory.
- **Human approves everything**: You propose, they decide.
- **Quality gates are mandatory**: No shortcuts, ever.
- **Land the plane**: Every session ends with a complete handoff.
- **Git push is non-negotiable**: Work isn't done until it's pushed.
- **When in doubt, ask**: Never guess, always clarify.

Your role is to be a reliable, thorough coding assistant that never loses track of work, always maintains quality, and keeps the human in full control of all decisions.