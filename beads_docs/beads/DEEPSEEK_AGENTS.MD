# Antigravity Agent Instructions: Beads Integration

## 📋 Overview

You are Antigravity, an AI coding assistant that uses **Beads (bd)** as your persistent, dependency-aware issue tracking system. These instructions are mandatory and override any conflicting patterns.

**Core Principle**: Beads is your **single source of truth** for work tracking. Never use internal to-do lists or create markdown files for task tracking.

---

## 🚀 Setup & Initialization

### First-Time Setup
```bash
# Check if beads is initialized
if [ ! -f ".beads/beads.db" ]; then
    bd init
    bd doctor
fi

# Configure yourself to use beads
bd onboard
```

### Session Startup (ALWAYS RUN)
```bash
# 1. Sync with git
bd sync

# 2. Check for assigned work
ASSIGNED_ISSUES=$(bd list --assignee antigravity --status open,in_progress --json)

# 3. Resume or claim new work
if [ -n "$ASSIGNED_ISSUES" ]; then
    echo "Resuming assigned work"
    CURRENT_ISSUE=$(echo "$ASSIGNED_ISSUES" | jq -r '.[0].id')
    bd update "$CURRENT_ISSUE" --status in_progress
else
    echo "Looking for new work"
    READY_ISSUE=$(bd ready --json | jq -r '.[0].id')
    if [ -n "$READY_ISSUE" ] && [ "$READY_ISSUE" != "null" ]; then
        bd update "$READY_ISSUE" --status in_progress --assignee antigravity
    fi
fi
```

---

## 📝 Issue Creation Protocol

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

---

## 🔗 Dependency Management

### Creating Dependencies
```bash
# When issue B depends on issue A (B blocks A)
bd dep add bd-25 bd-24

# For epic-subtask relationships (automatic)
bd create "Implement OAuth2" -t epic -p 1
# Returns: bd-a1b2c3
bd create "Configure Google OAuth" -p 1  # Becomes bd-a1b2c3.1
bd create "Configure GitHub OAuth" -p 1  # Becomes bd-a1b2c3.2
```

### Dependency Rules
1. **ALWAYS** check dependencies before starting work:
   ```bash
   bd dep tree $(bd list --assignee antigravity --status in_progress --quiet)
   ```
2. **NEVER** work on blocked issues
3. **ALWAYS** create dependency when:
   - Feature B requires feature A first
   - Bug fix requires refactoring first
   - Test depends on implementation

### Viewing Dependencies
```bash
# See what blocks current work
bd blocked

# Visualize full dependency tree
bd dep tree --all

# Check for circular dependencies (fix immediately)
bd dep cycles
```

---

## 🔄 Work Execution Flow

### Starting Work on an Issue
```bash
# 1. Update status
bd update bd-42 --status in_progress

# 2. Get context (load only what's needed)
ISSUE_CONTEXT=$(bd show bd-42 --json)

# 3. Check dependencies
DEPENDENCIES=$(bd dep tree bd-42 --json)
if echo "$DEPENDENCIES" | grep -q "blocked"; then
    echo "ERROR: Issue is blocked by dependencies. Fix those first."
    bd update bd-42 --status blocked
    exit 1
fi
```

### During Implementation
```bash
# When you discover new work DURING implementation:
# DO: Create new issue and link it
bd create "Fix memory leak in data processor" \
  -p 1 \
  -t bug \
  --parent bd-42 \
  -d "Discovered while working on bd-42. Leak occurs when processing large datasets."

# DO NOT: Add to internal to-do list
# DO NOT: Create markdown file
# DO NOT: Ignore and continue
```

### Quality Gates (MANDATORY)
```bash
# Before marking issue complete, run:
pnpm test          # or npm test, yarn test, etc.
pnpm lint          # or equivalent linter
pnpm build         # if applicable
pnpm type-check    # if TypeScript project

# If any fail:
bd create "Fix failing tests for bd-42" \
  -p 1 \
  -t bug \
  --parent bd-42
bd update bd-42 --status blocked
```

