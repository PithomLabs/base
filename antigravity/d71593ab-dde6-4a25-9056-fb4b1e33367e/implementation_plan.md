# Mention Notification System Implementation Plan

## Goal Description
Implement a simplified notification system for user mentions in tickets and memos. This involves a new database table, backend logic to populate it, and a frontend update to display these notifications and show an unread count.

## Proposed Changes

### Database
#### [NEW] `store/migration/mysql/LATEST__notifications.sql` (and other dialects)
- Create `notifications` table with columns:
  - `id` (Primary Key)
  - `initiator_id` (INT)
  - `receiver_id` (INT)
  - `ticket_url` (TEXT)
  - `created_ts` (BIGINT)
  - `is_read` (BOOLEAN)

### Backend
#### [NEW] `store/notification.go`
- Define `Notification` struct.
- Define `FindNotification` and `UpdateNotification` structs.
- Implement `CreateNotification`, `ListNotifications`, `UpdateNotification` methods.

#### [NEW] `server/router/api/v1/notification_service.go`
- Implement `RegisterNotificationRoutes`.
- Implement `ListNotifications` handler.
- Implement `PatchNotification` handler (for marking as read).

#### [MODIFY] `server/router/api/v1/v1.go`
- Register notification routes.

#### [MODIFY] `server/router/api/v1/ticket_service.go`
- Update `dispatchTicketMentions` to create `Notification` entries instead of `Inbox`.
- Trigger SSE notification.

#### [MODIFY] `server/router/api/v1/memo_service.go`
- Update `dispatchMemoMentions` to create `Notification` entries.

#### [MODIFY] `server/router/api/v1/notification_hub.go`
- Update `Notification` struct to align with store model if needed or keep using view model.
- Ensure SSE sends correct data structure for the new frontend listener.

### Frontend
#### [MODIFY] `web/src/types/proto/api/v1/notification_service.ts`
- _(Actually, we likely use plain JSON types or manually defined types since proto gen isn't fully automated in this environment/user instructions might prefer manual typescript definitions)_ -> Create `web/src/types/notification.ts`.

#### [MODIFY] `web/src/pages/Notifications.tsx`
- Replace `Inbox` logic with `Notification` table list.
- columns: Initiator, Ticket URL (link), Time, Status (Read/Unread).

#### [MODIFY] `web/src/components/Navigation.tsx`
- Add unread count badge to Notifications icon.
- Subscribe to SSE events to update this count real-time.

## Verification Plan

### Automated Tests
- Create a test script `test_notification.sh` using `curl` to:
    1. Create a ticket with a mention.
    2. Check if notification exists via API.
    3. Mark notification as read.

### Manual Verification
1. **Mention Flow**:
    - Log in as User A.
    - Create a ticket and mention User B (e.g. `@userb`).
    - Log in as User B.
    - Check Notifications page -> Should see entry.
    - Check Sidebar -> Should see unread badge.
    - Click notification -> Should go to ticket.
    - Badge count should decrease (if logic marks read on click/visit).
