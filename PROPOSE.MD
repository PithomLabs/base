# PROPOSE.MD: Basic Jira-like Ticket Management System

This proposal outlines the "bare essential" extension of the Memos ticket system to a production-grade, lightweight Jira API.

## 1. Design Philosophy
-   **Simplicity**: Leverage existing Memos architecture (user system, database patterns).
-   **SQLite-first**: Optimized for the single-file database.
-   **Production-grade**: Proper error handling, validation, and clean UI/UX.
-   **Reuse**: Leverage existing Memo system for comments.

## 2. Data Model Changes (SQLite)

We will modify the existing `tickets` table.

### 2.1 Table: `tickets` (Update)
Add `type` to categorize work.

```sql
ALTER TABLE tickets ADD COLUMN type TEXT NOT NULL DEFAULT 'TASK'; -- Enum: 'TASK', 'BUG', 'STORY'
```

*Note: `tags` can be managed within the description (Memo) or added as a column. For robust filtration, we will add a JSON `tags` column.*

```sql
ALTER TABLE tickets ADD COLUMN tags TEXT NOT NULL DEFAULT '[]'; -- JSON Array
```

### 2.2 Comments (Native Memo Integration)
Since every ticket description is linked to a Memo, we will use the existing **Memo Relations** system for comments.
-   **Comments**: A comment is a regular Memo.
-   **Relation**: Linked via `memo_relation` table with type `COMMENT`.
-   **No new table required**.

## 3. Backend Implementation
Location: `server/router/api/v1/` and `store/`

### 3.1 Ticket Service Enhancements
-   **Update Create/Update API**: Support `type` and `tags`.
-   **Filtering**: Enhance `ListTickets` to filter by:
    -   `status` (Array of status)
    -   `priority`
    -   `type`
    -   `creator_id / assignee_id`
    -   `tags`

*Note: Comment operations will use existing Memo APIs (`CreateMemo` with relation).*

## 4. Frontend Implementation
Location: `web/src/pages/Tickets.tsx` (and new components)

### 4.1 Views
We will introduce a **View Switcher** (Kanban | List/Queue).

#### A. Queue View (Enhanced List)
-   **Filters Bar**: Dropdowns for Type, Status, Priority, Assignee.
-   **Sorting**: Sort by Created, Updated, Priority.
-   **Compact Rows**: Quick actions (Assign, change status).

#### B. Kanban Board
-   **Columns**: Open, In Progress, Closed (Mapped from Status).
-   **Drag & Drop**: Using `@dnd-kit` or similar lightweight library (or standard HTML5 DnD if minimizing deps).
-   **Card Design**: Title, Priority Badge, Assignee Avatar, Type Icon.

### 4.2 Ticket Detail
-   **Slide-over / Modal**: Clicking a ticket opens a detailed view.
-   **Metadata**: Editable fields (Status, Priority, Type, Tags).
-   **Description**: Render existing Memo-linked description (Markdown).
-   **Comments Section**: 
    -   **View**: Reuse `MemoView` component to render linked memos.
    -   **Create**: Reuse `MemoEditor` with `parentMemoName` prop to automatically create `COMMENT` relation.

## 5. Technical Stack Details
-   **Kanban DnD**: Use `@dnd-kit/core` (already in `package.json`) for accessible, lightweight drag-and-drop.
-   **State Management**: Local React state for Board; specific updates trigger `fetchTickets`.
-   **Reusability**: Heavily leveraging `MemoView` and `MemoEditor` components to ensure UI consistency.

## 6. Evaluation & Verification

### 6.1 Success Criteria
-   [ ] User can create 'Bug', 'Task', 'Story'.
-   [ ] User can drag a ticket from 'Open' to 'In Progress' on Kanban board.
-   [ ] User can filter the Queue to show only "High Priority Bugs".
-   [ ] User can comment on a ticket (persisted as linked Memo).
-   [ ] SQLite migration runs successfully without data loss.

### 6.2 detailed Test Plan

#### A. Backend Unit Tests
I will add `store/test/ticket_test.go` mirroring `memo_test.go` to verify:
-   CRUD operations for tickets with `Type` and `Tags`.
-   Filtering logic (Status list, Priority, Type).

#### B. Manual Verification (via Browser Helper)
I will perform the following manual checks once the server is running:
1.  **Frontend Flow**:
    -   Navigate to `/tickets`.
    -   Click "New Ticket" -> Select Type "Bug" -> Properties: High Priority -> Save.
    -   Verify ticket appears in "OPEN" column.
    -   Drag ticket to "IN PROGRESS".
    -   Refresh page -> Verify status persists.
2.  **Comment Integration**:
    -   Click the ticket.
    -   Typing "This is a bug fix" in the editor -> Save.
    -   Verify the comment appears in the list.

#### C. API Verification (Curl)
-   `POST /api/v1/tickets` (Create)
-   `PATCH /api/v1/tickets/:id` (Update Status)
-   `GET /api/v1/tickets?status=OPEN&type=BUG` (Filter)