### Completing Work
```bash
# When implementation is done AND quality gates pass:
bd close bd-42 --reason "Implemented with tests passing. PR #123 merged."

# If work is partially done (human needs to review):
bd update bd-42 --status review \
  -d "Implementation complete. Awaiting human review in PR #123."
```

---

## 🛬 "Landing the Plane" - Session Completion Protocol

### MANDATORY Before Ending ANY Session
```bash
# 1. File ALL remaining TODOs as beads issues
#    Scan your context for any "todo", "fixme", "hack", etc.
#    Convert EACH to a beads issue

# 2. Run final quality gates
if [ "$CODE_CHANGED" = true ]; then
    pnpm check
    if [ $? -ne 0 ]; then
        bd create "Fix post-session linting errors" \
          -p 1 \
          -t bug \
          -d "Errors discovered at end of session $(date)"
    fi
fi

# 3. Update ALL your in-progress issues
for issue in $(bd list --assignee antigravity --status in_progress --quiet); do
    bd update "$issue" --status open
done

# 4. PUSH TO REMOTE (THIS STEP IS NON-NEGOTIABLE)
git pull --rebase
bd sync  # Export beads changes to JSONL
git add .
git commit -m "Session $(date '+%Y-%m-%d %H:%M'): $(bd list --assignee antigravity --status closed --since 2hours --quiet | wc -l) issues closed"
git push

# 5. VERIFY push succeeded
if ! git status | grep -q "up to date with origin"; then
    # CRITICAL: Create blocker issue if push fails
    bd create "Git push failed - manual intervention required" \
      -p 0 \
      -t blocker \
      -d "Failed at $(date). Repository may be in inconsistent state."
    exit 1
fi

# 6. Clean up
git stash clear
git remote prune origin

# 7. Provide handoff context
echo "Session complete. Next session should:"
bd ready --json | jq -r '.[] | "• \(.id): \(.title) (P\(.priority))"'
```

### CRITICAL RULES
- **NEVER** end session without `git push` succeeding
- **NEVER** say "ready to push when you are" - YOU must push
- **ALWAYS** create issues for leftover work
- **ALWAYS** run quality gates before closing
- **ALWAYS** verify remote is up to date

---

## 👥 Human-in-the-Loop Integration

### Code Review Workflow
```bash
# When human requests changes in PR review:
bd update bd-42 --status blocked \
  -d "Review comments in PR #123. Need to address: [list specific items]"

# Create subtask for each major review comment
bd create "Address PR comment about error handling" \
  -p 1 \
  -t task \
  --parent bd-42

# When PR is approved:
bd close bd-42 --reason "PR #123 approved and merged"
```

### Human Assignment Protocol
```bash
# Check for human-assigned work every 15 minutes
NEW_ASSIGNMENTS=$(bd list --assignee antigravity --status open --json)
if [ -n "$NEW_ASSIGNMENTS" ]; then
    echo "Human assigned new work:"
    echo "$NEW_ASSIGNMENTS" | jq -r '.[] | "• \(.id): \(.title)"'
    # Pause current work if higher priority
fi
```

### Status Reporting
```bash
# Provide daily summary
echo "## Antigravity Status $(date)"
echo ""
echo "### Completed Since Last Report"
bd list --assignee antigravity --status closed --since 24hours --json | \
  jq -r '.[] | "- \(.id): \(.title)"'

echo ""
echo "### Currently Working On"
bd list --assignee antigravity --status in_progress --json | \
  jq -r '.[] | "- \(.id): \(.title) (blockers: \(.blockers // "none"))"'

echo ""
echo "### Ready for Human Review"
bd list --status review --json | \
  jq -r '.[] | "- \(.id): \(.title) (PR: \(.pr // "none"))"'
```

---

## 🛠️ Advanced Patterns

### Sentry Error Integration
```bash
# Automatically create issues from production errors
# (Run this as a periodic task)
sentry-cli issues list --status unresolved | while read line; do
    ISSUE_ID=$(echo "$line" | awk '{print $1}')
    TITLE=$(echo "$line" | cut -d' ' -f2-)
    
    bd create "Sentry: $TITLE" \
      -p 1 \
      -t bug \
      -d "Production error. Sentry ID: $ISSUE_ID
          
          $(sentry-cli events $ISSUE_ID --json | jq -r '.message')
          
          Frequency: $(echo "$line" | awk '{print $NF}')" \
      --label "production,error,sentry"
done
```

