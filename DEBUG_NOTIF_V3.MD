# Debugging Notification V3

## Issue
User reports "page not found 404" when clicking a ticket link (`/tickets/40`) from notifications.

## Investigation findings
1.  **Backend**: `server/router/api/v1/ticket_service.go` was MISSING `GET /tickets/:id`.
2.  **Frontend**: `web/src/router/index.tsx` was MISSING a route definition for `/tickets/:id` and had no detailed view component.

## Changes Applied
1.  **Backend**:
    - Implemented `GetTicket` handler in `server/router/api/v1/ticket_service.go` (appended to end).
    - Registered `g.GET("/tickets/:id", s.GetTicket)` in `RegisterTicketRoutes`.
    - Verified compilation.
2.  **Frontend**:
    - Created `web/src/pages/TicketDetail.tsx` to fetch and display single ticket.
    - Added `/tickets/:id` route in `web/src/router/index.tsx` mapping to `TicketDetail`.

## Verification
- Backend compiles: Success.
- Frontend code exists: Success.
