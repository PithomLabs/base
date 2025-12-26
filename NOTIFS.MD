# Design: Notification System with @Mentions

## Overview
The goal is to implement a notification system where users are notified when they are mentioned (via `@nickname`) in:
1.  **Ticket Descriptions** (which are Memos).
2.  **Ticket Comments** (which are related Memos).

## Current System Analysis
-   **Notifications**: Handled via the `inbox` table and `Inbox` struct in `store/inbox.go`.
-   **Triggers**: Currently, only `CreateMemoComment` triggers a notification (to the author of the parent memo).
-   **Activity**: Notifications are linked to an `Activity` (e.g., `ActivityTypeMemoComment`).
-   **Limitations**:
    -   No "Mention" activity type exists.
    -   No "Mention" inbox message type exists.
    -   No Parsing logic for `@nickname` exists in the backend.

## Proposed Design

### 1. Mention Parsing Logic
**Location**: `server/router/api/v1/memo_service.go`
**Functions**: `CreateMemo`, `UpdateMemo` (if we want edits to trigger?), `CreateMemoComment`.

**Logic**:
1.  Extract content from the memo.
2.  Use Regex `/@([^\s@]+)/` (simplified) to find all potential nicknames.
3.  Query `store.ListUsers` with definitions to find matching users.
4.  Filter out:
    -   The creator themselves (no self-notification).
    -   Duplicate mentions.

### 2. Data Model Changes

We need to distinguishable types for "Comment" vs "Mention".

#### Option A: Clean Implementation (Requires Proto Re-generation)
-   **`proto/store/activity.proto`**: Add `ActivityTypeMemoMention`.
-   **`proto/store/inbox.proto`**: Add `InboxMessage.Type.MENTION`.
-   **`store/activity.go`**: Add `ActivityTypeMemoMention` constant.

#### Option B: Zero-Dependency (Reuse Existing)
-   Reuse `ActivityTypeMemoComment` and `InboxMessage.Type.MEMO_COMMENT`.
-   **Pros**: No `protoc` needed.
-   **Cons**: UI calls it a "Comment", but technically a mention in a comment IS a comment notification. A mention in a main ticket body might look weird ("Alice commented on Ticket..." when Alice created the ticket? No, "Alice created a ticket..."?).
-   **Recommendation**: Use Option A if possible. If strict environment limitations exist, use Option B but be aware of UI labelling.

### 3. Implementation Steps

#### Step 1: Implement Parsing Helper
Create `ExtractMentions(content string) ([]string)` helper.

#### Step 2: Implement User Resolution
Query users by nickname.

#### Step 3: Trigger Notifications
In `CreateMemo` and `CreateMemoComment`:
```go
// After saving memo...
nicknames := ExtractMentions(memo.Content)
for _, nickname := range nicknames {
    user := FindUserByNickname(nickname)
    if user != nil && user.ID != currentUser.ID {
        // Create Activity & Inbox
        // ...
    }
}
```

### 4. Ticket Integration
Since Tickets use Memos for descriptions:
-   **Ticket Creation**: The description is a Memo link (e.g., `/m/uid`). The actual content is in the Memo. The `CreateTicket` flow currently doesn't create the memo content, the frontend does via `MemoEditor`.
-   **Flow**:
    1.  User types description in `MemoEditor`.
    2.  `MemoEditor` calls `POST /api/v1/memos`.
    3.  **Our Hook catches it here**.
    4.  Notification sent.
    5.  User saves Ticket linking to that Memo.

**Result**: It works out of the box if we implement it in `CreateMemo`.

## Code Example (Conceptual)

```go
func (s *APIV1Service) DispatchMentions(ctx context.Context, memo *v1pb.Memo) {
    mentions := parseMentions(memo.Content)
    for _, nick := range mentions {
        targetUser, err := s.Store.GetUser(ctx, &store.FindUser{Nickname: &nick})
        if targetUser == nil { continue }
        
        // Create Activity
        activity, _ := s.Store.CreateActivity(...)
        
        // Create Inbox
        s.Store.CreateInbox(ctx, &store.Inbox{
            ReceiverID: targetUser.ID,
            Message: &storepb.InboxMessage{
                Type: storepb.InboxMessage_MEMO_COMMENT, // or MENTION
                ActivityId: &activity.ID,
            },
        })
    }
}
```