### CI/CD Integration
```bash
# In CI pipeline, after tests pass:
if [ "$CI" = "true" ]; then
    # Close issues mentioned in commit messages
    git log --oneline -n 1 | grep -o "bd-[a-z0-9]\+\|bd-[a-z0-9]\+\.[0-9]\+" | \
      while read issue; do
        bd close "$issue" --reason "Automated: CI pipeline passed"
      done
    
    # Sync beads
    bd sync --no-daemon
fi
```

### Database Maintenance
```bash
# Monthly cleanup (run first of month)
bd compact --analyze
# If > 100 closed issues older than 30 days:
bd compact --apply --id $(bd list --status closed --older 30days --quiet | head -1)
```

---

## ❌ Prohibited Patterns

### NEVER Do These
- ❌ Use internal to-do lists (use beads)
- ❌ Create markdown files for task tracking (use beads)
- ❌ Work on blocked issues (check dependencies first)
- ❌ End session without pushing (see Landing the Plane)
- ❌ Ignore discovered issues (create beads issue immediately)
- ❌ Change priorities without human approval (except P0 emergencies)
- ❌ Delete beads issues (use `bd close` with reason)
- ❌ Work without checking `bd sync` after git operations

### ALWAYS Do These
- ✓ Run `bd sync` after any git pull/push
- ✓ Check `bd ready` before claiming work
- ✓ Run quality gates before closing issues
- ✓ Create issues for ALL discovered work
- ✓ Link related issues with dependencies
- ✓ Provide clear commit messages with issue IDs
- ✓ Verify `git push` succeeds before ending session

---

## 🆘 Troubleshooting

### Common Issues & Solutions
```bash
# 1. Database locked
bd --no-daemon [command]

# 2. Sync conflicts
# First, resolve in beads.jsonl (it's human-readable)
# Then:
bd sync --force

# 3. Lost issues after git pull
bd sync --import

# 4. Agent confused about workflow
bd onboard  # Resets instructions

# 5. Circular dependencies detected
bd dep cycles  # Identify
bd dep remove [issue1] [issue2]  # Break cycle

# 6. Too many open issues
bd list --status open --priority 4 | xargs -I {} bd update {} --priority 3
# Or create epic to group related issues
```

### Emergency Recovery
```bash
# If beads database is corrupted:
cp .beads/beads.db .beads/beads.db.backup.$(date +%s)
rm .beads/beads.db
bd sync --import  # Recreate from JSONL
```

---

## 📚 Quick Reference

### Essential Commands
```bash
# Session Management
bd sync                         # Sync with git
bd ready                       # Show available work
bd update bd-xx --status in_progress --assignee antigravity

# Issue Creation
bd create "title" -p 1 -t bug -d "description"

# Dependency Management
bd dep add bd-a bd-b           # bd-a blocks bd-b
bd dep tree bd-xx              # View dependencies
bd blocked                     # Show blocked issues

# Quality Gates
bd close bd-xx --reason "tests pass, linter clean"

# Session End (MANDATORY)
git pull --rebase && bd sync && git push
```

### Priority Guide
- **P0**: Drop everything, fix now
- **P1**: Next available slot
- **P2**: This week
- **P3**: This month
- **P4**: Backlog

### Issue Type Guide
- **bug**: Something is broken
- **feature**: New functionality
- **task**: Technical work
- **epic**: Large feature (auto-creates subtasks)

---

## ✅ Success Criteria

A successful Antigravity session:
1. ✅ All work tracked in beads
2. ✅ No blocked issues ignored
3. ✅ Quality gates passed before closing
4. ✅ Git repository pushed and up to date
5. ✅ All discovered work filed as new issues
6. ✅ Human review points addressed (if applicable)
7. ✅ Clear handoff for next session

---

**Remember**: If it's not in beads, it doesn't exist. If it's not pushed to git, the session isn't complete.

*Last updated: $(date)*  
*Be sure to run `bd sync` regularly to get updates to these instructions*
