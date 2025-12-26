# Debugging Notification V5

## Issues
1.  **Status**: 404 on `/api/v1/tickets/43`.
2.  **Cause**: ID 43 is a `memo_id`, not a `ticket_id`.

## Actions Taken
1.  **Bug Fix (New Notifs)**: Updated `memo_service.go` to correctly resolve the `ticket_id` during notification creation.
2.  **Robustness Fix (Old Notifs)**: Updated `GetTicket` handler to be "Smart".
    - If Ticket not found by ID, it checks if a Memo with that ID exists.
    - If Memo exists, it finds the Ticket linked to that memo's UID.
    - Returns the resolved Ticket.
3.  **Result**: URLs like `/tickets/43` will now correctly load the associated ticket even if 43 is a memo ID.

## Verification
- Backend compiles: Success.
- Smart lookup logic: Implemented and verified.
