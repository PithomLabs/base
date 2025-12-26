# Memo-Based Ticket Description Implementation Plan

## Goal

Replace the simple textarea description field in the Ticket form with the existing Memos markdown editor. When a ticket is created or updated, the description content is:

1. Saved as a new Memo record (using existing gRPC API)
2. The Memo URL/name is stored in the ticket's `description` field (instead of raw content)

## User-Facing Flow

1. User opens ticket create/edit dialog
2. User types description using the full Memos markdown editor (same as composing a memo)
3. On submit:
   - A new Memo is created with the description content
   - The Memo URL (e.g., `/memo/abc123`) is saved to the ticket's `description` field
4. When viewing/editing a ticket, the system fetches the memo by its URL and displays its content in the editor

## Proposed Changes

---

### Frontend

#### [MODIFY] [Tickets.tsx](file:///home/chaschel/Documents/ibm/go/memos-0244/web/src/pages/Tickets.tsx)

1. **Import MemoEditor**: Add import for the `MemoEditor` component and memo store
2. **Replace description textarea** with `MemoEditor` component
3. **Modify `handleCreateOrUpdate`**:
   - Before creating/updating ticket, call `memoStore.createMemo()` with the editor content
   - Extract the memo name (e.g., `memos/abc123`) from the created memo
   - Convert to URL format (e.g., `/m/abc123`) and save as `description`
4. **When editing**: Fetch the memo content by parsing the URL in `description` and populate the editor

#### Key Implementation Details

```typescript
// 1. Create memo with description content
const memo = await memoStore.createMemo({
  content: editorContent,
  visibility: Visibility.PRIVATE,
});

// 2. Save memo URL as description
const memoUrl = `/m/${memo.uid}`;
await fetch("/api/v1/tickets", {
  body: JSON.stringify({ title, description: memoUrl, status, priority }),
});
```

---

### Backend

No backend changes required. The existing ticket API already accepts and stores strings in the `description` field.

---

## Verification Plan

### Manual Testing

1. **Start servers**:
   ```bash
   # Backend
   ./memos --mode dev --driver sqlite --data ./data
   
   # Frontend
   cd web && npm run dev
   ```

2. **Test ticket creation with memo**:
   - Navigate to `http://localhost:3001/tickets`
   - Click "New Ticket"
   - Enter a title
   - Type markdown content in the editor (e.g., `# Heading\n- bullet point`)
   - Click Create
   - Verify ticket is created
   - Check backend logs for memo creation
   - Verify `description` field contains a memo URL (e.g., `/m/xyz`)

3. **Test ticket editing**:
   - Click on an existing ticket to edit
   - Verify the editor loads the memo content
   - Modify content and save
   - Verify a new memo is created with updated content

4. **Verify memo exists**:
   - Navigate to the memo URL directly (e.g., `http://localhost:3001/m/xyz`)
   - Confirm the memo content matches what was entered

---

## Implementation Decisions

Based on user feedback:

1. **When editing a ticket**: ✅ **Option B** - Update the existing memo (overwrite)
2. **Visibility**: ✅ **PUBLIC** 
3. **View Full Description link**: ✅ **Yes** - Added to table
