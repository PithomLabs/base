# SSE Notification System Walkthrough

## Summary

Implemented a real-time notification system using **Server-Sent Events (SSE)** with the [datastar-go](https://github.com/starfederation/datastar-go) SDK. When a user is @mentioned in a memo comment, they receive an **instant toast popup** notification without refreshing the page.

## Changes Made

### Backend

| File | Change |
|------|--------|
| [notification_hub.go](file:///home/chaschel/Documents/ibm/go/tix-gemini-master/server/router/api/v1/notification_hub.go) | **NEW** - SSE connection manager, notification renderer, and stream handler |
| [v1.go](file:///home/chaschel/Documents/ibm/go/tix-gemini-master/server/router/api/v1/v1.go) | Added SSE endpoint `/api/v1/notifications/stream` |
| [memo_service.go](file:///home/chaschel/Documents/ibm/go/tix-gemini-master/server/router/api/v1/memo_service.go) | Integrated `NotifyUser()` call in `dispatchMemoMentions` |

### Frontend

| File | Change |
|------|--------|
| [index.html](file:///home/chaschel/Documents/ibm/go/tix-gemini-master/web/index.html) | Added Datastar script, SSE container, toast CSS |

## How It Works

```
User A mentions @UserB in memo    SSE Connection    Toast Popup
          │                            │              (top-right)
          ▼                            │                  │
   dispatchMemoMentions ──────────────►│                  │
          │                            │                  │
   NotificationHub.NotifyUser() ──────►│ ─────────────────►
          │                            │
   sse.PatchElements(toastHTML) ──────►│
```

## Testing

1. **Restart backend**: `./bin/memos`
2. Open two browser tabs logged in as different users (e.g., `ibm2100` and `ading`)
3. As `ibm2100`, create a memo comment with `@ading`
4. **Expected**: `ading's` browser shows a blue toast popup in top-right:
   > 🔔 ibm2100 mentioned you in a memo
5. Click toast to navigate to `/notifications`

## Design Doc

See [SSE_NOTIFICATION_DESIGN.MD](file:///home/chaschel/Documents/ibm/go/tix-gemini-master/SSE_NOTIFICATION_DESIGN.MD) for architecture details.
