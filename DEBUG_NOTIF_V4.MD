# Debugging Notification V4

## Issues
1.  **Display**: Notifications show `Initiator: [int]` instead of `Initiator: [username]`.
2.  **Routing**: Clicking ticket links results in a 404.

## Investigation findings
1.  **Display**: The backend was only sending the ID. I updated the response struct to include `InitiatorName` and populated it by fetching the user during conversion.
2.  **Routing**: Found that in `v1.go`, the gRPC-gateway was registered with a wildcard `/api/v1/*` *before* the custom Echo routes for tickets and notifications. This caused standard Echo requests to `/api/v1/tickets/:id` to be incorrectly handled by the gateway handler, resulting in a 404.

## Changes Applied
1.  **Backend**:
    - Updated `NotificationResponse` in `notification_service.go`.
    - Updated `convertNotificationFromStore` to fetch the username/nickname.
    - **CRITICAL**: Reordered route registration in `v1.go` to put Echo routes before the gRPC-gateway wildcard.
2.  **Frontend**:
    - Updated `Notification` type in `types/notification.ts`.
    - Updated `Notifications.tsx` to display `initiatorName`.

## Verification
- Backend compiles: Success.
- Logic verified: Yes.
