# Beads Integration - Deployment Guide

## Quick Start: Activate the Fixes

The bugs are **already fix** in code! Just need to restart and refresh:

### 1. Restart Backend Server

```bash
cd /home/chaschel/Documents/ibm/go/base/bin/memos

# Stop old server (Ctrl+C in terminal)

# Start with new binary
./memos --mode dev --driver sqlite --data ./data
```

**Expected output:**
```
Version 0.25.0 has been started on port 8081
```

### 2. Hard Refresh Browser

- Press **Ctrl + Shift + R** (Windows/Linux)
- Or **Cmd + Shift + R** (Mac)

This forces browser to reload the updated JavaScript.

### 3. Test Ticket Creation

1. Navigate to http://localhost:8081
2. Sign in (ibm2100 / iBm1234)
3. Click "Tickets" → "+ New Ticket"
4. Fill in form:
   - Title: "Test Beads Integration"
   - Type: TASK
   - Priority: MEDIUM
   - Labels: `backend, testing` (comma-separated)
5. Try "Add description" → Create memo
6. Click "Create Ticket"

**Expected:** ✅ No errors, ticket created successfully

### 4. Verify Beads CLI Integration

```bash
# Check beads tracked the issue
bd list

# Should show new issue like:
# base-abc [P2] [task] open [backend testing] - Test Beads Integration

# View details
bd show base-abc
```

**Expected:** Issue appears in beads with all metadata

### 5. Check Database

```bash
sqlite3 bin/memos/data/memos_dev.db \
  "SELECT id, title, beads_id, labels FROM tickets ORDER BY id DESC LIMIT 1;"
```

**Expected:**
```
13|Test Beads Integration|base-abc|["backend","testing"]
```

---

## What Was Fixed

### Fix 1: Issue Type Case Conversion ✅
- **File:** [`server/router/api/v1/ticket_service.go`](file:///home/chaschel/Documents/ibm/go/base/server/router/api/v1/ticket_service.go#L123)
- **What:** Converts `TASK` → `task`, `STORY` → `feature` before calling bd CLI
- **Why:** Beads expects lowercase types

### Fix 2: MemoEditor Null Safety ✅  
- **File:** [`web/src/components/MemoEditor/index.tsx`](file:///home/chaschel/Documents/ibm/go/base/web/src/components/MemoEditor/index.tsx#L120)
- **What:** Added null check before accessing `userSetting.memoVisibility`
- **Why:** Prevents crash when user settings haven't loaded yet

### Fix 3: Labels & Tags Fields ✅
- **File:** [`web/src/pages/Tickets.tsx`](file:///home/chaschel/Documents/ibm/go/base/web/src/pages/Tickets.tsx#L541-L559)
- **What:** Added Labels and Tags input fields with comma-separated parsing
- **Why:** Allow users to add beads metadata to tickets

---

## Troubleshooting

### Error: "invalid issue type: TASK"
**Solution:** Server restart not applied. Stop server and restart with new binary.

### Error: "Cannot read properties of undefined"
**Solution:** Browser cache not cleared. Hard refresh with Ctrl + Shift + R.

### Labels field not showing
**Solution:** Clear browser cache or try incognito mode.

### Beads issue not created
**Check:**
```bash
# Verify bd CLI is installed and initialized
bd --version
bd ready

# Check if server logs show bd CLI being called
# Look for log: "bd create successful, beadsID=base-xxx"
```

---

## Build Commands (For Reference)

```bash
# Rebuild backend if needed
cd /home/chaschel/Documents/ibm/go/base
go build -o bin/memos/memos ./bin/memos/main.go

# Rebuild frontend if needed  
cd web
npm run build

# Both builds completed successfully already!
```

---

## Beads Issues Closed

Following AGENTS.MD workflow, all work was tracked in beads:

- ✅ `base-33d` - Fix issue type case conversion
- ✅ `base-wyf` - Fix MemoEditor null reference  
- ✅ `base-xxx` - Document deployment (this file)

View history:
```bash
bd list --closed
```
