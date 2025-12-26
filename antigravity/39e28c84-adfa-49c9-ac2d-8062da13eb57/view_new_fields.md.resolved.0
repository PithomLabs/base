# How to See New Fields in Ticket Form

## ✅ Changes Made

Added to ticket form:
- **Labels** field (Beads tags) - comma-separated input
- **Tags** field (Legacy) - comma-separated input

Location in code: [`web/src/pages/Tickets.tsx` line 541-558](file:///home/chaschel/Documents/ibm/go/base/web/src/pages/Tickets.tsx#L541-L558)

## 🔄 How to See the Changes

### Option 1: Hard Refresh Browser (Fastest)

1. Open your browser with the ticket form
2. Press **Ctrl + Shift + R** (or **Cmd + Shift + R** on Mac)
3. This forces reload from server, bypassing cache

### Option 2: Use Dev Server (Hot Reload)

```bash
cd web
npm run dev
```

Then open: http://localhost:3001

Dev server has hot reload - changes appear automatically as you edit files.

### Option 3: Rebuild Production

```bash
cd web
npm run build
cd ..
go run ./bin/memos/main.go --mode dev --port 8081
```

Then refresh browser.

---

## Expected Form Fields (After Refresh)

You should now see:

```
Title
Type  
Status
Priority
Assignee
Memo URL (Description)
Labels                    ← NEW!
Tags (Legacy)            ← NEW!
[Cancel] [Create Ticket]
```

**Labels example:** `backend, frontend, security`  
**Tags example:** `important, urgent, review-needed`

---

## Verify Changes Were Saved

```bash
grep -A 5 "Labels Field" web/src/pages/Tickets.tsx
```

Should show:
```tsx
{/* Labels Field (NEW - Beads labels) */}
<div>
  <label className="block text-sm font-medium mb-1">Labels</label>
  <Input
    value={labels.join(", ")}
    ...
```

---

## Testing

1. Refresh browser (Ctrl + Shift + R)
2. Click "+ New Ticket" button
3. You should see Labels and Tags fields
4. Enter labels: `backend, security`  
5. Create ticket
6. Check database:
   ```bash
   sqlite3 bin/memos/data/memos_dev.db \
     "SELECT id, title, labels, tags FROM tickets ORDER BY id DESC LIMIT 1;"
   ```

Should show: `["backend","security"]`
