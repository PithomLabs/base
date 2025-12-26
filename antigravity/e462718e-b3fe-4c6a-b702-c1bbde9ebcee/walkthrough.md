# Walkthrough - Enforce Memo for Ticket Creation & Mention Notifications

## Enforce Memo for Ticket Creation

I have implemented the validation rule requiring tickets to have a memo description and improved the error messaging.

### Changes

#### Backend
**File**: `store/ticket.go`
- Added validation logic in `Validate()`:
```go
	if len(t.Description) < 3 || t.Description[:3] != "/m/" {
		return errors.New("description must be a valid memo link starting with /m/")
	}
```

#### Frontend
**File**: `web/src/pages/Tickets.tsx`
- Added pre-submission check validation.
- **Improved Error Handling**: Updated the API response handling to display the actual server error message instead of a generic failure.

### Verification Results
Tests passed for valid/invalid descriptions (`store/test/ticket_test.go`).

---

## Implement Mention Notifications

I have implemented the system to notify users when they are mentioned (`@nickname`) in tickets or comments.

### Changes

#### Backend
**File**: `server/router/api/v1/memo_service.go`
- Added `dispatchMemoMentions` helper function that uses regex/parsing to find `@nickname` mentions.
- Integrated this into `CreateMemo` (triggers when Ticket Description is created) and `CreateMemoComment` (triggers when commenting on a ticket).
- **Strategy**: Reuses the existing `MEMO_COMMENT` activity/inbox type to avoid protocol buffer changes ("Zero-Dependency").

#### Frontend
**File**: `web/src/components/Inbox/MemoCommentMessage.tsx`
- **Smart Redirection**: When clicking a notification, the frontend now checks if the related memo is actually a **Ticket Description**.
- If it is, the user is redirected to the **Ticket View** (`/tickets?id=X`) instead of the Memo View.

### Verification Results

#### Manual Verification Plan
1.  **Ticket Creation Mention**: User A creates a ticket with description `@UserB check this`. User B receives a notification. Clicking it opens the Ticket.
2.  **Comment Mention**: User A comments `@UserB` on a ticket. User B receives a notification. Clicking it opens the Ticket (via the description check).
