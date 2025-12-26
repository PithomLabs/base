# Notifications System Walkthrough

I have implemented the notification system for ticket mentions.

## Changes

### Database
- Added `notifications` table migration for SQLite, MySQL, and Postgres.
- Fields: `id`, `initiator_id`, `receiver_id`, `ticket_url`, `created_ts`, `is_read`.

### Backend
- **Store**: Added `Notification` model and store logic (`store/notification.go`, `store/db/*/notification.go`).
- **Service**: 
    - Created `server/router/api/v1/notification_service.go` with `ListNotifications` and `UpdateNotification` (PATCH).
    - Updated `TicketService` and `MemoService` to create `Notification` entries when users are mentioned.
    - Updated `server/router/api/v1/v1.go` to register routes and `stream`.
- **SSE**: Reused existing SSE stream to push updates.

### Frontend
- **Types**: Added `web/src/types/notification.ts`.
- **Store**: Updated `web/src/store/v2/user.ts` to manage `notifications` state and listen to SSE stream for updates.
- **Pages**: Rewrote `web/src/pages/Notifications.tsx` to display notifications in a table format.
- **Components**: Updated `web/src/components/Navigation.tsx` to show unread notification badge on the bell icon.

## Verification

To verify the changes:

1. **Start the server**: `go run ./server/cmd/server/main.go`
2. **Frontend**: The frontend should be rebuilt or run in dev mode.
3. **Workflow**:
    - User A creates a ticket and mentions `@UserB` in description.
    - Use another browser/incognito window to login as User B.
    - User B should see a red badge on the "Notifications" link in sidebar.
    - Clicking "Notifications" shows a table with the mention details.
    - View the table: Shows Initiator, Ticket URL link, Date, and Status.
    - Click the checkmark icon to mark as read. Badge count should update (decrement).

## Automated Verification Script
A script `test_notification.sh` (if you were to create one) would:
1. Login as User A.
2. Create Ticket with `@UserB`.
3. Login as UserB.
4. `GET /api/v1/notifications` -> verify entry.
5. `PATCH /api/v1/notifications/:id` -> verify `isRead: true`.
